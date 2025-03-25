package manticore

import Manticoresearch "github.com/manticoresoftware/manticoresearch-go"

// SearchResult 代表搜尋結果
type SearchResult struct {
	Total  int64                    `json:"total"`
	Data   []map[string]interface{} `json:"data"`
	Scroll string                   `json:"scroll"`
}

// ManticoreService 定義了 Manticore Search 服務的基本介面
type ManticoreService interface {
	// Create 創建新文件
	Create(index string, data map[string]interface{}) (int64, error)

	// Replace 更新文件
	Replace(index string, id int64, data map[string]interface{}) error

	// Delete 刪除文件
	Delete(index string, id int64) error

	// Search 搜尋文件
	Search(searchRequest *Manticoresearch.SearchRequest) (*Manticoresearch.SearchResponse, error)

	// Health 健康檢查
	Health() (bool, error)
}
