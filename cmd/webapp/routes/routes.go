package routes

import (
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/swingfox/image-poller/internal/image"
	"github.com/swingfox/image-poller/internal/user"
	"github.com/swingfox/image-poller/internal/userrole"
	"net/http"
	"strconv"
	"time"
)

type ServiceResponse struct {
	Timestamp time.Time `json:"timestamp"`
	Status    int       `json:"status"`
	Error     string    `json:"error"`
	Message   string    `json:"message"`
	Path      string    `json:"path"`
}

//go:generate mockgen -destination=../../../mocks/mock_userrole.go -package=mocks -source routes.go
type UserRole interface {
	CreateRole(request userrole.Request)
	GetRole(ID string)
	UpdateRole(ID string)
	DeleteRole(ID string)
}

//go:generate mockgen -destination=../../../mocks/mock_user.go -package=mocks -source routes.go
type User interface {
	CreateUser(request user.Request)
	GetUser(ID string)
	UpdateUser(ID string)
	DeleteUser(ID string)
}

//go:generate mockgen -destination=../../../mocks/mock_image_provider.go -package=mocks -source routes.go
type ImageProvider interface {
	CreateImage(request image.ImageData)
	GetImages(limit int) (*image.ImageResponse, error)
	GetImage(ID string) (*image.ImageData, error)
	UpdateImage(ID string, data image.ImageData) (*image.ImageData, error)
	DeleteImage(ID string) (int64, error)
}

type Handler struct {
	ImageService    ImageProvider
	UserService     User
	UserRoleService UserRole
}

func (hndlr *Handler) Set(router *mux.Router) {

	// User Routes
	router.HandleFunc("/users", hndlr.CreateUser).Methods("POST")
	router.HandleFunc("/users", hndlr.GetUser).Methods("GET")
	router.HandleFunc("/users/{id}", hndlr.UpdateUser).Methods("PATCH")
	router.HandleFunc("/users/{id}", hndlr.DeleteUser).Methods("DELETE")

	// Image Routes
	router.HandleFunc("/images", hndlr.CreateImage).Methods("POST")
	router.HandleFunc("/images", hndlr.GetImages).Methods("GET")
	router.HandleFunc("/images/{id}", hndlr.GetImage).Methods("GET")
	router.HandleFunc("/images/{id}", hndlr.UpdateImage).Methods("PATCH")
	router.HandleFunc("/images/{id}", hndlr.DeleteImage).Methods("DELETE")

	// User role Routes
	router.HandleFunc("/userrole", hndlr.CreateRole).Methods("POST")
	router.HandleFunc("/userrole/{id}", hndlr.GetRole).Methods("GET")
	router.HandleFunc("/userrole/{id}", hndlr.UpdateRole).Methods("PATCH")
	router.HandleFunc("/userrole/{id}", hndlr.DeleteRole).Methods("DELETE")

}

func (hndlr *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {

}

func (hndlr *Handler) GetUser(w http.ResponseWriter, r *http.Request) {

}

func (hndlr *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {

}

func (hndlr *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {

}

func (hndlr *Handler) CreateImage(w http.ResponseWriter, r *http.Request) {
	var imageData image.ImageData
	err := json.NewDecoder(r.Body).Decode(&imageData)
	if err != nil {
		log.Error("Error decoding request body: ", err)
		methodBadRequestHandler(w, r)
	}
	hndlr.ImageService.CreateImage(imageData)
	if err != nil {
		log.Error("Error calling CreateImage with data: ", err)
		methodNotFoundHandler(w, r)
	} else {
		response := ServiceResponse{
			Timestamp: time.Now().UTC(),
			Status:    200,
			Path:      r.URL.Path,
		}
		errorResponse, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.Write(errorResponse)
	}
}

func (hndlr *Handler) GetImages(w http.ResponseWriter, r *http.Request) {
	var imageResponse *image.ImageResponse
	var imageErr error
	// extract optional limit query param from URL
	if limit, ok := r.URL.Query()["limit"]; ok {
		value, err := strconv.Atoi(limit[0])
		if err != nil {
			log.Error("Limit query param parse error.")
		}
		imageResponse, imageErr = hndlr.ImageService.GetImages(value)
		if imageErr != nil {
			log.Error(imageErr)
		}
	} else {
		imageResponse, imageErr = hndlr.ImageService.GetImages(5)
		if imageErr != nil {
			log.Error(imageErr)
		}
	}

	writeJsonResponse(w, imageResponse)
}

func (hndlr *Handler) GetImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		log.Info("GetImage: ID is missing in parameters")
	}
	imageMetadata, err := hndlr.ImageService.GetImage(id)

	if err != nil {
		log.Error("Error calling GetImage with ID: " + id)
		methodNotFoundHandler(w, r)
	} else {
		writeJsonResponse(w, imageMetadata)
	}
}

func (hndlr *Handler) UpdateImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		log.Info("UpdateImage: ID is missing in parameters")
	}
	var imageData image.ImageData
	err := json.NewDecoder(r.Body).Decode(&imageData)
	if err != nil {
		log.Error("Error decoding request body with ID: " + id)
		methodBadRequestHandler(w, r)
	}
	imageMetadata, err := hndlr.ImageService.UpdateImage(id, imageData)
	if err != nil {
		log.Error("Error calling UpdateImage with ID: " + id)
		methodNotFoundHandler(w, r)
	} else {
		writeJsonResponse(w, imageMetadata)
	}
}

func (hndlr *Handler) DeleteImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		log.Info("DeleteImage: ID is missing in parameters")
	}
	imageMetadata, err := hndlr.ImageService.DeleteImage(id)
	if err != nil {
		log.Error("Error calling DeleteImage with ID: " + id)
		methodNotFoundHandler(w, r)
	} else {
		writeJsonResponse(w, imageMetadata)
	}
}

func (hndlr *Handler) CreateRole(w http.ResponseWriter, r *http.Request) {

}

func (hndlr *Handler) UpdateRole(w http.ResponseWriter, r *http.Request) {

}

func (hndlr *Handler) GetRole(w http.ResponseWriter, r *http.Request) {

}

func (hndlr *Handler) DeleteRole(w http.ResponseWriter, r *http.Request) {

}

// see http.HandlerFunc
func methodNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	serviceCallErrorHandler(http.StatusNotFound, w, r)
}

// see http.HandlerFunc
func methodBadRequestHandler(w http.ResponseWriter, r *http.Request) {
	serviceCallErrorHandler(http.StatusBadRequest, w, r)
}

// see http.HandlerFunc
func methodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	serviceCallErrorHandler(http.StatusMethodNotAllowed, w, r)
}

// see http.HandlerFunc
func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	serviceCallErrorHandler(http.StatusNotFound, w, r)
}

func serviceCallErrorHandler(errorCode int, w http.ResponseWriter, r *http.Request) {
	error := ServiceResponse{
		Timestamp: time.Now().UTC(),
		Status:    errorCode,
		Error:     http.StatusText(errorCode),
		Message:   http.StatusText(errorCode),
		Path:      r.URL.Path,
	}
	errorResponse, _ := json.Marshal(error)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errorCode)
	w.Write(errorResponse)
}

func writeJsonResponse(w http.ResponseWriter, response interface{}) {
	jsonResponse, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
