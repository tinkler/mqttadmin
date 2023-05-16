package status

import (
	"net/http"

	errzpb "github.com/tinkler/mqttadmin/errz/v1"
)

func HttpError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}
	switch s := err.(type) {
	case *Status:
		if s.Code == http.StatusOK {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":1,"message":"` + s.Message + `"}`))
			return true
		}
		http.Error(w, s.Error(), int(s.Code))
		return true
	case *errzpb.ValidateError:
		http.Error(w, s.Message, http.StatusBadRequest)
		return true
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return true
	}
}
