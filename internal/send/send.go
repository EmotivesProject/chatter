package send

import (
	"fmt"
	"os"

	"github.com/EmotivesProject/common/logger"
	"github.com/EmotivesProject/common/notification"
)

func MessageNotification(fromUser, newUsername, content string) {
	notif := notification.Notification{
		Username:   newUsername,
		Type:       notification.Message,
		Title:      "New message!",
		Message:    fmt.Sprintf("%s messaged you: %s", fromUser, content),
		Link:       fmt.Sprintf("%smessenger?talking-to=%s", os.Getenv("EMOTIVES_URL"), fromUser),
		UsernameTo: &newUsername,
	}

	_, err := notification.SendEvent(os.Getenv("NOTIFICATION_URL"), os.Getenv("NOTIFICATION_AUTH"), notif)
	if err != nil {
		logger.Error(err)
	}
}
