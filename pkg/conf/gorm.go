package conf

import "gorm.io/gorm"

func NewGormConfig() *gorm.Config {
	return &gorm.Config{}
}
