// Code generated by github.com/tinkler/mqttadmin; DO NOT EDIT.
package route
import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/tinkler/mqttadmin/pkg/jsonz/sjson"
	"github.com/tinkler/mqttadmin/pkg/model/role"
	"github.com/tinkler/mqttadmin/pkg/status"
)

func RoutesRole(m chi.Router) {
	m.Route("/role", func(r chi.Router) {
		
		r.Post("/role/save-role", func(w http.ResponseWriter, r *http.Request) {
			m := Model[*role.Role, any]{}
			err := sjson.Bind(r, &m)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			res := Res[*role.Role,any]{Data:m.Data}
			err = m.Data.SaveRole(r.Context())
			
			if status.HttpError(w, err) {
				return
			}
			if sjson.HttpWrite(w, res) {
				return
			}

		})
	})
}
