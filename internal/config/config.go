package config

import (
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	ListenPort  int    `mapstructure:"listen_port" yaml:"listen_port" json:"-"`
	Timeout     int    `json:"-"`
	Version     string `json:"version"`
	Environment string `json:"-" mapstructure:"environment"`
	LogLevel    string `mapstructure:"log_level" yaml:"log_level" json:"-"`

	BaseURL string `mapstructure:"base_url" yaml:"base_url" json:"-"`

	// DB Settings
	RedisAddress string `json:"-" mapstructure:"redis_address"`
	// DB Settings
	DBUserName    string `json:"-" mapstructure:"db_username"`
	DBPassword    string `json:"-" mapstructure:"db_password"`
	DBHost        string `json:"-" mapstructure:"db_host"`
	DBName        string `json:"-" mapstructure:"db_name"`
	DBAutoMigrate bool   `json:"-" mapstructure:"db_auto_migrate"`

	// MySQLCatalog Settings
	// This should always be set in a secret
	MySQLCatalogDBPassword string `json:"-" mapstructure:"mysql_catalog_db_password"`

	// General Auth
	AuthSessionHandlerKey string `json:"-" mapstructure:"auth_session_handler_key"`

	// Auth settings
	AuthIDP             string `json:"-" mapstructure:"auth_idp"`
	OIDCIssuer          string `json:"-" mapstructure:"oidc_issuer"`
	OIDCRedirectURI     string `json:"-" mapstructure:"oidc_redirect_uri"`
	OIDCClientID        string `json:"-" mapstructure:"oidc_client_id"`
	OIDCCLientTLSVerify bool   `json:"-" mapstructure:"oidc_client_tls_verify"`
	OIDCAudience        string `json:"-" mapstructure:"oidc_audience"`

	// Kubernetes settings
	K8sInCluster               bool `json:"-" mapstructure:"k8s_in_cluster"`
	K8sDataSinkIntervalSeconds int  `json:"-" mapstructure:"k8s_data_sink_interval_seconds"`

	// AWS settings
	AWSRegion     string `json:"-" mapstructure:"aws_region"`
	ReportsBucket string `json:"-" mapstructure:"reports_bucket"`
}

func (c Config) IsProduction() bool {
	return strings.ToUpper(c.Environment) == "PRODUCTION"
}

func Load(version string, cfgFile string) *Config {
	// SET CONFIG DEFAULTS
	c := &Config{
		Environment:                "Development",
		Version:                    version,
		ListenPort:                 8080,
		Timeout:                    2000,
		BaseURL:                    "http://localhost:3000",
		AuthSessionHandlerKey:      "auth-session",
		OIDCIssuer:                 "{REPLACE_ME}",
		OIDCRedirectURI:            "http://localhost:8080/authorization-code/callback",
		OIDCClientID:               "{REPLACE_ME}",
		OIDCCLientTLSVerify:        false, // Zitadel cloud's self-signed cert is not trusted by default, for example
		OIDCAudience:               "{REPLACE_ME}",
		K8sInCluster:               true,
		RedisAddress:               "redis-master.redis:6379",
		DBUserName:                 "postgres",
		DBPassword:                 "postgres1011",
		DBHost:                     "postgres-postgresql.default.svc.cluster.local",
		DBName:                     "khub",
		DBAutoMigrate:              true,
		K8sDataSinkIntervalSeconds: 5,
		MySQLCatalogDBPassword:     "khub1011",
	}

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
	}
	viper.SetEnvPrefix("khub")
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	_ = viper.BindEnv("LISTEN_PORT")
	_ = viper.BindEnv("TIMEOUT")
	_ = viper.BindEnv("ENVIRONMENT")
	_ = viper.BindEnv("LOG_LEVEL")
	_ = viper.BindEnv("BASE_URL")
	_ = viper.BindEnv("AUTH_SESSION_HANDLER_KEY")
	_ = viper.BindEnv("OIDC_ISSUER")
	_ = viper.BindEnv("OIDC_REDIRECT_URI")
	_ = viper.BindEnv("OIDC_CLIENT_ID")
	_ = viper.BindEnv("OIDC_CLIENT_SECRET")
	_ = viper.BindEnv("OIDC_AUDIENCE")
	_ = viper.BindEnv("REDIS_ADDRESS")
	_ = viper.BindEnv("DB_USERNAME")
	_ = viper.BindEnv("DB_PASSWORD")
	_ = viper.BindEnv("DB_HOST")
	_ = viper.BindEnv("DB_NAME")
	_ = viper.BindEnv("DB_AUTO_MIGRATE")
	_ = viper.BindEnv("K8S_IN_CLUSTER")
	_ = viper.BindEnv("K8S_DATA_SYNC_INTERVAL_SECONDS")
	_ = viper.BindEnv("AWS_REGION")
	_ = viper.BindEnv("REPORTS_BUCKET")
	_ = viper.BindEnv("MYSQL_CATALOG_DB_PASSWORD")

	_ = viper.ReadInConfig()
	viper.AutomaticEnv()
	_ = viper.Unmarshal(&c)
	return c
}
