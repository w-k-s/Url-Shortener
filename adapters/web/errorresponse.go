package web

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func SendError(w http.ResponseWriter, e Err) {

	encoder := json.NewEncoder(w)

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(httpStatusCode(e.Code()))
	err := encoder.Encode(map[string]interface{}{
		"code":    e.Code(),
		"message": e.Error(),
		"domain":  e.Domain(),
		"fields":  e.Fields(),
	})
	if err != nil {
		http.Error(
			w,
			fmt.Sprintf("Error encoding %v. Cause: %s", e, err),
			http.StatusInternalServerError,
		)
	}
}

func httpStatusCode(e err.Code) int {
	switch e {
	case ShortenURLValidation:
		fallthrough
	case RetrieveFullURLValidation:
		fallthrough
	case ShortenURLShortIdInUse:
		return http.StatusBadRequest
	case RetrieveFullURLNotFound:
		fallthrough
	case RedirectionFullURLNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
