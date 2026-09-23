package mercadolibremarketplace

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type Notification struct {
	ID            string `json:"_id,omitempty"`
	Resource      string `json:"resource"`
	UserID        int64  `json:"user_id"`
	Topic         string `json:"topic"`
	ApplicationID int64  `json:"application_id"`
	Attempts      int    `json:"attempts"`
	Sent          string `json:"sent"`
	Received      string `json:"received"`
}

type FetchJob struct {
	Schema           string `json:"schema"`
	Topic            string `json:"topic"`
	Resource         string `json:"resource"`
	DeduplicationKey string `json:"deduplication_key"`
	ReceivedAt       string `json:"received_at"`
}

var resourcePattern = regexp.MustCompile(`^/(items/ML[A-Z][0-9]{6,20}|orders/[0-9]{3,20}|shipments/[0-9]{3,20}|questions/[0-9]{3,20})$`)

func NormalizeNotification(profile Profile, raw []byte) (FetchJob, error) {
	if profile.Decision != "PROVEN" || !profile.NotificationOriginControlProven || !profile.NotificationTopicsApproved || !profile.ReconciliationApproved {
		return FetchJob{}, errors.New("notification profile is not proven")
	}
	if len(raw) == 0 || len(raw) > 1<<20 {
		return FetchJob{}, errors.New("notification size is invalid")
	}
	var notification Notification
	decoder := json.NewDecoder(strings.NewReader(string(raw)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&notification); err != nil {
		return FetchJob{}, err
	}
	if !resourcePattern.MatchString(notification.Resource) || notification.Attempts < 1 || notification.Attempts > 100 {
		return FetchJob{}, errors.New("notification envelope is invalid")
	}
	if notification.UserID <= 0 || notification.ApplicationID <= 0 || notification.UserID != parsePositive(profile.SellerID) || notification.ApplicationID != parsePositive(profile.ApplicationID) {
		return FetchJob{}, errors.New("notification identity mismatch")
	}
	allowed := false
	for _, topic := range profile.ApprovedNotificationTopics {
		if topic == notification.Topic {
			allowed = true
		}
	}
	if !allowed || !((notification.Topic == "items" && strings.HasPrefix(notification.Resource, "/items/")) || (notification.Topic == "orders_v2" && strings.HasPrefix(notification.Resource, "/orders/")) || (notification.Topic == "shipments" && strings.HasPrefix(notification.Resource, "/shipments/")) || (notification.Topic == "questions" && strings.HasPrefix(notification.Resource, "/questions/"))) {
		return FetchJob{}, errors.New("notification topic/resource is not approved")
	}
	if _, err := time.Parse(time.RFC3339Nano, notification.Sent); err != nil {
		return FetchJob{}, errors.New("notification sent timestamp is invalid")
	}
	sum := sha256.Sum256(raw)
	return FetchJob{Schema: "elite-mercadolibre-fetch-job/v1", Topic: notification.Topic, Resource: notification.Resource, DeduplicationKey: hex.EncodeToString(sum[:]), ReceivedAt: time.Now().UTC().Format(time.RFC3339Nano)}, nil
}

func WriteFetchJob(job FetchJob, outputDirectory string) error {
	abs, err := filepath.Abs(outputDirectory)
	if err != nil {
		return err
	}
	if _, err := os.Stat(abs); !os.IsNotExist(err) {
		return errors.New("output directory must not exist")
	}
	stage := abs + ".stage"
	if err := os.Mkdir(stage, 0o700); err != nil {
		return err
	}
	committed := false
	defer func() {
		if !committed {
			_ = os.RemoveAll(stage)
		}
	}()
	data, _ := json.MarshalIndent(job, "", "  ")
	if err := os.WriteFile(filepath.Join(stage, "FETCH_JOB.json"), append(data, '\n'), 0o600); err != nil {
		return err
	}
	if err := os.Rename(stage, abs); err != nil {
		return err
	}
	committed = true
	return nil
}

func parsePositive(value string) int64 {
	var result int64
	for _, char := range value {
		if char < '0' || char > '9' {
			return -1
		}
		result = result*10 + int64(char-'0')
	}
	return result
}
