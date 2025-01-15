package main

import (
	"DeliFood/backend/handlers"
	"DeliFood/backend/pkg/logger"
	"context"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func MenuPageHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("frontend/menu.html")
	if err != nil {
		log.Fatal(err)
	}
	t.ExecuteTemplate(w, "menu.html", nil)
}

func main() {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize logger
	logger := logger.NewLogger()

	// Log application start
	logger.Info("Application started", map[string]interface{}{
		"module": "main",
		"status": "success",
	})

	//working with server side
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./frontend/assets/"))))

	http.HandleFunc("/", handlers.MainPageHandler)
	http.HandleFunc("/contact", handlers.ContactUsHandler)
	http.HandleFunc("/menu", handlers.MenuHandler)

	// Start HTTP server
	server := &http.Server{
		Addr: ":9078",
	}

	go func() {
		log.Println("Listening on :9078")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not start server: %s\n", err)
		}
	}()

	//Gracefully shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan
	log.Println("Shutting down the server!")

	cancel()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	err := server.Shutdown(shutdownCtx)
	if err != nil {
		logger.Error("Could not shut down the server!!!", map[string]interface{}{"Error:": err})
	}

	logger.Warn("Server gracefully stopped", map[string]interface{}{"Module": "main"})
}
