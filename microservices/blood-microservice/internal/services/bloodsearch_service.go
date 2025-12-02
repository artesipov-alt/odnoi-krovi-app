package services

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	bloodsearchv1 "github.com/artesipov-alt/odnoi-krovi-app/microservices/blood-microservice/gen/api/bloodsearch/v1"
	"github.com/artesipov-alt/odnoi-krovi-app/microservices/blood-microservice/internal/repositories"

	"go.uber.org/zap"
)

// BloodSearchService implements the BloodSearchPool service
type BloodSearchService struct {
	repo repositories.BloodSearchRepository
	log  *zap.Logger
}

// NewBloodSearchService creates a new instance of BloodSearchService
func NewBloodSearchService(repo repositories.BloodSearchRepository, log *zap.Logger) *BloodSearchService {
	if log == nil {
		log = zap.NewNop()
	}
	return &BloodSearchService{
		repo: repo,
		log:  log,
	}
}

// AddPet adds a new pet blood search request to the database
func (s *BloodSearchService) AddPet(ctx context.Context, req *bloodsearchv1.PetRow) (*bloodsearchv1.PetRowStatus, error) {
	s.log.Info("AddPet called", zap.String("pet_id", req.PetId))

	// Validate request
	if req.PetId == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("pet_id is required"))
	}

	if req.PetType == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("pet_type is required"))
	}

	if req.BloodGroup == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("blood_group is required"))
	}

	if req.BloodVolume <= 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("blood_volume must be positive"))
	}

	if len(req.Regions) == 0 {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("at least one region must be specified"))
	}

	// Set default status if not provided
	if req.Status == "" {
		req.Status = "active"
	}

	// Check if pet already exists
	exists, err := s.repo.Exists(ctx, req.PetId)
	if err != nil {
		s.log.Error("Failed to check if pet exists", zap.Error(err))
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to check pet existence"))
	}

	if exists {
		return nil, connect.NewError(connect.CodeAlreadyExists, errors.New("pet with this ID already exists"))
	}

	// Add request to repository
	err = s.repo.AddRequest(ctx, req)
	if err != nil {
		s.log.Error("Failed to add pet request", zap.Error(err))
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to add pet request"))
	}

	// Create response
	response := &bloodsearchv1.PetRowStatus{
		PetId:  req.PetId,
		Status: req.Status,
	}

	s.log.Info("Pet added successfully", zap.String("pet_id", req.PetId))
	return response, nil
}

// GetPets retrieves pet blood search requests based on criteria
func (s *BloodSearchService) GetPets(ctx context.Context, req *bloodsearchv1.GetPetRows) (*bloodsearchv1.PetRows, error) {
	s.log.Info("GetPets called",
		zap.String("pet_type", req.PetType),
		zap.String("blood_group", req.BloodGroup),
		zap.Int("regions_count", len(req.Regions)))

	// Get requests from repository
	petRows, err := s.repo.GetRequestsByCriteria(ctx, req)
	if err != nil {
		s.log.Error("Failed to get pet requests", zap.Error(err))
		return nil, connect.NewError(connect.CodeInternal, errors.New("failed to get pet requests"))
	}

	// Create response
	response := &bloodsearchv1.PetRows{
		Pets: petRows,
	}

	s.log.Info("GetPets completed", zap.Int("pets_count", len(petRows)))
	return response, nil
}
