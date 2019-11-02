package web

import (
	"encoding/json"
	domain "github.com/w-k-s/short-url/domain"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendError(t *testing.T) {

	w := httptest.NewRecorder()

	inputErr := domain.NewError(
		10300,
		"test",
		"message",
		map[string]string{"FIELD": "VALUE"},
	)

	sendError(w, inputErr)

	resp := w.Result()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Send Error sends wrong status code, got: %d, want: %d.", resp.StatusCode, http.StatusBadRequest)
	}

	var outputErr map[string]interface{}
	decoder := json.NewDecoder(resp.Body)
	decoder.Decode(&outputErr)

	outCode := domain.Code(outputErr["code"].(float64))
	outDomain := outputErr["domain"].(string)
	outError := outputErr["message"].(string)
	outFields := outputErr["fields"].(map[string]interface{})

	if outCode != inputErr.Code() {
		t.Errorf("Incorrect err.Code, got: %d, want: %d.", outCode, inputErr.Code())
	}

	if outDomain != inputErr.Domain() {
		t.Errorf("Incorrect err.Domain, got: %s, want: %s.", outDomain, inputErr.Domain())
	}

	if outError != inputErr.Error() {
		t.Errorf("Incorrect err.Message, got: %s, want: %s.", outError, inputErr.Error())
	}

	if len(outFields) != len(inputErr.Fields()) {
		t.Errorf("Incorrect err.Fields length, got: %d, want: %d.", len(outFields), len(inputErr.Fields()))
	}

	if outFields["FIELD"] != inputErr.Fields()["FIELD"] {
		t.Error("Incorrect err.Fields keys")
	}
}
