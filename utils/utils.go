package utils

import (
	"net/http"
	"strconv"
)

func ParseLimitOffset(r *http.Request) (int, int) {
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")
	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)
	if limit == 0 {
		limit = 100
	}
	return limit, offset
}
