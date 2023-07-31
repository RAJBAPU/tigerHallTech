package firebase

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	models "simpl_pr/model"
	"strings"

	"cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

const (
	ConstFRONT_SLASH = "/"
	ConstEMPTY       = ""
	url2F            = "%2F"
	urlAltMediaToken = "?alt=media&token="
	uuidFormat       = "%x-%x-%x-%x-%x"
)

func initializeFirebase() (*firebase.App, error) {
	opt := option.WithCredentialsFile("./firebase/firebase_creds.json") // Replace with the path to your Firebase Admin SDK credentials file
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}
	return app, nil
}

// Function to retrieve the image from Firebase Storage
func GetImageFromFirebase(filePath, imageName string) ([]byte, error) {
	// Initialize Firebase
	app, err := initializeFirebase()
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize Firebase: %v", err)
	}

	// Create a new Firebase Storage client
	client, err := app.Storage(context.Background())
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize Firebase Storage client: %v", err)
	}

	// Get a reference to the image in Firebase Storage
	bucketName := "your-firebase-storage-bucket-name" // Replace with your Firebase Storage bucket name
	objectName := fmt.Sprintf("%s/%s", filePath, imageName)
	obj, _ := client.Bucket(bucketName)
	objj := obj.Object(objectName)
	// Download the image from Firebase Storage
	imageData, err := objj.NewReader(context.Background())
	if err != nil {
		return nil, fmt.Errorf("Failed to get image from Firebase Storage: %v", err)
	}
	defer imageData.Close()

	// Read the image data into a byte slice
	imageBytes, err := io.ReadAll(imageData)
	if err != nil {
		return nil, fmt.Errorf("Failed to read image data: %v", err)
	}

	return imageBytes, nil
}
func UploadToFireBaseAndGetAccessUrl(content []byte, contentType, filepath string) (url string, err error) {
	configs := models.GetAllConfigs()
	bucket := configs["bucket"]
	url, err = firebaseUploadWithBucket(content, contentType, filepath, bucket, "")
	firebaseBaseUrl := configs["firebaseBaseUrl"]
	url = firebaseBaseUrl + url
	return
}

// firebaseUploadWithBucket ...
func firebaseUploadWithBucket(content []byte, contentType, filepath, bucket, downloadToken string) (url string, err error) {
	// get the firebase app
	firebaseApp, err := initializeFirebase()
	if err != nil {
		return
	}
	client, err := firebaseApp.Storage(context.Background())
	if err != nil {
		fmt.Println("clienterr", err)
		return
	}

	var bucketHandle *storage.BucketHandle

	if bucket == "" {
		bucketHandle, err = client.DefaultBucket()
	} else {
		bucketHandle, err = client.Bucket(bucket)
	}
	fmt.Println("bucketHandle", bucketHandle, err)

	// get the bucket
	if err != nil {
		return
	}

	// generate uuid and send it in write request which will be used as download token
	uuid := GenerateUUID()

	// if a download token is passed then use it as the download token
	if downloadToken != "" {
		uuid = downloadToken
	}

	// create a firebase bucket object with the given filepath
	obj := bucketHandle.Object(filepath)

	// create a writer for the object
	w := obj.NewWriter(context.Background())
	w.ContentType = contentType
	w.Metadata = map[string]string{"firebaseStorageDownloadTokens": uuid} // download token for getting download url without authorization
	// w.ChunkSize = 0
	w.Write(content)
	if err = w.Close(); err != nil {
		return
	}

	url = strings.Replace(filepath, ConstFRONT_SLASH, url2F, -1) + urlAltMediaToken + uuid
	return
}

// GenerateUUID ... generate a random uuid
func GenerateUUID() (uuid string) {

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	uuid = fmt.Sprintf(uuidFormat, b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return
}
