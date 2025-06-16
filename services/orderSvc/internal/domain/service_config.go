package domain

// ServiceConfig holds configuration for services
type ServiceConfig struct {
	// Configuration for connections to other services
	InventoryServiceAddr string
	ProductServiceAddr   string
	PaymentServiceAddr   string
	NotificationServiceAddr string

	// Local service configuration
	ServerPort    string
	DatabaseURI   string
	JWTSecret     string
	DefaultLocationID string
}
