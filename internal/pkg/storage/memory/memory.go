package memory

import (
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/seggga/observability/internal/pkg/model"
	"golang.org/x/crypto/bcrypt"
)

type mem struct {
	linkStor  map[string]model.Link
	linkMutex sync.RWMutex

	userStore map[string]model.User
	userMutex sync.RWMutex
}

// New creates a new in-memory repository (a map with predefined data)
func New() *mem {
	// users
	userStore := make(map[string]model.User, 2)
	pass1, err := bcrypt.GenerateFromPassword([]byte("pass1"), 0)
	if err != nil {
		log.Fatal("error creating initial users")
	}
	userStore["user1"] = model.User{
		ID:       uuid.New(),
		Name:     "user1",
		PassHash: string(pass1),
	}

	pass2, err := bcrypt.GenerateFromPassword([]byte("pass2"), 0)
	if err != nil {
		log.Fatal("error creating initial users")
	}
	userStore["user2"] = model.User{
		ID:       uuid.New(),
		Name:     "user2",
		PassHash: string(pass2),
	}

	// links
	linkStor := make(map[string]model.Link, 2)
	// predefine link #1
	linkStor["asdf"] = model.Link{
		Short:       "asdf",
		Long:        "https://google.com",
		Count:       0,
		Description: "a very useful link",
		Owner:       userStore["user1"].ID,
	}
	// predefine link #2
	linkStor["qwerty"] = model.Link{
		Short:       "qwerty",
		Long:        "https://yandex.ru",
		Count:       0,
		Description: "one more useful link",
		Owner:       userStore["user2"].ID,
	}

	return &mem{
		linkStor:  linkStor,
		userStore: userStore,
	}
}

// CreateLink adds new link data to the map
func (m *mem) CreateLink(link *model.Link) error {
	m.linkMutex.Lock()
	m.linkStor[link.Short] = *link
	m.linkMutex.Unlock()

	return nil
}

// GetLink retrieves specified link from the storage
func (m *mem) GetLink(short string) (*model.Link, error) {
	m.linkMutex.RLock()
	link := m.linkStor[short]
	m.linkMutex.RUnlock()

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
	m.linkMutex.RLock()
	_, ok := m.linkStor[short]
	m.linkMutex.RUnlock()

	return ok
}

// DeleteLink deletes link data from the map
func (m *mem) DeleteLink(short string) error {
	m.linkMutex.Lock()
	delete(m.linkStor, short)
	m.linkMutex.Unlock()

	return nil
}
