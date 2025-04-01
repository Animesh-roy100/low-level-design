package main

import "fmt"

// Factory Design Pattern

// This is an interface that defines the SendNotification method
type Notification interface {
	SendNotification(message string)
}

// This is a Concrete class that implements the Notification interface
type EmailNotification struct{}

func (en *EmailNotification) SendNotification(message string) {
	fmt.Println("Email notification sent: ", message)
}

type SMSNotification struct{}

func (sn *SMSNotification) SendNotification(message string) {
	fmt.Println("SMS notification sent: ", message)
}

type PushNotification struct{}

func (pn *PushNotification) SendNotification(message string) {
	fmt.Println("Push notification sent: ", message)
}

// NotificationFactory function is the factory method that
// creates and returns a notification object based on the input parameter
func NotificationFactory(notificationType string) Notification {
	switch notificationType {
	case "Email":
		return &EmailNotification{}
	case "SMS":
		return &SMSNotification{}
	case "Push":
		return &PushNotification{}
	default:
		return nil
	}
}

func main() {
	smsNotification := NotificationFactory("SMS")
	smsNotification.SendNotification("Hello, SMS!")

	pushNotification := NotificationFactory("Push")
	pushNotification.SendNotification("Hello, PUSH!")

	emailNotification := NotificationFactory("Email")
	emailNotification.SendNotification("Hello, Email!")
}
