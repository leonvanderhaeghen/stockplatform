package application

import (
	"fmt"
	"sync"

	"github.com/leonvanderhaeghen/stockplatform/services/supplierSvc/internal/domain"
)

// adapterRegistry is a registry of supplier adapters
type adapterRegistry struct {
	adapters map[string]domain.SupplierAdapter
	mu       sync.RWMutex
}

// NewAdapterRegistry creates a new adapter registry
func NewAdapterRegistry() domain.AdapterRegistry {
	return &adapterRegistry{
		adapters: make(map[string]domain.SupplierAdapter),
	}
}

// Register registers a new supplier adapter
func (r *adapterRegistry) Register(adapter domain.SupplierAdapter) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := adapter.Name()
	if name == "" {
		return fmt.Errorf("adapter name cannot be empty")
	}

	if _, exists := r.adapters[name]; exists {
		return fmt.Errorf("adapter with name %s already registered", name)
	}

	r.adapters[name] = adapter
	return nil
}

// Get returns an adapter by name
func (r *adapterRegistry) Get(name string) (domain.SupplierAdapter, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	adapter, exists := r.adapters[name]
	if !exists {
		return nil, fmt.Errorf("adapter %s not found", name)
	}
	return adapter, nil
}

// List returns all registered adapters
func (r *adapterRegistry) List() []domain.SupplierAdapter {
	r.mu.RLock()
	defer r.mu.RUnlock()

	adapters := make([]domain.SupplierAdapter, 0, len(r.adapters))
	for _, adapter := range r.adapters {
		adapters = append(adapters, adapter)
	}
	return adapters
}
