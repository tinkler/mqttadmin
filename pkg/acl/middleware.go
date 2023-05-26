package acl

import (
	"context"
	"net/http"
	"strings"

	"github.com/tinkler/mqttadmin/pkg/qm"
	"github.com/tinkler/mqttadmin/pkg/status"
)

const (
	QueueRoleAdd    = "queue_role_add"
	QueueRoleRemove = "queue_role_remove"
	QueueTokenCheck = "queue_token_check"
)

type middlewareKey string

var (
	authMiddlewareKey middlewareKey = "auth"
)

// builtin roles
const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

type authorization struct {
	claims   HasuraClaims
	deviceID string
}

// HasRole check if the user has the roles
// Auto add admin role to check if the user is admin
func HasRole(ctx context.Context, roles ...string) bool {
	if ctx == nil {
		return false
	}
	roles = append(roles, RoleAdmin)
	if auth, ok := ctx.Value(authMiddlewareKey).(*authorization); ok {
		bMap := make(map[string]bool)
		for _, r := range auth.claims.Roles {
			bMap[r] = true
		}
		for _, r := range roles {
			if bMap[r] {
				return true
			}
		}
	}
	return false
}

func GetAllRoles(ctx context.Context) []string {
	if ctx == nil {
		return nil
	}
	if auth, ok := ctx.Value(authMiddlewareKey).(*authorization); ok {
		return auth.claims.Roles
	}
	return nil
}

func GetUserID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if auth, ok := ctx.Value(authMiddlewareKey).(*authorization); ok {
		return auth.claims.UserID
	}
	return ""
}

func AddRole(userID string, roleID ...string) error {
	_, err := qm.Publish(QueueRoleAdd, userID+":"+strings.Join(roleID, ","))
	return err
}

func RemoveRole(userID string, roleID ...string) error {
	_, err := qm.Publish(QueueRoleRemove, userID+":"+strings.Join(roleID, ","))
	return err
}

type AuthConfig struct {
	NoNeedAuth bool
}

func WrapAuth(c AuthConfig) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
			ctx := r.Context()
			if token != "" {
				ctx, _ = checkJwtToken(ctx, token)
			} else if !c.NoNeedAuth {
				status.HttpError(w, status.StatusUnauthorized())
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
