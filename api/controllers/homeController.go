package controllers

import (
	"net/http"

	"github.com/lucassilveira96/rest-api-jwt-go-lang/api/responses"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome To This Awesome API")

}
