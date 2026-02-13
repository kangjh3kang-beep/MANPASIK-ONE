package clients

import (
	"context"
	"log"
)

// InProcessNotificationClient is a simple notification client for in-process use
type InProcessNotificationClient struct{}

// NewInProcessNotificationClient creates a new in-process notification client
func NewInProcessNotificationClient() *InProcessNotificationClient {
	return &InProcessNotificationClient{}
}

// SendNotification logs the notification (in-process stub)
func (c *InProcessNotificationClient) SendNotification(ctx context.Context, userID, notifType, title, body, priority, channel string) error {
	log.Printf("[Notification] → user=%s type=%s title=%s body=%s", userID, notifType, title, body)
	return nil
}
