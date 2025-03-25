package manticore

import (
	"context"
	"fmt"

	Manticoresearch "github.com/manticoresoftware/manticoresearch-go"
)

// Create 實現文件創建
func (c *manticore) Create(table string, data map[string]interface{}) (int64, error) {
	ctx := context.Background()
	req := Manticoresearch.NewInsertDocumentRequest(table, data)

	successRes, httpRes, err := c.apiClient.IndexAPI.Insert(ctx).InsertDocumentRequest(*req).Execute()
	if err != nil {
		return 0, fmt.Errorf("create document failed: %w", err)
	}

	if httpRes.StatusCode != 200 {
		return 0, fmt.Errorf("create document failed with status code: %d", httpRes.StatusCode)
	}
	if successRes != nil && successRes.Id != nil {
		return *successRes.Id, nil
	}

	return 0, fmt.Errorf("failed to get document ID from response")
}

// Read 實現文件讀取
/*func (c *manticore) Read(table string, id uint64) (*Document, error) {
	ctx := context.Background()

	// 使用 SQL 查詢來獲取文檔
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = %d", table, id)
	_, httpRes, err := c.apiClient.UtilsAPI.Sql(ctx).Body(query).Execute()
	if err != nil {
		return nil, fmt.Errorf("read document failed: %w", err)
	}

	if httpRes.StatusCode != 200 {
		return nil, fmt.Errorf("read document failed with status code: %d", httpRes.StatusCode)
	}

	// TODO: 從響應中解析文檔
	return &Document{
		ID:    id,
		Index: table,
		Data:  make(map[string]interface{}),
	}, nil
}*/

// Replace 實現文件更新，如果文檔不存在，則創建新文檔，如果文檔存在，則更新文檔
func (c *manticore) Replace(table string, id int64, data map[string]interface{}) error {
	ctx := context.Background()

	req := Manticoresearch.NewInsertDocumentRequest(table, data)
	req.SetId(id)

	// 執行 replace 操作
	_, httpRes, err := c.apiClient.IndexAPI.Replace(ctx).InsertDocumentRequest(*req).Execute()
	if err != nil {
		return fmt.Errorf("replace document failed: %w", err)
	}

	if httpRes.StatusCode != 200 {
		return fmt.Errorf("replace document failed with status code: %d", httpRes.StatusCode)
	}

	return nil
}

// Delete 實現文件刪除
func (c *manticore) Delete(table string, id int64) error {
	ctx := context.Background()
	req := Manticoresearch.NewDeleteDocumentRequest(table)
	req.SetId(id)

	_, httpRes, err := c.apiClient.IndexAPI.Delete(ctx).DeleteDocumentRequest(*req).Execute()
	if err != nil {
		return fmt.Errorf("delete document failed: %w", err)
	}

	if httpRes.StatusCode != 200 {
		return fmt.Errorf("delete document failed with status code: %d", httpRes.StatusCode)
	}

	return nil
}

// Search 實現文件搜尋
func (c *manticore) Search(searchRequest *Manticoresearch.SearchRequest) (*Manticoresearch.SearchResponse, error) {
	ctx := context.Background()

	searchRes, httpRes, err := c.apiClient.SearchAPI.Search(ctx).SearchRequest(*searchRequest).Execute()
	if err != nil {
		return nil, fmt.Errorf("search documents failed: %w", err)
	}

	if httpRes.StatusCode != 200 {
		return nil, fmt.Errorf("search documents failed with status code: %d", httpRes.StatusCode)
	}

	return searchRes, nil
}
