package user

import "context"

type middlewareKey string

var (
	authMiddlewareKey middlewareKey = "auth"
)

// builtin roles
const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

type authMiddleware struct {
	claims   *HasuraClaims
	deviceID string
}

func HasRole(ctx context.Context, role string) bool {
	if ctx == nil {
		return false
	}
	if auth, ok := ctx.Value(authMiddlewareKey).(*authMiddleware); ok {
		for _, r := range auth.claims.Roles {
			if r == role {
				return true
			}
		}
	}
	return false
}

func GetUserID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if auth, ok := ctx.Value(authMiddlewareKey).(*authMiddleware); ok {
		return auth.claims.UserID
	}
	return ""
}
