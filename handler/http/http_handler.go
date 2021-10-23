package http_handler

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rehandwi03/stockbit-question-2/service"
	"log"
	"net/http"
	"strconv"
)

type httpHandler struct {
	service service.Service
}

func NewMovieHttpHandler(router *mux.Router, service service.Service) {
	handler := httpHandler{service: service}
	router.HandleFunc("/movies", handler.fetch).Methods("GET")
	router.HandleFunc("/movies/{id}", handler.getById).Methods("GET")
}

func (h *httpHandler) getById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		log.Println("id query param can't be null")
		h.resError(w, "id query param can't be null", http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	res, err := h.service.GetByID(ctx, id)
	if err != nil {
		log.Printf("erorr fetch in handler: %v", err)
		h.resError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.resSuccess(w, "success", http.StatusOK, res)
	return
}

func (h *httpHandler) fetch(w http.ResponseWriter, r *http.Request) {
	var pagination, searchWord string
	pagination = r.URL.Query().Get("pagination")
	if pagination == "" {
		pagination = "1"
	}
	searchWord = r.URL.Query().Get("searchword")
	if searchWord == "" {
		log.Println("searchword query param can't be null")
		h.resError(w, "searchword query param can't be null", http.StatusBadRequest)
		return
	}

	params := map[string]interface{}{
		"query": map[string]interface{}{
			"s":    searchWord,
			"page": pagination,
		},
	}

	ctx := context.Background()

	res, err := h.service.Fetch(ctx, params)
	if err != nil {
		log.Printf("erorr fetch in handler: %v", err)
		h.resError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.resSuccess(w, "success", http.StatusOK, res)
	return
}

func (h *httpHandler) resSuccess(w http.ResponseWriter, message string, httpCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	response := map[string]interface{}{
		"message": message,
		"status":  "success",
		"code":    strconv.Itoa(httpCode),
		"data":    data,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Fatalf("error: %v", err)
	}
}

func (h *httpHandler) resError(w http.ResponseWriter, message string, httpCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	response := map[string]interface{}{
		"message": message,
		"status":  "failed",
		"code":    strconv.Itoa(httpCode),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Fatalf("error: %v", err)
	}
}
