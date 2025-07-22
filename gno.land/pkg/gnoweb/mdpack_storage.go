package gnoweb

import (
	"sync"
	"time"

	"github.com/gnolang/gno/gno.land/pkg/gnoweb/markdown"
)

// MDPackStorage stores extracted files temporarily
type MDPackStorage struct {
	mu       sync.RWMutex
	files    map[string]map[string]*storedFile // path -> filename -> file
	ttl      time.Duration
	cleanupInterval time.Duration
	stopCleanup    chan struct{}
}

type storedFile struct {
	file      *markdown.MDPackFile
	timestamp time.Time
}

// NewMDPackStorage creates a new MDPack storage with configurable TTL
func NewMDPackStorage() *MDPackStorage {
	return NewMDPackStorageWithTTL(30*time.Second, 10*time.Second) // 30 second TTL, cleanup every 10 seconds
}

// NewMDPackStorageWithTTL creates storage with custom TTL and cleanup interval
func NewMDPackStorageWithTTL(ttl, cleanupInterval time.Duration) *MDPackStorage {
	storage := &MDPackStorage{
		files:           make(map[string]map[string]*storedFile),
		ttl:             ttl,
		cleanupInterval: cleanupInterval,
		stopCleanup:     make(chan struct{}),
	}
	
	// Start cleanup goroutine
	go storage.cleanup()
	
	return storage
}

// Stop stops the cleanup goroutine
func (s *MDPackStorage) Stop() {
	close(s.stopCleanup)
}

// AddFile adds a file to storage
func (s *MDPackStorage) AddFile(path string, file *markdown.MDPackFile) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if s.files[path] == nil {
		s.files[path] = make(map[string]*storedFile)
	}
	
	s.files[path][file.Name] = &storedFile{
		file:      file,
		timestamp: time.Now(),
	}
	
	// Debug logging
	// fmt.Printf("MDPackStorage: Added file %s to path %s\n", file.Name, path)
}

// GetFile retrieves a file from storage
func (s *MDPackStorage) GetFile(path, filename string) (*markdown.MDPackFile, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	// Debug logging  
	// fmt.Printf("MDPackStorage: Looking for file %s at path %s\n", filename, path)
	// fmt.Printf("MDPackStorage: Available paths: %v\n", s.getPathsDebug())
	
	if pathFiles, ok := s.files[path]; ok {
		if stored, ok := pathFiles[filename]; ok {
			// Check if file is not too old
			if time.Since(stored.timestamp) < s.ttl {
				// fmt.Printf("MDPackStorage: Found file %s\n", filename)
				return stored.file, true
			}
		}
	}
	
	// fmt.Printf("MDPackStorage: File not found\n")
	return nil, false
}

func (s *MDPackStorage) getPathsDebug() []string {
	var paths []string
	for p := range s.files {
		paths = append(paths, p)
	}
	return paths
}

// cleanup removes old files periodically
func (s *MDPackStorage) cleanup() {
	ticker := time.NewTicker(s.cleanupInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			s.mu.Lock()
			now := time.Now()
			
			// Remove files older than TTL
			for path, pathFiles := range s.files {
				for filename, stored := range pathFiles {
					if now.Sub(stored.timestamp) > s.ttl {
						delete(pathFiles, filename)
					}
				}
				
				// Remove empty path entries
				if len(pathFiles) == 0 {
					delete(s.files, path)
				}
			}
			
			s.mu.Unlock()
		case <-s.stopCleanup:
			return
		}
	}
}

// RequestScopedStorage provides request-scoped file storage
type RequestScopedStorage struct {
	files map[string]*markdown.MDPackFile
	mu    sync.Mutex
}

// NewRequestScopedStorage creates a new request-scoped storage
func NewRequestScopedStorage() *RequestScopedStorage {
	return &RequestScopedStorage{
		files: make(map[string]*markdown.MDPackFile),
	}
}

// AddFile adds a file to the request storage
func (r *RequestScopedStorage) AddFile(path string, file *markdown.MDPackFile) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	// For request-scoped storage, we ignore the path and just use filename
	r.files[file.Name] = file
}

// GetFile retrieves a file from the request storage
func (r *RequestScopedStorage) GetFile(path, filename string) (*markdown.MDPackFile, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	file, ok := r.files[filename]
	return file, ok
}