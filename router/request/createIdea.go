package request

import (
	"errors"
	"strings"
)

type CreateIdea struct {
	IdeaId             int      `json:"idea_id"`
	ItineraryName      string   `json:"itinerary_name"`
	AttractionName     string   `json:"attraction_name"`
	Tags               []string `json:"tags"`
	WildMode           string   `json:"wild_mode"`
	AttractionLocation string   `json:"attraction_location"`
	ExperienceDuration float64  `json:"experience_duration"`
}

func (r *CreateIdea) Validate() error {
	if r == nil {
		return errors.New("nil request")
	}
	if r.IdeaId <= 0 {
		return errors.New("idea_id must be greater than zero")
	}

	if strings.TrimSpace(r.ItineraryName) == "" {
		return errors.New("itinerary_name is required")
	}

	if strings.TrimSpace(r.AttractionName) == "" {
		return errors.New("empty attraction_name")
	}

	if len(r.Tags) == 0 {
		return errors.New("empty tags")
	}

	if strings.TrimSpace(r.WildMode) == "" {
		return errors.New("empty wild_mode")
	}

	if strings.TrimSpace(r.AttractionLocation) == "" {
		return errors.New("empty attraction_location")
	}

	if r.ExperienceDuration <= 0 {
		return errors.New("experience_duration smaller than or equal to zero")
	}

	return nil
}
