package controllers

import (
	"encoding/json"
	"fmt"
	"http_monitoring/api/models"
	"http_monitoring/api/utils/formaterror"
	"http_monitoring/api/utils/responses"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (server *Server) CreateCall(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	epc := models.EndPointCalls{}
	err = json.Unmarshal(body, &epc)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	epc.Prepare()
	err = epc.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	// uid, err := auth.ExtractTokenID(r)
	// if err != nil {
	// 	responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
	// 	return
	// }
	// if uid != epc.OwnerID {
	// 	responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
	// 	return
	// }
	epcCreated, err := epc.SaveCall(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, epcCreated.ID))
	responses.JSON(w, http.StatusCreated, epcCreated)
}

func (server *Server) GetCalls(w http.ResponseWriter, r *http.Request) {

	epc := models.EndPointCalls{}

	epcs, err := epc.FindAllCalls(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, epcs)
}
func (server *Server) GetCallsByTime(w http.ResponseWriter, r *http.Request) {

	epc := models.EndPointCalls{}
	vars := mux.Vars(r)

	epcs, err := epc.FindCallsByTime(server.DB,vars["urlId"],vars["StartTime"],vars["EndTime"])
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, epcs)
}

func (server *Server) GetCall(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	epcid, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	epc := models.EndPointCalls{}

	CallReceived, err := epc.FindCallByID(server.DB, uint32(epcid))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, CallReceived)
}



