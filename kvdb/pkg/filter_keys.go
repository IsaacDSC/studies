package pkg

import (
	"encoding/base64"
	"errors"
	"fmt"
	"slices"
)

// Filter keys
type FilterKeys []string

func (fk FilterKeys) Contains(key string) bool {
	return slices.Contains(fk, key)
}

// Multiple values
type FilterKey interface {
	filterKey() // método privado - só nosso tipo implementa
	String() string
}

type filterKeyImpl string

func (f filterKeyImpl) filterKey()     {}
func (f filterKeyImpl) String() string { return string(f) }

func NewFilterKey(s string) FilterKey {
	return filterKeyImpl(s)
}

type KeyPair string // filterKey:filterValue

func NewKeyPair(filterKey FilterKey, filterValue any) (KeyPair, error) {
	value, err := ToStr(filterValue)
	if err != nil {
		return "", fmt.Errorf("failed to marshal filter value: %w", err)
	}

	b64Value := base64.StdEncoding.EncodeToString([]byte(value))
	return KeyPair(fmt.Sprintf("%s:%s", filterKey.String(), b64Value)), nil
}

type MultipleValues map[KeyPair][]Key

func (k MultipleValues) Set(keyPair KeyPair, value Key) {
	k[keyPair] = append(k[keyPair], value)
}

var ErrMultipleValuesNotFound = errors.New("multiple values not found")

func (k MultipleValues) Get(keyPair KeyPair) ([]Key, error) {
	keys, ok := k[keyPair]
	if !ok {
		return nil, ErrMultipleValuesNotFound
	}

	return keys, nil
}

// Unique values
type SingleValue map[Key]string

func (s SingleValue) Set(key Key, value string) {
	s[key] = value
}

var ErrSingleValueNotFound = errors.New("single value not found")

func (s SingleValue) Get(key Key) (string, error) {
	if data, ok := s[key]; ok {
		return data, nil
	}

	return "", ErrSingleValueNotFound
}
