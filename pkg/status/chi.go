package status

import "net/http"

func HttpError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}
	if s, ok := err.(*Status); ok {
		if s.Code == http.StatusOK {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"code":1,"message":"` + s.Message + `"}`))
			return true
		}
		http.Error(w, s.Error(), int(s.Code))
		return true
	}
	http.Error(w, err.Error(), http.StatusInternalServerError)
	return true
}
