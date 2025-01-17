package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"todo/internal/db"
	"todo/internal/web/router"

	"github.com/joho/godotenv"
)

const defaultPort = "8088"

func main() {

	err := godotenv.Load()
	if err != nil {
		panic(fmt.Errorf("load env file err: %v", err))
	}

	fmt.Printf("selected environment = [%v]", os.Getenv("ENVIRONMENT"))

	dB1 := db.InitDB()
	defer dB1.Close()

	dB2 := db.InitDB()
	defer dB2.Close()

	appRouter := router.BuildRouter(dB1, dB2)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: appRouter,
	}

	fmt.Printf("web-api listening on :%v", port)

	if err := server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			fmt.Println("Server shut down")
		} else {
			log.Fatal("Server shut down unexpectedly!")
		}
	}
}
