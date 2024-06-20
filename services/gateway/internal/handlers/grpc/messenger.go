package grpc_handler

import (
	"context"

	"messenger.api/go/api"
)

func (g *grpcHandler) SendMessage(ctx context.Context, req *api.SendMessageRequest) (*api.SendMessageResponse, error) {
	return g.service.Messenger().SendMessage(ctx, req)
}

func (g *grpcHandler) ReadMessage(ctx context.Context, req *api.ReadMessageRequest) (*api.ReadMessageResponse, error) {
	return g.service.Messenger().ReadMessage(ctx, req)
}

func (g *grpcHandler) GetDialog(ctx context.Context, req *api.GetDialogRequest) (*api.GetDialogResponse, error) {
	return g.service.Messenger().GetDialog(ctx, req)
}

func (g *grpcHandler) GetDialogs(ctx context.Context, req *api.GetDialogsRequest) (*api.GetDialogsResponse, error) {
	return g.service.Messenger().GetDialogs(ctx, req)
}

func (g *grpcHandler) GetMessages(ctx context.Context, req *api.GetMessagesRequest) (*api.GetMessagesResponse, error) {
	return g.service.Messenger().GetMessages(ctx, req)
}

func (g *grpcHandler) GetUnreadDialogsCounter(ctx context.Context, req *api.GetUnreadDialogsCounterRequest) (*api.GetUnreadDialogsCounterResponse, error) {
	return g.service.Messenger().GetUnreadDialogsCounter(ctx, req)
}
