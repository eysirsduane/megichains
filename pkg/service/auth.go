package service

import (
	"megichains/pkg/biz"
	"megichains/pkg/entity"
	"megichains/pkg/global"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type AuthService struct {
	db      *gorm.DB
	asecret string
	aexpire int64
	issuer  string
	rsecret string
	rexpire int64
}

func NewAuthService(db *gorm.DB, asecret string, aexpire int64, rsecret string, rexpire int64, issuer string) *AuthService {
	return &AuthService{db: db, asecret: asecret, aexpire: aexpire, rsecret: rsecret, rexpire: rexpire, issuer: issuer}
}

func (s *AuthService) Login(username, password string) (atoken, rtoken string, err error) {
	user := &entity.User{}
	err = s.db.Model(&entity.User{}).Where("username = ?", username).First(user).Error
	if err != nil {
		logx.Errorf("auth service login canot find username failed, un:%v, err:%v", username, err)
		err = biz.UserUsernameOrPasswordIncorrect
		return
	}

	if !global.CheckPassword(user.Password, password) {
		logx.Errorf("auth service login failed, err:%v", err)
		err = biz.UserUsernameOrPasswordIncorrect
		return
	}

	atoken, err = global.GenerateToken(user.Id, user.Username, s.asecret, s.aexpire, s.issuer)
	if err != nil {
		logx.Errorf("auth service login generate token failed, err:%v", err)
		err = biz.UserLoginGenerateTokenFailed
		return
	}

	rtoken, err = global.GenerateRefreshToken(user.Id, user.Username, s.rsecret, s.rexpire, s.issuer)
	if err != nil {
		logx.Errorf("auth service login generate token failed, err:%v", err)
		err = biz.UserLoginGenerateTokenFailed
		return
	}

	return
}

func (s *AuthService) RefreshToken(otoken string, secrect string) (atoken, rtoken string, err error) {
	claims, err := global.ParseToken(otoken, secrect)
	if err != nil {
		logx.Errorf("auth service refresh token parse otoken failed, err:%v", err)
		return
	}

	// 根据 refresh token 生成新的 token
	atoken, err = global.GenerateToken(claims.UserID, claims.Username, s.asecret, s.aexpire, s.issuer)
	if err != nil {
		logx.Errorf("auth service refresh token generate new token failed, err:%v", err)
		return
	}

	rtoken, err = global.GenerateRefreshToken(claims.UserID, claims.Username, s.rsecret, s.rexpire, s.issuer)
	if err != nil {
		logx.Errorf("auth service refresh token generate new refresh token failed, err:%v", err)
		return
	}

	return
}
