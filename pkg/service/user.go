package service

import (
	"megichains/pkg/biz"
	"megichains/pkg/converter"
	"megichains/pkg/entity"
	"megichains/pkg/global"
	"strings"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) Create(username, password string) (success bool, err error) {
	hash, err := global.HashPassword(password)
	if err != nil {
		logx.Errorf("user service generate hash password failed, username:%v, password:%v, err:%v", username, password, hash)
		err = biz.GenerateHashPasswordFailed
		return
	}

	user := &entity.User{Id: strings.ReplaceAll(uuid.NewString(), "-", ""), Username: username, Password: hash}
	err = s.db.Model(&entity.User{}).Create(user).Error
	if err != nil {
		logx.Errorf("user service create user failed, username:%v, password:%v, err:%v", username, password, err)
		err = biz.UserCreateFailed
		return
	}

	success = true

	return
}

func (s *UserService) Get(username string) (userinfo *converter.UserInfo, err error) {
	user := &entity.User{}
	err = s.db.Model(&entity.User{}).Where("username = ?", username).First(user).Error
	if err != nil {
		logx.Errorf("user service get user info failed, username:%v,  err:%v", username, err)
		err = biz.UserCreateFailed
		return
	}

	userinfo = &converter.UserInfo{}
	copier.Copy(userinfo, user)

	return
}
