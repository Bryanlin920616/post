package service

import (
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
func (s *IdeaService) SearchIdeas(query string, searchAfter string, limit int32) (*manticore.SearchResult, error) {
	return s.client.Search(s.index, query, searchAfter, limit)
}
