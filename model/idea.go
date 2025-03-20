package model

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/arwoosa/post/pkg/manticore"
)

// IdeaData 定義了 idea 表的資料結構
type IdeaData struct {
	ID                 uint64  `json:"id"`
	IdeaID             uint64  `json:"idea_id"`
	ItineraryName      string  `json:"itinerary_name"`
	AttractionName     string  `json:"attraction_name"`
	Tags               string  `json:"tags"`
	WildMode           string  `json:"wild_mode"`
	AttractionLocation string  `json:"attraction_location"`
	ExperienceDuration float64 `json:"experience_duration"`
}

// ToMap 將 IdeaData 轉換為 map
func (d *IdeaData) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":                  d.ID,
		"idea_id":             d.IdeaID,
		"itinerary_name":      d.ItineraryName,
		"attraction_name":     d.AttractionName,
		"tags":                d.Tags,
		"wild_mode":           d.WildMode,
		"attraction_location": d.AttractionLocation,
		"experience_duration": d.ExperienceDuration,
	}
}

// IdeaFromDocument 從 Document 轉換為 IdeaData
func IdeaFromDocument(doc *manticore.Document) (*IdeaData, error) {
	if doc == nil {
		return nil, fmt.Errorf("document is nil")
	}

	id, _ := getUint64(doc.Data, "id")
	ideaID, _ := getUint64(doc.Data, "idea_id")
	experienceDuration, _ := getFloat64(doc.Data, "experience_duration")

	data := &IdeaData{
		ID:                 id,
		IdeaID:             ideaID,
		ItineraryName:      getString(doc.Data, "itinerary_name"),
		AttractionName:     getString(doc.Data, "attraction_name"),
		Tags:               getString(doc.Data, "tags"),
		WildMode:           getString(doc.Data, "wild_mode"),
		AttractionLocation: getString(doc.Data, "attraction_location"),
		ExperienceDuration: experienceDuration,
	}

	return data, nil
}

// 輔助函數
func getString(data map[string]interface{}, key string) string {
	if val, ok := data[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func getUint64(data map[string]interface{}, key string) (uint64, error) {
	if val, ok := data[key]; ok {
		switch v := val.(type) {
		case uint64:
			return v, nil
		case int64:
			return uint64(v), nil
		case string:
			return strconv.ParseUint(v, 10, 64)
		}
	}
	return 0, fmt.Errorf("invalid value for key: %s", key)
}
func getFloat64(data map[string]interface{}, key string) (float64, error) {
	if val, ok := data[key]; ok {
		switch v := val.(type) {
		case float64:
			return v, nil
		case float32:
			return float64(v), nil
		case string:
			return strconv.ParseFloat(v, 64)
		}
	}
	return 0, fmt.Errorf("invalid value for key: %s", key)
}

// GetTags 將 tags 字串轉換為字串切片
func (d *IdeaData) GetTags() []string {
	if d.Tags == "" {
		return nil
	}
	return strings.Split(d.Tags, ",")
}

// SetTags 將字串切片轉換為 tags 字串
func (d *IdeaData) SetTags(tags []string) {
	if len(tags) == 0 {
		d.Tags = ""
		return
	}
	d.Tags = strings.Join(tags, ",")
}
