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

	dB := db.InitDB()
	defer dB.Close()

	appRouter := router.BuildRouter(dB)

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
