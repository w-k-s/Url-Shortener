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
	bytes, err := json.Marshal(Error{Message: error})
	if err != nil {
		panic(err)
	}
	http.Error(w, string(bytes), code)
}
