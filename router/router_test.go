package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/94peter/microservice/apitool"
	apiErr "github.com/94peter/microservice/apitool/err"
	"github.com/arwoosa/post/router/request"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetApis(t *testing.T) {
	tests := []struct {
		name string
		want []apitool.GinAPI
	}{
		{
			name: "test GetApis",
			want: []apitool.GinAPI{
				&idea{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetApis()
			if len(got) != len(tt.want) {
				t.Errorf("GetApis() = %v, want %v", got, tt.want)
			}
			for i, api := range got {
				if fmt.Sprintf("%T", api) != fmt.Sprintf("%T", tt.want[i]) {
					t.Errorf("GetApis() = %v, want %v", api, tt.want[i])
				}
			}
		})
	}
}
func TestIdeaGetHandlers(t *testing.T) {
	m := &idea{}

	handlers := m.GetHandlers()

	// Test that the function returns five handlers
	if len(handlers) != 5 {
		t.Errorf("expected 5 handlers, got %d", len(handlers))
	}

	// Test that the first handler has the correct path and method
	if handlers[0].Path != "/idea" || handlers[0].Method != "GET" {
		t.Errorf("expected first handler to have path '/idea' and method 'GET', got path '%s' and method '%s'", handlers[0].Path, handlers[0].Method)
	}

	// Test that the second handler has the correct path and method
	if handlers[1].Path != "/idea" || handlers[1].Method != "POST" {
		t.Errorf("expected second handler to have path '/idea' and method 'POST', got path '%s' and method '%s'", handlers[1].Path, handlers[1].Method)
	}

	// Test that the third handler has the correct path and method
	if handlers[2].Path != "/idea/:id" || handlers[2].Method != "PUT" {
		t.Errorf("expected third handler to have path '/idea/:id' and method 'PUT', got path '%s' and method '%s'", handlers[2].Path, handlers[2].Method)
	}

	// Test that the forth handler has the correct path and method
	if handlers[3].Path != "/idea/:id" || handlers[3].Method != "DELETE" {
		t.Errorf("expected forth handler to have path '/idea/:id' and method 'DELETE', got path '%s' and method '%s'", handlers[3].Path, handlers[3].Method)
	}

	// Test that the fifth handler has the correct path and method
	if handlers[4].Path != "/idea/autocomplete" || handlers[4].Method != "GET" {
		t.Errorf("expected fifth handler to have path '/idea/autocomplete' and method 'GET', got path '%s' and method '%s'", handlers[4].Path, handlers[4].Method)
	}
}
func TestCreateIdea(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name            string
		requestBody     *request.CreateIdea
		mockInsertFunc  func(query string) error
		mockInsertError error
		statusCode      int
		response        *request.CreateIdeaResponse
	}{
		{
			name:           "bind error",
			requestBody:    nil,
			mockInsertFunc: nil,
			statusCode:     http.StatusBadRequest,
		},
		{
			name: "invalid request body",
			requestBody: &request.CreateIdea{
				BaseIdea: request.BaseIdea{
					MongoId: -1,
				},
			},
			mockInsertFunc: nil,
			statusCode:     http.StatusBadRequest,
		},
		{
			name: "valid request",
			requestBody: &request.CreateIdea{
				BaseIdea: request.BaseIdea{
					MongoId:            1,
					ItineraryName:      "自行車地獄之旅",
					AttractionName:     "河濱公園",
					Tags:               []string{"新手", "情侶", "一日遊", "文化體驗", "期間限定"},
					WildMode:           "露營",
					AttractionLocation: "台北, 台灣",
					HostMessage:        "這是一個很棒的行程，歡迎參加",
					ExperienceDuration: 4.0,
				},
			},
			mockInsertFunc: func(query string) error {
				// 模擬 SQL 插入成功
				return nil
			},
			statusCode: http.StatusOK,
			response: &request.CreateIdeaResponse{
				BaseIdeaResponse: request.BaseIdeaResponse{
					BaseIdea: request.BaseIdea{
						MongoId:            1,
						ItineraryName:      "自行車地獄之旅",
						AttractionName:     "河濱公園",
						Tags:               []string{"新手", "情侶", "一日遊", "文化體驗", "期間限定"},
						WildMode:           "露營",
						AttractionLocation: "台北, 台灣",
						HostMessage:        "這是一個很棒的行程，歡迎參加",
						ExperienceDuration: 4.0,
					},
				},
			},
		},
		{
			name: "server error",
			requestBody: &request.CreateIdea{
				BaseIdea: request.BaseIdea{
					MongoId:            1,
					ItineraryName:      "自行車地獄之旅",
					AttractionName:     "河濱公園",
					Tags:               []string{"新手", "情侶", "一日遊", "文化體驗", "期間限定"},
					WildMode:           "露營",
					AttractionLocation: "台北, 台灣",
					HostMessage:        "這是一個很棒的行程，歡迎參加",
					ExperienceDuration: 4.0,
				},
			},
			mockInsertFunc: func(query string) error {
				// 模擬 SQL 插入失敗
				return nil
			},
			statusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		// defer func ()  {
		// 	// reset mock
		// }
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			var requestData *bytes.Buffer
			if test.requestBody != nil {
				data, _ := json.Marshal(test.requestBody)
				requestData = bytes.NewBuffer(data)
			} else {
				requestData = bytes.NewBuffer([]byte{})
			}

			c.Request, _ = http.NewRequest("POST", "/ideas", requestData)
			c.Request.Header.Set("Content-Type", "application/json")

			/* setup mock

			mockDatabaseInsert := func(query string) error {
				if test.mockInsertFunc != nil {
					return test.mockInsertFunc(query)
				}
				return test.mockInsertError
			}

			*/

			idea := &idea{}
			idea.SetErrorHandler(func(c *gin.Context, err error) {
				if apiErr, ok := err.(apiErr.ApiError); ok {
					c.JSON(apiErr.GetStatus(), gin.H{
						"error": apiErr.Error(),
					})
				}
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
			})
			idea.createIdea(c)

			assert.Equal(t, test.statusCode, w.Code)
			if test.response != nil {
				assert.Equal(t, test.response, w.Result())
			}
		})
	}
}
func TestUpdateIdea(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name            string
		mongoId         string
		requestBody     *request.UpdateIdea
		mockUpdateFunc  func(query string) error
		mockUpdateError error
		statusCode      int
		response        *request.UpdateIdeaResponse
	}{
		{
			name:        "bind error",
			mongoId:     "1",
			requestBody: nil,
			statusCode:  http.StatusBadRequest,
		},
		{
			name:    "invalid request body",
			mongoId: "1",
			requestBody: &request.UpdateIdea{
				BaseIdea: request.BaseIdea{
					MongoId: -1,
				},
			},
			statusCode: http.StatusBadRequest,
		},
		{
			name:    "valid request",
			mongoId: "1",
			requestBody: &request.UpdateIdea{
				BaseIdea: request.BaseIdea{
					MongoId:            1,
					ItineraryName:      "更新後的行程名稱",
					AttractionName:     "更新後的景點名稱",
					Tags:               []string{"更新標籤一", "更新標籤二"},
					WildMode:           "更新後的野放模式",
					AttractionLocation: "更新後的景點縣市",
					HostMessage:        "這是團主想說的話",
					ExperienceDuration: 3.5,
				},
			},
			mockUpdateFunc: func(query string) error {
				// 模擬 SQL 更新成功
				return nil
			},
			statusCode: http.StatusOK,
			response: &request.UpdateIdeaResponse{
				BaseIdeaResponse: request.BaseIdeaResponse{
					BaseIdea: request.BaseIdea{
						MongoId:            1,
						ItineraryName:      "更新後的行程名稱",
						AttractionName:     "更新後的景點名稱",
						Tags:               []string{"更新標籤一", "更新標籤二"},
						WildMode:           "更新後的野放模式",
						AttractionLocation: "更新後的景點縣市",
						HostMessage:        "這是團主想說的話",
						ExperienceDuration: 3.5,
					},
				},
			},
		},
		{
			name:    "not found error",
			mongoId: "999",
			requestBody: &request.UpdateIdea{
				BaseIdea: request.BaseIdea{
					MongoId:            999,
					ItineraryName:      "更新後的行程名稱",
					AttractionName:     "更新後的景點名稱",
					Tags:               []string{"更新標籤一", "更新標籤二"},
					WildMode:           "更新後的野放模式",
					AttractionLocation: "更新後的景點縣市",
					HostMessage:        "這是團主想說的話",
					ExperienceDuration: 3.5,
				},
			},
			mockUpdateFunc: func(query string) error {
				// 模擬找不到資料
				return fmt.Errorf("idea not found")
			},
			statusCode: http.StatusNotFound,
		},
		{
			name:    "server error",
			mongoId: "1",
			requestBody: &request.UpdateIdea{
				BaseIdea: request.BaseIdea{
					MongoId:            1,
					ItineraryName:      "更新後的行程名稱",
					AttractionName:     "更新後的景點名稱",
					Tags:               []string{"更新標籤一", "更新標籤二"},
					WildMode:           "更新後的野放模式",
					AttractionLocation: "更新後的景點縣市",
					HostMessage:        "這是團主想說的話",
					ExperienceDuration: 3.5,
				},
			},
			mockUpdateFunc: func(query string) error {
				// 模擬伺服器錯誤
				return fmt.Errorf("internal server error")
			},
			statusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		// defer func ()  {
		// 	// reset mock
		// }
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			var requestData *bytes.Buffer
			if test.requestBody != nil {
				data, _ := json.Marshal(test.requestBody)
				requestData = bytes.NewBuffer(data)
			} else {
				requestData = bytes.NewBuffer([]byte{})
			}

			c.Request, _ = http.NewRequest("PUT", "/idea/"+test.mongoId, requestData)
			c.Request.Header.Set("Content-Type", "application/json")
			c.Params = []gin.Param{{Key: "id", Value: test.mongoId}}

			/* setup mock

			mockDatabaseUpdate := func(query string) error {
				if test.mockUpdateFunc != nil {
					return test.mockUpdateFunc(query)
				}
				return test.mockUpdateError
			}

			*/

			idea := &idea{}
			idea.SetErrorHandler(func(c *gin.Context, err error) {
				if apiErr, ok := err.(apiErr.ApiError); ok {
					c.JSON(apiErr.GetStatus(), gin.H{
						"error": apiErr.Error(),
					})
				}
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
			})
			idea.updateIdea(c)

			assert.Equal(t, test.statusCode, w.Code)
			if test.response != nil {
				assert.Equal(t, test.response, w.Result())
			}
		})
	}
}
func TestDeleteIdea(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name            string
		mongoId         string
		mockDeleteFunc  func(query string) error
		mockDeleteError error
		statusCode      int
	}{
		{
			name:    "valid request",
			mongoId: "1",
			mockDeleteFunc: func(query string) error {
				// 模擬 SQL 刪除成功
				return nil
			},
			statusCode: http.StatusNoContent,
		},
		{
			name:    "not found error",
			mongoId: "999",
			mockDeleteFunc: func(query string) error {
				// 模擬找不到資料
				return fmt.Errorf("idea not found")
			},
			statusCode: http.StatusNotFound,
		},
		{
			name:    "server error",
			mongoId: "1",
			mockDeleteFunc: func(query string) error {
				// 模擬伺服器錯誤
				return fmt.Errorf("internal server error")
			},
			statusCode: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			c.Request, _ = http.NewRequest("DELETE", "/idea/"+test.mongoId, nil)
			c.Params = []gin.Param{{Key: "id", Value: test.mongoId}}

			/* setup mock

			mockDatabaseDelete := func(query string) error {
				if test.mockDeleteFunc != nil {
					return test.mockDeleteFunc(query)
				}
				return test.mockDeleteError
			}

			*/

			idea := &idea{}
			idea.SetErrorHandler(func(c *gin.Context, err error) {
				if apiErr, ok := err.(apiErr.ApiError); ok {
					c.JSON(apiErr.GetStatus(), gin.H{
						"error": apiErr.Error(),
					})
				}
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
			})
			idea.deleteIdea(c)

			assert.Equal(t, test.statusCode, w.Code)
		})
	}
}
func TestAutocomplete(t *testing.T) {
	// TODO: Write test
}
