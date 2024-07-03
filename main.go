package main

import (
	"aurora/initialize"
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/acheong08/endless"
	"github.com/joho/godotenv"

	"fmt"
    	"net/http"
    	"time"
)
func keepAwake() {
    for {
        resp, err := http.Get("https://duck2api-xnk4.onrender.com")//need deploy to render and get it's URL first
        if err != nil {
            fmt.Printf("Failed to send keep awake request: %v\n", err)
        } else {
            fmt.Println("Keep awake request sent. Status:", resp.Status)
            resp.Body.Close()
        }
        time.Sleep(10 * time.Minute)
    }
}
//go:embed web/*
var staticFiles embed.FS

func main() {
	go keepAwake()
	_ = godotenv.Load(".env")
	gin.SetMode(gin.ReleaseMode)
	router := initialize.RegisterRouter()
	subFS, err := fs.Sub(staticFiles, "web")
	if err != nil {
		log.Fatal(err)
	}
	router.StaticFS("/web", http.FS(subFS))
	host := os.Getenv("SERVER_HOST")
	port := os.Getenv("SERVER_PORT")
	tlsCert := os.Getenv("TLS_CERT")
	tlsKey := os.Getenv("TLS_KEY")

	if host == "" {
		host = "0.0.0.0"
	}
	if port == "" {
		port = os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
	}

	if tlsCert != "" && tlsKey != "" {
		_ = endless.ListenAndServeTLS(host+":"+port, tlsCert, tlsKey, router)
	} else {
		_ = endless.ListenAndServe(host+":"+port, router)
	}
}
