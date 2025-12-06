package router

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func jsonResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func jsonError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

func strToInt(s string, d *int64) *int64 {
	if v, err := strconv.ParseInt(s, 10, 64); err == nil {
		return &v
	}
	return d
}
func strToInts(s []string) []int64 {
	vs := make([]int64, 0, len(s))
	for _, s := range s {
		if v := strToInt(s, nil); v != nil {
			vs = append(vs, *v)
		}
	}
	return vs
}
func intPtr(i int64) *int64 {
	return &i
}
