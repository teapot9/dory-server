package controller

import (
	"encoding/json"
	"net/http"
	"service"
	"structures"
)

func Reinitialize(w http.ResponseWriter, r *http.Request){

	HandlePreflight(w, "POST")
	if r.Method == "OPTIONS" {
		return
	}

	//Decoding JSON
	var user structures.User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, "Could not decode JSON data !", http.StatusBadRequest)
		return
	}

	//Check that all fields we need are present
	if user.Username == "" {
		http.Error(w, "Missing username in payload.", http.StatusBadRequest)
		return
	}

	if user.Token == "" {
		http.Error(w, "Missing token in payload.", http.StatusBadRequest)
		return
	}

	if user.NewPassword == "" {
		http.Error(w, "Missing new password in payload.", http.StatusBadRequest)
		return
	}

	//Contact service
	serviceError := service.ReinitializePassword(user.Username, user.Token, user.NewPassword)

	if serviceError.Error != nil {
		http.Error(w, serviceError.Error.Error(), serviceError.HttpCode)
		return
	}

	w.WriteHeader(200)

}