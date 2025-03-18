package request

import (
	"errors"
	"strings"
)

// BaseIdea 包含所有共用的欄位
type BaseIdea struct {
	MongoId            int      `json:"mongo_id"`
	ItineraryName      string   `json:"itinerary_name"`
	AttractionName     string   `json:"attraction_name"`
	Tags               []string `json:"tags"`
	WildMode           string   `json:"wild_mode"`
	AttractionLocation string   `json:"attraction_location"`
	HostMessage        string   `json:"host_message"`
	ExperienceDuration float64  `json:"experience_duration"`
}
type CreateIdea struct {
	BaseIdea
}
type UpdateIdea struct {
	BaseIdea
}

// Validate 驗證基礎欄位
func (b *BaseIdea) Validate() error {
	if b == nil {
		return errors.New("nil request")
	}
	if b.MongoId <= 0 {
		return errors.New("mongo_id must be greater than zero")
	}

	if strings.TrimSpace(b.ItineraryName) == "" {
		return errors.New("itinerary_name is required")
	}

	if strings.TrimSpace(b.AttractionName) == "" {
		return errors.New("empty attraction_name")
	}

	if len(b.Tags) == 0 {
		return errors.New("empty tags")
	}

	if strings.TrimSpace(b.WildMode) == "" {
		return errors.New("empty wild_mode")
	}

	if strings.TrimSpace(b.AttractionLocation) == "" {
		return errors.New("empty attraction_location")
	}

	if strings.TrimSpace(b.HostMessage) == "" {
		return errors.New("empty host_message")
	}

	if b.ExperienceDuration <= 0 {
		return errors.New("experience_duration smaller than or equal to zero")
	}

	return nil
}
