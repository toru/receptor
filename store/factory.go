package store

import (
	"errors"

	"github.com/toru/dexter/subscription"
)

// Store is an interface that Storage Engines must implement.
type Store interface {
	Name() string
	Subscriptions() []subscription.Subscription
	CreateSubscription(sub *subscription.Subscription) error
	NumSubscriptions() int
}

// GetStore returns a Storage Engine based on the given name.
func GetStore(name string) (Store, error) {
	switch name {
	case "memory":
		s, err := NewMemoryStore()
		if err != nil {
			return nil, err
		}
		return s, nil
	case "mysql", "mariadb":
		return nil, errors.New("work in progress")
	default:
		return nil, errors.New("unknown store")
	}
}
