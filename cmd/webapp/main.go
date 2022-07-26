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

	// Listen at specified port
	http.ListenAndServe(":"+serverPort, r)
}

// Initialize Initalize is the main init method to be used for services initialization
func Initialize() *routes.Handler {

	imageService := newImageProviderService()

	return &routes.Handler{
		ImageService:    imageService,
		UserService:     user.UserService{},
		UserRoleService: userrole.UserRoleService{},
	}
}

// Initialize image service
func newImageProviderService() *image.Service {
	imageConfig := image.Config{
		ImageProviderHost:   config.Registry.GetString("IMAGE_PROVIDER.HOST"),
		ImageProviderAPIKey: config.Registry.GetString("IMAGE_PROVIDER.API_KEY"),
		Limit:               config.Registry.GetInt("IMAGE_PROVIDER.LIMIT"),
		StorageName:         config.Registry.GetString("STORAGE.CLOUD_NAME"),
		StorageAPIKey:       config.Registry.GetString("STORAGE.API_KEY"),
		StorageSecret:       config.Registry.GetString("STORAGE.API_SECRET"),
	}
	imageService := &image.Service{
		ImageConfig: imageConfig,
	}
	return imageService
}
