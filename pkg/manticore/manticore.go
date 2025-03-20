package manticore

// Document 代表一個 Manticore Search 文件
type Document struct {
	ID    int64                  `json:"id"`
	Index string                 `json:"index"`
	Data  map[string]interface{} `json:"data"`
}

// SearchResult 代表搜尋結果
type SearchResult struct {
	Total       int64      `json:"total"`
	Documents   []Document `json:"documents"`
	SearchAfter string     // 用於分頁的 token
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
	Search(table string, query string, searchAfter string, limit int32) (*SearchResult, error)

	// Health 健康檢查
	Health() (bool, error)
}
