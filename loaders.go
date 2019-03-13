package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"os"
	"trustwallet.com/blockatlas/platform/binance"
)

var loaders = map[string]func(gin.IRouter){
	"binance": binance.Setup,
}

func loadPlatforms(router gin.IRouter) {
	enabled := viper.GetStringSlice("platforms")

	for _, ns := range enabled {
		loader := loaders[ns]
		if loader == nil {
			fmt.Fprintf(os.Stderr, "Platform does not exist: %s\n", ns)
			os.Exit(1)
		}

		loader(router)
		fmt.Printf("Loaded /%s\n", ns)
	}
}
