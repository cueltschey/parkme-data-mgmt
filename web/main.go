package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	var (
		port       = flag.Int("port", 8080, "Port to listen on")
		dbPath     = flag.String("db", "./sbc.db", "Path to SQLite database file")
		staticPath = flag.String("static", "../web", "Path to web files")
	)

	flag.Parse()

	listenPort := *port
	listenAddr := ":" + strconv.Itoa(listenPort)

	dbConn, err := db.Open(*dbPath)
	if err != nil {
		log.Fatalf("❌ failed to open DB: %v", err)
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

  // insecure login endpoint
	router.POST("/login", api.LoginHandler(dbConn))

  // Requests to get information
	router.POST("/search", api.AuthMiddleware(), api.Search(dbConn))
	router.POST("/client", api.AuthMiddleware(), api.GetClientHandler(dbConn))
	router.POST("/agent", api.AuthMiddleware(), api.GetAgentHandler(dbConn))
	router.POST("/obligee", api.AuthMiddleware(), api.GetObligeeHandler(dbConn))
	router.POST("/overview", api.AuthMiddleware(), api.GetOverviewHandler(dbConn))
	router.POST("/details", api.AuthMiddleware(), api.GetDetailsHandler(dbConn))
	router.POST("/notes", api.AuthMiddleware(), api.GetNotesHandler(dbConn))

  // static file server
	router.GET("/*filepath", api.StaticFileHandler(*staticPath, gin.Dir(*staticPath, false)))

	srv := &http.Server{
		Addr:    listenAddr,
		Handler: router,
	}

	go func() {
		var err error
		log.Printf("🚀 Starting HTTP server on %s", listenAddr)
		err = srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Forced shutdown: %v", err)
	}
	log.Println("✅ Server stopped cleanly")
}

