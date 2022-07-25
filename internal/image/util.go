package image

import (
	"context"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"sync"
)

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

	for v := range ch {
		imagesData = append(imagesData, v)
	}
	return imagesData
}

func convertImageDataToDBObject(imageData []ImageData) []interface{} {
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
