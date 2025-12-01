package main

import (
	"log"

	"github.com/internalWizzard/idp-server/internal/auth/infrastructure/middleware"
	oidcprovicer "github.com/internalWizzard/idp-server/internal/auth/infrastructure/oidc"
	monitorhttp "github.com/internalWizzard/idp-server/internal/monitor/delivery/http"
	"github.com/internalWizzard/idp-server/pkg/httpserver"
)

func main() {
	srv := httpserver.New()

	e := srv.Echo()

	authProvider := oidcprovicer.NewOIDCProvider("http://localhost:8080", "master")
	authMW := middleware.NewAuthMiddleware(authProvider)

	// API v1 group
	apiV1 := e.Group("/api/v1")

	// Secured group
	securedV1 := apiV1.Group("", authMW.Handle)

	// Register module routes
	monitorhttp.RegisterRoutes(apiV1, securedV1)

	if err := srv.Start(":8000"); err != nil {
		log.Fatalf("http server failed: %v", err)
	}
}
