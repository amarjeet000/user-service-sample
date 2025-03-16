package server

import (
	"encoding/json"
	"net/http"
)

func RespondWithData(w http.ResponseWriter, r *http.Request, httpStatus int, obj any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(httpStatus)
	if obj != nil {
		switch obj := obj.(type) {
		case []byte:
			_, err := w.Write(obj)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			return
		default:
			jsonRes, err := json.Marshal(obj)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			_, err = w.Write(jsonRes)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			return
		}
	}
}
