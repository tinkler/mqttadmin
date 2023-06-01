package acl

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/tinkler/mqttadmin/errz"
	errzpb "github.com/tinkler/mqttadmin/errz/v1"
	"github.com/tinkler/mqttadmin/pkg/kv"
	"github.com/tinkler/mqttadmin/pkg/logger"
	"github.com/tinkler/mqttadmin/pkg/qm"
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

func GetJwtToken(accountID, deviceID string, roles []string) (string, error) {
	t := jwt.New()
	if len(roles) == 0 {
		roles = []string{"user"}
	}
	t.Set(jwt.SubjectKey, "https://mqtt.sfs.ink")
	ec := HasuraClaims{
		Roles:  roles,
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
	setShortToken(accountID, deviceID, token[:10], expireTime)
	return token, nil
}

// shortTokenCache 短token缓存
// 用于存储短token，key为短token，value为deviceID
// 用于校验短token是否存在
// 当短token存在时，校验短token是否匹配deviceID
// 如果短token的值前面带-符号，表示设备需要刷新token
type shortTokenCache map[string]string

func (tc shortTokenCache) SetShortToken(deviceID string, shortToken string) {
	// clear the old short token
	for k := range tc {
		if tc[k] == deviceID {
			delete(tc, k)
		}
	}
	tc[shortToken] = deviceID
}

func (tc shortTokenCache) AddRemoveFlag(deviceID string) {
	for k := range tc {
		if tc[k] == deviceID {
			tc[k] = "-" + tc[k]
		}
	}
}

func (tc shortTokenCache) GetDeviceID(shortToken string) (deviceID string) {
	return tc[shortToken]
}

func (tc shortTokenCache) GetShortToken(deviceID string) (shortToken string) {
	for k := range tc {
		if tc[k] == deviceID {
			return k
		}
	}
	return ""
}

func (tc shortTokenCache) String() string {
	byt, _ := json.Marshal(tc)
	return string(byt)
}

func getDeviceID(accountID, shortToken string) string {
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
	return tc[shortToken]
}

func setShortToken(accountID, deviceID string, shortToken string, expireTime time.Time) {
	kstr := shortTokenKvKey(accountID)
	tv, ok := kv.Get(kstr)
	if !ok || tv == "" {
		tc := make(shortTokenCache)
		tc.SetShortToken(deviceID, shortToken)
		kv.Set(kstr, tc.String(), expireTime)
	} else {
		tc := make(shortTokenCache)
		if err := json.Unmarshal([]byte(tv), &tc); err != nil {
			logger.Error(err)
			tc := make(shortTokenCache)
			tc.SetShortToken(deviceID, shortToken)
			kv.Set(kstr, tc.String(), expireTime)
		} else {
			tc.SetShortToken(deviceID, shortToken)
			kv.Set(kstr, tc.String(), expireTime)
		}
	}
}

// 标记设备需要重新登录
func SetDeviceRemoveFlag(accountID, deviceID string) {
	kstr := shortTokenKvKey(accountID)
	tv, ok := kv.Get(kstr)
	if !ok || tv == "" {
		return
	}

	tc := make(shortTokenCache)
	if err := json.Unmarshal([]byte(tv), &tc); err != nil {
		logger.Error(err)
		return
	}
	tc.AddRemoveFlag(deviceID)
	kv.Set(kstr, tc.String(), time.Now().Add(TokenExpiredDuration))
}

// 标记所有设备需要重新登录
func SetAllDeviceRemoveFlag(accountID string) {
	kstr := shortTokenKvKey(accountID)
	tv, ok := kv.Get(kstr)
	if !ok || tv == "" {
		return
	}

	tc := make(shortTokenCache)
	if err := json.Unmarshal([]byte(tv), &tc); err != nil {
		logger.Error(err)
		return
	}
	for k := range tc {
		tc[k] = "-" + tc[k]
	}
	kv.Set(kstr, tc.String(), time.Now().Add(TokenExpiredDuration))
}

func shortTokenKvKey(accountID string) string {
	return fmt.Sprintf("auth:%s:", accountID)
}

// CheckDeviceToken 检查设备token
// 如果token中包含设备ID，检查设备ID是否匹配
func CheckDeviceToken(ctx context.Context, deviceToken string, tokenStr string) (context.Context, error) {
	if len(tokenStr) <= 10 || deviceToken == "" {
		return ctx, errz.ErrAuth(errzpb.AuthError_TOKEN_INVALID)
	}
	token, auth, err := parseTokenStr(tokenStr, ctx)
	if err != nil {
		return ctx, err
	}
	if deviceIdInt, deviceOk := token.Get(TokenDeviceKey); deviceOk {
		if deviceID, ok := deviceIdInt.(string); !ok {
			logger.Debug("%T", deviceIdInt)
			return ctx, errz.ErrAuth(errzpb.AuthError_TOKEN_INVALID)
		} else {
			if deviceToken[:10] != deviceID {
				// the token is not owned by the device
				return ctx, errz.ErrAuth(errzpb.AuthError_TOKEN_INVALID)
			}
			var cacheDeviceID string
			if Aclm().mode == AM_SERVER {
				cacheDeviceID = getDeviceID(auth.claims.UserID, tokenStr[:10])
			} else {
				cacheDeviceID, err = qm.PublishAndReceive(QueueTokenCheck, auth.claims.UserID+":"+tokenStr[:10])
				if err != nil {
					logger.Error(err)
					return ctx, errz.ErrInternal(err)
				}
			}
			if strings.TrimPrefix(cacheDeviceID, "-") != deviceToken {
				// the token is not owned by the device
				return ctx, errz.ErrAuth(errzpb.AuthError_TOKEN_INVALID)
			}

			// 存设备id
			auth.deviceID = deviceToken
			return context.WithValue(ctx, authMiddlewareKey, auth), nil
		}
	} else {
		return ctx, errz.ErrAuth(errzpb.AuthError_TOKEN_INVALID)
	}
}

func checkJwtToken(ctx context.Context, tokenStr string) (context.Context, error) {
	if len(tokenStr) <= 10 {
		return ctx, errz.ErrAuth(errzpb.AuthError_TOKEN_INVALID)
	}
	_, auth, err := parseTokenStr(tokenStr, ctx)
	if err != nil {
		return ctx, err
	}

	var cacheDeviceID string
	if Aclm().mode == AM_SERVER {
		cacheDeviceID = getDeviceID(auth.claims.UserID, tokenStr[:10])
	} else {
		cacheDeviceID, err = qm.PublishAndReceive(QueueTokenCheck, auth.claims.UserID+":"+tokenStr[:10])
		if err != nil {
			return ctx, errz.ErrInternal(err)
		}
	}
	if strings.HasPrefix(cacheDeviceID, "-") || cacheDeviceID == "" {
		return ctx, errz.ErrAuth(errzpb.AuthError_TOKEN_EXPIRED)
	}

	return context.WithValue(ctx, authMiddlewareKey, auth), nil
}

func parseTokenStr(token string, ctx context.Context) (jwt.Token, *authorization, error) {
	verifiedToken, err := jwt.Parse([]byte(token), jwt.WithKey(jwa.HS256, keyStr))
	if err != nil {
		switch err.Error() {
		case `"exp" not satisfied`:
			return nil, nil, errz.ErrAuth(errzpb.AuthError_TOKEN_INVALID)
		}
		logger.Error(err)
		return nil, nil, errz.ErrAuth(errzpb.AuthError_TOKEN_INVALID)
	}
	v1, expireTimeOk := verifiedToken.Get(jwt.ExpirationKey)
	if !expireTimeOk {
		return nil, nil, errz.ErrAuth(errzpb.AuthError_TOKEN_INVALID)
	}
	if expireTime, expirTimeOk := v1.(time.Time); !expirTimeOk {
		return nil, nil, errz.ErrAuth(errzpb.AuthError_TOKEN_INVALID)
	} else {
		if expireTime.Before(time.Now()) {
			return nil, nil, errz.ErrAuth(errzpb.AuthError_TOKEN_INVALID)
		}
	}

	v2, hasuraOk := verifiedToken.Get(TokenHasuraKey)
	if !hasuraOk {
		return nil, nil, errz.ErrAuth(errzpb.AuthError_TOKEN_INVALID)
	}
	var hasuraClaims HasuraClaims
	if hasuraData, hasuraOk := v2.(map[string]interface{}); !hasuraOk {
		return nil, nil, errz.ErrAuth(errzpb.AuthError_TOKEN_INVALID)
	} else {
		if rolesInter, ok := hasuraData["x-hasura-allowed-roles"].([]interface{}); ok {
			roles := make([]string, len(rolesInter))
			for i := range rolesInter {
				roles[i] = rolesInter[i].(string)
			}
			hasuraClaims.Roles = roles
		} else {
			logger.Debug("%T", hasuraData["x-hasura-allowed-roles"])
			return nil, nil, errz.ErrAuth(errzpb.AuthError_TOKEN_INVALID)
		}
		if role, ok := hasuraData["x-hasura-default-role"].(string); ok {
			hasuraClaims.Role = role
		} else {
			return nil, nil, errz.ErrAuth(errzpb.AuthError_TOKEN_INVALID)
		}
		if id, ok := hasuraData["x-hasura-user-id"].(string); ok {
			hasuraClaims.UserID = id
		} else {
			return nil, nil, errz.ErrAuth(errzpb.AuthError_TOKEN_INVALID)
		}
	}

	auth := &authorization{
		claims: hasuraClaims,
	}
	return verifiedToken, auth, nil
}
