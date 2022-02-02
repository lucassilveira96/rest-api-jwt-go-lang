package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lucassilveira96/rest-api-jwt-go-lang/api/auth"
	"github.com/lucassilveira96/rest-api-jwt-go-lang/api/models"
	"github.com/lucassilveira96/rest-api-jwt-go-lang/api/responses"
	"github.com/lucassilveira96/rest-api-jwt-go-lang/api/utils/formaterror"
)

func (server *Server) CreateStock(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	stock := models.Stock{}
	err = json.Unmarshal(body, &stock)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	stock.Prepare()
	err = stock.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if uid != uint32(stock.ID) {
		//responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		//return
	}
	stockCreated, err := stock.SaveStock(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Lacation", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, stockCreated.ID))
	responses.JSON(w, http.StatusCreated, stockCreated)
}

func (server *Server) GetAllStocks(w http.ResponseWriter, r *http.Request) {

	stock := models.Stock{}

	stocks, err := stock.FindAllStocks(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, stocks)
}

func (server *Server) GetStock(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	stock := models.Stock{}

	stockReceived, err := stock.FindStockByID(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, stockReceived)
}

func (server *Server) UpdateStock(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Check if the Stock id is valid
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	//CHeck if the auth token is valid and  get the user id from it
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the Stock exist
	stock := models.Stock{}
	err = server.DB.Debug().Model(models.Stock{}).Where("id = ?", pid).Take(&stock).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Stock not found"))
		return
	}

	// If a user attempt to update a Stock not belonging to him
	if uid != uint32(stock.ID) {
		//responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		//return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Start processing the request data
	stockUpdate := models.Stock{}
	err = json.Unmarshal(body, &stockUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	stockUpdate.Prepare()
	err = stockUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	stockUpdate.ID = stock.ID //this is important to tell the model the Stock id to update, the other update field are set above

	stockUpdated, err := stockUpdate.UpdateAStock(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, stockUpdated)
}

func (server *Server) DeleteStock(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Is a valid Stock id given to us?
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	// Is this user authenticated?
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the Stock exist
	stock := models.Stock{}
	err = server.DB.Debug().Model(models.Stock{}).Where("id = ?", pid).Take(&stock).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	// Is the authenticated user, the owner of this Stock?
	if uid != uint32(stock.ID) {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	_, err = stock.DeleteAStock(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	responses.JSON(w, http.StatusNoContent, "")
}
