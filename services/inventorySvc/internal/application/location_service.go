package application

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/leonvanderhaeghen/stockplatform/services/inventorySvc/internal/domain"
)

// LocationService handles business logic for store locations
type LocationService struct {
	repo   domain.LocationRepository
	logger *zap.Logger
}

// NewLocationService creates a new location service
func NewLocationService(repo domain.LocationRepository, logger *zap.Logger) *LocationService {
	return &LocationService{
		repo:   repo,
		logger: logger.Named("location_service"),
	}
}

// CreateLocation creates a new store location
func (s *LocationService) CreateLocation(ctx context.Context, name, locationType, address, city, state, postalCode, country string) (*domain.StoreLocation, error) {
	s.logger.Info("Creating store location",
		zap.String("name", name),
		zap.String("type", locationType),
	)

	// Check if location with this name already exists
	existingLocation, err := s.repo.GetByName(ctx, name)
	if err != nil && !errors.Is(err, domain.ErrNotFound) {
		return nil, err
	}

	if existingLocation != nil {
		s.logger.Warn("Store location with this name already exists",
			zap.String("name", name),
		)
		return nil, errors.New("store location with this name already exists")
	}

	location := domain.NewStoreLocation(name, locationType, address, city, state, postalCode, country)
	if err := s.repo.Create(ctx, location); err != nil {
		return nil, err
	}

	return location, nil
}

// GetLocation retrieves a store location by ID
func (s *LocationService) GetLocation(ctx context.Context, id string) (*domain.StoreLocation, error) {
	s.logger.Debug("Getting store location", zap.String("id", id))
	
	location, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	if location == nil {
		return nil, errors.New("store location not found")
	}
	
	return location, nil
}

// GetLocationByName retrieves a store location by name
func (s *LocationService) GetLocationByName(ctx context.Context, name string) (*domain.StoreLocation, error) {
	s.logger.Debug("Getting store location by name", zap.String("name", name))
	
	location, err := s.repo.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}
	
	if location == nil {
		return nil, errors.New("store location not found")
	}
	
	return location, nil
}

// UpdateLocation updates an existing store location
func (s *LocationService) UpdateLocation(ctx context.Context, location *domain.StoreLocation) error {
	s.logger.Info("Updating store location",
		zap.String("id", location.ID),
		zap.String("name", location.Name),
	)
	
	return s.repo.Update(ctx, location)
}

// DeleteLocation marks a store location as inactive
func (s *LocationService) DeleteLocation(ctx context.Context, id string) error {
	s.logger.Info("Deleting store location", zap.String("id", id))
	
	// Get the location first
	location, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	
	if location == nil {
		return errors.New("store location not found")
	}
	
	// Set inactive instead of hard delete
	location.IsActive = false
	return s.repo.Update(ctx, location)
}

// ListLocations returns all store locations with pagination and filters
func (s *LocationService) ListLocations(ctx context.Context, limit, offset int, includeInactive bool) ([]*domain.StoreLocation, error) {
	s.logger.Debug("Listing store locations",
		zap.Int("limit", limit),
		zap.Int("offset", offset),
		zap.Bool("include_inactive", includeInactive),
	)
	
	return s.repo.List(ctx, limit, offset, includeInactive)
}

// ListLocationsByType returns all store locations of a specific type
func (s *LocationService) ListLocationsByType(ctx context.Context, locationType string, limit, offset int) ([]*domain.StoreLocation, error) {
	s.logger.Debug("Listing store locations by type",
		zap.String("type", locationType),
		zap.Int("limit", limit),
		zap.Int("offset", offset),
	)
	
	return s.repo.ListByType(ctx, locationType, limit, offset)
}
