package handlers

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/antontuzov/minimalytics/internal/storage"
)

type Handler struct {
	store storage.Storage
}

func NewHandler(store storage.Storage) *Handler {
	return &Handler{store: store}
}

func (h *Handler) APIHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/api/stats/daily":
		h.dailyStatsHandler(w, r)
	case "/api/stats/unique-visits":
		h.uniqueVisitsHandler(w, r)
	case "/api/stats/top-pages":
		h.topPagesHandler(w, r)
	case "/api/stats/referrers":
		h.referrerStatsHandler(w, r)
	case "/api/stats/devices":
		h.deviceStatsHandler(w, r)
	case "/api/stats/browsers":
		h.browserStatsHandler(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h *Handler) dailyStatsHandler(w http.ResponseWriter, r *http.Request) {
	stats, err := h.store.GetDailyStats()
	h.handleStatsResponse(w, stats, err)
}

func (h *Handler) uniqueVisitsHandler(w http.ResponseWriter, r *http.Request) {
	stats, err := h.store.GetUniqueVisits()
	h.handleStatsResponse(w, stats, err)
}

func (h *Handler) topPagesHandler(w http.ResponseWriter, r *http.Request) {
	stats, err := h.store.GetTopPages()
	h.handleStatsResponse(w, stats, err)
}

func (h *Handler) referrerStatsHandler(w http.ResponseWriter, r *http.Request) {
	stats, err := h.store.GetReferrers()
	h.handleStatsResponse(w, stats, err)
}

func (h *Handler) deviceStatsHandler(w http.ResponseWriter, r *http.Request) {
	stats, err := h.store.GetDevices()
	h.handleStatsResponse(w, stats, err)
}

func (h *Handler) browserStatsHandler(w http.ResponseWriter, r *http.Request) {
	stats, err := h.store.GetBrowsers()
	h.handleStatsResponse(w, stats, err)
}

func (h *Handler) handleStatsResponse(w http.ResponseWriter, data interface{}, err error) {
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(data)
}

// Dashboard handler
func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("dashboard.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}
