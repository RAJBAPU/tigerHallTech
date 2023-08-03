package service

import "time"

type SignUpInput struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm" `
}

type SignInInput struct {
	Email    string `json:"email"  binding:"required"`
	Password string `json:"password"  binding:"required"`
}

type UserResponse struct {
	Id    int    `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
}

type TigerDetails struct {
	TigerId              int         `json:"tigerId,omitempty"`
	Name                 string      `json:"name,omitempty"`
	Dob                  time.Time   `json:"dob,omitempty"`
	LastSteenTimeStamp   time.Time   `json:"lastSteenTimeStamp,omitempty"`
	LastSteenCoordinates Coordinates `json:"lastSteenCoordinates,omitempty"`
	Image                string      `json:"image,omitempty"`
}

type TigerDetailsWithPaginationResponse struct {
	TigerDetails []TigerDetails `json:"tigerDetails"`
	TotalTigers  int            `json:"totalTigers"`
}

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type TigerSightingResponse struct {
	TigerId         int               `json:"tigerId"`
	SightingDetails []SightingDetails `json:"sightingDetails"`
	Total           int               `json:"total"`
}

type SightingDetails struct {
	LastSteenTimeStamp   time.Time   `json:"lastSteenTimeStamp"`
	LastSteenCoordinates Coordinates `json:"lastSteenCoordinates"`
	Image                string      `json:"image"`
}

type NotificationMessage struct {
	Name               string
	Email              string
	TigerName          string
	Latitude           float64
	Longitude          float64
	LastSteenTimeStamp time.Time
	EmailFrom          string
}
