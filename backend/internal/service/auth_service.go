package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	apperrors "github.com/your-org/project-budget-tracker/backend/internal/errors"
	"github.com/your-org/project-budget-tracker/backend/internal/models"
)

type AuthService struct {
	db        *gorm.DB
	jwtSecret string
}

func NewAuthService(db *gorm.DB, jwtSecret string) *AuthService {
	return &AuthService{
		db:        db,
		jwtSecret: jwtSecret,
	}
}

// Claims represents JWT claims
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// Register creates a new user
func (s *AuthService) Register(email, password, name string) (*models.User, string, error) {
	// Check if user already exists
	var existingUser models.User
	if err := s.db.Where("email = ?", email).First(&existingUser).Error; err == nil {
		return nil, "", apperrors.ErrAlreadyExists("User")
	}

	// Hash password
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, "", apperrors.ErrInternal(err)
	}

	// Create user
	user := &models.User{
		Email:        email,
		PasswordHash: hashedPassword,
		Name:         name,
		Role:         "member", // default role
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, "", apperrors.ErrDatabaseError(err)
	}

	// Generate token
	token, err := s.generateToken(user)
	if err != nil {
		return nil, "", apperrors.ErrInternal(err)
	}

	return user, token, nil
}

// Login authenticates a user
func (s *AuthService) Login(email, password string) (*models.User, string, error) {
	var user models.User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", apperrors.ErrInvalidCredentials()
		}
		return nil, "", apperrors.ErrDatabaseError(err)
	}

	// Check password
	if !checkPassword(password, user.PasswordHash) {
		return nil, "", apperrors.ErrInvalidCredentials()
	}

	// Generate token
	token, err := s.generateToken(&user)
	if err != nil {
		return nil, "", apperrors.ErrInternal(err)
	}

	return &user, token, nil
}

// GetUserByID retrieves a user by ID
func (s *AuthService) GetUserByID(userID string) (*models.User, error) {
	id, err := uuid.Parse(userID)
	if err != nil {
		return nil, apperrors.ErrInvalidInput(err)
	}

	var user models.User
	if err := s.db.First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound("User")
		}
		return nil, apperrors.ErrDatabaseError(err)
	}

	return &user, nil
}

// ValidateToken validates a JWT token and returns claims
func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, apperrors.ErrInvalidToken()
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, apperrors.ErrInvalidToken()
	}

	return claims, nil
}

// generateToken generates a JWT token for a user
func (s *AuthService) generateToken(user *models.User) (string, error) {
	claims := &Claims{
		UserID: user.ID.String(),
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

// hashPassword hashes a password using bcrypt
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

// checkPassword checks if password matches hash
func checkPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
