package services

import (
	"errors"
	"time"

	"github.com/clarkzhu2020/aidecms/config"
	"github.com/clarkzhu2020/aidecms/internal/app/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserService 处理用户相关的业务逻辑
type UserService struct {
	DB *gorm.DB
}

// NewUserService 创建用户服务实例
func NewUserService() *UserService {
	return &UserService{
		DB: config.DB,
	}
}

// Register 注册新用户
func (s *UserService) Register(username, email, password, firstName, lastName string) (*models.User, error) {
	// 检查用户名是否已存在
	var existingUser models.User
	if result := s.DB.Where("username = ?", username).First(&existingUser); result.Error == nil {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	if result := s.DB.Where("email = ?", email).First(&existingUser); result.Error == nil {
		return nil, errors.New("邮箱已存在")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 创建新用户
	user := &models.User{
		Username:  username,
		Email:     email,
		Password:  string(hashedPassword),
		FirstName: firstName,
		LastName:  lastName,
	}

	if err := s.DB.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// Login 用户登录
func (s *UserService) Login(usernameOrEmail, password string) (*models.User, string, error) {
	var user models.User

	// 查找用户（通过用户名或邮箱）
	result := s.DB.Where("username = ? OR email = ?", usernameOrEmail, usernameOrEmail).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, "", errors.New("用户不存在")
		}
		return nil, "", result.Error
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", errors.New("密码错误")
	}

	// 更新最后登录时间
	now := time.Now()
	user.LastLogin = &now
	s.DB.Save(&user)

	// 生成JWT令牌
	token, err := s.GenerateJWT(user.ID)
	if err != nil {
		return nil, "", err
	}

	return &user, token, nil
}

// GetUserByID 通过ID获取用户
func (s *UserService) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	if err := s.DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}
	return &user, nil
}

// UpdateProfile 更新用户资料
func (s *UserService) UpdateProfile(id uint, firstName, lastName, avatar string) (*models.User, error) {
	var user models.User
	if err := s.DB.First(&user, id).Error; err != nil {
		return nil, err
	}

	// 更新字段
	user.FirstName = firstName
	user.LastName = lastName
	if avatar != "" {
		user.Avatar = avatar
	}

	if err := s.DB.Save(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// GenerateJWT 生成JWT令牌
func (s *UserService) GenerateJWT(userID uint) (string, error) {
	// 创建JWT声明
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * time.Duration(config.JWT.ExpiresIn)).Unix(),
		"iat":     time.Now().Unix(),
	}

	// 创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名令牌
	tokenString, err := token.SignedString([]byte(config.JWT.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseJWT 解析JWT令牌
func (s *UserService) ParseJWT(tokenString string) (uint, error) {
	// 解析令牌
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWT.SecretKey), nil
	})

	if err != nil {
		return 0, err
	}

	if !token.Valid {
		return 0, errors.New("无效的令牌")
	}

	// 获取声明
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("无效的令牌声明")
	}

	// 获取用户ID
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("无效的用户ID")
	}

	return uint(userID), nil
}
