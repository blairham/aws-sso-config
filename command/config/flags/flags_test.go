package flags

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListFlag(t *testing.T) {
	flag := NewListFlag()

	// Test flag properties
	assert.Equal(t, "list", flag.GetFlagName())
	assert.Equal(t, "l", flag.GetShortFlag())
	assert.Equal(t, "List all configuration variables and their values", flag.GetDescription())
	assert.Equal(t, "List all configuration variables and their values", flag.GetUsage())
}

func TestFlagRegistry(t *testing.T) {
	registry := NewFlagRegistry()

	// Test GetAllFlags
	flags := registry.GetAllFlags()
	assert.Equal(t, 1, len(flags), "Expected exactly one flag in registry")

	// Test GetFlagByName
	listFlag := registry.GetFlagByName("list")
	assert.NotNil(t, listFlag, "Expected to find list flag")
	assert.Equal(t, "list", listFlag.GetFlagName())

	// Test non-existent flag
	nonExistent := registry.GetFlagByName("nonexistent")
	assert.Nil(t, nonExistent, "Expected nil for non-existent flag")
}
