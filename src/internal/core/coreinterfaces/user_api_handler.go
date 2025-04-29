package coreinterfaces

import (
	"net/http"
)

type UserAPIHandler interface {
	LoginHandler(w http.ResponseWriter, r *http.Request)
	AuthMeHandler(w http.ResponseWriter, r *http.Request)
	CategoryHandler(w http.ResponseWriter, r *http.Request)
	ProductHandler(w http.ResponseWriter, r *http.Request)
}
