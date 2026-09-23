package mercadolibremarketplace

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func notificationJSON(topic, resource string, user, application int64) []byte {
	value := Notification{ID: "provider-id", Resource: resource, UserID: user, Topic: topic, ApplicationID: application, Attempts: 1, Sent: "2026-08-26T10:00:00Z", Received: "2026-08-26T10:00:01Z"}
	result, _ := json.Marshal(value)
	return result
}

func TestNotificationProducesDurableFetchJob(t *testing.T) {
	job, err := NormalizeNotification(provenProfile(), notificationJSON("orders_v2", "/orders/123456789", 123456789, 987654321))
	if err != nil {
		t.Fatal(err)
	}
	if job.Resource != "/orders/123456789" || job.Topic != "orders_v2" || len(job.DeduplicationKey) != 64 {
		t.Fatal("invalid fetch job")
	}
	output := filepath.Join(t.TempDir(), "job")
	if err := WriteFetchJob(job, output); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(output, "FETCH_JOB.json"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "987654321") || strings.Contains(string(data), "provider-id") {
		t.Fatal("fetch job leaked envelope identifiers")
	}
}

func TestQuestionNotificationProducesAllowlistedFetchJob(t *testing.T) {
	job, err := NormalizeNotification(provenProfile(), notificationJSON("questions", "/questions/11751825075", 123456789, 987654321))
	if err != nil {
		t.Fatal(err)
	}
	if job.Resource != "/questions/11751825075" || job.Topic != "questions" || len(job.DeduplicationKey) != 64 {
		t.Fatalf("invalid question fetch job: %+v", job)
	}
}

func TestNotificationIdentityTopicAndResourceFailClosed(t *testing.T) {
	if _, err := NormalizeNotification(provenProfile(), notificationJSON("orders_v2", "/orders/123456789", 111111111, 987654321)); err == nil {
		t.Fatal("wrong seller passed")
	}
	if _, err := NormalizeNotification(provenProfile(), notificationJSON("unknown", "/orders/123456789", 123456789, 987654321)); err == nil {
		t.Fatal("unknown topic passed")
	}
	if _, err := NormalizeNotification(provenProfile(), notificationJSON("orders_v2", "https://evil.example/x", 123456789, 987654321)); err == nil {
		t.Fatal("untrusted resource passed")
	}
}

func TestNotificationRequiresOriginControl(t *testing.T) {
	profile := provenProfile()
	profile.NotificationOriginControlProven = false
	if _, err := NormalizeNotification(profile, notificationJSON("items", "/items/MLA1234567890", 123456789, 987654321)); err == nil {
		t.Fatal("unproven origin control passed")
	}
}
