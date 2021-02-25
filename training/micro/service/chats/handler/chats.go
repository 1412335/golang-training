package handler

import (
	"context"
	"encoding/json"
	"sort"
	"strings"
	"time"

	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"

	chats "chats/proto"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

var (
	ErrMissingUserIds = errors.BadRequest("MISSING_ID", "missing userIds")
	ErrDatabase       = errors.InternalServerError("DATABASE_ERROR", "Connecting database failed")
)

// model
type Chat struct {
	ID        string
	UserIds   string `gorm:"uniqueIndex"`
	CreatedAt time.Time
}

func (c *Chat) Serialize() (*chats.Chat, error) {
	var userIds []string
	if err := json.Unmarshal([]byte(c.UserIds), &userIds); err != nil {
		return nil, err
	}
	return &chats.Chat{
		Id:        c.ID,
		UserIds:   userIds,
		CreatedAt: timestamppb.New(c.CreatedAt),
	}, nil
}

// handler
type Chats struct {
	DB *gorm.DB
}

func (h *Chats) CreateChat(ctx context.Context, req *chats.CreateChatRequest, rsp *chats.CreateChatResponse) error {
	logger.Info("chats.CreateChat request")
	if len(req.UserIds) < 2 {
		return ErrMissingUserIds
	}

	// sort user IDs & encode
	sort.Strings(req.UserIds)
	bytes, err := json.Marshal(req.UserIds)
	if err != nil {
		logger.Errorf("encode userIds failed: %v", err)
		return errors.InternalServerError("ENCODE_ERROR", "encode error")
	}

	// construct chat model
	chat := &Chat{
		ID:        uuid.New().String(),
		UserIds:   string(bytes),
		CreatedAt: time.Now(),
	}

	// write to db
	if err := h.DB.Create(chat).Error; err != nil && strings.Contains(err.Error(), "idx_chats_user_ids") {
		// get existing chat
		var existingChat Chat
		if err := h.DB.Where(&Chat{UserIds: string(bytes)}).First(&existingChat).Error; err != nil {
			logger.Errorf("get existing chat failed: %v", err)
			return ErrDatabase
		}
		chat = &existingChat
	} else if err != nil {
		logger.Errorf("create chat failed: %v", err)
		return ErrDatabase
	}

	// response
	if rsp.Chat, err = chat.Serialize(); err != nil {
		logger.Errorf("decode userIds failed: %v", err)
		return errors.InternalServerError("DECODE_ERROR", "decode error")
	}

	return nil
}
