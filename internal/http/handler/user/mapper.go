package user

import "transaction/internal/user/domain"

func ToResponse(user *domain.User) Response {
	return Response{
		ID:     user.ID,
		Name:   user.Name,
		Email:  user.Email,
		APIKey: nil,
	}
}

func ToCreateResponse(user *domain.User, apiKey string) Response {
	return Response{
		ID:     user.ID,
		Name:   user.Name,
		Email:  user.Email,
		APIKey: &apiKey,
	}
}
