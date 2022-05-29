package src

import (
	"log"
	"os"
	"strings"

	"github.com/Koubae/goweb/src/config"
	"github.com/Koubae/goweb/src/mvc/controllers"
)

func AppInit() {
	log.Println("-------------------------- Initialzing App ...")
	var arrow []string = []string{"="}
	log.Printf("* INIT %s> Seting App Environment Variable\n", strings.Join(arrow, ""))
	err := config.Env()
	if err != nil {
		panic(err)
	}
	arrow = append(arrow, "=")
	log.Printf("* INIT %s> Seting App Environment Variable", strings.Join(arrow, ""))
}

func AppRun() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"

	}

	r := controllers.Routes()
	r.Run(":" + port)
}
