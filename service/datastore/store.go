// Package datastore exposes an extremely simple key-val datastore.
package datastore

// Store implements Datastore interface
type Store struct {
	Data map[string]interface{}
}

func (s *Store) Get(key string) interface{} {
	return s.Data[key]
}

func (s *Store) Set(key string, val interface{}) error {
	s.Data[key] = val
	return nil
}

func InitStore() *Store {
	return &Store{
		Data: map[string]interface{}{},
	}
}
