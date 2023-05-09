package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	log.Println("Authenticate")
	err := app.readJSON(w, r, &requestPayload)

	if err != nil {
		_ = app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// validate the user
	user, err := app.Models.User.GetByEmail(requestPayload.Email)

	log.Println("User: ", user)
	log.Println("Error: ", err)

	if err != nil {
		_ = app.errorJSON(w, errors.New("Invalid credentials"), http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)

	if err != nil || !valid {
		_ = app.errorJSON(w, errors.New("Invalid credentials"), http.StatusBadRequest)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	log.Println("response: ", payload)
	log.Println("Niice")
	_ = app.writeJSON(w, http.StatusAccepted, payload)

}
