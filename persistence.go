package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
)

type Storage struct {
	File *os.File
	Data map[string]string
	mu   sync.RWMutex
}

func NewStorage(filePath string) (*Storage, error) {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	// Initialize storage
	storage := &Storage{
		File: file,
		Data: make(map[string]string),
	}

	// Load existing data
	if err := storage.loadExistingData(); err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *Storage) loadExistingData() error {
	info, err := s.File.Stat()
	if err != nil {
		return err
	}

	if info.Size() == 0 {
		return nil
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	scanner := bufio.NewScanner(s.File)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) != 2 {
			return errors.New("invalid data format")
		}
		s.Data[parts[0]] = parts[1]
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func (s *Storage) SaveEntry(key, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry := fmt.Sprintf("%s\t%s\n", key, value)
	if _, err := s.File.WriteString(entry); err != nil {
		return err
	}

	s.Data[key] = value
	return nil
}

func (s *Storage) Get(key string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, exists := s.Data[key]
	if !exists {
		return ""
	}
	return value
}

func (s *Storage) Close() error {
	return s.File.Close()
}
