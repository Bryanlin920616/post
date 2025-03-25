package service

import (
	"fmt"
	"net/url"

	"github.com/arwoosa/post/model"
	"github.com/arwoosa/post/pkg/manticore"
)

// IdeaService 提供 idea 表的特定操作
type IdeaService struct {
	client manticore.ManticoreService
	index  string
}

// NewIdeaService 創建新的 IdeaService 實例
func NewIdeaService(client manticore.ManticoreService) *IdeaService {
	return &IdeaService{
		client: client,
		index:  "idea",
	}
}

// CreateIdea 創建新的 idea
func (s *IdeaService) CreateIdea(data *model.IdeaData) (int64, error) {
	return s.client.Create(s.index, data.ToMap())
}

// UpdateIdea 更新指定的 idea
func (s *IdeaService) ReplaceIdea(id int64, data *model.IdeaData) error {
	return s.client.Replace(s.index, id, data.ToMap())
}

// DeleteIdea 刪除指定的 idea
func (s *IdeaService) DeleteIdea(id int64) error {
	return s.client.Delete(s.index, id)
}

// SearchIdeas 搜尋 ideas
func (s *IdeaService) SearchIdeas(query string, scroll string, limit int32) (*model.SearchResponse, error) {
	filters := decodeQuery(query)

	// 使用查詢工廠創建搜尋請求
	factory := NewQueryFactory()
	searchRequest, err := factory.CreateSearchRequest(filters, "idea")
	if err != nil {
		return nil, fmt.Errorf("創建搜尋請求失敗: %w", err)
	}

	if scroll != "" {
		options := searchRequest.GetOptions()
		if options == nil {
			options = make(map[string]interface{})
		}
		options["scroll"] = scroll
		searchRequest.SetOptions(options)
	}

	if limit > 0 {
		searchRequest.SetLimit(limit)
	}

	// 執行搜尋
	result, err := s.client.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("執行搜尋失敗: %w", err)
	}

	return model.FromManticoreResponse(result), nil
}

// decodeQuery 解析 query 參數
func decodeQuery(query string) map[string]interface{} {
	filters := make(map[string]interface{})

	// 解析 URL Query String
	params, _ := url.ParseQuery(query)

	for key, values := range params {
		value := values[0] // 只取第一個值
		filters[key] = value
	}

	return filters
}
