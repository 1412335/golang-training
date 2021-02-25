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
	ErrMissingUserIds = errors.BadRequest("MISSING_USER_ID", "missing userIds")
	ErrMissingChatId  = errors.BadRequest("MISSING_CHAT_ID", "missing chat id")
	ErrMissingAuthor  = errors.BadRequest("MISSING_AUTHOR", "missing author")
	ErrChatNotFound   = errors.NotFound("CHAT_NOT_FOUND", "chat not found")
	ErrAuthorNotFound = errors.NotFound("AUTHOR_NOT_FOUND", "author not found in chat users")
	ErrDatabase       = errors.InternalServerError("DATABASE_ERROR", "Connecting database failed")
)

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// model
type Chat struct {
	ID        string
	UserIds   string `gorm:"uniqueIndex"`
	CreatedAt time.Time
	Messages  []Message
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

// model
type Message struct {
	ID     string
	Author string
	ChatID string
	Text   string
	SendAt time.Time
}

func (c *Message) Serialize() (*chats.Message, error) {
	return &chats.Message{
		Id:     c.ID,
		Author: c.Author,
		ChatId: c.ChatID,
		Text:   c.Text,
		SendAt: timestamppb.New(c.SendAt),
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
		CreatedAt: time.Now().Round(time.Microsecond),
		// https://stackoverflow.com/questions/60433870/saving-time-time-in-golang-to-postgres-timestamp-with-time-zone-field
		// postgres: auto round to microseconds (>5) eg: 100100500 => 100100
		// go: Round(time.Microsecond) (>=5) eg: 100100500 => 100101
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

func (h *Chats) CreateMessage(ctx context.Context, req *chats.CreateMessageRequest, rsp *chats.CreateMessageResponse) error {
	logger.Info("chats.CreateMessage request")
	if len(req.ChatId) == 0 {
		return ErrMissingChatId
	}
	if len(req.Author) == 0 {
		return ErrMissingAuthor
	}

	// lookup chat
	var chat Chat
	if err := h.DB.Where(&Chat{ID: req.ChatId}).First(&chat).Error; err == gorm.ErrRecordNotFound {
		return ErrChatNotFound
	} else if err != nil {
		logger.Errorf("lookup chat failed: %v", err)
		return ErrDatabase
	}

	// check author in chat
	var userIds []string
	if err := json.Unmarshal([]byte(chat.UserIds), &userIds); err != nil {
		logger.Errorf("decode userIds failed: %v", err)
		return errors.InternalServerError("DECODE_ERROR", "decode error")
	}
	if !contains(userIds, req.Author) {
		return ErrAuthorNotFound
	}

	// construct chat model
	msg := &Message{
		ID:     uuid.New().String(),
		ChatID: req.ChatId,
		Author: req.Author,
		Text:   req.Text,
		SendAt: time.Now().Round(time.Microsecond),
		// https://stackoverflow.com/questions/60433870/saving-time-time-in-golang-to-postgres-timestamp-with-time-zone-field
	}

	// write to db
	if err := h.DB.Create(msg).Error; err != nil {
		logger.Errorf("create message failed: %v", err)
		return ErrDatabase
	}

	// response
	var err error
	if rsp.Message, err = msg.Serialize(); err != nil {
		logger.Errorf("decode msg failed: %v", err)
		return errors.InternalServerError("DECODE_ERROR", "decode error")
	}

	return nil
}

func (h *Chats) ListMessage(ctx context.Context, req *chats.ListMessageRequest, rsp *chats.ListMessageResponse) error {
	logger.Info("chats.ListMessage request")
	if len(req.ChatId) == 0 {
		return ErrMissingChatId
	}

	// build query
	q := h.DB.Where(&Chat{ID: req.ChatId}).Preload("Messages", func(db *gorm.DB) *gorm.DB {
		limit := 25
		if req.Limit != nil {
			limit = int(req.Limit.Value)
		}
		// one day ago
		after := time.Now().Add(-24 * time.Hour)
		if req.After != nil && req.After.IsValid() {
			after = req.After.AsTime()
		}
		return db.Where("send_at > ?", after).Order("send_at desc").Limit(limit)
	})

	// lookup chat & its messages
	var chat Chat
	if err := q.First(&chat).Error; err == gorm.ErrRecordNotFound {
		return ErrChatNotFound
	} else if err != nil {
		logger.Errorf("lookup chat error: %v", err)
		return ErrDatabase
	}

	rsp.Messages = make([]*chats.Message, len(chat.Messages))
	for i, msg := range chat.Messages {
		var err error
		rsp.Messages[i], err = msg.Serialize()
		if err != nil {
			logger.Errorf("decode msg failed: %v", err)
			return errors.InternalServerError("DECODE_ERROR", "decode error")
		}
	}
	return nil
}
