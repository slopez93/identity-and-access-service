package config

type Config struct {
	Server              ServerConfig
	AWS                 AWSConfig
	Cognito             CognitoConfig
	Database            DatabaseConfig
	AuthorizationScopes AuthorizationScopesConfig
}

type ServerConfig struct {
	GRPCPort int `env:"GRPC_PORT" envDefault:"50051"`
}

type AWSConfig struct {
	Region string `env:"AWS_REGION,required"`
}

type CognitoConfig struct {
	Issuer       string `env:"COGNITO_ISSUER,required"`
	UserPoolID   string `env:"COGNITO_USER_POOL_ID,required"`
	UserClientID string `env:"COGNITO_USER_CLIENT_ID,required"`

	ServiceClients CognitoServiceClientsConfig `env:"COGNITO_SERVICE_CLIENTS"`
}

type CognitoServiceClientConfig struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type CognitoServiceClientsConfig map[string]CognitoServiceClientConfig

type DatabaseConfig struct {
	Host     string `env:"DB_HOST,required"`
	Port     int    `env:"DB_PORT" envDefault:"5432"`
	Name     string `env:"DB_NAME,required"`
	User     string `env:"DB_USER,required"`
	Password string `env:"DB_PASSWORD,required"`
	SSLMode  string `env:"DB_SSL_MODE" envDefault:"disable"`
}

type AuthorizationScopesConfig struct {
	TransactionRead    string `env:"TRANSACTION_READ_SCOPE,required"`
	TransactionExecute string `env:"TRANSACTION_EXECUTE_SCOPE,required"`
	TransactionCreate  string `env:"TRANSACTION_CREATE_SCOPE,required"`

	PaymentRead    string `env:"PAYMENT_READ_SCOPE,required"`
	PaymentExecute string `env:"PAYMENT_EXECUTE_SCOPE,required"`

	AccountRead   string `env:"ACCOUNT_READ_SCOPE,required"`
	AccountUpdate string `env:"ACCOUNT_UPDATE_SCOPE,required"`
}
