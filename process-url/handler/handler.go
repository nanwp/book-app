package handler

import (
	"byfood-interview/process-url/service"
	"encoding/json"
	"net/http"
)

type processReq struct {
	URL       string `json:"url"`
	Operation string `json:"operation"`
}

type processResp struct {
	ProcessedURL string `json:"processed_url"`
}

type errResp struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// CleanupURL Godoc
// @Summary Cleanup a URL
// @Description Cleanup a URL by applying the specified operation
// @Tags URLs
// @Accept json
// @Produce json
// @Param request body processReq true "URL Cleanup Request"
// @Success 200 {object} processResp
// @Failure 400 {object} errResp
// @Router /api/v1/process-url [post]
func ProcessURLHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSON(w, http.StatusMethodNotAllowed, errResp{Error: "method not allowed"})
			return
		}

		var req processReq
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, errResp{Error: "invalid JSON body"})
			return
		}

		if req.URL == "" {
			writeJSON(w, http.StatusBadRequest, errResp{Error: "url is required"})
			return
		}

		op, err := service.ParseOperation(req.Operation)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, errResp{Error: err.Error()})
			return
		}

		out, err := service.ProcessURL(req.URL, op)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, errResp{Error: err.Error()})
			return
		}

		writeJSON(w, http.StatusOK, processResp{ProcessedURL: out})
	}
}
