package config

import (
	"github.com/spf13/viper"
	"github.com/utkarsh-pro/notion-gister/pkg/utils"
)

const EnvPrefix = "GISTER"

func Setup() {
	viper.SetConfigName(".gister")
	viper.AddConfigPath(".")

	viper.SetEnvPrefix(EnvPrefix)
	viper.BindEnv("mail.smtp.host", prefixed("SMTP_HOST"))
	viper.BindEnv("mail.smtp.port", prefixed("SMTP_PORT"))
	viper.BindEnv("mail.smtp.username", prefixed("SMTP_USERNAME"))
	viper.BindEnv("mail.smtp.password", prefixed("SMTP_PASSWORD"))
	viper.AutomaticEnv()

	utils.PanicIfError(viper.ReadInConfig())
}

func prefixed(key string) string {
	return EnvPrefix + "_" + key
}
