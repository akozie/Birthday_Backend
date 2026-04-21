package utils

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, code int, obj interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	enc := json.NewEncoder(w)
	enc.Encode(obj)
}

func SetCookie(w http.ResponseWriter, name string, value string) {
	cookie := http.Cookie{
		Name:  name,
		Value: value,
		Path:  "/",
	}
	http.SetCookie(w, &cookie)
}
