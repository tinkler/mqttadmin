package sjson

import (
	"io"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func init() {
	json.RegisterExtension(&snakedNamedExtension{})
}

func Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

func Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func Bind(r *http.Request, v any) error {
	reqByt, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(reqByt, v)
}
