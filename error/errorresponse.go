package error

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func SendError(w http.ResponseWriter, code int, e Err) {

	encoder := json.NewEncoder(w)

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(code)
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
