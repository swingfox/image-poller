package image

import (
	"encoding/json"
	"fmt"
	"github.com/cloudinary/cloudinary-go/v2"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Config struct {
	ImageProviderHost   string
	ImageProviderAPIKey string
	StorageName         string
	StorageAPIKey       string
	StorageSecret       string
	Limit               int
}

type Service struct {
	ImageConfig Config
	Client      http.Client
}

func (is *Service) GetImages(limit int) (resp *ImageResponse, err error) {
	if limit > is.ImageConfig.Limit {
		log.Info("Image Exceeded the hard limit.")
		limit = is.ImageConfig.Limit
	}

	// get Images from Pexel
	var imagesResponse ProviderResponse
	imagesResponse, err = is.getImagesFromProvider(limit)
	if err != nil {
		return resp, err
	}

	// get Images from Cloudinary
	var cld *cloudinary.Cloudinary
	cld, cldErr := is.getStorageClient()
	if cldErr != nil {
		return resp, err
	}

	// get the response data
	imagesData := getImagesData(imagesResponse, cld)

	log.Info("Finished Processing Images...")
	return &ImageResponse{
		Limit:     len(imagesData),
		ImageData: imagesData,
	}, nil
}

func (is *Service) getStorageClient() (*cloudinary.Cloudinary, error) {
	var cld, cldErr = cloudinary.NewFromParams(is.ImageConfig.StorageName, is.ImageConfig.StorageAPIKey, is.ImageConfig.StorageSecret)
	if cldErr != nil {
		log.Fatalf("Failed to intialize Cloudinary, %v", cldErr)
		return nil, fmt.Errorf("failed to intialize Cloudinary, %v", cldErr)
	}
	return cld, nil
}

func (is *Service) getImagesFromProvider(limit int) (ProviderResponse, error) {
	query := fmt.Sprintf(is.ImageConfig.ImageProviderHost, strconv.Itoa(limit))

	client := http.Client{}
	req, err := http.NewRequest("GET", query, nil)
	req.Header.Add("Authorization", is.ImageConfig.ImageProviderAPIKey)
	resp, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("error sending request to the Image Provider: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = fmt.Errorf("error reading request from Image Provider:  %v", err)
	}

	var imageResponse ProviderResponse
	err = json.Unmarshal(body, &imageResponse)
	if err != nil {
		err = fmt.Errorf("error reading request from Image Provider:  %v", err)
	}
	return imageResponse, err
}

func (is *Service) GetImage(ID string) {
	//TODO implement me
	panic("implement me")
}

func (is *Service) UpdateImage(ID string) {
	//TODO implement me
	panic("implement me")
}

func (is *Service) CreateImage(request Request) {
	//TODO implement me
	panic("implement me")
}
