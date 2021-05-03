package api

import (
	"net/http"
	"strconv"
)

func findSkip(r *http.Request) int64 {
	skipParam := r.URL.Query().Get("skip")
	if skipParam == "" {
		skipParam = "0"
	}

	limit, err := strconv.ParseInt(skipParam, 10, 64)
	if err != nil {
		return 0
	}

	return limit
}
