package core

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/osamikoyo/music-and-marks/services/user/config"
	"github.com/osamikoyo/music-and-marks/services/user/entity"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmptyFields      = errors.New("empty fields")
	ErrJwtFailed        = errors.New("failed to generate token")
	ErrDoublePasswords  = errors.New("passwords is simular")
	ErrComparePasswords = errors.New("old password is wrong")
	ErrParseToken       = errors.New("failed parse jwt token")
	ErrInvalidToken     = errors.New("invalid jwt token")
	ErrGetNewToken      = errors.New("fialed to get new token")
	ErrInternal         = errors.New("internal error")
)

type Repository interface {
	CreateUser(ctx context.Context, user *entity.User) error
	UpdateUser(ctx context.Context, update *entity.User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	GetUser(ctx context.Context, id uuid.UUID) (*entity.User, error)
	CheckUser(ctx context.Context, password, email string) (string, error)
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
}

type UserCore struct {
	repo    Repository
	timeout time.Duration
	cfg     *config.Config
}

type RefreshTokenClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func (uc *UserCore) context() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), uc.timeout)
}
func newJwtRefreshKey(uid, key string, dur time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"uid": uid,
		"exp": time.Now().Add(dur).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(key))
}

func newJwtAccessKey(uid,ref, key string, dur time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"uid": uid,
		"ref": ref,
		"exp": time.Now().Add(dur).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(key))
}

func NewUserCore(repo Repository,timeout time.Duration, cfg *config.Config) *UserCore {
	return &UserCore{
		repo:    repo,
		timeout: timeout,
		cfg:     cfg,
	}
}

func (uc *UserCore) RegisterUser(username, password, email string) (*entity.TokenPair, error) {
	if len(username) == 0 || len(password) == 0 {
		return nil, ErrEmptyFields
	}

	user := entity.NewUser(username, password, email)

	ctx, cancel := uc.context()
	defer cancel()

	if err := uc.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	reftoken, err := newJwtRefreshKey(user.ID.String(), uc.cfg.JwtKey, uc.cfg.RTokenTTL)
	if err != nil {
		return nil, ErrJwtFailed
	}

	access, err := newJwtAccessKey(user.ID.String(), reftoken, uc.cfg.JwtKey, uc.cfg.ATokenTTL)
	if err != nil{
		return nil, ErrJwtFailed
	}

	return &entity.TokenPair{
		RefreshToken: reftoken,
		AccessToken:  access,
	}, nil
}

func (uc *UserCore) LoginUser(password, email string) (*entity.TokenPair, error) {
	if len(password) == 0 || len(email) == 0 {
		return nil, ErrEmptyFields
	}

	ctx, cancel := context.WithTimeout(context.Background(), uc.timeout)
	defer cancel()

	id, err := uc.repo.CheckUser(ctx, password, email)
	if err != nil {
		return nil, err
	}

	reftoken, err := newJwtRefreshKey(id, uc.cfg.JwtKey, uc.cfg.RTokenTTL)
	if err != nil {
		return nil, err
	}

	access, err := newJwtAccessKey(id, reftoken, uc.cfg.JwtKey, )

	return &entity.TokenPair{
		RefreshToken: reftoken,
		AccessToken:  access,
	}, nil
}

func (uc *UserCore) ChangePassword(id uuid.UUID, old, new string) error {
	if old == new {
		return ErrDoublePasswords
	}

	ctx, cancel := uc.context()
	defer cancel()

	user, err := uc.repo.GetUser(ctx, id)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(old))
	if err != nil {
		return ErrComparePasswords
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(new), bcrypt.DefaultCost)
	if err != nil {
		return ErrInternal
	}

	user.Password = string(hash)

	if err = uc.repo.UpdateUser(ctx, user); err != nil {
		return err
	}

	return nil
}

func (uc *UserCore) Refresh(refreshToken string) (string, error) {
	if len(refreshToken) == 0 {
		return "", ErrEmptyFields
	}

	token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (any, error) {
		return uc.cfg.JwtKey, nil
	})
	if err != nil {
		return "", ErrParseToken
	}

	refreshClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", ErrInvalidToken
	}

	expiresAt := time.Now().Add(uc.cfg.ATokenTTL)

	claims := Claims{
		UserID: refreshClaims["uid"].(string),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "music-and-marks",
		},
	}

	accessTkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessTknStr, err := accessTkn.SignedString(uc.cfg.JwtKey)
	if err != nil {
		return "", ErrGetNewToken
	}

	return accessTknStr, nil
}

func (uc *UserCore) IncLike(uid uuid.UUID) error {
	if len(uid) == 0 {
		return ErrEmptyFields
	}

	ctx, cancel := uc.context()
	defer cancel()

	user, err := uc.repo.GetUser(ctx, uid)
	if err != nil {
		return err
	}

	user.Likes++

	return uc.repo.UpdateUser(ctx, user)
}

func (uc *UserCore) DecLike(uid uuid.UUID) error {
	if len(uid) == 0 {
		return ErrEmptyFields
	}

	ctx, cancel := uc.context()
	defer cancel()

	user, err := uc.repo.GetUser(ctx, uid)
	if err != nil {
		return err
	}

	user.Likes--

	return uc.repo.UpdateUser(ctx, user)
}

func (uc *UserCore) IncReview(uid uuid.UUID) error {
	if len(uid) == 0 {
		return ErrEmptyFields
	}

	ctx, cancel := uc.context()
	defer cancel()

	user, err := uc.repo.GetUser(ctx, uid)
	if err != nil {
		return err
	}

	user.Likes--

	return uc.repo.UpdateUser(ctx, user)
}

func (uc *UserCore) DecReview(uid uuid.UUID) error {
	if len(uid) == 0 {
		return ErrEmptyFields
	}

	ctx, cancel := uc.context()
	defer cancel()

	user, err := uc.repo.GetUser(ctx, uid)
	if err != nil {
		return err
	}

	user.Reviews++

	return uc.repo.UpdateUser(ctx, user)
}

func (uc *UserCore) GetUserByID(uid uuid.UUID) (*entity.User, error) {
	if len(uid) == 0 {
		return nil, ErrEmptyFields
	}

	ctx, cancel := uc.context()
	defer cancel()

	return uc.repo.GetUser(ctx, uid)
}

func (uc *UserCore) GetUserByUsername(username string) (*entity.User, error) {
	if len(username) == 0 {
		return nil, ErrEmptyFields
	}

	ctx, cancel := uc.context()
	defer cancel()

	return uc.repo.GetUserByUsername(ctx, username)
}

func (uc *UserCore) DeleteUser(uid uuid.UUID) error {
	if len(uid) == 0 {
		return ErrEmptyFields
	}

	ctx, cancel := uc.context()
	defer cancel()

	return uc.repo.DeleteUser(ctx, uid)
}
