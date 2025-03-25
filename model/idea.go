package model

import (
	"strconv"
	"strings"

	manticoresearch "github.com/manticoresoftware/manticoresearch-go"
)

// IdeaData 定義了 idea 表的資料結構
type IdeaData struct {
	ID                 uint64  `json:"id"`
	Name               string  `json:"name"`
	Rewilding_name     string  `json:"rewilding_name"`
	Rewilding_mode     string  `json:"rewilding_mode"`
	Rewilding_location string  `json:"rewilding_location"`
	Tags               string  `json:"tags"`
	Host_message       string  `json:"host_message"`
	Experience_hours   float64 `json:"experience_hours"`
}

// ToMap 將 IdeaData 轉換為 map
func (d *IdeaData) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":                 d.ID,
		"name":               d.Name,
		"rewilding_name":     d.Rewilding_name,
		"rewilding_mode":     d.Rewilding_mode,
		"rewilding_location": d.Rewilding_location,
		"tags":               d.Tags,
		"host_message":       d.Host_message,
		"experience_hours":   d.Experience_hours,
	}
}

// IdeaResponse 用於 API 回傳的資料結構
type IdeaResponse struct {
	ID                 uint64   `json:"id"`
	Name               string   `json:"name"`
	RewildingName      string   `json:"rewilding_name"`
	RewildingMode      string   `json:"rewilding_mode"`
	RewildingLocation  string   `json:"rewilding_location"`
	HostMessage        string   `json:"host_message"`
	ExperienceDuration float64  `json:"experience_duration"`
	Tags               []string `json:"tags"`
}

// SearchResponse 搜尋結果的回傳格式
type SearchResponse struct {
	Data   []IdeaResponse `json:"data"`
	Total  int64          `json:"total"`
	Scroll string         `json:"scroll"`
}

// FromManticoreResponse 從 Manticore 的搜尋結果轉換為 API 回傳格式
func FromManticoreResponse(result *manticoresearch.SearchResponse) *SearchResponse {
	ideas := make([]IdeaResponse, 0)
	if result != nil && result.Hits != nil && result.Hits.Hits != nil {
		for _, hit := range result.Hits.Hits {
			source := hit["_source"].(map[string]interface{})
			idea := IdeaResponse{
				ID:                 uint64(hit["_id"].(float64)),
				Name:               getString(source, "name"),
				RewildingName:      getString(source, "rewilding_name"),
				RewildingMode:      getString(source, "rewilding_mode"),
				RewildingLocation:  getString(source, "rewilding_location"),
				HostMessage:        getString(source, "host_message"),
				ExperienceDuration: getFloat64(source, "experience_hours"),
				Tags:               strings.Split(getString(source, "tags"), ","),
			}
			ideas = append(ideas, idea)
		}
	}

	total := int64(0)
	if result != nil && result.Hits != nil && result.Hits.Total != nil {
		total = int64(*result.Hits.Total)
	}

	scroll := ""
	if result != nil && result.Scroll != nil {
		scroll = *result.Scroll
	}

	return &SearchResponse{
		Data:   ideas,
		Total:  total,
		Scroll: scroll,
	}
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

func getFloat64(data map[string]interface{}, key string) float64 {
	if val, ok := data[key]; ok {
		switch v := val.(type) {
		case float64:
			return v
		case float32:
			return float64(v)
		case string:
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				return f
			}
		}
	}
	return 0
}
