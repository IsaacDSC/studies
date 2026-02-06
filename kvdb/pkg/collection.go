package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"sync"
)

type Collections map[string]*Collection

type Collection struct {
	mu         sync.RWMutex
	FilterKeys FilterKeys
	uniqueData SingleValue
	multiData  MultipleValues
}

func NewCollection() *Collection {
	return &Collection{
		FilterKeys: FilterKeys{},
		uniqueData: make(SingleValue),
		multiData:  make(MultipleValues),
	}
}

var ErrMaxFilterKeysReached = errors.New("max filter keys reached (max is 5)")
var ErrFilterKeyAlreadyExists = errors.New("filter key already exists")

func (c *Collection) CreateFilterKey(filterKey string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.FilterKeys) >= 5 {
		return ErrMaxFilterKeysReached
	}

	if slices.Contains(c.FilterKeys, filterKey) {
		return ErrFilterKeyAlreadyExists
	}

	c.FilterKeys = append(c.FilterKeys, filterKey)

	return nil
}

var ErrCollectionNotFound = errors.New("collection not found")

func (c *Collection) Set(key string, value any) error {
	valueStr, err := ToStr(value)
	if err != nil {
		return fmt.Errorf("failed to convert value to string: %w", err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	if !IsPrimitive(value) {
		var data map[string]any
		if err := json.Unmarshal([]byte(valueStr), &data); err != nil {
			return fmt.Errorf("failed to unmarshal data: %w", err)
		}

		for k, v := range data {
			if c.FilterKeys.Contains(k) {
				keypair, err := NewKeyPair(NewFilterKey(k), v)
				if err != nil {
					return fmt.Errorf("failed to create key pair: %w", err)
				}
				c.multiData.Set(keypair, Key(key))
			}
		}
	}

	c.uniqueData.Set(Key(key), valueStr)

	return nil
}

func (c *Collection) FindOne(key Key) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.uniqueData.Get(key)
}

func (c *Collection) FindAll(filterKey FilterKey, filterValue string) ([]string, error) {
	keypair, err := NewKeyPair(filterKey, filterValue)
	if err != nil {
		return nil, fmt.Errorf("failed to create key pair: %w", err)
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	keys, err := c.multiData.Get(keypair)
	if err != nil {
		return nil, fmt.Errorf("failed to get keys: %w", err)
	}

	var data []string
	for _, key := range keys {
		d, _ := c.uniqueData.Get(key)
		data = append(data, d)
	}

	return data, nil
}

func (c *Collection) Scan() ([]string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var data []string
	for _, value := range c.uniqueData {
		data = append(data, value)
	}

	return data, nil
}
