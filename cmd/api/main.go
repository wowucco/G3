package main

import (
	"github.com/spf13/viper"
	"github.com/wowucco/G3/config"
	"github.com/wowucco/G3/server"
	"log"
)

func main()  {
	if err := config.Init(); err != nil {
		log.Fatalf("%s", err.Error())
	}

	app := server.NewApp()

	if err := app.Run(viper.GetString("port")); err != nil {
		log.Fatalf("%s", err.Error())
	}
}
