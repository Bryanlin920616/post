package manticore

import (
	"context"
	"errors"
	"fmt"

	Manticoresearch "github.com/manticoresoftware/manticoresearch-go"
	"github.com/spf13/viper"
)

type manticore struct {
	apiClient *Manticoresearch.APIClient
}
type option func(*manticore)

// NewManticore 創建新的 Manticore Search 客戶端
func NewManticore(opts ...option) (ManticoreService, error) {
	// 創建配置
	url := viper.GetString("manticore.url")
	if url == "" {
		return nil, errors.New("manticore.url is empty")
	}
	configuration := Manticoresearch.NewConfiguration()
	configuration.Servers[0].URL = url
	apiClient := Manticoresearch.NewAPIClient(configuration)

	manticore := &manticore{
		apiClient: apiClient,
	}

	for _, opt := range opts {
		opt(manticore)
	}
	return manticore, nil
}

// Health 實現健康檢查
func (c *manticore) Health() (bool, error) {
	// 使用 SQL 查詢來檢查服務是否正常運行
	ctx := context.Background()
	_, _, err := c.apiClient.UtilsAPI.Sql(ctx).Body("SHOW STATUS").Execute()
	if err != nil {
		return false, fmt.Errorf("health check failed: %w", err)
	}
	return true, nil
}
