package grpc

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/justcgh9/discord-clone-friends/internal/app/grpc/client"
	"github.com/justcgh9/discord-clone-friends/internal/models"
	"github.com/justcgh9/discord-clone-proto/gen/go/friends"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Friends interface {
	AcceptRequest(ctx context.Context, userID, requesterID string) error
	DenyRequest(ctx context.Context, userID, requesterID string) error
	RemoveFriend(ctx context.Context, userID, friendID string) error
	ListFriends(ctx context.Context, userID string) ([]models.Friend, error)
	BlockUser(ctx context.Context, userID, targetID string) error
	SendRequest(ctx context.Context, fromUserID, toUserID string) error
}

type AuthValidator interface {
	Verify(ctx context.Context, token string) <-chan client.LoginResult
}

type serverAPI struct {
	friends.UnimplementedFriendServiceServer
	svc  Friends
	auth AuthValidator
}

func RegisterServer(srv *grpc.Server, svc Friends, auth AuthValidator) {
	friends.RegisterFriendServiceServer(srv, &serverAPI{svc: svc, auth: auth})
}

func extractToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("missing metadata")
	}
	values := md.Get("authorization")
	log.Println(values)
	if len(values) == 0 {
		return "", errors.New("missing authorization token")
	}
	return values[0], nil
}

func (s *serverAPI) verifyUser(ctx context.Context) (*client.LoginResult, error) {
	token, err := extractToken(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "auth error: %v", err)
	}

	res := <-s.auth.Verify(ctx, token)
	if res.Err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "auth failed: %v", res.Err)
	}
	return &res, nil
}

func (s *serverAPI) SendRequest(ctx context.Context, req *friends.FriendRequest) (*friends.FriendResponse, error) {
	if req == nil || req.FromUserId == "" || req.ToUserId == "" {
		return nil, status.Error(codes.InvalidArgument, "missing from_user_id or to_user_id")
	}

	res, err := s.verifyUser(ctx)
	if err != nil {
		return nil, err
	}
	if req.FromUserId != strconv.Itoa(int(res.User.User.UserId)) {
		return nil, status.Error(codes.PermissionDenied, "cannot act on behalf of another user")
	}

	err = s.svc.SendRequest(ctx, req.FromUserId, req.ToUserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "send request failed: %v", err)
	}
	return &friends.FriendResponse{Message: "Friend request sent"}, nil
}

func (s *serverAPI) AcceptRequest(ctx context.Context, req *friends.FriendAction) (*friends.FriendResponse, error) {
	if req == nil || req.UserId == "" || req.TargetId == "" {
		return nil, status.Error(codes.InvalidArgument, "missing user_id or target_id")
	}

	res, err := s.verifyUser(ctx)
	if err != nil {
		return nil, err
	}
	if req.UserId != strconv.Itoa(int(res.User.User.UserId)) {
		return nil, status.Error(codes.PermissionDenied, "cannot act on behalf of another user")
	}

	err = s.svc.AcceptRequest(ctx, req.UserId, req.TargetId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "accept request failed: %v", err)
	}
	return &friends.FriendResponse{Message: "Friend request accepted"}, nil
}

func (s *serverAPI) DenyRequest(ctx context.Context, req *friends.FriendAction) (*friends.FriendResponse, error) {
	if req == nil || req.UserId == "" || req.TargetId == "" {
		return nil, status.Error(codes.InvalidArgument, "missing user_id or target_id")
	}

	res, err := s.verifyUser(ctx)
	if err != nil {
		return nil, err
	}
	if req.UserId != strconv.Itoa(int(res.User.User.UserId)) {
		return nil, status.Error(codes.PermissionDenied, "cannot act on behalf of another user")
	}

	err = s.svc.DenyRequest(ctx, req.UserId, req.TargetId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "deny request failed: %v", err)
	}
	return &friends.FriendResponse{Message: "Friend request denied"}, nil
}

func (s *serverAPI) RemoveFriend(ctx context.Context, req *friends.FriendAction) (*friends.FriendResponse, error) {
	if req == nil || req.UserId == "" || req.TargetId == "" {
		return nil, status.Error(codes.InvalidArgument, "missing user_id or target_id")
	}

	res, err := s.verifyUser(ctx)
	if err != nil {
		return nil, err
	}
	if req.UserId != strconv.Itoa(int(res.User.User.UserId)) {
		return nil, status.Error(codes.PermissionDenied, "cannot act on behalf of another user")
	}

	err = s.svc.RemoveFriend(ctx, req.UserId, req.TargetId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "remove friend failed: %v", err)
	}
	return &friends.FriendResponse{Message: "Friend removed"}, nil
}

func (s *serverAPI) BlockUser(ctx context.Context, req *friends.FriendAction) (*friends.FriendResponse, error) {
	if req == nil || req.UserId == "" || req.TargetId == "" {
		return nil, status.Error(codes.InvalidArgument, "missing user_id or target_id")
	}
	if req.UserId == req.TargetId {
		return nil, status.Error(codes.InvalidArgument, "cannot block yourself")
	}

	res, err := s.verifyUser(ctx)
	if err != nil {
		return nil, err
	}
	if req.UserId != strconv.Itoa(int(res.User.User.UserId)) {
		return nil, status.Error(codes.PermissionDenied, "cannot act on behalf of another user")
	}

	err = s.svc.BlockUser(ctx, req.UserId, req.TargetId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "block user failed: %v", err)
	}
	return &friends.FriendResponse{Message: "User blocked"}, nil
}

func (s *serverAPI) ListFriends(ctx context.Context, req *friends.UserID) (*friends.FriendList, error) {
	if req == nil || req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "missing user_id")
	}

	res, err := s.verifyUser(ctx)
	if err != nil {
		return nil, err
	}
	if req.UserId != strconv.Itoa(int(res.User.User.UserId)) {
		return nil, status.Error(codes.PermissionDenied, "cannot view another user's friends")
	}

	list, err := s.svc.ListFriends(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list friends failed: %v", err)
	}

	resp := &friends.FriendList{}
	for _, f := range list {
		resp.Friends = append(resp.Friends, &friends.Friend{
			Id:     f.ID,
			Handle: f.Handle,
			Status: string(f.Status),
		})
	}
	return resp, nil
}
