package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Port         int    `mapstructure:"port"`
	HostKeyFile  string `mapstructure:"host_key_file"`
	AuthKeyFile  string `mapstructure:"auth_key_file"`
	PasswordHash string `mapstructure:"password_hash"`
}

var cfg Config

func init() {
	_ = godotenv.Load()

	v := viper.NewWithOptions(
		viper.ExperimentalBindStruct(),
		viper.ExperimentalFinder(),
	)

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	err := v.Unmarshal(&cfg)

	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	repo := NewServerRepository("./servers.json")

	server := NewSSHServer(
		fmt.Sprintf("0.0.0.0:%d", cfg.Port),
		cfg.HostKeyFile,
		cfg.AuthKeyFile,
		cfg.PasswordHash,
		repo,
	)

	if err := server.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
	fmt.Println("Server stopped")
}
