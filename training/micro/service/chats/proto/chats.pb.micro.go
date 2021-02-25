// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: proto/chats.proto

package chats

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	math "math"
)

import (
	context "context"
	api "github.com/micro/micro/v3/service/api"
	client "github.com/micro/micro/v3/service/client"
	server "github.com/micro/micro/v3/service/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ api.Endpoint
var _ context.Context
var _ client.Option
var _ server.Option

// Api Endpoints for Chats service

func NewChatsEndpoints() []*api.Endpoint {
	return []*api.Endpoint{}
}

// Client API for Chats service

type ChatsService interface {
	CreateChat(ctx context.Context, in *CreateChatRequest, opts ...client.CallOption) (*CreateChatResponse, error)
}

type chatsService struct {
	c    client.Client
	name string
}

func NewChatsService(name string, c client.Client) ChatsService {
	return &chatsService{
		c:    c,
		name: name,
	}
}

func (c *chatsService) CreateChat(ctx context.Context, in *CreateChatRequest, opts ...client.CallOption) (*CreateChatResponse, error) {
	req := c.c.NewRequest(c.name, "Chats.CreateChat", in)
	out := new(CreateChatResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Chats service

type ChatsHandler interface {
	CreateChat(context.Context, *CreateChatRequest, *CreateChatResponse) error
}

func RegisterChatsHandler(s server.Server, hdlr ChatsHandler, opts ...server.HandlerOption) error {
	type chats interface {
		CreateChat(ctx context.Context, in *CreateChatRequest, out *CreateChatResponse) error
	}
	type Chats struct {
		chats
	}
	h := &chatsHandler{hdlr}
	return s.Handle(s.NewHandler(&Chats{h}, opts...))
}

type chatsHandler struct {
	ChatsHandler
}

func (h *chatsHandler) CreateChat(ctx context.Context, in *CreateChatRequest, out *CreateChatResponse) error {
	return h.ChatsHandler.CreateChat(ctx, in, out)
}
