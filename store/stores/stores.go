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

type StoreInfo struct {
	Key  string `json:"key"`
	User string `json:"owner"`
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
			pkgloggers.StoreLogger.Printf("Store found for the passed in key %v.", Key)
			return store, true
		}
	}
	pkgloggers.StoreLogger.Printf("Store not found for the passed in key %v.", Key)
	return nil, false
}

func indexOf(key int) (int, bool) {
	for i := 0; i <= len(Stores)-1; i++ {
		if key == Stores[i].Key {
			pkgloggers.StoreLogger.Printf("key %v and index %v", key, i)
			return i, true
		}
	}

	pkgloggers.StoreLogger.Printf("Store item not found for key %v", key)
	return -1, false
}

func RemoveIndex(index int) []*Store {
	pkgloggers.StoreLogger.Printf("Remove item from slice")

	newStores := append(Stores[:index], Stores[index+1:]...)
	return newStores

}
