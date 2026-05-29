package push

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"sync"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

var (
	fcmOnce   sync.Once
	fcmClient *messaging.Client
	fcmErr    error
)

func fcmEnabled() bool {
	return os.Getenv("FCM_ENABLED") == "true" || os.Getenv("FIREBASE_CREDENTIALS_JSON") != "" || os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") != ""
}

func getFCMClient(ctx context.Context) (*messaging.Client, error) {
	if !fcmEnabled() {
		return nil, errFCMDisabled
	}
	fcmOnce.Do(func() {
		var opts []option.ClientOption
		if raw := os.Getenv("FIREBASE_CREDENTIALS_JSON"); raw != "" {
			opts = append(opts, option.WithCredentialsJSON([]byte(raw)))
		}
		app, err := firebase.NewApp(ctx, nil, opts...)
		if err != nil {
			fcmErr = err
			return
		}
		fcmClient, fcmErr = app.Messaging(ctx)
	})
	return fcmClient, fcmErr
}

var errFCMDisabled = &pushError{"FCM not configured"}

type pushError struct{ msg string }

func (e *pushError) Error() string { return e.msg }

func sendFCM(ctx context.Context, token, title, body, link string) error {
	client, err := getFCMClient(ctx)
	if err != nil {
		return err
	}
	msg := &messaging.Message{
		Token: token,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Webpush: &messaging.WebpushConfig{
			Notification: &messaging.WebpushNotification{
				Icon: "/logo/android-chrome-192x192.png",
			},
			FCMOptions: &messaging.WebpushFCMOptions{
				Link: link,
			},
		},
	}
	_, err = client.Send(ctx, msg)
	if err != nil {
		if messaging.IsRegistrationTokenNotRegistered(err) || messaging.IsInvalidArgument(err) {
			_ = revokeToken(token)
		}
		return err
	}
	return nil
}

func logFCMConfigStatus() {
	if !fcmEnabled() {
		log.Println("push: FCM disabled (set FCM_ENABLED=true and FIREBASE_CREDENTIALS_JSON)")
		return
	}
	ctx := context.Background()
	_, err := getFCMClient(ctx)
	if err != nil {
		log.Printf("push: FCM init failed: %v", err)
		return
	}
	log.Println("push: FCM ready")
}

func debugFCMCredentials() {
	raw := os.Getenv("FIREBASE_CREDENTIALS_JSON")
	if raw == "" {
		return
	}
	var m map[string]interface{}
	if json.Unmarshal([]byte(raw), &m) != nil {
		log.Println("push: FIREBASE_CREDENTIALS_JSON is not valid JSON")
	}
}
