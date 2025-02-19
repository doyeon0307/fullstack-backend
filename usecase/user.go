package usecase

import (
	"context"
	"time"

	"github.com/doyeon0307/tickit-backend/common"
	"github.com/doyeon0307/tickit-backend/domain"
	"github.com/doyeon0307/tickit-backend/dto"
	"github.com/doyeon0307/tickit-backend/models"
	"github.com/doyeon0307/tickit-backend/service"
)

type userUsecase struct {
	userRepo domain.UserRepository
}

func NewUserUsecase(repo domain.UserRepository) domain.UserUsecase {
	return &userUsecase{
		userRepo: repo,
	}
}

func (u userUsecase) GetProfile(id string) (*dto.KakaoProfile, error) {
	model, err := u.userRepo.GetById(context.Background(), id)
	if err != nil {
		return nil, err
	}

	profile := &dto.KakaoProfile{
		NickName: model.Name,
	}
	return profile, nil
}

func (u userUsecase) CreateUser(idToken string, accessToken string) (string, error) {
	oauthId, err := service.GetOAuthIdFromKakao(idToken)
	if err != nil {
		return "", &common.AppError{
			Code:    common.ErrNotFound,
			Message: "ID Token: 카카오로부터 사용자 정보를 불러오는데 실패했습니다",
			Err:     err,
		}
	}

	info, err := service.GetUserInfoFromKakao(accessToken)
	if err != nil {
		return "", &common.AppError{
			Code:    common.ErrNotFound,
			Message: "Access Token: 카카오로부터 사용자 정보를 불러오는데 실패했습니다",
			Err:     err,
		}
	}
	name := info.NickName

	user := &models.User{
		OAuthId: oauthId,
		Name:    name,
	}

	id, err := u.userRepo.Create(context.Background(), user)
	return id, err
}

func (u userUsecase) DeleteUser(id string) error {
	return u.userRepo.Delete(context.Background(), id)
}

func (u *userUsecase) GetUserByOAuthId(oauthId string) (*models.User, error) {
	user, err := u.userRepo.GetByOAuthId(context.Background(), oauthId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *userUsecase) SaveRefreshToken(userId string, refreshToken string, expiryTime time.Time) error {
	err := u.userRepo.SaveRefreshToken(context.Background(), userId, refreshToken, expiryTime)
	return err
}

func (u *userUsecase) ValidateStoredRefreshToken(userId string, refreshToken string) (bool, error) {
	storedToken, err := u.userRepo.GetRefreshToken(context.Background(), userId)
	if err != nil {
		return false, &common.AppError{
			Code:    common.ErrServer,
			Message: "Refresh Token 조회에 실패했습니다",
			Err:     err,
		}
	}

	return storedToken == refreshToken, nil
}

func (u *userUsecase) WithdrawUser(userId string) error {
	return u.userRepo.DeleteUser(context.Background(), userId)
}

func (u *userUsecase) Logout(userId string) error {
	return u.userRepo.RemoveRefreshToken(context.Background(), userId)
}
