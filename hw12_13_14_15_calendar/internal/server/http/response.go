package internalhttp

import (
	"encoding/json"
	"net/http"
)

type DataResponse map[string]interface{}

// SendResponse Отправить ответ клиенту.
func SendResponse(w http.ResponseWriter, data DataResponse, caption string) {
	response, errEncode := json.Marshal(data)
	if errEncode != nil {
		logg.Errorf("Failed to serialize data to get %s. Error: %v", caption, errEncode)
		http.Error(w, errEncode.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application-json")
	_, errWrite := w.Write(response)
	if errWrite != nil {
		logg.Errorf("Failed to send %s data. Error: %v", caption, errWrite)
		http.Error(w, errEncode.Error(), http.StatusInternalServerError)
		return
	}
}
