package manticore

import (
	"testing"

	"github.com/spf13/viper"
)

func TestNewManticore(t *testing.T) {
	// Test case: Successful creation of Manticore client
	viper.Set("manticore.url", "http://localhost:9308/")

	client, err := NewManticore()
	if err != nil {
		t.Errorf("Expected nil error, got %v", err)
	}
	if client == nil {
		t.Errorf("Expected non-nil client, got nil")
	}

	// Test case: Failed creation of Manticore client due to empty URL
	viper.Set("manticore.url", "")
	_, err = NewManticore()
	if err == nil {
		t.Errorf("Expected non-nil error, got nil")
	}
}
