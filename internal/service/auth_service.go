package service

import (
	"blog/internal/model/dto/request"
	dto "blog/internal/model/dto/response"
	"blog/internal/model/entity"
	"blog/internal/repository"
	"blog/pkg/captcha"
	"blog/pkg/email"
	bizerrors "blog/pkg/errors"
	"blog/pkg/jwt"
	"blog/pkg/logger"
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

// AuthService 认证服务接口
type AuthService interface {
	// Login 用户登录
	Login(req *request.LoginRequest) (*dto.LoginResponse, error)
	// RefreshToken 刷新 Token
	RefreshToken(token string) (string, error)
	// SendEmailCode 发送邮箱验证
	SendEmailCode(email string) error
	// Register 用户注册
	Register(r *request.RegisterRequest) error
	// Logout 用户登出
	Logout(id uint) error

	// ResetPassword 重置密码（邮箱验证码）
	ResetPassword(req *request.ResetPasswordRequest) error
}

// authService 认证服务实现
type authService struct {
	userRepo    repository.UserRepository
	userService UserService
	redis       *redis.Client
}

func (s *authService) Logout(id uint) error {
	// 当前项目未实现 token 黑名单/撤销机制，因此这里默认无操作成功返回。
	_ = id
	return nil
}

func (s *authService) ResetPassword(req *request.ResetPasswordRequest) error {
	ctx := context.Background()
	codeKey := fmt.Sprintf("email_code:%s", req.Email)
	cacheCode, err := s.redis.Get(ctx, codeKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return bizerrors.New(bizerrors.CodeInvalidParam, "验证码已过期或未发送，请重新获取")
		}
		return bizerrors.NewWithErr(bizerrors.CodeInternalError, "获取验证码缓存失败", err)
	}

	if req.EmailCaptcha != cacheCode {
		return bizerrors.New(bizerrors.CodeInvalidParam, "验证码输入错误，请重新核对")
	}

	// 校验通过后，删除验证码，避免重复使用
	_ = s.redis.Del(ctx, codeKey).Err()

	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return err
	}
	if user == nil {
		return bizerrors.ErrUserNotFound
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.userRepo.Update(user.ID, map[string]interface{}{
		"password": string(hashedPassword),
	})
}

// NewAuthService 创建认证服务
func NewAuthService(
	userRepo repository.UserRepository,
	userService UserService,
	redisClient *redis.Client,
) AuthService {
	return &authService{
		userRepo:    userRepo,
		userService: userService,
		redis:       redisClient,
	}
}

func (s *authService) Login(req *request.LoginRequest) (*dto.LoginResponse, error) {
	// 1. 校验验证码
	if !captcha.Verify(req.CaptchaID, req.Captcha) {
		return nil, bizerrors.ErrInvalidCaptcha
	}
	// 2. 校验用户
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		logger.Errorf("查找用户名失败: username=%s, err=%v", req.Username, err)
		return nil, bizerrors.NewWithErr(bizerrors.CodeInternalError, "查找用户名失败:", err)
	}
	if user == nil {
		return nil, bizerrors.ErrInvalidCredentials
	}

	// 3.检查用户状态 (1正常 2禁用)
	if user.Status != 1 {
		return nil, bizerrors.ErrUserDisabled
	}

	// 4.验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		logger.Errorf("密码验证失败: username=%s, err=%v", req.Username, err)
		return nil, bizerrors.NewWithErr(bizerrors.CodeInvalidCredentials, "登录失败：密码错误", err)
	}

	// 5.生成 JWT
	token, err := jwt.GenerateToken(user.ID, user.Username)
	if err != nil {
		logger.Error("JWT令牌生成失败:", zap.Error(err))
		return nil, bizerrors.NewWithErr(bizerrors.CodeInvalidToken, "无效的令牌", err)
	}

	// 构建响应
	userResp := s.userService.GetUserResponse(user)
	return &dto.LoginResponse{
		Token: token,
		User:  *userResp,
	}, nil
}

func (s *authService) Register(req *request.RegisterRequest) error {
	// 检查用户名是否存在
	exists, err := s.userRepo.ExistsByUsername(req.Username)
	if err != nil {
		return err
	}
	if exists {
		return bizerrors.ErrUserAlreadyExists
	}

	// 检查邮箱是否存在
	exists, err = s.userRepo.ExistsByEmail(req.Email)
	if err != nil {
		return err
	}
	if exists {
		return bizerrors.New(bizerrors.CodeUserAlreadyExists, "邮箱已被注册")
	}

	// 验证邮箱验证码
	ctx := context.Background()
	codeKey := fmt.Sprintf("email_code:%s", req.Email)
	cacheCode, err := s.redis.Get(ctx, codeKey).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return bizerrors.New(bizerrors.CodeInvalidParam, "验证码已过期或未发送，请重新获取")
		}
		return bizerrors.NewWithErr(bizerrors.CodeInternalError, "获取验证码缓存失败", err)
	}

	if req.EmailCaptcha != cacheCode {
		return bizerrors.New(bizerrors.CodeInvalidParam, "验证码输入错误，请重新核对")
	}
	if delErr := s.redis.Del(ctx, codeKey); delErr != nil {
		logger.Info("删除验证码失败")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 创建用户
	user := &entity.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
		Status:   1, // 默认状态正常
		Role:     1, // 默认普通用户
	}
	return s.userRepo.Create(user)
}

func (s *authService) RefreshToken(token string) (string, error) {
	newToken, err := jwt.RefreshToken(token)
	if err != nil {
		return "", err
	}
	return newToken, nil
}
func (s *authService) SendEmailCode(emailStr string) error {
	// 1. 生成6位随机数字
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := fmt.Sprintf("%06d", rnd.Intn(1000000))

	// 2. 存储到 Redis (有效期5分钟)
	ctx := context.Background()
	key := fmt.Sprintf("email_code:%s", emailStr)
	err := s.redis.Set(ctx, key, code, 5*time.Minute)
	if err != nil {
		return bizerrors.NewWithErr(bizerrors.CodeInternalError, "缓存验证码失败", err.Err())
	}

	// 3. 发送邮件
	subject := "Ryne 博客验证码"
	body := fmt.Sprintf("<h1>您的验证码是: %s</h1><p>有效期5分钟，请勿泄露给他人。</p>", code)
	if err := email.SendEmail(emailStr, subject, body); err != nil {
		return bizerrors.NewWithErr(bizerrors.CodeInternalError, "发送邮件失败", err)
	}

	return nil
}
