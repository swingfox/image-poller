package image

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cloudinary/cloudinary-go/v2"
	log "github.com/sirupsen/logrus"
	"github.com/swingfox/image-poller/internal/persistence"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	// save image info to DB
	imageCollection := persistence.GetCollection("images")
	result, err := imageCollection.InsertMany(context.TODO(), convertImageDataToDBObject(imagesData))

	if err != nil {
		log.Error("error saving db", err)
	} else {
		log.Info("successfully inserted to DB", result)
	}

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

func (is *Service) GetImage(ID string) (resp *ImageData, err error) {
	// save image info to DB
	imageCollection := persistence.GetCollection("images")
	opts := options.FindOneOptions{}
	// find one ImageData by ID
	result := imageCollection.FindOne(context.TODO(), bson.D{{"_id", ID}}, &opts)

	if err != nil {
		// if query did not match any documents
		if err == mongo.ErrNoDocuments {
			log.Error("GetImage: Query for "+ID+" did not match any documents", err)
			return nil, err
		} else {
			log.Error("GetImage: Error on FindOne", err)
			return nil, err
		}
	}

	return createImageData(result)
}

func (is *Service) UpdateImage(ID string, data ImageData) (resp *ImageData, err error) {
	// update image info to DB
	imageCollection := persistence.GetCollection("images")
	filter := bson.D{{"_id", ID}}
	update := bson.D{{"$set", convertImageDataToDocument(data)}}
	// find one and update ImageData by ID
	res, err := imageCollection.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		// if query did not match any documents
		if err == mongo.ErrNoDocuments {
			log.Error("UpdateImage: Query for "+ID+" did not match any documents", err)
			return nil, err
		} else {
			log.Error("UpdateImage: Error on FindOneAndUpdate", err)
			return nil, err
		}
	}
	log.Info(fmt.Sprintf("updated document %v", res))

	return &data, nil
}

func (is *Service) DeleteImage(ID string) (int64, error) {
	// update image info to DB
	imageCollection := persistence.GetCollection("images")
	filter := bson.D{{"_id", ID}}
	update := bson.D{{"$set", bson.D{{"isDeleted", true}}}}
	// find one and update ImageData by ID
	res, err := imageCollection.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		// if query did not match any documents
		if err == mongo.ErrNoDocuments {
			log.Error("DeleteImage: Query for "+ID+" did not match any documents", err)
			return 0, err
		} else {
			log.Error("DeleteImage: Error on FindOneAndUpdate", err)
			return 0, err
		}
	}
	log.Info(fmt.Sprintf("deleted document %v", res))
	return res.MatchedCount, nil
}

func (is *Service) CreateImage(request Request) {
	//TODO implement me
	panic("implement me")
}
