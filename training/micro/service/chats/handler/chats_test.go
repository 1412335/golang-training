package handler

import (
	chats "chats/proto"
	"context"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
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

func TestChats_ListMessage(t *testing.T) {
	// handler
	h := testHandler(t)
	require.NotNil(t, h)

	tests := []struct {
		name   string
		chatId string
		err    error
	}{
		{
			name: "MissingChatID",
			err:  ErrMissingChatId,
		},
		{
			name:   "ChatNotFound",
			chatId: "a",
			err:    ErrChatNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &chats.ListMessageRequest{
				ChatId: tt.chatId,
			}
			rsp := &chats.ListMessageResponse{}
			err := h.ListMessage(context.TODO(), req, rsp)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				require.Nil(t, rsp.Messages)
			} else {
				require.NoError(t, err)
				require.NotNil(t, rsp.Messages)
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

	// mockup messages
	messReq1 := &chats.CreateMessageRequest{
		ChatId: rsp.Chat.Id,
		Author: rsp.Chat.UserIds[0],
		Text:   "mess 1",
	}
	messRsp1 := &chats.CreateMessageResponse{}
	err = h.CreateMessage(context.TODO(), messReq1, messRsp1)
	require.NoError(t, err)
	require.NotNil(t, messRsp1.Message)

	messReq2 := &chats.CreateMessageRequest{
		ChatId: rsp.Chat.Id,
		Author: rsp.Chat.UserIds[1],
		Text:   "mess 2",
	}
	messRsp2 := &chats.CreateMessageResponse{}
	err = h.CreateMessage(context.TODO(), messReq2, messRsp2)
	require.NoError(t, err)
	require.NotNil(t, messRsp2.Message)

	// test
	tests2 := []struct {
		name   string
		chatId string
		after  *timestamppb.Timestamp
		limit  *wrapperspb.Int32Value
		count  int
		err    error
	}{
		{
			name:   "Valid",
			chatId: rsp.Chat.Id,
		},
		{
			name:   "ValidWithLimit",
			chatId: rsp.Chat.Id,
			limit:  wrapperspb.Int32(1),
		},
		{
			name:   "ValidWithAfterNextOneDay",
			chatId: rsp.Chat.Id,
			after:  timestamppb.New(time.Now().Add(24 * time.Hour)),
			count:  0,
		},
		{
			name:   "ValidWithAfterTwoDayBefore",
			chatId: rsp.Chat.Id,
			after:  timestamppb.New(time.Now().Add(-2 * 24 * time.Hour)),
			count:  2,
		},
	}
	for _, tt := range tests2 {
		t.Run(tt.name, func(t *testing.T) {
			req := &chats.ListMessageRequest{
				ChatId: tt.chatId,
				After:  tt.after,
				Limit:  tt.limit,
			}
			rsp := &chats.ListMessageResponse{}
			err := h.ListMessage(context.TODO(), req, rsp)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				require.Nil(t, rsp.Messages)
			} else {
				require.NoError(t, err)
				require.NotNil(t, rsp.Chat)
				// require.Equal(t, rsp.Chat.UserIds, tt.req.UserIds)
				// require.True(t, rsp.Chat.CreatedAt.AsTime().Equal(chat.CreatedAt.AsTime()))
				require.NotNil(t, rsp.Messages)
				if tt.limit != nil {
					require.Len(t, rsp.Messages, int(tt.limit.Value))
				}
				if tt.after != nil {
					require.Len(t, rsp.Messages, int(tt.count))
				}
				for _, msg := range rsp.Messages {
					require.NotNil(t, msg)
					switch msg.Id {
					case messRsp1.Message.Id:
						require.Equal(t, msg.ChatId, messRsp1.Message.ChatId)
						require.Equal(t, msg.Author, messRsp1.Message.Author)
						require.Equal(t, msg.Text, messRsp1.Message.Text)
						require.True(t, msg.SendAt.AsTime().Equal(messRsp1.Message.SendAt.AsTime()))
					case messRsp2.Message.Id:
						require.Equal(t, msg.ChatId, messRsp2.Message.ChatId)
						require.Equal(t, msg.Author, messRsp2.Message.Author)
						require.Equal(t, msg.Text, messRsp2.Message.Text)
						require.True(t, msg.SendAt.AsTime().Equal(messRsp2.Message.SendAt.AsTime()))
					default:
						t.Errorf("Unexpected message: %v", msg.Id)
					}
				}
			}
		})
	}
}
