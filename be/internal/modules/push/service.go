package push

import (
	"context"
	"log"
	"os"
	"time"

	"be/internal/db"
	"be/internal/modules/users"

	"github.com/google/uuid"
)

type RegisterDeviceInput struct {
	FCMToken   string
	DeviceType string
	UserAgent  string
}

func RegisterDevice(userID int, in RegisterDeviceInput) (*UserDevice, error) {
	now := time.Now()
	d := &UserDevice{
		ID:         uuid.New(),
		UserID:     userID,
		FCMToken:   in.FCMToken,
		DeviceType: in.DeviceType,
		UserAgent:  in.UserAgent,
		LastSeenAt: now,
		CreatedAt:  now,
	}
	if err := upsertDevice(d); err != nil {
		return nil, err
	}
	_ = touchUserLastSeen(userID, now)
	return d, nil
}

func Heartbeat(userID int, fcmToken string) error {
	now := time.Now()
	if err := touchDevice(fcmToken, userID, now); err != nil {
		return err
	}
	return touchUserLastSeen(userID, now)
}

func RevokeDevice(userID int, fcmToken string) error {
	return revokeDevice(fcmToken, userID)
}

func touchUserLastSeen(userID int, at time.Time) error {
	return db.DB.Model(&users.User{}).Where("id = ?", userID).Update("last_seen_at", at).Error
}

func TouchUserLastSeen(userID int) {
	_ = touchUserLastSeen(userID, time.Now())
}

type pushContent struct {
	kind  string
	title string
	body  string
}

func contentForInactive(days int, _ string) pushContent {
	switch {
	case days >= 7:
		return pushContent{kind: KindMiss7Day, title: "Hani 💕", body: ""}
	case days >= 3:
		return pushContent{
			kind:  KindMiss3Day,
			title: "Hani 💕",
			body:  "Lâu quá không thấy anh... em lo quá 🥺",
		}
	case days >= 1:
		return pushContent{
			kind:  KindMiss1Day,
			title: "Hani ❤️",
			body:  "Anh bận hả? Em nhớ anh đó 🥺",
		}
	default:
		return pushContent{}
	}
}

func SendToUser(ctx context.Context, userID int, kind, title, body string) error {
	tokens, err := activeTokensForUser(userID)
	if err != nil {
		return err
	}
	if len(tokens) == 0 {
		return nil
	}

	link := envOr("PUSH_APP_URL", "https://hani.app/")
	var sendErr error
	sent := 0
	for _, token := range tokens {
		if err := sendFCM(ctx, token, title, body, link); err != nil {
			log.Printf("push: send user=%d token=%s… err=%v", userID, token[:min(8, len(token))], err)
			sendErr = err
			continue
		}
		sent++
	}
	if sent == 0 && sendErr != nil {
		return sendErr
	}

	n := &Notification{
		ID:     uuid.New(),
		UserID: userID,
		Kind:   kind,
		Title:  title,
		Body:   body,
		SentAt: time.Now(),
	}
	return insertNotification(n)
}

// SendTestToUser sends a one-off push without recording re-engagement history.
func SendTestToUser(ctx context.Context, userID int, title, body string) error {
	if !fcmEnabled() {
		return errFCMDisabled
	}
	tokens, err := activeTokensForUser(userID)
	if err != nil {
		return err
	}
	if len(tokens) == 0 {
		return &pushError{"no registered device — bật thông báo trong Settings trước"}
	}

	link := envOr("PUSH_APP_URL", "https://hani.app/")
	var lastErr error
	sent := 0
	for _, token := range tokens {
		if err := sendFCM(ctx, token, title, body, link); err != nil {
			lastErr = err
			continue
		}
		sent++
	}
	if sent == 0 && lastErr != nil {
		return lastErr
	}
	return nil
}

func RunReengagement(ctx context.Context, now time.Time) {
	if !fcmEnabled() {
		return
	}

	loc := koreaLocation()
	day := now.In(loc)

	cutoff := now.Add(-24 * time.Hour)
	inactive, err := usersInactiveSince(cutoff)
	if err != nil {
		log.Printf("push: reengage query: %v", err)
		return
	}

	for _, row := range inactive {
		days := int(now.Sub(row.LastSeenAt).Hours() / 24)
		content := contentForInactive(days, "")
		if content.kind == "" {
			continue
		}

		already, err := sentKindToday(row.UserID, content.kind, day)
		if err != nil || already {
			continue
		}

		u, err := users.GetUserByIDService(itoa(row.UserID))
		if err != nil || !u.IsActive {
			continue
		}

		if content.kind == KindMiss7Day {
			title, body, err := generateMiss7Day(ctx, u.Name)
			if err != nil {
				log.Printf("push: ai miss_7d user=%d: %v", row.UserID, err)
				content.body = "Lâu quá rồi không thấy anh... em nhớ anh lắm 🥺"
			} else {
				content.title = title
				content.body = body
			}
		}

		if err := SendToUser(ctx, row.UserID, content.kind, content.title, content.body); err != nil {
			log.Printf("push: reengage user=%d kind=%s: %v", row.UserID, content.kind, err)
		}
	}
}

func koreaLocation() *time.Location {
	loc, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		return time.FixedZone("KST", 9*3600)
	}
	return loc
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
