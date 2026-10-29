package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	authenticationv1 "github.com/example/identity-access-service/gen/proto/authentication/v1"
	authorizationv1 "github.com/example/identity-access-service/gen/proto/authorization/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

const (
	defaultHTTPAddress = ":8080"
	defaultGRPCAddress = "localhost:50051"
)

type transaction struct {
	ID       string  `json:"id"`
	OwnerID  string  `json:"owner_id"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	Status   string  `json:"status"`
}

type app struct {
	authentication authenticationv1.AuthenticationServiceClient
	authorization  authorizationv1.AuthorizationServiceClient
	transactions   map[string]transaction
}

func main() {
	grpcAddress := getenv("IDENTITY_ACCESS_GRPC_ADDRESS", defaultGRPCAddress)
	httpAddress := getenv("HTTP_ADDRESS", defaultHTTPAddress)

	conn, err := grpc.NewClient(
		grpcAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("create gRPC client: %v", err)
	}
	defer conn.Close()

	application := &app{
		authentication: authenticationv1.NewAuthenticationServiceClient(conn),
		authorization:  authorizationv1.NewAuthorizationServiceClient(conn),
		transactions: map[string]transaction{
			"tx-001": {
				ID:       "tx-001",
				OwnerID:  "992e01e4-cf93-4052-ab73-0d87719420e1",
				Amount:   125.50,
				Currency: "EUR",
				Status:   "COMPLETED",
			},
			"tx-002": {
				ID:       "tx-002",
				OwnerID:  "another-user",
				Amount:   42.25,
				Currency: "EUR",
				Status:   "PENDING",
			},
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /transactions/{id}", application.getTransaction)
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	server := &http.Server{
		Addr:              httpAddress,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("example HTTP client listening on %s", httpAddress)
	log.Printf("Identity & Access gRPC address: %s", grpcAddress)

	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

func (a *app) getTransaction(w http.ResponseWriter, r *http.Request) {
	token, err := bearerToken(r.Header.Get("Authorization"))
	if err != nil {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	validateResponse, err := a.authentication.ValidateToken(
		ctx,
		&authenticationv1.ValidateTokenRequest{Token: token},
	)
	if err != nil {
		handleAuthenticationError(w, err)
		return
	}

	principal := validateResponse.GetPrincipal()
	if principal == nil {
		writeError(w, http.StatusUnauthorized, "token did not resolve to a principal")
		return
	}

	transactionID := r.PathValue("id")
	tx, ok := a.transactions[transactionID]
	if !ok {
		writeError(w, http.StatusNotFound, "transaction not found")
		return
	}

	authorizeResponse, err := a.authorization.AuthorizeTransaction(
		ctx,
		&authorizationv1.AuthorizeTransactionRequest{
			Principal: principal,
			Action:    authorizationv1.TransactionAction_READ,
			Resource: &authorizationv1.TransactionResource{
				Id:      tx.ID,
				OwnerId: tx.OwnerID,
			},
		},
	)
	if err != nil {
		log.Printf("authorize transaction: %v", err)
		writeError(w, http.StatusInternalServerError, "authorization service failed")
		return
	}

	if !authorizeResponse.GetAllowed() {
		writeJSON(w, http.StatusForbidden, map[string]string{
			"error":  "forbidden",
			"reason": authorizeResponse.GetReason(),
		})
		return
	}

	writeJSON(w, http.StatusOK, tx)
}

func bearerToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("missing Authorization header")
	}

	parts := strings.Fields(header)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return "", errors.New("Authorization header must be: Bearer <token>")
	}

	if strings.TrimSpace(parts[1]) == "" {
		return "", errors.New("bearer token is empty")
	}

	return parts[1], nil
}

func handleAuthenticationError(w http.ResponseWriter, err error) {
	switch status.Code(err) {
	case codes.Unauthenticated, codes.InvalidArgument:
		writeError(w, http.StatusUnauthorized, "invalid access token")
	default:
		log.Printf("validate token: %v", err)
		writeError(w, http.StatusInternalServerError, "authentication service failed")
	}
}

func writeError(w http.ResponseWriter, statusCode int, message string) {
	writeJSON(w, statusCode, map[string]string{"error": message})
}

func writeJSON(w http.ResponseWriter, statusCode int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(value); err != nil {
		log.Printf("encode HTTP response: %v", err)
	}
}

func getenv(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}
