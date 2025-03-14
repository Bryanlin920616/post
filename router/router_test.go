package router

import (
	"fmt"
	"testing"

	"github.com/94peter/microservice/apitool"
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
	// TODO: Write test
}
func TestUpdateIdea(t *testing.T) {
	// TODO: Write test
}
func TestDeleteIdea(t *testing.T) {
	// TODO: Write test
}
func TestAutocomplete(t *testing.T) {
	// TODO: Write test
}
