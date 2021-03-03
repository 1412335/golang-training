package handler

import (
	"context"
	"encoding/json"
	"fw/pkg/audit"
	pb "fw/users/proto"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        string
	FirstName string
	LastName  string
	ValidFrom time.Time
	ValidTo   time.Time
	Active    bool
	Password  string
	Email     string `gorm:"uniqueIndex"`
	CreatedAt time.Time
	UpdatedAt time.Time
	audit     *audit.Audit `gorm:"-" json:"-"`
}

func (u *User) sanitize() *pb.User {
	return &pb.User{
		Id:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		ValidFrom: timestamppb.New(u.ValidFrom),
		ValidTo:   timestamppb.New(u.ValidTo),
		Active:    u.Active,
		Email:     u.Email,
		CreatedAt: timestamppb.New(u.CreatedAt),
		UpdatedAt: timestamppb.New(u.UpdatedAt),
	}
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	logger.Infof("before create")
	return nil
}

func (u *User) AfterCreate(tx *gorm.DB) error {
	if u.audit == nil {
		return nil
	}
	ctx := tx.Statement.Context
	err := u.sendUserAudit(ctx, "Users", "AfterCreate", "insert", "user", u.ID)
	if err != nil {
		logger.Errorf("Call AfterCreate failed: %v", err)
	}
	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
	logger.Infof("before update")
	return nil
}

// Updating data in same transaction
func (u *User) AfterUpdate(tx *gorm.DB) error {
	if u.audit == nil {
		return nil
	}
	ctx := tx.Statement.Context
	err := u.sendUserAudit(ctx, "Users", "AfterUpdate", "insert", "user", u.ID)
	if err != nil {
		logger.Errorf("Call AfterUpdate failed: %v", err)
	}
	return nil
}

func (u *User) sendUserAudit(ctx context.Context, serviceName, actionFunc, actionType string, objectName string, iObjectId string) error {
	bytes, err := json.Marshal(u)
	if err != nil {
		return errors.InternalServerError("ENCODE_ERROR", "encode user error")
	}
	return u.audit.Send(ctx, serviceName, actionFunc, actionType, objectName, iObjectId, bytes)
}
