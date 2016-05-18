package xrmgo

import (
	"fmt"
	"time"

	"github.com/satori/go.uuid"
)

func guid() string {
	return fmt.Sprintf("urn:uuid:%s", uuid.NewV4())
}

func contracted(action string) string {
	return "http://schemas.microsoft.com/xrm/2011/Contracts/Services/IOrganizationService/" + action
}

func toCurrentTime(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

func toTomorrowTime(t time.Time) string {
	return t.UTC().Add(time.Hour * 24).Format(time.RFC3339)
}
