package synk

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	log "github.com/mneira10/synk/logger"

	"github.com/spf13/viper"
)

// Currently only supports Cloudflare's R2
type R2ConfigData struct {
	Type            string `validate:"required"`
	BucketName      string `validate:"required"`
	Url             string `validate:"required"`
	AccountId       string `validate:"required"`
	AccessKeyId     string `validate:"required"`
	AccessKeySecret string `validate:"required"`
}

func GetConfiguration(path string) R2ConfigData {
	log.WithFields(log.Fields{"configPath": path}).Info("Getting configuration...")
	viper.SetConfigName(".synk.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Could not find .synk.yaml config file in:", path)
			os.Exit(1)
		} else {
			fmt.Println("Could not read the .synk.yaml config file.")
			os.Exit(1)
		}
	}

	log.Info("Unmarshalling data...")

	var config R2ConfigData

	err := viper.Unmarshal(&config)
	if err != nil {
		log.Fatal("unable to decode into struct, %v", err)
		os.Exit(1)
	}

	log.Info("Validating struct data...")

	validate := validator.New()
	if err := validate.Struct(&config); err != nil {
		log.WithFields(log.Fields{"validateErrMsg": err}).Fatal("Missing required yaml attributes")
		os.Exit(1)
	}

	return config

}
