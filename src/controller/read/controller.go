package read

import (
	"log"
	"net/http"
	"../../../localization"
	cnf "../../../configuration"
	database "../../../lib/core/database"
)

type ReadController struct {
	Response *http.ResponseWriter
	Request * http.Request
	UserId int64
	Token string
	StatusCode int
	ResponseType string
	Result string
	Configuration *cnf.Configuration
	DB *database.Database
	Logger *log.Logger
	DeviceToken string
	DeviceType string
	Language localization.Localization
}


