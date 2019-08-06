package web

import (
	"encoding/json"
	"fmt"
	"github.com/w-k-s/short-url/domain"
	"github.com/w-k-s/short-url/domain/urlshortener/usecase"
	"net/http"
)

func sendResponse(w http.ResponseWriter, status int, body interface{}) {

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	err := encoder.Encode(body)

	if err != nil {
		sendEncodingError(w, body, err)
		return
	}
}

func sendError(w http.ResponseWriter, e domain.Err) {

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
		sendEncodingError(w, e, err)
	}
}

func httpStatusCode(e domain.Code) int {
	switch e {
	case usecase.ShortenURLValidation:
		fallthrough
	case usecase.RetrieveFullURLValidation:
		fallthrough
	case usecase.ShortenURLShortIDInUse:
		return http.StatusBadRequest
	case usecase.RetrieveFullURLNotFound:
		fallthrough
	case usecase.RedirectionFullURLNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

func sendEncodingError(w http.ResponseWriter, encodee interface{}, err error) {
	http.Error(
		w,
		fmt.Sprintf("Error encoding '%v'. Cause: %s", encodee, err),
		http.StatusInternalServerError,
	)
}
