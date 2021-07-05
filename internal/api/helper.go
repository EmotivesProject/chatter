package api

import (
	"net/http"
	"strconv"
)

const (
	bit = 64
	ten = 10
)

func findSkip(r *http.Request) int64 {
	skipParam := r.URL.Query().Get("skip")
	if skipParam == "" {
		skipParam = "0"
	}

	limit, err := strconv.ParseInt(skipParam, ten, bit)
	if err != nil {
		return 0
	}

	return limit
}
