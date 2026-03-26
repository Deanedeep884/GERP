package cli

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Config binds the GERP matrix operational endpoints.
type Config struct {
	GraphQLEndpoint string `mapstructure:"graphql_endpoint"`
	TemporalHost    string `mapstructure:"temporal_host"`
	SpannerDB       string `mapstructure:"spanner_db"`
}

var ActiveConfig Config

// InitConfig directs Viper to bind to the localized user system configurations or env overrides.
func InitConfig() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("🚨 CRITICAL: Root matrix bounds failed, cannot resolve home directory:", err)
		os.Exit(1)
	}

	viper.AddConfigPath(home)
	viper.SetConfigType("yaml")
	viper.SetConfigName(".gerp")

	viper.SetEnvPrefix("GERP")
	viper.AutomaticEnv()

	// Establish sane defaults for local sandbox development
	viper.SetDefault("graphql_endpoint", "http://localhost:8080/query")
	viper.SetDefault("temporal_host", "localhost:7233")
	viper.SetDefault("spanner_db", "projects/gerp-local-dev/instances/gerp-instance/databases/gerp-db")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Println("🚨 CRITICAL: Matrix configuration binding error:", err)
		}
	}

	if err := viper.Unmarshal(&ActiveConfig); err != nil {
		fmt.Println("🚨 CRITICAL: Matrix configuration deserialization fault:", err)
		os.Exit(1)
	}
}
