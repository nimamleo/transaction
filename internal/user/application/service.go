package application

import (
	"context"

	"transaction/internal/user/domain"
	"transaction/pkg/hash"
)

type Service struct {
	userRepository domain.Repository
	apiKeyRepo     domain.APIKeyRepository
}

func NewService(userRepository domain.Repository, apiKeyRepo domain.APIKeyRepository) *Service {
	return &Service{
		userRepository: userRepository,
		apiKeyRepo:     apiKeyRepo,
	}
}

type CreateUserResult struct {
	User   *domain.User
	APIKey string
}

func (s *Service) CreateUser(ctx context.Context, name, email string) (*CreateUserResult, error) {
	exists, err := s.userRepository.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, domain.ErrEmailAlreadyExists
	}

	user := domain.New(name, email)

	if err := s.userRepository.Create(ctx, user); err != nil {
		return nil, err
	}

	apiKey, err := domain.NewAPIKey(user.ID)
	if err != nil {
		return nil, err
	}

	if err := s.apiKeyRepo.Create(ctx, apiKey); err != nil {
		return nil, err
	}

	return &CreateUserResult{
		User:   user,
		APIKey: apiKey.PlainAPIKey,
	}, nil
}

func (s *Service) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	return s.userRepository.GetByID(ctx, id)
}

func (s *Service) GetUserIDByAPIKey(ctx context.Context, apiKey string) (string, error) {
	hashedKey := hash.Hash(apiKey)
	return s.apiKeyRepo.GetUserIDByAPIKey(ctx, hashedKey)
}
