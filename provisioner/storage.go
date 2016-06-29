package main

import (
	"fmt"
	"net/url"
	"strings"
)

type Storage interface {
	Put(id string, update StatusMsg) error
	Get(id string) (*StatusMsg, error)
	List() ([]StatusMsg, error)
}

func NewStorage(spec string) (Storage, error) {
	conn, err := url.Parse(spec)
	if err != nil {
		return nil, err
	}

	switch strings.ToUpper(conn.Scheme) {
	case "MEMORY":
		return NewMemoryStorage(), nil
	case "CONSUL":
		return NewConsulStorage(spec)
	default:
		return nil, fmt.Errorf("Unknown storage scheme specified, '%s'", conn.Scheme)
	}
}

type MemoryStorage struct {
	Data map[string]StatusMsg
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		Data: make(map[string]StatusMsg),
	}
}

func (s *MemoryStorage) Put(id string, update StatusMsg) error {
	s.Data[id] = update
	return nil
}

func (s *MemoryStorage) Get(id string) (*StatusMsg, error) {
	m, ok := s.Data[id]
	if !ok {
		return nil, nil
	}
	return &m, nil
}

func (s *MemoryStorage) List() ([]StatusMsg, error) {
	r := make([]StatusMsg, len(s.Data))
	i := 0
	for _, v := range s.Data {
		r[i] = v
		i += 1
	}
	return r, nil
}
