package models


type FriendStatus string

const (
	StatusPending  FriendStatus = "PENDING"
	StatusAccepted FriendStatus = "ACCEPTED"
	StatusBlocked  FriendStatus = "BLOCKED"
)

type Friendship struct {
	UserID   string
	FriendID string
	Status   FriendStatus
}