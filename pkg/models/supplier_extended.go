package models

import "time"

// ListAdaptersResponse represents the response from listing supplier adapters
type ListAdaptersResponse struct {
	Adapters []*SupplierAdapter `json:"adapters"`
	Count    int32              `json:"count"`
}

// SupplierAdapter represents a supplier integration adapter
type SupplierAdapter struct {
	Name         string            `json:"name"`
	DisplayName  string            `json:"display_name"`
	Description  string            `json:"description"`
	Version      string            `json:"version"`
	Capabilities *AdapterCapabilities `json:"capabilities,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
	IsActive     bool              `json:"is_active"`
}

// AdapterCapabilities represents the capabilities of a supplier adapter
type AdapterCapabilities struct {
	SupportsProductSync   bool     `json:"supports_product_sync"`
	SupportsInventorySync bool     `json:"supports_inventory_sync"`
	SupportsOrderSync     bool     `json:"supports_order_sync"`
	SupportedFormats      []string `json:"supported_formats"`
	MaxBatchSize          int32    `json:"max_batch_size"`
	RateLimitPerMinute    int32    `json:"rate_limit_per_minute"`
	RequiredConfig        []string `json:"required_config"`
	OptionalConfig        []string `json:"optional_config"`
}

// TestConnectionResponse represents the response from testing an adapter connection
type TestConnectionResponse struct {
	Success      bool              `json:"success"`
	Message      string            `json:"message"`
	ErrorCode    string            `json:"error_code,omitempty"`
	ConnectionInfo map[string]string `json:"connection_info,omitempty"`
	TestedAt     time.Time         `json:"tested_at"`
}

// SyncResponse represents the response from sync operations (products/inventory)
type SyncResponse struct {
	JobID         string            `json:"job_id"`
	Status        SyncStatus        `json:"status"`
	Message       string            `json:"message"`
	RecordsTotal  int32             `json:"records_total"`
	RecordsSuccess int32             `json:"records_success"`
	RecordsFailed int32             `json:"records_failed"`
	StartedAt     time.Time         `json:"started_at"`
	CompletedAt   *time.Time        `json:"completed_at,omitempty"`
	ErrorDetails  []SyncError       `json:"error_details,omitempty"`
	DryRun        bool              `json:"dry_run"`
}

// SyncStatus represents the status of a sync operation
type SyncStatus string

const (
	SyncStatusPending    SyncStatus = "pending"
	SyncStatusRunning    SyncStatus = "running"
	SyncStatusCompleted  SyncStatus = "completed"
	SyncStatusFailed     SyncStatus = "failed"
	SyncStatusCancelled  SyncStatus = "cancelled"
)

// SyncError represents an error during sync operations
type SyncError struct {
	RecordID    string `json:"record_id,omitempty"`
	Message     string `json:"message"`
	ErrorCode   string `json:"error_code"`
	LineNumber  int32  `json:"line_number,omitempty"`
}
