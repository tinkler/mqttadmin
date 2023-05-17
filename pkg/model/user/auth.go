package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/tinkler/mqttadmin/errz"
	"github.com/tinkler/mqttadmin/pkg/db"
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
	return "v1.user"
}

func (a *Auth) Signin(ctx context.Context) (*Auth, error) {
	u := Auth{}
	se := db.GetDB(ctx).Where("username = ?", a.Username).First(&u)
	if se.Error != nil {
		if errors.Is(se.Error, gorm.ErrRecordNotFound) {
			return nil, status.Ok("user not found")
		}
		return nil, se.Error
	}

	if err := a.Compare(u.Password); err != nil {
		return nil, status.Ok("invalid password")
	}

	if a.DeviceToken == "" {
		a.DeviceToken = uuid.New().String()
	}

	// when login success, clear the token
	ClearShortTokenKV(a.ID)

	token, err := getJwtToken(a.ID, a.DeviceToken)
	if err != nil {
		return nil, err
	}
	a.Password = ""
	a.Token = token

	return a, nil
}

// QuickSignin quick signin with password
func (a *Auth) QuickSignin(ctx context.Context) error {
	if a.DeviceToken == "" {
		return errz.ErrVdM("device_token", "device token is empty", "登录设备号不能为空")
	}
	if a.Token == "" {
		return errz.ErrVdM("token", "token is empty", "登录令牌不能为空")
	}
	_, err := CheckJwtToken(context.TODO(), a.Token)
	return err
}

func (a *Auth) Signup(ctx context.Context) (*Auth, error) {
	if a.Username == "" {
		return nil, status.Ok("username is empty")
	}
	if a.Password == "" {
		return nil, status.Ok("password is empty")
	}

	se := db.GetDB(ctx).Where("username = ?", a.Username).First(&Auth{})
	if se.Error == nil {
		return nil, status.Ok(ErrMsgNameAreadyExist)
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
	ClearShortTokenKV(r.ID)

	token, err := getJwtToken(r.ID, r.DeviceToken)
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
