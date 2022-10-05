package server

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/superkruger/thunderdrone/internal/routes"
	"log"
	"net/http"
)

func Start(context context.Context, routes []routes.Routable, port string) {
	router := gin.Default()
	applyCors(router)

	api := router.Group("/api")

	for _, route := range routes {
		route.RegisterRoutes(api)
	}

	fmt.Println("Listening on port " + port)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	select {
	case <-context.Done():
		log.Println(port + " context done.")
	}
	log.Println(port + " Server exiting")
}

func applyCors(r *gin.Engine) {
	corsConfig := cors.DefaultConfig()
	//hot reload CORS
	corsConfig.AllowOrigins = []string{"http://localhost:3000"}
	corsConfig.AllowCredentials = true
	r.Use(cors.New(corsConfig))
}
