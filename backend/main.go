package main

import (
	"github.com/gagarin/backend/api"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	if err := viperSetup(); err != nil {
		zap.S().Panic(err)
	}

	// Setting up logger
	if viper.GetBool("dev") {
		zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))
	} else {
		zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
	}
	defer zap.L().Sync()

	api.CreateApi(viper.GetString("app.address"), viper.GetString("app.port")).ConfigureApp().Run()

}

func viperSetup() error {
	viper.SetConfigFile("app.json")
	viper.AddConfigPath(".")
	viper.SetDefault("app.port", "3033")
	viper.SetDefault("app.address", "")
	viper.SetDefault("dev", true)

	return viper.ReadInConfig()
}
