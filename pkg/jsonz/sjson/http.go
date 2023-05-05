package sjson

import (
	"net/http"

	"github.com/tinkler/mqttadmin/pkg/status"
)

const (
	// ContentType is the content type for json
	ContentType = "application/json"
)

// HttpWrite writes the v to the http.ResponseWriter
func HttpWrite(w http.ResponseWriter, v any) bool {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", ContentType)
	byt, err := Marshal(map[string]interface{}{
		"code":    0,
		"message": "success",
		"data":    v,
	})
	if err != nil {
		status.HttpError(w, err)
		return true
	}
	_, err = w.Write(byt)
	if err != nil {
		status.HttpError(w, err)
		return true
	}
	return true
}
