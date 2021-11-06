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

// CreateLink adds data to the map
func (m *mem) CreateLink(link *model.Link) error {
	m.storage[link.Short] = *link
	return nil
}

// GetLink retrieves specified link from the storage
func (m *mem) GetLink(short string) (*model.Link, error) {
	link := m.storage[short]
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
	_, ok := m.storage[short]
	if ok {
		return true
	}
	return false
}
