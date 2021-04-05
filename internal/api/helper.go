package api

import (
	"net/http"
	"strconv"
)

func findLimit(r *http.Request) int64 {
	limitParam := r.URL.Query().Get("limit")
	if limitParam == "" {
		limitParam = "50"
	}
	limit, err := strconv.ParseInt(limitParam, 10, 64)
	if err != nil {
		return 50
	}
	return limit
}
