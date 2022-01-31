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

func (server *Server) CreateProductCategory(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	productCategory := models.ProductCategory{}
	err = json.Unmarshal(body, &productCategory)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	productCategory.Prepare()
	err = productCategory.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if uid != uint32(productCategory.ID) {
		//responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		//return
	}
	productCategoryCreated, err := productCategory.SaveProductCategory(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Lacation", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, productCategoryCreated.ID))
	responses.JSON(w, http.StatusCreated, productCategoryCreated)
}

func (server *Server) GetAllProductCategories(w http.ResponseWriter, r *http.Request) {

	productCategory := models.ProductCategory{}

	productCategories, err := productCategory.FindAllProductCategory(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, productCategories)
}

func (server *Server) GetProductCategory(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	productCategory := models.ProductCategory{}

	productCategoryReceived, err := productCategory.FindProductCategoryByID(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, productCategoryReceived)
}

func (server *Server) UpdateProductCategory(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Check if the productCategory id is valid
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

	// Check if the productCategory exist
	productCategory := models.ProductCategory{}
	err = server.DB.Debug().Model(models.ProductCategory{}).Where("id = ?", pid).Take(&productCategory).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Product Category not found"))
		return
	}

	// If a user attempt to update a product category' not belonging to him
	if uid != uint32(productCategory.ID) {
		//responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		//return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Start processing the request data
	productCategoryUpdate := models.ProductCategory{}
	err = json.Unmarshal(body, &productCategoryUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	productCategoryUpdate.Prepare()
	err = productCategoryUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	productCategoryUpdate.ID = productCategory.ID //this is important to tell the model the ProductCategory id to update, the other update field are set above

	productCategoryUpdated, err := productCategoryUpdate.UpdateAProductCategory(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, productCategoryUpdated)
}

func (server *Server) DeleteProductCategory(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Is a valid productCategory id given to us?
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

	// Check if the product category exist
	productCategory := models.ProductCategory{}
	err = server.DB.Debug().Model(models.ProductCategory{}).Where("id = ?", pid).Take(&productCategory).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	// Is the authenticated user, the owner of this product category?
	if uid != uint32(productCategory.ID) {
		//responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		//return
	}
	_, err = productCategory.DeleteAProductCategory(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	responses.JSON(w, http.StatusNoContent, "")
}
