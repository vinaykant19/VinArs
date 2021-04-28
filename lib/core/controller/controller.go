package controller

import (
	"log"
	"net/http"
	cnf "../../../configuration"
)

/*
 * Base controller
 * @author: Vinaykant
 */
type Controller struct {
	Response *http.ResponseWriter
	Request * http.Request
	StatusCode int
	ResponseType string
	Result string
	Configuration *cnf.Configuration
	Logger *log.Logger
}


