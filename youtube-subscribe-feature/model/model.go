package model

import (
	"time"

	"github.com/google/uuid"
)

type NotificationType string
type NotificationPreferences string

const (
	NewVideo      NotificationType = "NewVideo"
	CommunityPost NotificationType = "CommunityPost"
	LiveStream    NotificationType = "LiveStream"
	Other         NotificationType = "Other"
)

const (
	Email NotificationPreferences = "Email"
	Push  NotificationPreferences = "Push"
)

type User struct {
	UserID             uuid.UUID
	Email              string
	SubscribedChannels []string
}

type Channel struct {
	ChannelID   uuid.UUID
	Channel     string
	Subscribers []string // array of userIds
}

type Subscription struct {
	SubscriptionID          uuid.UUID
	UserID                  uuid.UUID
	ChannelID               uuid.UUID
	CreatedAt               time.Time
	UpdatedAt               time.Time
	NotificationPreferences NotificationPreferences
}

type Notification struct {
	NotificationID   uuid.UUID
	UserID           uuid.UUID
	ChannelID        uuid.UUID
	NotificationType NotificationType
	SentAt           time.Time
}
