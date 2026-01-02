package main

import (
	"fmt"
	"sync"
	"time"
)

type MemberType int

const (
	CUSTOMER MemberType = iota
	SELLER
	DELIVERY
)

type OrderEventType int

const (
	DELIVERED OrderEventType = iota
	SHIPPED
	ORDERED
)

type NotificationChannelType int

const (
	EMAIL NotificationChannelType = iota
	SMS
)

type Member struct {
	MemberID   string
	MemberType MemberType
	Name       string
}

func NewMember(memberType MemberType, name string) *Member {
	return &Member{
		MemberID:   "",
		MemberType: memberType,
		Name:       name,
	}
}

func GetMemberByID(memberId string) *Member { return nil }

type Notification struct {
	TimeStamp time.Time
	OrderID   string
	MemberID  string
	Channel   NotificationChannelType
	EventType OrderEventType
	Message   string
}

func NewNotification(orderId, memberId, message string, channel NotificationChannelType, eventType OrderEventType) *Notification {
	return &Notification{
		TimeStamp: time.Now(),
		OrderID:   orderId,
		MemberID:  memberId,
		Message:   message,
		Channel:   channel,
		EventType: eventType,
	}
}

type NotificationChannel interface {
	Send(notification *Notification)
	GetType() NotificationChannelType
}

type SMSNotificationChannel struct{}

func (s *SMSNotificationChannel) Send(notification *Notification)  { fmt.Println(notification) }
func (s *SMSNotificationChannel) GetType() NotificationChannelType { return EMAIL }

type EmailNotificationChannel struct{}

func (e *EmailNotificationChannel) Send(notification *Notification)  { fmt.Println(notification) }
func (e *EmailNotificationChannel) GetType() NotificationChannelType { return SMS }

type NotificationChannelFactory struct {
	NotificationChannels map[NotificationChannelType]NotificationChannel
	mu                   sync.Mutex
}

func (n *NotificationChannelFactory) GetNotificationChannel(notificationType NotificationChannelType) NotificationChannel {
	n.mu.Lock()
	defer n.mu.Unlock()

	if channel, ok := n.NotificationChannels[notificationType]; ok {
		return channel
	}

	var channel NotificationChannel

	switch notificationType {
	case EMAIL:
		channel = &EmailNotificationChannel{}
	case SMS:
		channel = &SMSNotificationChannel{}
	}

	n.NotificationChannels[notificationType] = channel

	return channel
}

type Subscription struct {
	SubscriptionID string
	Member         *Member
	Channels       []NotificationChannel
}

func NewSubscription(subscriptionId string, member *Member, channels []NotificationChannel) *Subscription {
	return &Subscription{
		SubscriptionID: subscriptionId,
		Member:         member,
		Channels:       append([]NotificationChannel{}, channels...),
	}
}

func (s *Subscription) AddChannel(channel NotificationChannel) {
	for _, ch := range s.Channels {
		if ch == channel {
			return
		}
	}
	s.Channels = append(s.Channels, channel)
}

func (s *Subscription) RemoveChannel(channel NotificationChannel) {
	for i, ch := range s.Channels {
		if ch == channel {
			s.Channels = append(s.Channels[:i], s.Channels[i+1:]...)
			return
		}
	}
}

func (s *Subscription) Update(orderId string, event OrderEventType) {

}
