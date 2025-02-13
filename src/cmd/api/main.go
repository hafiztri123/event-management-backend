package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)


func main(){

	router := chi.NewRouter()
	applyMiddleware(router)
	healthRouteInit(router)
	startServer(router)



}

func applyMiddleware(router *chi.Mux) {
	log.Println("[OK] Apply Middleware")
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Timeout(60 * time.Second))
	
}

func healthRouteInit(router *chi.Mux) {
	log.Println("[OK] Health route initialization")
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
}

func startServer(router *chi.Mux) {
	log.Println("[OK] Starting server...")
	server := &http.Server{
		Addr: ":8080",
		Handler: router,
	}

	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sig:= make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func(){
		<-sig
		log.Println("[OK] Shuting Down...")

		shutdownCtx,_ := context.WithTimeout(serverCtx, 30*time.Second)

		go func ()  {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded{
				log.Fatal("[FAIL] graceful shutdown timed out... forcing exit")
			}
		}()

		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed{
		log.Fatal(err)
	}

	<-serverCtx.Done()


}