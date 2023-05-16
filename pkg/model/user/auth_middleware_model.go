package user

type middlewareKey string

var (
	authMiddlewareKey middlewareKey = "auth"
)

type authMiddleware struct {
	claims   *HasuraClaims
	deviceID string
}
