package service

import (
	"fmt"
	"strings"

	openapi "github.com/manticoresoftware/manticoresearch-go"
)

// FieldType 定義欄位類型
type FieldType int

const (
	FullText   FieldType = iota // 全文搜索，跨欄位
	TextMatch                   // 單欄位文本匹配
	Attribute                   // 屬性欄位，支援 EQUAL 和 IN
	RangeField                  // 數值範圍查詢
)

// QueryType 定義查詢類型
type QueryType int

const (
	Match QueryType = iota
	Equal
	In
	RangeQuery
)

// FieldDefinition 定義欄位特性
type FieldDefinition struct {
	Type FieldType
	Name string
}

// Filter 定義過濾條件
type Filter struct {
	Field    string
	Type     FieldType
	Value    interface{}
	Operator string
}

// QueryFactory 用於創建不同類型的查詢條件
type QueryFactory struct {
	fieldDefinitions map[string]FieldDefinition
}

// NewQueryFactory 創建新的查詢工廠
func NewQueryFactory() *QueryFactory {
	return &QueryFactory{
		fieldDefinitions: map[string]FieldDefinition{
			"keyword": {
				Type: FullText,
				Name: "*",
			},
			"tags": {
				Type: TextMatch,
				Name: "tags",
			},
			"rewilding_mode": {
				Type: Attribute,
				Name: "rewilding_mode",
			},
			"rewilding_location": {
				Type: Attribute,
				Name: "rewilding_location",
			},
			"experience_hours": {
				Type: RangeField,
				Name: "experience_hours",
			},
		},
	}
}

// CreateSearchRequest 根據過濾條件創建搜尋請求
func (f *QueryFactory) CreateSearchRequest(filters map[string]interface{}, index string) (*openapi.SearchRequest, error) {
	searchRequest := openapi.NewSearchRequest(index)
	query := openapi.NewSearchQuery()
	boolFilter := openapi.NewBoolFilter()

	// 設定選項
	options := map[string]interface{}{
		"field_weights": map[string]int{
			"name":               6,
			"rewilding_mode":     4,
			"rewilding_location": 4,
			"rewilding_name":     3,
			"tags":               2,
			"host_message":       1,
		},
		"scroll": true,
	}
	searchRequest.SetOptions(options)

	// 處理各種過濾條件
	must := make([]openapi.QueryFilter, 0)
	// should := make([]*openapi.QueryFilter, 0)

	for key, value := range filters {
		// 檢查欄位定義是否存在
		fieldDef, exists := f.fieldDefinitions[key]
		if !exists {
			return nil, fmt.Errorf("未知的欄位: %s", key)
		}

		// 根據欄位類型處理查詢條件
		switch fieldDef.Type {
		case FullText:
			if filter, err := f.handleFullText(value); err == nil {
				must = append(must, *filter)
			}

		case TextMatch:
			if filter, err := f.handleTextMatch(key, value); err == nil {
				must = append(must, *filter)
			}

		case Attribute:
			if filter, err := f.handleAttribute(key, value); err == nil {
				must = append(must, *filter)
			}

		case RangeField:
			if filter, err := f.handleRange(key, value); err == nil {
				// TODO: 如果有設定must，should會被忽略，所以一樣要放在must的地方
				must = append(must, *filter)
			}
		}
	}

	// 設置布林查詢條件
	boolFilter.SetMust(must)
	// if len(should) > 0 {
	// 	boolFilter.SetShould(should)
	// }

	query.SetBool(*boolFilter)
	searchRequest.SetQuery(*query)

	// 設置排序
	sort := map[string]string{
		"id": "asc",
	}
	searchRequest.SetSort(sort)

	return searchRequest, nil
}

// handleFullText 處理全文搜索
func (f *QueryFactory) handleFullText(value interface{}) (*openapi.QueryFilter, error) {
	if keyword, ok := value.(string); ok {
		return &openapi.QueryFilter{
			Match: map[string]interface{}{
				"*": map[string]interface{}{
					"query":    keyword,
					"operator": "and",
				},
			},
		}, nil
	}
	return nil, fmt.Errorf("全文搜索需要字串值")
}

// handleTextMatch 處理文本匹配
func (f *QueryFactory) handleTextMatch(key string, value interface{}) (*openapi.QueryFilter, error) {
	if valueStr, ok := value.(string); ok {
		// 解析包含邏輯運算符的查詢
		if strings.Contains(valueStr, "&") || strings.Contains(valueStr, "|") {
			// 處理 AND 運算
			if strings.Contains(valueStr, "&") {
				andParts := strings.Split(valueStr, "&")
				andFilters := make([]openapi.QueryFilter, 0)

				for _, part := range andParts {
					// 清理括號
					part = strings.Trim(part, "()")
					// 處理 OR 運算
					if strings.Contains(part, "|") {
						orValues := strings.Split(part, "|")
						orFilters := make([]*openapi.QueryFilter, 0)
						for _, v := range orValues {
							orFilters = append(orFilters, &openapi.QueryFilter{
								Match: map[string]interface{}{
									key: map[string]interface{}{
										"query":    v,
										"operator": "and",
									},
								},
							})
						}
						if len(orFilters) > 0 {
							andFilters = append(andFilters, openapi.QueryFilter{
								Bool: &openapi.BoolFilter{
									Should: orFilters,
								},
							})
						}
					} else {
						andFilters = append(andFilters, openapi.QueryFilter{
							Match: map[string]interface{}{
								key: map[string]interface{}{
									"query":    part,
									"operator": "and",
								},
							},
						})
					}
				}

				if len(andFilters) > 0 {
					return &openapi.QueryFilter{
						Bool: &openapi.BoolFilter{
							Must: andFilters,
						},
					}, nil
				}
			} else {
				// 只有 OR 運算
				orValues := strings.Split(valueStr, "|")
				orFilters := make([]*openapi.QueryFilter, 0)
				for _, v := range orValues {
					orFilters = append(orFilters, &openapi.QueryFilter{
						Match: map[string]interface{}{
							key: map[string]interface{}{
								"query":    v,
								"operator": "and",
							},
						},
					})
				}
				if len(orFilters) > 0 {
					return &openapi.QueryFilter{
						Bool: &openapi.BoolFilter{
							Should: orFilters,
						},
					}, nil
				}
			}
		} else {
			// 單一值
			return &openapi.QueryFilter{
				Match: map[string]interface{}{
					key: map[string]interface{}{
						"query":    valueStr,
						"operator": "and",
					},
				},
			}, nil
		}
	}
	return nil, fmt.Errorf("文本匹配需要字串值")
}

// handleAttribute 處理屬性查詢
func (f *QueryFactory) handleAttribute(key string, value interface{}) (*openapi.QueryFilter, error) {
	if valueStr, ok := value.(string); ok {
		// 解析包含邏輯運算符的查詢
		if strings.Contains(valueStr, "&") || strings.Contains(valueStr, "|") {
			// 處理 AND 運算
			if strings.Contains(valueStr, "&") {
				andParts := strings.Split(valueStr, "&")
				andFilters := make([]openapi.QueryFilter, 0)

				for _, part := range andParts {
					// 清理括號
					part = strings.Trim(part, "()")
					// 處理 OR 運算
					if strings.Contains(part, "|") {
						orValues := strings.Split(part, "|")
						andFilters = append(andFilters, openapi.QueryFilter{
							In: map[string]interface{}{
								key: orValues,
							},
						})
					} else {
						andFilters = append(andFilters, openapi.QueryFilter{
							Equals: map[string]interface{}{
								key: part,
							},
						})
					}
				}

				if len(andFilters) > 0 {
					return &openapi.QueryFilter{
						Bool: &openapi.BoolFilter{
							Must: andFilters,
						},
					}, nil
				}
			} else {
				// 只有 OR 運算
				orValues := strings.Split(valueStr, "|")
				return &openapi.QueryFilter{
					In: map[string]interface{}{
						key: orValues,
					},
				}, nil
			}
		} else if strings.Contains(valueStr, ",") {
			// 多選
			orValues := strings.Split(valueStr, ",")
			return &openapi.QueryFilter{
				In: map[string]interface{}{
					key: orValues,
				},
			}, nil
		} else {
			// 單選
			return &openapi.QueryFilter{
				Equals: map[string]interface{}{
					key: valueStr,
				},
			}, nil
		}
	}
	return nil, fmt.Errorf("屬性查詢需要字串值")
}

// handleRange 處理範圍查詢
func (f *QueryFactory) handleRange(key string, value interface{}) (*openapi.QueryFilter, error) {
	if valueStr, ok := value.(string); ok {
		// 解析包含邏輯運算符的查詢
		if strings.Contains(valueStr, "&") || strings.Contains(valueStr, "|") {
			// 處理 AND 運算
			if strings.Contains(valueStr, "&") {
				andParts := strings.Split(valueStr, "&")
				rangeFilters := make([]openapi.QueryFilter, 0)

				for _, part := range andParts {
					// 清理括號
					part = strings.Trim(part, "()")
					// 處理 OR 運算
					if strings.Contains(part, "|") {
						orRanges := strings.Split(part, "|")
						orFilters := make([]*openapi.QueryFilter, 0)
						for _, r := range orRanges {
							if filter, err := f.createRangeFilter(key, r); err == nil {
								orFilters = append(orFilters, filter)
							}
						}
						if len(orFilters) > 0 {
							rangeFilters = append(rangeFilters, openapi.QueryFilter{
								Bool: &openapi.BoolFilter{
									Should: orFilters,
								},
							})
						}
					} else {
						if filter, err := f.createRangeFilter(key, part); err == nil {
							rangeFilters = append(rangeFilters, openapi.QueryFilter{
								Bool: &openapi.BoolFilter{
									Should: []*openapi.QueryFilter{filter},
								},
							})
						}
					}
				}

				if len(rangeFilters) > 0 {
					return &openapi.QueryFilter{
						Bool: &openapi.BoolFilter{
							Must: rangeFilters,
						},
					}, nil
				}
			} else {
				// 只有 OR 運算
				orRanges := strings.Split(valueStr, "|")
				rangeFilters := make([]*openapi.QueryFilter, 0)
				for _, r := range orRanges {
					if filter, err := f.createRangeFilter(key, r); err == nil {
						rangeFilters = append(rangeFilters, filter)
					}
				}
				if len(rangeFilters) > 0 {
					return &openapi.QueryFilter{
						Bool: &openapi.BoolFilter{
							Should: rangeFilters,
						},
					}, nil
				}
			}
		} else {
			// 單一範圍
			if filter, err := f.createRangeFilter(key, valueStr); err == nil {
				return &openapi.QueryFilter{
					Bool: &openapi.BoolFilter{
						Should: []*openapi.QueryFilter{filter},
					},
				}, nil
			}
		}
	}
	return nil, fmt.Errorf("範圍查詢需要字串值")
}

// createRangeFilter 創建範圍過濾器
func (f *QueryFactory) createRangeFilter(field, rangeStr string) (*openapi.QueryFilter, error) {
	// 處理 [3,6] 格式
	if strings.HasPrefix(rangeStr, "[") && strings.HasSuffix(rangeStr, "]") {
		rangeStr = strings.Trim(rangeStr, "[]")
		parts := strings.Split(rangeStr, ",")
		if len(parts) == 2 {
			return &openapi.QueryFilter{
				Range: map[string]interface{}{
					field: map[string]interface{}{
						"gte": parseFloat(parts[0]),
						"lte": parseFloat(parts[1]),
					},
				},
			}, nil
		}
	}

	// 處理 <3 格式
	if strings.HasPrefix(rangeStr, "<") {
		value := strings.TrimPrefix(rangeStr, "<")
		return &openapi.QueryFilter{
			Range: map[string]interface{}{
				field: map[string]interface{}{
					"lt": parseFloat(value),
				},
			},
		}, nil
	}

	// 處理 >3 格式
	if strings.HasPrefix(rangeStr, ">") {
		value := strings.TrimPrefix(rangeStr, ">")
		return &openapi.QueryFilter{
			Range: map[string]interface{}{
				field: map[string]interface{}{
					"gt": parseFloat(value),
				},
			},
		}, nil
	}

	return nil, fmt.Errorf("不支援的範圍格式: %s", rangeStr)
}

// parseFloat 將字串轉換為浮點數
func parseFloat(s string) float64 {
	var result float64
	fmt.Sscanf(s, "%f", &result)
	return result
}
