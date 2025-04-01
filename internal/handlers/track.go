package handlers

import (
	"net/http"
)

func (h *Handler) TrackHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		http.Error(w, "path parameter is required", http.StatusBadRequest)
		return
	}

	referrer := r.URL.Query().Get("referrer")
	userAgent := r.UserAgent()
	ip := r.RemoteAddr

	if err := h.store.TrackPageView(path, referrer, userAgent, ip); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
