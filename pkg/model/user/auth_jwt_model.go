package user

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/tinkler/mqttadmin/errz"
	errzpb "github.com/tinkler/mqttadmin/errz/v1"
	"github.com/tinkler/mqttadmin/pkg/kv"
	"github.com/tinkler/mqttadmin/pkg/logger"
)

const (
	TokenExpiredDuration = time.Hour * 24 * 30
	TokenDeviceKey       = "dev"
	TokenHasuraKey       = "https://hasura.io/jwt/claims"
)

var (
	keyStr = []byte("d1489e2df216240887f41014a7ac8fd6")
)

type HasuraClaims struct {
	Roles  []string `json:"x-hasura-allowed-roles"`
	Role   string   `json:"x-hasura-default-role"`
	UserID string   `json:"x-hasura-user-id"`
}

func getJwtToken(accountID, deviceID string) (string, error) {
	t := jwt.New()
	t.Set(jwt.SubjectKey, "https://mqtt.sfs.ink")
	ec := HasuraClaims{
		Roles:  []string{"user"},
		Role:   "user",
		UserID: accountID,
	}
	expireTime := time.Now().Add(TokenExpiredDuration)
	t.Set(TokenHasuraKey, ec)
	t.Set(jwt.ExpirationKey, expireTime)
	// 当区分设备时加入设备UUID前10位值
	if len(deviceID) >= 10 {
		t.Set(TokenDeviceKey, deviceID[:10])
	}

	signed, err := jwt.Sign(t, jwt.WithKey(jwa.HS256, keyStr))
	if err != nil {
		return "", err
	}

	token := string(signed)
	setShortTokenKV(accountID, deviceID, token[:10], expireTime)
	return token, nil
}

// shortTokenCache 短token缓存
// 用于存储短token，key为accountID，value为deviceID
// 用于校验短token是否和服务器一致
type shortTokenCache map[string]string

func (tc shortTokenCache) Set(deviceID string, shortToken string) {
	tc[deviceID] = shortToken
}

func (tc shortTokenCache) Get(deviceID string) string {
	return tc[deviceID]
}

func (tc shortTokenCache) String() string {
	byt, _ := json.Marshal(tc)
	return string(byt)
}

func getShortTokenKV(accountID, deviceID string) string {
	kstr := shortTokenKvKey(accountID)
	tv, ok := kv.Get(kstr)
	if !ok || tv == "" {
		return ""
	}

	tc := make(shortTokenCache)
	if err := json.Unmarshal([]byte(tv), &tc); err != nil {
		logger.Error(err)
		return ""
	}
	return tc[accountID]
}

func setShortTokenKV(accountID, deviceID string, token string, expireTime time.Time) {
	kstr := shortTokenKvKey(accountID)
	tv, ok := kv.Get(kstr)
	if !ok || tv == "" {
		tc := make(shortTokenCache)
		tc.Set(deviceID, token)
		kv.Set(kstr, tc.String(), expireTime)
	} else {
		tc := make(shortTokenCache)
		if err := json.Unmarshal([]byte(tv), &tc); err != nil {
			logger.Error(err)
			tc := make(shortTokenCache)
			tc.Set(deviceID, token)
			kv.Set(kstr, tc.String(), expireTime)
		} else {
			tc.Set(deviceID, token)
			kv.Set(kstr, tc.String(), expireTime)
		}
	}
}

func ClearShortTokenKV(accountID string) {
	kstr := shortTokenKvKey(accountID)
	kv.Set(kstr, "", time.Now())
}

func shortTokenKvKey(accountID string) string {
	return fmt.Sprintf("auth:%s:", accountID)
}

func CheckJwtToken(ctx context.Context, token string) (context.Context, error) {
	if len(token) <= 10 {
		return ctx, errz.ErrAuth(errzpb.AuthError_TOKEN_INVALID)
	}
	verifiedToken, err := jwt.Parse([]byte(token), jwt.WithKey(jwa.HS256, keyStr))
	if err != nil {
		switch err.Error() {
		case `"exp" not satisfied`:
			return ctx, errz.ErrAuth(errzpb.AuthError_TOKEN_EXPIRED)
		}
		logger.Error(err)
		return ctx, errz.ErrAuth(errzpb.AuthError_TOKEN_INVALID)
	}
	v1, expireTimeOk := verifiedToken.Get(jwt.ExpirationKey)
	if !expireTimeOk {
		return ctx, errz.ErrAuth(errzpb.AuthError_TOKEN_INVALID)
	}
	if expireTime, expirTimeOk := v1.(time.Time); !expirTimeOk {
		return ctx, errz.ErrAuth(errzpb.AuthError_TOKEN_INVALID)
	} else {
		if expireTime.Before(time.Now()) {
			return ctx, errz.ErrAuth(errzpb.AuthError_TOKEN_EXPIRED)
		}
	}

	v2, hasuraOk := verifiedToken.Get(TokenHasuraKey)
	if !hasuraOk {
		return ctx, errz.ErrAuth(errzpb.AuthError_TOKEN_INVALID)
	}
	var hasuraClaims HasuraClaims
	if hasuraData, hasuraOk := v2.(map[string]interface{}); !hasuraOk {
		return ctx, errz.ErrAuth(errzpb.AuthError_TOKEN_INVALID)
	} else {
		if rolesInter, ok := hasuraData["x-hasura-allowed-roles"].([]interface{}); ok {
			roles := make([]string, len(rolesInter))
			for i := range rolesInter {
				roles[i] = rolesInter[i].(string)
			}
			hasuraClaims.Roles = roles
		} else {
			logger.Debug("%T", hasuraData["x-hasura-allowed-roles"])
			return ctx, errz.ErrAuth(errzpb.AuthError_TOKEN_INVALID)
		}
		if role, ok := hasuraData["x-hasura-default-role"].(string); ok {
			hasuraClaims.Role = role
		} else {
			return ctx, errz.ErrAuth(errzpb.AuthError_TOKEN_INVALID)
		}
		if id, ok := hasuraData["x-hasura-user-id"].(string); ok {
			hasuraClaims.UserID = id
		} else {
			return ctx, errz.ErrAuth(errzpb.AuthError_TOKEN_INVALID)
		}
	}

	auth := authMiddleware{
		claims: &hasuraClaims,
	}
	if deviceIdInt, deviceOk := verifiedToken.Get(TokenDeviceKey); deviceOk {
		if deviceID, ok := deviceIdInt.(string); !ok {
			logger.Debug("%T", deviceIdInt)
			return ctx, errz.ErrAuth(errzpb.AuthError_TOKEN_INVALID)
		} else {
			// 不存在,令牌被刷新
			shortToken := getShortTokenKV(hasuraClaims.UserID, deviceID)
			if shortToken != "" {
				return ctx, errz.ErrAuth(errzpb.AuthError_TOKEN_EXPIRED)
			}
			// 被重新登录刷新
			if shortToken != token[len(token)-10:] {
				return ctx, errz.ErrAuth(errzpb.AuthError_TOKEN_EXPIRED)
			}

			// 存设备id
			auth.deviceID = deviceID
		}
	} else {
		// 不存在,令牌被刷新
		shortToken := getShortTokenKV(hasuraClaims.UserID, "")
		if shortToken != "" {
			return ctx, errz.ErrAuth(errzpb.AuthError_TOKEN_EXPIRED)
		}
		// 被重新登录刷新
		if shortToken != token[len(token)-10:] {
			return ctx, errz.ErrAuth(errzpb.AuthError_TOKEN_EXPIRED)
		}
	}

	return context.WithValue(ctx, authMiddlewareKey, &auth), nil
}
