package gnoweb

import (
	"testing"
	"time"

	"github.com/gnolang/gno/gno.land/pkg/gnoweb/markdown"
	"github.com/stretchr/testify/assert"
)

func TestMDPackStorageTTL(t *testing.T) {
	// Create storage with very short TTL for testing
	storage := NewMDPackStorageWithTTL(100*time.Millisecond, 50*time.Millisecond)
	defer storage.Stop()
	
	// Add a file
	file := &markdown.MDPackFile{
		Name:     "test.txt",
		Content:  []byte("test content"),
		MimeType: "text/plain",
	}
	storage.AddFile("/test", file)
	
	// Should be retrievable immediately
	retrieved, ok := storage.GetFile("/test", "test.txt")
	assert.True(t, ok)
	assert.Equal(t, "test content", string(retrieved.Content))
	
	// Wait for TTL to expire
	time.Sleep(150 * time.Millisecond)
	
	// Should no longer be retrievable
	_, ok = storage.GetFile("/test", "test.txt")
	assert.False(t, ok, "file should have expired")
}

func TestMDPackStorageCleanup(t *testing.T) {
	// Create storage with short TTL and cleanup interval
	storage := NewMDPackStorageWithTTL(100*time.Millisecond, 50*time.Millisecond)
	defer storage.Stop()
	
	// Add multiple files
	for i := 0; i < 5; i++ {
		file := &markdown.MDPackFile{
			Name:    string(rune('a' + i)) + ".txt",
			Content: []byte("content"),
		}
		storage.AddFile("/test", file)
	}
	
	// Verify all files are present
	storage.mu.RLock()
	initialCount := len(storage.files["/test"])
	storage.mu.RUnlock()
	assert.Equal(t, 5, initialCount)
	
	// Wait for cleanup to run
	time.Sleep(200 * time.Millisecond)
	
	// Verify files have been cleaned up
	storage.mu.RLock()
	finalCount := len(storage.files["/test"])
	storage.mu.RUnlock()
	assert.Equal(t, 0, finalCount, "all files should have been cleaned up")
}