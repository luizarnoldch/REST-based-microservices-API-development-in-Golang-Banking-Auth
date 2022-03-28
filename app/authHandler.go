package app

import (
	"net/http"
"encoding/json"
"github.com/luizarnoldch/REST-based-microservices-API-development-in-Golang-Banking-Auth/domain"
"github.com/luizarnoldch/REST-based-microservices-API-development-in-Golang-Banking-Auth/service"
"github.com/luizarnoldch/REST-based-microservices-API-development-in-Golang-Banking-Lib/logger"
)

type AuthHandler struct {
	service service.AuthHandler
}

func (h AuthHandler) NotImplementedHandler(w http.ResponseWriter, r *http.Request){
	writeResponse(w,http.StatusOK,"Handler not implemented ...")
}

func (h AuthHandler) Login(w http.ResponseWriter, r *http.Request){
	var loginRequest dto.LoginRequest
	if err:=json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		logger.Error("Error while decoding login request: ", err.Error())
		w.WriteHeader(http.StatusBadRequest)
	} else {
		token, appErr := h.service.Login(loginRequest)
		if appErr != nil {
			writeResponse(w,appErr.Code, appErr.AsMessage())
		} else {
			writeResponse(w,http.StatusOK,*token)
		}
	}
}

func (h AuthHandler) Verify (w http.ResponseWriter, r *http.Request){
	urlParams := make(map[string]string)

	for k := range r.URL.Query() {
		urlParams[k] = r.URL.Query().Get(k)
	}

	if urlParams["token"] != "" {
		appErr := h.service.Verify(urlParams)
		if appErr != nil {
			writeResponse(w,appErr.Code, notAuthorizedResponse(appErr.Message))
		} else {
			writeResponse(w,http.StatusOK,authorizedResponse())
		}
	} else {
		writeResponse(w,http.StatusForbidden,notAuthorizedResponse("missing token"))
	}
}

func (h AuthHandler) Refresh(w http.ResponseWriter, r *http.Request){
	var refreshRequest dto.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&refreshRequest); err != nil {
		logger.Error("Error while decoding refresh roken request: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
	} else {
		token, appErr := h.service.Refresh(refreshRequest)
		if appErr != nil {
			writeResponse(w,appErr.Code, appErr.AsMessage())
		} else {
			writeResponse(w,http.StatusOK,*token)
		}
	}
}

func notAuthorizedResponse(msg string) map[string]interface{}{
	return map[string]interface{}{
		"isAuthorized": false,
		"message": msg,
	}
}

func authorizedResponse() map[string]bool{
	return map[string]bool{"isAuthorize": true}
}

func writeResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if err:= json.NewEncoder(w).Encode(data); err != nil {
		panic(err)
	}
}