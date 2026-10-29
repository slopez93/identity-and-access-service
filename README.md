# Identity & Access Service — Go

Example implementation of the architecture discussed for a financial platform.

## Modules

- `authentication`: validates Cognito access tokens locally through Cognito JWKS and centralizes M2M token issuance with Cognito `GetClientToken`.
- `authorization`: typed `AuthorizeTransaction` contract, RBAC in PostgreSQL, required OAuth scopes from service configuration, and optional ownership checks.

The code intentionally keeps the modules self-contained. gRPC controllers live under each module at `infrastructure/controllers/grpc`.

## Flow

```text
Authentication

gRPC Controller
      -> ValidateTokenUseCase / GetServiceTokenUseCase
      -> domain/interfaces
      -> infrastructure/cognito

Authorization

gRPC Controller
      -> AuthorizeTransactionUseCase
      -> domain/interfaces + RequiredScopeProvider
      -> PostgreSQL / configuration
```

## Project structure

```text
api/proto/                         protobuf source of truth
cmd/server/                        composition root / dependency injection
internal/
  shared/
    config/                        config.go + loader.go
    database/                      pgx pool creation
    domain/                        Principal, Scope, PrincipalType
  authentication/
    domain/
      interfaces/
    application/
    infrastructure/
      controllers/grpc/
      cognito/
  authorization/
    domain/
      interfaces/
      services/
    application/
    infrastructure/
      controllers/grpc/
      config/
      postgres/
migrations/
```

## Configuration

Copy `.env.example` to `.env`.

Configuration is loaded with:

- `github.com/caarlos0/env/v11`
- `github.com/joho/godotenv`

`COGNITO_SERVICE_CLIENTS` is a generic JSON map, so adding a new M2M caller does not require changing `config.go`:

```json
{
  "payments-service": {
    "client_id": "...",
    "client_secret": "...",
    "allowed_scopes": ["financial-api/transactions.execute"]
  }
}
```

For production, inject secrets from your EKS secret-management mechanism rather than committing them.

## Protobuf generation

The `.proto` files under `api/proto/` define the gRPC contract, and generated Go stubs under `gen/proto/` are committed so Go clients can import them from this module. Install Buf to regenerate the stubs after changing a proto:

```bash
brew install bufbuild/buf/buf
make generate
```

`buf.gen.yaml` uses remote Buf plugins, so local `protoc-gen-go` binaries are not required.

## Run locally

```bash
cp .env.example .env
make db-up
make generate
go mod tidy
go test ./...
make run
```

## Authentication notes

`ValidateToken` does not call Cognito per request. `CognitoTokenVerifier` downloads the User Pool JWKS and `keyfunc` maintains the keys in memory and refreshes them. Validation includes:

- JWT signature using the JWKS public key
- issuer
- expiration
- `token_use == access`
- `client_id`

Scopes are read from Cognito's signed `scope` claim and normalized into `Principal.Scopes`.

For M2M, `GetServiceToken(service_id, scopes)` looks up the registered Cognito app client in configuration, checks the requested scopes against `allowed_scopes`, and calls Cognito `GetClientToken`.

In this first version `service_id` is supplied by the caller and is acceptable only under the explicit trust assumption that these gRPC endpoints are private inside EKS. A stronger version should derive the caller identity from workload identity / mTLS / service-account identity instead of trusting a free-form `service_id`.

## Authorization notes

The consumer does not send a free-form resource or permission string.

`AuthorizeTransaction` fixes the domain to `TRANSACTION`. Its protobuf enum is mapped internally to `READ`, `EXECUTE`, or `APPROVE`. The `Authorizer` then checks:

1. the permission `(resource, action)` exists;
2. the required scope configured for that permission is present in `Principal.Scopes`;
3. the principal has that permission through its roles;
4. ownership if the permission requires it.

Business rules such as account balance, daily limits or transaction state remain in the domain service that owns the transaction.
