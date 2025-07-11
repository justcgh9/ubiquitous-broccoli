package friends

import (
    "context"
    "fmt"

    "github.com/justcgh9/discord-clone-friends/internal/models"
)

type RequestsManager interface {
    SendRequest(ctx context.Context, fromUserID, toUserID string) error
    AcceptRequest(ctx context.Context, userID, requesterID string) error
    DenyRequest(ctx context.Context, userID, requesterID string) error
}

type FriendsManager interface {
    RemoveFriend(ctx context.Context, userID, friendID string) error
    ListFriends(ctx context.Context, userID string) ([]models.Friendship, error)
}

type BlockManager interface {
    BlockUser(ctx context.Context, userID, targetID string) error
}

type Service struct {
    requests RequestsManager
    friends  FriendsManager
    block    BlockManager
}

func NewService(req RequestsManager, fri FriendsManager, blk BlockManager) *Service {
    return &Service{
        requests: req,
        friends:  fri,
        block:    blk,
    }
}

func (s *Service) SendRequest(ctx context.Context, fromUserID, toUserID string) error {
    if fromUserID == toUserID {
        return fmt.Errorf("cannot send request to self")
    }
    return s.requests.SendRequest(ctx, fromUserID, toUserID)
}

func (s *Service) AcceptRequest(ctx context.Context, userID, requesterID string) error {
    return s.requests.AcceptRequest(ctx, userID, requesterID)
}

func (s *Service) DenyRequest(ctx context.Context, userID, requesterID string) error {
    return s.requests.DenyRequest(ctx, userID, requesterID)
}

func (s *Service) RemoveFriend(ctx context.Context, userID, friendID string) error {
    return s.friends.RemoveFriend(ctx, userID, friendID)
}

func (s *Service) ListFriends(ctx context.Context, userID string) ([]models.Friend, error) {
    friendships, err := s.friends.ListFriends(ctx, userID)
    if err != nil {
        return nil, err
    }

    var friends []models.Friend
    for _, f := range friendships {
        friends = append(friends, models.Friend{
            ID:     f.FriendID,
            Status: f.Status,
        })
    }
    return friends, nil
}

func (s *Service) BlockUser(ctx context.Context, userID, targetID string) error {
    if userID == targetID {
        return fmt.Errorf("cannot block yourself")
    }
    return s.block.BlockUser(ctx, userID, targetID)
}
