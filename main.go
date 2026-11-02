package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/guvanch999/abrayly-mebel-backend/cmd"
	"github.com/guvanch999/abrayly-mebel-backend/docs"
	"github.com/guvanch999/abrayly-mebel-backend/internal"
	"github.com/guvanch999/abrayly-mebel-backend/internal/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

//	@title						Abraýly Mebel API
//	@version					1.0
//	@description				REST API of the Abraýly Mebel furniture store: trilingual (tm/ru/en) catalog with a category tree, product attributes and images, home-page banners and an admin panel backend.
//	@BasePath					/
//	@contact.name				API Support
//	@contact.email				guvanch.d.1999@gmail.com
//	@securityDefinitions.apikey	Bearer
//	@in							header
//	@name						Authorization
//	@description				Admin access token: "Bearer <access_token>" obtained from POST /auth/admin/login.

func main() {
	server := cmd.NewServer()

	if server.Config.IsDevelopment() || server.Config.AppUrl == "" {
		docs.SwaggerInfo.Schemes = []string{"http"}
		docs.SwaggerInfo.Host = fmt.Sprintf("localhost:%v", server.Config.Port)
	} else {
		docs.SwaggerInfo.Schemes = []string{"https"}
		host := strings.TrimPrefix(strings.TrimPrefix(server.Config.AppUrl, "https://"), "http://")
		docs.SwaggerInfo.Host = strings.TrimSuffix(strings.TrimSuffix(host, "/v1"), "/")
	}

	server.Echo.GET("/swagger/*", echoSwagger.EchoWrapHandler())

	middleware.RegisterMiddlewares(server)
	internal.InitRouters(server)

	server.Echo.Static("/uploads", "uploads")

	signalCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := server.Echo.Start(fmt.Sprintf(":%v", server.Config.Port)); err != nil && !errors.Is(err, http.ErrServerClosed) {
			server.Echo.Logger.Fatal("shutting down the server. error: " + err.Error())
		}
	}()

	// Wait for SIGINT/SIGTERM (Ctrl+C, systemd stop, docker stop) and shut
	// down gracefully with a 10 second timeout.
	<-signalCtx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Echo.Shutdown(ctx); err != nil {
		server.Echo.Logger.Fatal(err)
	}
	if sqlDB, err := server.Database.DB(); err == nil {
		_ = sqlDB.Close()
	}
}
