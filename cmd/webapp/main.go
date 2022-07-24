package main

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/swingfox/image-poller/cmd/webapp/config"
	"github.com/swingfox/image-poller/cmd/webapp/routes"
	"github.com/swingfox/image-poller/internal/image"
	"github.com/swingfox/image-poller/internal/user"
	"github.com/swingfox/image-poller/internal/userrole"
	"net/http"
)

func main() {
	config.Set()

	routesHandler := Initialize()

	r := mux.NewRouter()
	routesHandler.Set(r)
	serverPort := config.Registry.GetString("SERVER_PORT")
	log.Info("Server listening on port " + serverPort)
	http.ListenAndServe(":"+serverPort, r)
}

func Initialize() *routes.Handler {
	return &routes.Handler{
		ImageService:    image.ImageService{},
		UserService:     user.UserService{},
		UserRoleService: userrole.UserRoleService{},
	}
}
