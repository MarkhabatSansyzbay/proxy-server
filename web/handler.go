package web

import (
	"encoding/json"
	"errors"
	"net/http"

	"proxyserver/internal/models"
	"proxyserver/internal/service"
)

type Handler struct {
	proxy service.ReverseProxy
}

func NewHandler() *Handler {
	return &Handler{
		proxy: service.NewProxy(),
	}
}

func (h *Handler) InitRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", h.homePage)

	return mux
}

func (h *Handler) homePage(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	case http.MethodPost:
		var req *models.ClientMessage

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		proxyResponse, err := h.proxy.ResponseToMessage(req)
		if err != nil {
			if errors.Is(err, service.ErrMethod) {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		if err := json.NewEncoder(w).Encode(proxyResponse); err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}
