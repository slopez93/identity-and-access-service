package main

import (
	"context"
	"fmt"
	"log"
	"net"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"google.golang.org/grpc"

	authenticationv1 "github.com/example/identity-access-service/gen/proto/authentication/v1"
	authorizationv1 "github.com/example/identity-access-service/gen/proto/authorization/v1"
	authapp "github.com/example/identity-access-service/internal/authentication/application"
	authcognito "github.com/example/identity-access-service/internal/authentication/infrastructure/cognito"
	authgrpc "github.com/example/identity-access-service/internal/authentication/infrastructure/controllers/grpc"
	authzapp "github.com/example/identity-access-service/internal/authorization/application"
	authzconfig "github.com/example/identity-access-service/internal/authorization/infrastructure/config"
	authzgrpc "github.com/example/identity-access-service/internal/authorization/infrastructure/controllers/grpc"
	authzpostgres "github.com/example/identity-access-service/internal/authorization/infrastructure/postgres"
	sharedconfig "github.com/example/identity-access-service/internal/shared/config"
	"github.com/example/identity-access-service/internal/shared/database"
)

func main() {
	ctx := context.Background()

	cfg, err := sharedconfig.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.NewPostgresPool(ctx, cfg.Database)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	permissionRepo := authzpostgres.NewPermissionRepository(db)
	principalPermissionRepo := authzpostgres.NewPermissionChecker(db)
	requiredScopes := authzconfig.NewRequiredScopeProvider(cfg.AuthorizationScopes)
	authorizeTransaction := authzapp.NewAuthorizeTransactionUseCase(requiredScopes, permissionRepo, principalPermissionRepo)
	authorizationController := authzgrpc.NewController(authorizeTransaction)

	serviceClientIDs := make(map[string]string, len(cfg.Cognito.ServiceClients))
	serviceClients := make(map[string]authcognito.ServiceClient, len(cfg.Cognito.ServiceClients))
	for serviceID, service := range cfg.Cognito.ServiceClients {
		serviceClientIDs[service.ClientID] = serviceID
		serviceClients[serviceID] = authcognito.ServiceClient{ClientID: service.ClientID, ClientSecret: service.ClientSecret}
	}

	tokenVerifier, err := authcognito.NewTokenVerifier(ctx, authcognito.TokenVerifierConfig{
		Issuer:           cfg.Cognito.Issuer,
		UserClientID:     cfg.Cognito.UserClientID,
		ServiceClientIDs: serviceClientIDs,
	})
	if err != nil {
		log.Fatal(err)
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(cfg.AWS.Region))
	if err != nil {
		log.Fatal(err)
	}
	cognitoClient := cognitoidentityprovider.NewFromConfig(awsCfg)
	serviceTokenProvider := authcognito.NewServiceTokenProvider(cognitoClient, serviceClients)
	validateToken := authapp.NewValidateTokenUseCase(tokenVerifier)
	getServiceToken := authapp.NewGetServiceTokenUseCase(serviceTokenProvider)
	authenticationController := authgrpc.NewController(validateToken, getServiceToken)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.GRPCPort))
	if err != nil {
		log.Fatal(err)
	}

	server := grpc.NewServer()
	authenticationv1.RegisterAuthenticationServiceServer(server, authenticationController)
	authorizationv1.RegisterAuthorizationServiceServer(server, authorizationController)

	log.Printf("gRPC server listening on %s", listener.Addr())
	if err := server.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
