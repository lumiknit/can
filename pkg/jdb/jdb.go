package jdb

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type Action string

const (
	ActionBASE Action = "BASE"
	ActionSET  Action = "SET"
	ActionDEL  Action = "DEL"
)

type Entry struct {
	Action Action      `json:"action"`
	Key    string      `json:"key"`
	Value  interface{} `json:"value,omitempty"`
}

type JDB struct {
	filename       string
	data           sync.Map
	file           *os.File
	compactRatio   float64
	baseEntryCount int
	totalEntries   int
	mu             sync.Mutex
}

func New(filename string) (*JDB, error) {
	jdb := &JDB{
		filename:     filename,
		compactRatio: 3.0,
	}

	if err := jdb.load(); err != nil {
		return nil, fmt.Errorf("failed to load database: %w", err)
	}

	return jdb, nil
}

func (jdb *JDB) load() error {
	file, err := os.OpenFile(jdb.filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	jdb.file = file

	// Read existing entries
	file.Seek(0, 0)
	scanner := bufio.NewScanner(file)
	entries := make([]Entry, 0)

	for scanner.Scan() {
		var entry Entry
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			continue // Skip malformed entries
		}
		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Process entries
	jdb.processEntries(entries)

	// Check if compaction is needed
	if jdb.needsCompaction() {
		return jdb.compact()
	}

	return nil
}

func (jdb *JDB) processEntries(entries []Entry) {
	for i, entry := range entries {
		switch entry.Action {
		case ActionBASE:
			if i == 0 {
				jdb.baseEntryCount = 1
			}
			if entry.Value != nil {
				if baseData, ok := entry.Value.(map[string]interface{}); ok {
					for k, v := range baseData {
						jdb.data.Store(k, v)
					}
				}
			}
		case ActionSET:
			jdb.data.Store(entry.Key, entry.Value)
		case ActionDEL:
			jdb.data.Delete(entry.Key)
		}
	}
	jdb.totalEntries = len(entries)
}

func (jdb *JDB) Get(key string) (interface{}, bool) {
	return jdb.data.Load(key)
}

func (jdb *JDB) Set(key string, value interface{}) error {
	jdb.mu.Lock()
	defer jdb.mu.Unlock()

	entry := Entry{
		Action: ActionSET,
		Key:    key,
		Value:  value,
	}

	if err := jdb.appendEntry(entry); err != nil {
		return err
	}

	jdb.data.Store(key, value)
	jdb.totalEntries++

	if jdb.needsCompaction() {
		return jdb.compact()
	}

	return nil
}

func (jdb *JDB) Del(key string) error {
	jdb.mu.Lock()
	defer jdb.mu.Unlock()

	entry := Entry{
		Action: ActionDEL,
		Key:    key,
	}

	if err := jdb.appendEntry(entry); err != nil {
		return err
	}

	jdb.data.Delete(key)
	jdb.totalEntries++

	if jdb.needsCompaction() {
		return jdb.compact()
	}

	return nil
}

func (jdb *JDB) appendEntry(entry Entry) error {
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("failed to marshal entry: %w", err)
	}

	if _, err := jdb.file.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("failed to write entry: %w", err)
	}

	return jdb.file.Sync()
}

func (jdb *JDB) needsCompaction() bool {
	if jdb.baseEntryCount == 0 {
		return false
	}
	return float64(jdb.totalEntries) >= float64(jdb.baseEntryCount)*jdb.compactRatio
}

func (jdb *JDB) compact() error {
	// Collect all current data
	allData := make(map[string]interface{})
	jdb.data.Range(func(key, value interface{}) bool {
		if k, ok := key.(string); ok {
			allData[k] = value
		}
		return true
	})

	// Close current file
	jdb.file.Close()

	// Create new file with BASE entry
	file, err := os.Create(jdb.filename)
	if err != nil {
		return fmt.Errorf("failed to create new file: %w", err)
	}
	jdb.file = file

	baseEntry := Entry{
		Action: ActionBASE,
		Key:    "",
		Value:  allData,
	}

	if err := jdb.appendEntry(baseEntry); err != nil {
		return fmt.Errorf("failed to write base entry: %w", err)
	}

	jdb.baseEntryCount = 1
	jdb.totalEntries = 1

	return nil
}

func (jdb *JDB) Close() error {
	if jdb.file != nil {
		return jdb.file.Close()
	}
	return nil
}
