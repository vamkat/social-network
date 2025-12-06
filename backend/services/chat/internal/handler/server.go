package handler

import (
	"social-network/services/chat/internal/application"
	"social-network/shared/gen-go/chat"

	_ "github.com/lib/pq"
)

type ChatHandler struct {
	chat.UnimplementedChatServiceServer
	Application *application.ChatService
	Port        string
}
