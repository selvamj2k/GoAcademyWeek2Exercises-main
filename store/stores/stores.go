package stores

import (
	"fmt"
	pkgloggers "loggers"
	"time"
)

type Store struct {
	Key    int
	Value  string
	User   string
	Writes int
	Reads  int
	Age    time.Time
}

func (s *Store) String() string {
	return fmt.Sprintf("Store key %v, value %v, user %v", s.Key, s.Value, s.User)
}

func NewStore(key int, value string, owner string, writes int, read int, age time.Time) *Store {
	store := Store{key, value, owner, writes, read, age}
	return &store
}

var Stores []*Store

func FindStore(Key int) (*Store, bool) {
	for i := 0; i <= len(Stores)-1; i++ {
		store := Stores[i]
		if store.Key == Key {
			pkgloggers.StoreLogger.Printf("Store found for the passed in key %v. update store", Key)
			return store, true
		}
	}
	pkgloggers.StoreLogger.Printf("Store not found for the passed in key %v. Add new", Key)
	return nil, false
}
