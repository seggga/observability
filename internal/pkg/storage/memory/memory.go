package memory

import (
	"fmt"
	"sync"

	"github.com/seggga/observability/internal/pkg/model"
)

type mem struct {
	storage map[string]model.Link
	mutex   sync.RWMutex
}

// New creates a new in-memory repository (a map with predefined data)
func New() *mem {
	mapStorage := make(map[string]model.Link)
	// predefine data 1
	mapStorage["asdf"] = model.Link{
		Short:       "asdf",
		Long:        "https://google.com",
		Count:       0,
		Description: "a very useful link",
	}
	// predefine data 2
	mapStorage["qwerty"] = model.Link{
		Short:       "qwerty",
		Long:        "https://yandex.ru",
		Count:       0,
		Description: "one more useful link",
	}

	return &mem{
		storage: mapStorage,
	}
}

// CreateLink adds new link data to the map
func (m *mem) CreateLink(link *model.Link) error {
	m.mutex.Lock()
	m.storage[link.Short] = *link
	m.mutex.Unlock()

	return nil
}

// GetLink retrieves specified link from the storage
func (m *mem) GetLink(short string) (*model.Link, error) {
	m.mutex.RLock()
	link := m.storage[short]
	m.mutex.RUnlock()

	// check for valid data
	if link.Long == "" {
		err := fmt.Errorf("long link is empty")
		return nil, err
	}

	return &link, nil
}

func (m *mem) Close() {
}

// IsSet checks if an element exist in the storage
func (m *mem) IsSet(short string) bool {
	m.mutex.RLock()
	_, ok := m.storage[short]
	m.mutex.RUnlock()

	return ok
}

// DeleteLink deletes link data from the map
func (m *mem) DeleteLink(short string) error {
	m.mutex.Lock()
	delete(m.storage, short)
	m.mutex.Unlock()

	return nil
}
