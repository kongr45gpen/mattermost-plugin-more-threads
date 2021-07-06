package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func (p *Plugin) SetupRouting() {
	root := mux.NewRouter()
	api := root.PathPrefix("/api/v0").Subrouter()
	api.Use(MattermostAuthorizationRequired)
	api.Handle("{anything:.*}", http.NotFoundHandler())
	api.NotFoundHandler = http.NotFoundHandler()

	autofollow := api.PathPrefix("/autofollow/{channelid:[A-Za-z0-9]+}").Subrouter()
	autofollow.Use(p.PermissionsToChannel)
	autofollow.HandleFunc("", p.getMoreThreads).Methods(http.MethodGet)
	autofollow.HandleFunc("", p.setAutofollow).Methods(http.MethodPut, http.MethodDelete)

	p.root = root
}

func (p *Plugin) setAutofollow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	channelId := vars["channelid"]
	userId := r.Header.Get("Mattermost-User-ID")

	newFollow := r.Method == http.MethodPut

	newFollowByte := byte(0)
	if newFollow {
		newFollowByte = byte(1)
	}

	p.API.KVSet(GetKVkey(userId, channelId), []byte{newFollowByte})

	w.WriteHeader(http.StatusOK)
}

type MoreThreadsResponse struct {
	AutoFollowing bool `json:"auto_following"`
}

func (p *Plugin) getMoreThreads(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	channelId := vars["channelid"]
	userId := r.Header.Get("Mattermost-User-ID")

	followBytes, err := p.API.KVGet(GetKVkey(userId, channelId))
	if err != nil {
		HandleError(w, err)
		return
	}

	autoFollowing := false
	if len(followBytes) >= 1 && followBytes[0] == 1 {
		autoFollowing = true
	}

	resp := MoreThreadsResponse{
		AutoFollowing: autoFollowing,
	}

	ReturnJSON(w, &resp, http.StatusOK)
}

func (p *Plugin) PermissionsToChannel(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		channelID := vars["channelid"]
		userID := r.Header.Get("Mattermost-User-ID")

		_, err := p.API.GetChannel(channelID)
		if err != nil {
			http.Error(w, "Not a channel", http.StatusBadRequest)
			return
		}

		if p.API.HasPermissionToChannel(userID, channelID, model.PERMISSION_READ_CHANNEL) {
			next.ServeHTTP(w, r)
			return
		}

		http.Error(w, "No Permissions to Channel", http.StatusForbidden)
	})
}

// Copy pasted tools. Really should be in a library.

// MattermostAuthorizationRequired checks if request is authorized.
func MattermostAuthorizationRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("Mattermost-User-ID")
		if userID != "" {
			next.ServeHTTP(w, r)
			return
		}

		http.Error(w, "Not authorized", http.StatusUnauthorized)
	})
}

// ReturnJSON writes the given pointer to object as json with a success response
func ReturnJSON(w http.ResponseWriter, pointerToObject interface{}, httpStatus int) {
	jsonBytes, err := json.Marshal(pointerToObject)
	if err != nil {
		HandleError(w, errors.Wrapf(err, "unable to marshal json"))
		return
	}

	w.WriteHeader(httpStatus)
	if _, err = w.Write(jsonBytes); err != nil {
		HandleError(w, err)
		return
	}
}

// HandleError logs the internal error and sends a generic error as JSON in a 500 response.
func HandleError(w http.ResponseWriter, internalErr error) {
	HandleErrorWithCode(w, http.StatusInternalServerError, "An internal error has occurred. Check app server logs for details.", internalErr)
}

// HandleErrorWithCode logs the internal error and sends the public facing error
// message as JSON in a response with the provided code.
func HandleErrorWithCode(w http.ResponseWriter, code int, publicErrorMsg string, internalErr error) {
	w.WriteHeader(code)

	details := ""
	if internalErr != nil {
		details = internalErr.Error()
	}

	loggedMsg, _ := json.Marshal(struct {
		Message string `json:"message"` // A public facing message providing details about the error.
		Details string `json:"details"` // More details, potentially sensitive, about the error.
	}{
		Message: publicErrorMsg,
		Details: details,
	})
	logrus.Warn(string(loggedMsg))

	responseMsg, _ := json.Marshal(struct {
		Error string `json:"error"` // A public facing message providing details about the error.
	}{
		Error: publicErrorMsg,
	})
	_, _ = w.Write(responseMsg)
}
