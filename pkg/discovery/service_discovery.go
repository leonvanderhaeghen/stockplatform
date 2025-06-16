package discovery

import (
	"fmt"
	"os"
	"strings"
)

// ServiceRegistry holds service addresses
type ServiceRegistry struct {
	services map[string]string
}

// NewServiceRegistry creates a new service registry
func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{
		services: make(map[string]string),
	}
}

// RegisterService registers a service with its address
func (sr *ServiceRegistry) RegisterService(name, address string) {
	sr.services[name] = address
}

// GetServiceAddress returns the address for a given service
func (sr *ServiceRegistry) GetServiceAddress(serviceName string) (string, error) {
	// First check if it's registered in memory
	if addr, exists := sr.services[serviceName]; exists {
		return addr, nil
	}

	// Check environment variables
	envKey := fmt.Sprintf("%s_ADDR", strings.ToUpper(strings.ReplaceAll(serviceName, "-", "_")))
	if addr := os.Getenv(envKey); addr != "" {
		sr.RegisterService(serviceName, addr)
		return addr, nil
	}

	// Fallback to default patterns
	defaultAddr := getDefaultServiceAddress(serviceName)
	if defaultAddr != "" {
		sr.RegisterService(serviceName, defaultAddr)
		return defaultAddr, nil
	}

	return "", fmt.Errorf("service %s not found", serviceName)
}

// getDefaultServiceAddress returns default addresses for known services
func getDefaultServiceAddress(serviceName string) string {
	defaults := map[string]string{
		"gateway-service":   "gateway-service:8080",
		"product-service":   "product-service:50053",
		"inventory-service": "inventory-service:50054",
		"order-service":     "order-service:50055",
		"user-service":      "user-service:50056",
		"supplier-service":  "supplier-service:50057",
	}

	return defaults[serviceName]
}

// ListServices returns all registered services
func (sr *ServiceRegistry) ListServices() map[string]string {
	result := make(map[string]string)
	for k, v := range sr.services {
		result[k] = v
	}
	return result
}

// DefaultRegistry is a global service registry instance
var DefaultRegistry = NewServiceRegistry()

// GetService is a convenience function to get service address from default registry
func GetService(serviceName string) (string, error) {
	return DefaultRegistry.GetServiceAddress(serviceName)
}

// RegisterService is a convenience function to register service in default registry
func RegisterService(name, address string) {
	DefaultRegistry.RegisterService(name, address)
}
