package handler

import (
	chats "chats/proto"
	"context"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func testHandler(t *testing.T) *Chats {
	const dsn = "postgresql://root:root@localhost:5432/chats?sslmode=disable"

	// connect db
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	require.NoError(t, err)

	// truncate table
	err = db.Exec("TRUNCATE TABLE chats CASCADE").Error
	require.NoError(t, err)

	// migration
	err = db.AutoMigrate(
		&Chat{},
		&Message{},
	)
	require.NoError(t, err)

	return &Chats{
		DB: db,
	}
}

func TestChats_Handler(t *testing.T) {
	h := testHandler(t)
	require.NotNil(t, h)
}

func TestChats_CreateChat(t *testing.T) {
	// handler
	h := testHandler(t)
	require.NotNil(t, h)

	tests := []struct {
		name string
		req  *chats.CreateChatRequest
		err  error
	}{
		{
			name: "MissingUserIDs",
			req:  &chats.CreateChatRequest{},
			err:  ErrMissingUserIds,
		},
		{
			name: "MissingUserIDs",
			req: &chats.CreateChatRequest{
				UserIds: []string{"1"},
			},
			err: ErrMissingUserIds,
		},
		{
			name: "Valid",
			req: &chats.CreateChatRequest{
				UserIds: []string{"1", "2"},
			},
		},
		{
			name: "ExistingUserIDs",
			req: &chats.CreateChatRequest{
				UserIds: []string{"2", "1"},
			},
		},
	}
	var chat *chats.Chat
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rsp := &chats.CreateChatResponse{}
			err := h.CreateChat(context.TODO(), tt.req, rsp)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				require.Nil(t, rsp.Chat)
			} else {
				require.NoError(t, err)
				require.NotNil(t, rsp.Chat)
				sort.Strings(tt.req.UserIds)
				require.Equal(t, rsp.Chat.UserIds, tt.req.UserIds)
				require.NotNil(t, rsp.Chat.CreatedAt)
				if chat == nil {
					chat = rsp.Chat
				} else {
					require.True(t, rsp.Chat.CreatedAt.AsTime().Equal(chat.CreatedAt.AsTime()))
				}
			}
		})
	}
}

func TestChats_CreateMessage(t *testing.T) {
	// handler
	h := testHandler(t)
	require.NotNil(t, h)

	tests := []struct {
		name   string
		chatId string
		author string
		text   string
		err    error
	}{
		{
			name: "MissingChatID",
			err:  ErrMissingChatId,
		},
		{
			name:   "MissingAuthor",
			chatId: "a",
			err:    ErrMissingAuthor,
		},
		{
			name:   "ChatNotFound",
			chatId: "a",
			author: "a",
			err:    ErrChatNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &chats.CreateMessageRequest{
				ChatId: tt.chatId,
				Author: tt.author,
				Text:   tt.text,
			}
			rsp := &chats.CreateMessageResponse{}
			err := h.CreateMessage(context.TODO(), req, rsp)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				require.Nil(t, rsp.Message)
			}
		})
	}

	// mockup chat
	req := &chats.CreateChatRequest{
		UserIds: []string{"1", "2"},
	}
	rsp := &chats.CreateChatResponse{}
	err := h.CreateChat(context.TODO(), req, rsp)
	require.NoError(t, err)
	require.NotNil(t, rsp.Chat)

	// test
	tests2 := []struct {
		name   string
		chatId string
		author string
		text   string
		err    error
	}{
		{
			name:   "AuthorNotFound",
			chatId: rsp.Chat.Id,
			author: "a",
			err:    ErrAuthorNotFound,
		},
		{
			name:   "Valid",
			chatId: rsp.Chat.Id,
			author: rsp.Chat.UserIds[0],
			text:   "text",
		},
	}
	for _, tt := range tests2 {
		t.Run(tt.name, func(t *testing.T) {
			req := &chats.CreateMessageRequest{
				ChatId: tt.chatId,
				Author: tt.author,
				Text:   tt.text,
			}
			rsp := &chats.CreateMessageResponse{}
			err := h.CreateMessage(context.TODO(), req, rsp)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				require.Nil(t, rsp.Message)
			} else {
				require.NoError(t, err)
				require.NotNil(t, rsp.Message)
				require.NotEmpty(t, rsp.Message.Id)
				require.Equal(t, rsp.Message.ChatId, tt.chatId)
				require.Equal(t, rsp.Message.Author, tt.author)
				require.Equal(t, rsp.Message.Text, tt.text)
				require.NotNil(t, rsp.Message.SendAt)
			}
		})
	}
}
