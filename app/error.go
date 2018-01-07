package app

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	Message string `json:"error"`
}

func (e Error) Error() string {
	return e.Message
}

func EncodeNewErrorJSON(w http.ResponseWriter, error string, code int) {
	
	encoder := json.NewEncoder(w)

	w.Header().Set("Content-Type","application/json;charset=utf-8")
	w.WriteHeader(code)
	err := encoder.Encode(Error{Message: error})
	if err != nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
	}
}
