package service

import (
	"errors"

	"eventix/internal/entity"
	"eventix/internal/repository"
	"eventix/pkg/utils"

	"gorm.io/gorm"
)

// Common errors for authentication operations
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailAlreadyExists = errors.New("email already registered")
)

// AuthService defines the interface for authentication operations.
// It abstracts the business logic for user registration and login.
type AuthService interface {
	// Register creates a new user account with hashed password
	Register(input *entity.RegisterInput) (*entity.User, error)
	// Login authenticates a user and returns a JWT token
	Login(input *entity.LoginInput) (string, error)
}

// authService is the implementation of AuthService.
type authService struct {
	userRepo repository.UserRepository
}

// NewAuthService creates a new instance of AuthService.
func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

// Register creates a new user account.
// Flow:
// 1. Check if email already exists
// 2. Hash the password using bcrypt
// 3. Create the user entity
// 4. Save to database via repository
func (s *authService) Register(input *entity.RegisterInput) (*entity.User, error) {
	// Step 1: Check if user with this email already exists
	existingUser, err := s.userRepo.FindByEmail(input.Email)
	if err == nil && existingUser != nil {
		return nil, ErrEmailAlreadyExists
	}
	// Only return error if it's not a "not found" error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Step 2: Hash the password for secure storage
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	// Step 3: Create user entity with hashed password
	user := &entity.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPassword,
		Role:     "user", // Default role for new users
	}

	// Step 4: Save user to database
	if err := s.userRepo.Save(user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login authenticates a user and returns a JWT token.
// Flow:
// 1. Find user by email
// 2. Verify password using bcrypt
// 3. Generate JWT token with user claims
// 4. Return the token
func (s *authService) Login(input *entity.LoginInput) (string, error) {
	// Step 1: Find user by email
	user, err := s.userRepo.FindByEmail(input.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	// Step 2: Verify password matches the stored hash
	if err := utils.CheckPassword(user.Password, input.Password); err != nil {
		return "", ErrInvalidCredentials
	}

	// Step 3: Generate JWT token with user ID and role
	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return "", err
	}

	// Step 4: Return the generated token
	return token, nil
}
