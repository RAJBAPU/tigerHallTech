package service

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"math"
	"simpl_pr/firebase"
	models "simpl_pr/model"
	"simpl_pr/persistence"
	util "simpl_pr/utils"
	"time"

	"strings"

	"github.com/astaxie/beego/orm"
	"github.com/nfnt/resize"
	"github.com/spf13/cast"
)

var messageChan = make(chan NotificationMessage, 100) // Buffered channel to handle messages

type Tiger struct {
	Configs      persistence.TgConfigPersistence
	User         persistence.TgUserPersistence
	TigerDetails persistence.TgTigerDetailsPersistence
	TigerImages  persistence.TgTigerImagesPersistence
}

func (tg *Tiger) PostTigerDetails(payload TigerDetails) (err error) {

	dobString := payload.Dob.Format("2006-01-02")

	tigerDetails, err := tg.TigerDetails.GetTgTigerDetails(payload.Name, dobString)
	if (err != nil && err != orm.ErrNoRows) || tigerDetails != nil {
		fmt.Println("Error in GetTgTigerDetails: ", err)
		return errors.New("tiger already present in database")
	}
	tigerDetails = &models.TgTigerDetails{}
	tigerDetails.Name = payload.Name
	tigerDetails.Dob = dobString
	tigerDetails.LastSteenTimeStamp = payload.LastSteenTimeStamp
	tigerDetails.Longitude = payload.LastSteenCoordinates.Longitude
	tigerDetails.Latitude = payload.LastSteenCoordinates.Latitude

	_, err = tg.TigerDetails.AddTgTigerDetails(tigerDetails)
	if err != nil {
		fmt.Println("Error in AddTgTigerDetails: ", err)
		return
	}
	return
}

func (tg *Tiger) GetAllTigers(page int, pageSize int) (response *TigerDetailsWithPaginationResponse, err error) {

	total, err := tg.TigerDetails.GetCountOfTigers()
	if err != nil {
		fmt.Println("Error in GetCountOfTigers: ", err)
		return
	}

	// Calculate the offset and limit for pagination
	offset := (page - 1) * pageSize
	limit := pageSize

	tigers, err := tg.TigerDetails.GetAllTigers(offset, limit)
	if err != nil {
		fmt.Println("Error in GetAllTigers: ", err)
		return
	}

	tigerResponse := &TigerDetails{}
	var tigerDetailsSlice []TigerDetails
	for _, tigerDetails := range tigers {
		tigerResponse.Name = tigerDetails.Name
		tigerResponse.Dob = cast.ToTime(tigerDetails.Dob)
		tigerResponse.LastSteenTimeStamp = tigerDetails.LastSteenTimeStamp
		tigerResponse.LastSteenCoordinates.Longitude = tigerDetails.Longitude
		tigerResponse.LastSteenCoordinates.Latitude = tigerDetails.Latitude
		tigerDetailsSlice = append(tigerDetailsSlice, *tigerResponse)
	}

	response = &TigerDetailsWithPaginationResponse{}
	response.TigerDetails = tigerDetailsSlice
	response.TotalTigers = total

	return

}

func (tg *Tiger) PostSightingDetails(request TigerDetails, user *models.TgUser) (errorCode int, err error) {
	tigerDetails, err := tg.TigerDetails.GetTgTigerDetailsById(request.TigerId)
	if err != nil {
		fmt.Println("Error in GetTgTigerDetailsById: ", err)
		return
	}

	if request.LastSteenCoordinates.Latitude > 180 || request.LastSteenCoordinates.Longitude > 180 || request.LastSteenCoordinates.Longitude < -180 || request.LastSteenCoordinates.Longitude < -180 {
		fmt.Println("Not Valid Coordinates")
		return 0, errors.New("invalid Coordinates")
	}
	// Calculate distance between new sighting and previous sighting
	distance := calculateDistance(request.LastSteenCoordinates.Latitude, request.LastSteenCoordinates.Longitude, tigerDetails.Latitude, tigerDetails.Longitude)
	if distance < 5.0 {
		fmt.Println("distance less than 5km")
		errorCode = 452
		return
	}

	// Decode the Base64 string to bytes
	imageData, err := base64.StdEncoding.DecodeString(request.Image)
	if err != nil {
		fmt.Println("Failed to decode Base64 data:", err)
		return
	}
	resizedImageData := resizeImage(imageData, 250, 200)
	filename := firebase.GenerateUUID()

	filePath := cast.ToString(tigerDetails.Id) + "/" + "tigerImage" + "/" + filename
	configs := tg.Configs.GetAllConfigs()

	tigerImageUrl, err := firebase.UploadToFireBaseAndGetAccessUrl(resizedImageData, "image/jpg", filePath, configs)
	if err != nil {
		fmt.Println("Error in UploadToFireBaseAndGetAccessUrl: ", err)
		return
	}

	tigerDetails.LastSteenTimeStamp = request.LastSteenTimeStamp
	tigerDetails.Latitude = request.LastSteenCoordinates.Latitude
	tigerDetails.Longitude = request.LastSteenCoordinates.Longitude
	err = tg.TigerDetails.UpdateTgTiger(tigerDetails, nil, "lastSteenTimeStamp", "latitude", "longitude")
	if err != nil {
		fmt.Println(" error in UpdateTgTiger ", err)
		return
	}

	tigerImage := &models.TgTigerImages{}
	tigerImage.TigerId = tigerDetails.Id
	tigerImage.Image = tigerImageUrl
	tigerImage.LastSteenTimeStamp = request.LastSteenTimeStamp
	tigerImage.Latitude = request.LastSteenCoordinates.Latitude
	tigerImage.Longitude = request.LastSteenCoordinates.Longitude
	tigerImage.SightedByUser = user.Id
	_, err = tg.TigerImages.AddTgTigerImages(tigerImage)
	if err != nil {
		fmt.Println(" error in AddTgTigerImages: ", err)
		return
	}

	user, err = tg.User.GetYpUserById(user.Id)
	if err != nil {
		fmt.Println(" error in GetYpUserById: ", err)
		return
	}
	var firstName = user.Name

	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[1]
	}

	message := NotificationMessage{
		Name:               firstName,
		Email:              user.Email,
		TigerName:          tigerDetails.Name,
		Latitude:           tigerDetails.Latitude,
		Longitude:          tigerDetails.Longitude,
		LastSteenTimeStamp: tigerImage.LastSteenTimeStamp,
		EmailFrom:          configs["EmailFrom"],
	}

	messageChan <- message

	return

}

const earthRadiusKm = 6371.0 // Earth's radius in kilometers

func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	// Convert latitude and longitude from degrees to radians
	lat1Rad := degreesToRadians(lat1)
	lon1Rad := degreesToRadians(lon1)
	lat2Rad := degreesToRadians(lat2)
	lon2Rad := degreesToRadians(lon2)

	// Haversine formula
	deltaLat := lat2Rad - lat1Rad
	deltaLon := lon2Rad - lon1Rad
	a := math.Pow(math.Sin(deltaLat/2), 2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Pow(math.Sin(deltaLon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := earthRadiusKm * c

	return distance
}

func degreesToRadians(deg float64) float64 {
	return deg * (math.Pi / 180)
}

func (tg *Tiger) GetSightingDetails(tigerId int, page int, pageSize int) (response *TigerSightingResponse, err error) {
	total, err := tg.TigerImages.GetCountOfTigerSightings(tigerId)
	if err != nil {
		fmt.Println("Error in GetCountOfTigerSightings: ", err)
		return
	}

	// Calculate the offset and limit for pagination
	offset := (page - 1) * pageSize
	limit := pageSize

	sightingDetails, err := tg.TigerImages.GetAllTigerSightings(tigerId, offset, limit)
	if err != nil {
		fmt.Println("Error in GetAllTigerSightings: ", err)
		return
	}
	configs := tg.Configs.GetAllConfigs()
	baseUrl := configs["firebaseBaseUrl"]
	sightingResponse := &SightingDetails{}
	var tigerDetailsSlice []SightingDetails
	for _, details := range sightingDetails {
		sightingResponse.LastSteenTimeStamp = details.LastSteenTimeStamp
		sightingResponse.LastSteenCoordinates.Longitude = details.Longitude
		sightingResponse.LastSteenCoordinates.Latitude = details.Latitude
		sightingResponse.Image = RemoveBaseURLInFirebaseUrl(details.Image, baseUrl)
		tigerDetailsSlice = append(tigerDetailsSlice, *sightingResponse)
	}

	response = &TigerSightingResponse{}
	response.TigerId = tigerId
	response.SightingDetails = tigerDetailsSlice
	response.Total = total
	return

}

func ProcessMessages() {
	for {
		message := <-messageChan // Wait for a message from the channel
		sendNotificationEmail(message)
	}
}

func sendNotificationEmail(message NotificationMessage) {
	timestamp := message.LastSteenTimeStamp.Add(-5 * time.Hour).Add(-30 * time.Minute)
	time := timestamp.Format("2006-01-02 15:04:05")
	// ? Send Email
	emailData := util.EmailData{
		Subject: "Tiger Has been spotted",
		Text:    "Tiger named" + message.TigerName + "was spotted at " + time,
	}

	err := util.SendEmail(message.EmailFrom, message.Email, &emailData)
	if err != nil {
		fmt.Println("Unable to send email: ", err)
	}

	fmt.Println("Sending notification email:", message)
}

func RemoveBaseURLInFirebaseUrl(imageUrl string, baseUrl string) string {
	return strings.Replace(imageUrl, baseUrl, "", 1)
}

func resizeImage(data []byte, width, height int) []byte {
	// Create a new image from the bytes.
	img, _, _ := image.Decode(bytes.NewReader(data))

	// Resize the image.
	resizedImg := resize.Resize(250, 200, img, resize.Lanczos3)

	// Encode the resized image to bytes.
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, resizedImg, nil)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return buf.Bytes()
}
