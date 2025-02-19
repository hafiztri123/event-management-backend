package handler

import "net/http"

type EventHandler interface{
	CreateEvent(w http.ResponseWriter, r *http.Request)
	UpdateEvent(w http.ResponseWriter, r *http.Request)
	DeleteEvent(w http.ResponseWriter, r *http.Request)
	GetEvent(w http.ResponseWriter, r *http.Request)
	ListEvents(w http.ResponseWriter, r *http.Request)
	SearchEvents(w http.ResponseWriter, r *http.Request)
}