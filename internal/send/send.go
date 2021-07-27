package send

import (
	"fmt"
	"os"

	"github.com/TomBowyerResearchProject/common/logger"
	"github.com/TomBowyerResearchProject/common/notification"
)

func MessageNotification(fromUser, newUsername string) {
	notif := notification.Notification{
		Username:   newUsername,
		Type:       notification.Message,
		Title:      "New message!",
		Message:    fmt.Sprintf("%s messaged you", fromUser),
		Link:       fmt.Sprintf("%smessenger?talking-to=%s", os.Getenv("EMOTIVES_URL"), fromUser),
		UsernameTo: &newUsername,
	}

	logger.Infof("%s", os.Getenv("NOTIFICATION_AUTH"))

	_, err := notification.SendEvent(os.Getenv("NOTIFICATION_URL"), os.Getenv("NOTIFICATION_AUTH"), notif)
	if err != nil {
		logger.Error(err)
	}
}
