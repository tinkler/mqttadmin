package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/tinkler/mqttadmin/errz"
	"github.com/tinkler/mqttadmin/pkg/acl"
	"github.com/tinkler/mqttadmin/pkg/db"
	"github.com/tinkler/mqttadmin/pkg/model/role"
	"github.com/tinkler/mqttadmin/pkg/status"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Auth struct {
	ID          string // UUID
	DeviceToken string `gorm:"-"` // UUID
	Username    string
	Password    string
	Token       string `gorm:"-"`
}

func (a *Auth) TableName() string {
	return "authv1.user"
}

// Signin sign in with username and password
// Require: Username, Password
// Optional: DeviceToken
func (a *Auth) Signin(ctx context.Context) (*Auth, error) {
	u := Auth{}
	se := db.DB().Where("username = ?", a.Username).First(&u)
	if se.Error != nil {
		if errors.Is(se.Error, gorm.ErrRecordNotFound) {
			return nil, status.OkCn("user not found", "用户不存在")
		}
		return nil, se.Error
	}

	if err := a.Compare(u.Password); err != nil {
		return nil, status.OkCn("invalid password", "密码错误")
	}
	a.ID = u.ID
	if a.DeviceToken == "" {
		a.DeviceToken = uuid.New().String()
	}

	// when login success, clear the token
	acl.SetDeviceRemoveFlag(a.ID, a.DeviceToken)

	user := &User{ID: a.ID}
	err := user.GetRoles(ctx)
	if err != nil {
		return nil, err
	}
	token, err := acl.GetJwtToken(a.ID, a.DeviceToken, user.GetRolesStrings())
	if err != nil {
		return nil, err
	}
	a.Password = ""
	a.Token = token

	return a, nil
}

// QuickSignin quick signin without password
// Require: DeviceToken, Token
func (a *Auth) QuickSignin(ctx context.Context) error {
	if a.DeviceToken == "" {
		return errz.ErrVdM("device_token", "device token is empty", "登录设备号不能为空")
	}
	if a.Token == "" {
		return errz.ErrVdM("token", "token is empty", "登录令牌不能为空")
	}
	var err error
	ctx, err = acl.CheckDeviceToken(ctx, a.DeviceToken, a.Token)
	if err != nil {
		return err
	}
	u := &User{ID: acl.GetUserID(ctx)}
	err = u.GetRoles(ctx)
	if err != nil {
		return err
	}
	a.ID = u.ID

	token, err := acl.GetJwtToken(a.ID, a.DeviceToken, u.GetRolesStrings())
	if err != nil {
		return err
	}
	a.Token = token
	return nil
}

func (a *Auth) Signup(ctx context.Context) (*Auth, error) {
	if a.Username == "" {
		return nil, status.OkCn("username is empty", "用户名不能为空")
	}
	if a.Password == "" {
		return nil, status.OkCn("password is empty", "密码不能为空")
	}

	se := db.DB().Where("username = ?", a.Username).First(&Auth{})
	if se.Error == nil {
		return nil, status.OkCn(ErrMsgNameAreadyExist, "用户名已存在")
	}
	if !errors.Is(gorm.ErrRecordNotFound, se.Error) {
		return nil, se.Error
	}

	hashed, err := a.Encrypt()
	if err != nil {
		return nil, err
	}
	se.Error = nil
	se = se.Begin()
	se.Omit("id").Create(&Auth{Username: a.Username, Password: hashed})
	if se.Error != nil {
		se.Rollback()
		return nil, se.Error
	}
	r := &Auth{Username: a.Username}
	se.First(r)
	if se.Error != nil {
		se.Rollback()
		return nil, se.Error
	}
	se.Commit()
	if se.Error != nil {
		return nil, se.Error
	}
	if r.DeviceToken == "" {
		r.DeviceToken = uuid.New().String()
	}

	// when login success, clear the token
	acl.SetDeviceRemoveFlag(r.ID, r.DeviceToken)

	u := &User{ID: r.ID}
	// new user must add role user
	err = addRole(u, &role.Role{Name: acl.RoleUser})
	if err != nil {
		return nil, err
	}

	err = u.GetRoles(ctx)
	if err != nil {
		return nil, err
	}
	token, err := acl.GetJwtToken(r.ID, r.DeviceToken, u.GetRolesStrings())
	if err != nil {
		return nil, err
	}
	r.Password = ""
	r.Token = token
	return r, nil
}

// encrypt password
func (a *Auth) Encrypt() (string, error) {
	if a.Password == "" {
		return "", nil
	}
	if len(a.Password) > 72 {
		return "", errors.New("password is too long")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}

// compare password
func (a *Auth) Compare(hashed string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(a.Password))
}
