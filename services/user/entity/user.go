package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/osamikoyo/music-and-marks/services/user/api/proto/gen/pb"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type (
	User struct {
		ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
		CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
		UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
		DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

		Username string `gorm:"size:50;uniqueIndex;not null" json:"username" validate:"required,alphanum,min=3,max=50"`
		Password string `gorm:"size:255;not null" json:"-" validate:"required,min=8"`
		Email    string `gorm:"size:255;uniqueIndex" json:"email,omitempty" validate:"omitempty,email"`
		Likes    int    `json:"likes"`
		Reviews  int    `json:"reciews"`
	}
)

func NewUser(username, password, email string) *User {
	return &User{
		Username: username,
		Password: password,
		Email:    email,
	}
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}

	if u.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}

	return nil
}

func (u *User) ToProto() *pb.User {
	return &pb.User{
		Email: u.Email,
		Username: u.Username,
		CreatedAt: timestamppb.New(u.CreatedAt),
		UpdatedAt: timestamppb.New(u.UpdatedAt),
		Likes: int64(u.Likes),
		Reviews: int64(u.Reviews),
		Id: u.ID.String(),
	}
}