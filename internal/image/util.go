package image

import (
	"context"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"sync"
)

// UploadToStorage - Upload image to Cloudinary (image data store).
func UploadToStorage(cld *cloudinary.Cloudinary, photo Photo, ch chan ImageData, wg *sync.WaitGroup) {
	defer wg.Done()
	ctx := context.Background()
	cldUploadResult, err := cld.Upload.Upload(ctx, photo.Src.Original, uploader.UploadParams{})
	if err != nil {
		log.Fatalf("Failed to upload file, %v\n", err)
	}

	ch <- ImageData{
		ID:   cldUploadResult.PublicID,
		Hits: 1,
		Uri:  cldUploadResult.URL,
	}
}

// getImagesData - Performs the asynchronous upload to the Cloudinary.
func getImagesData(imagesResponse ProviderResponse, cld *cloudinary.Cloudinary) []ImageData {
	var wg sync.WaitGroup
	ch := make(chan ImageData)
	imagesData := make([]ImageData, 0)
	for i := range imagesResponse.Photos {
		wg.Add(1)
		go UploadToStorage(cld, imagesResponse.Photos[i], ch, &wg)
	}

	go func() {
		// wait for all the workers to finish before collecting the results
		wg.Wait()
		// channel is closed only after the below for loop terminates
		close(ch)
	}()

	// traverse the channel and append it to the imagesData slice
	for v := range ch {
		imagesData = append(imagesData, v)
	}
	return imagesData
}

// convertSliceImageDataToDBObject - Converter used for converting []ImageData to MongoDB Document.
func convertSliceImageDataToDBObject(imageData []ImageData) []interface{} {
	dbObjects := make([]interface{}, 0)

	for _, data := range imageData {
		bson := bson.D{
			{"_id", data.ID},
			{"url", data.Uri},
			{"hits", data.Hits},
		}
		dbObjects = append(dbObjects, bson)
	}

	return dbObjects
}

// convertImageDataToDocument - Converter used for converting ImageData to MongoDB Document.
func convertImageDataToDocument(imageData ImageData) interface{} {
	document := bson.D{
		{"url", imageData.Uri},
		{"hits", imageData.Hits},
	}
	if imageData.Uri != "" {
		document = bson.D{{"url", imageData.Uri}}
	} else if imageData.Hits != 0 {
		document = bson.D{{"hits", imageData.Hits}}
	} else {
		document = bson.D{}
	}
	return document
}

// createImageData Convert SingleResult to Image Data model
func createImageData(result *mongo.SingleResult) (*ImageData, error) {
	var imageMetadata bson.M
	err := result.Decode(&imageMetadata)

	if err != nil {
		log.Error("Error on decoding ImageData ", err)
		return nil, err
	}

	var id = ""
	var hits int32 = 0
	var uri = ""

	if val, ok := imageMetadata["_id"]; ok {
		id = val.(string)
	}

	if val, ok := imageMetadata["hits"]; ok {
		hits = val.(int32)
	}

	if val, ok := imageMetadata["uri"]; ok {
		uri = val.(string)
	}

	return &ImageData{
		ID:   id,
		Hits: hits,
		Uri:  uri,
	}, nil
}
