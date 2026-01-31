package service

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/your-org/project-budget-tracker/backend/internal/dto"
	apperrors "github.com/your-org/project-budget-tracker/backend/internal/errors"
	"github.com/your-org/project-budget-tracker/backend/internal/models"
	"github.com/your-org/project-budget-tracker/backend/internal/repository"
)

// MemberService handles business logic for members
type MemberService struct {
	memberRepo *repository.MemberRepository
	db         *gorm.DB
}

// NewMemberService creates a new MemberService
func NewMemberService(db *gorm.DB) *MemberService {
	return &MemberService{
		memberRepo: repository.NewMemberRepository(db),
		db:         db,
	}
}

// CreateMember creates a new member
func (s *MemberService) CreateMember(req *dto.CreateMemberRequest) (*dto.MemberResponse, error) {
	// Check if email already exists
	existing, err := s.memberRepo.GetByEmail(req.Email)
	if err == nil && existing != nil {
		return nil, apperrors.ErrConflict("Member with this email already exists")
	}

	member := &models.Member{
		Name:       req.Name,
		Email:      req.Email,
		Role:       req.Role,
		HourlyRate: req.HourlyRate,
		Department: req.Department,
		UserID:     req.UserID,
	}

	if err := s.memberRepo.Create(member); err != nil {
		return nil, apperrors.ErrDatabaseError(err)
	}

	return s.toMemberResponse(member), nil
}

// GetMember retrieves a member by ID
func (s *MemberService) GetMember(id uuid.UUID) (*dto.MemberResponse, error) {
	member, err := s.memberRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound("Member")
		}
		return nil, apperrors.ErrDatabaseError(err)
	}

	return s.toMemberResponse(member), nil
}

// ListMembers retrieves members with pagination
func (s *MemberService) ListMembers(page, perPage int, search, department string) (*dto.MemberListResponse, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	params := repository.MemberListParams{
		Page:       page,
		PerPage:    perPage,
		Search:     search,
		Department: department,
	}

	members, total, err := s.memberRepo.List(params)
	if err != nil {
		return nil, apperrors.ErrDatabaseError(err)
	}

	// Convert to response
	memberResponses := make([]dto.MemberResponse, len(members))
	for i, member := range members {
		memberResponses[i] = *s.toMemberResponse(&member)
	}

	totalPages := int(total) / perPage
	if int(total)%perPage > 0 {
		totalPages++
	}

	return &dto.MemberListResponse{
		Members: memberResponses,
		Pagination: dto.Pagination{
			Page:       page,
			PerPage:    perPage,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

// UpdateMember updates a member
func (s *MemberService) UpdateMember(id uuid.UUID, req *dto.UpdateMemberRequest) (*dto.MemberResponse, error) {
	member, err := s.memberRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound("Member")
		}
		return nil, apperrors.ErrDatabaseError(err)
	}

	// Update fields
	if req.Name != nil {
		member.Name = *req.Name
	}
	if req.Email != nil {
		// Check if new email conflicts with existing member
		existing, err := s.memberRepo.GetByEmail(*req.Email)
		if err == nil && existing != nil && existing.ID != id {
			return nil, apperrors.ErrConflict("Member with this email already exists")
		}
		member.Email = *req.Email
	}
	if req.Role != nil {
		member.Role = req.Role
	}
	if req.HourlyRate != nil {
		member.HourlyRate = *req.HourlyRate
	}
	if req.Department != nil {
		member.Department = req.Department
	}

	if err := s.memberRepo.Update(member); err != nil {
		return nil, apperrors.ErrDatabaseError(err)
	}

	return s.toMemberResponse(member), nil
}

// DeleteMember deletes a member
func (s *MemberService) DeleteMember(id uuid.UUID) error {
	_, err := s.memberRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.ErrNotFound("Member")
		}
		return apperrors.ErrDatabaseError(err)
	}

	if err := s.memberRepo.Delete(id); err != nil {
		return apperrors.ErrDatabaseError(err)
	}

	return nil
}

// GetProjectMembers retrieves members assigned to a project
func (s *MemberService) GetProjectMembers(projectID uuid.UUID) ([]dto.ProjectMemberResponse, error) {
	// Verify project exists
	var project models.Project
	if err := s.db.First(&project, "id = ?", projectID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound("Project")
		}
		return nil, apperrors.ErrDatabaseError(err)
	}

	projectMembers, err := s.memberRepo.GetProjectMembers(projectID)
	if err != nil {
		return nil, apperrors.ErrDatabaseError(err)
	}

	responses := make([]dto.ProjectMemberResponse, len(projectMembers))
	for i, pm := range projectMembers {
		responses[i] = *s.toProjectMemberResponse(&pm)
	}

	return responses, nil
}

// AssignMemberToProject assigns a member to a project
func (s *MemberService) AssignMemberToProject(projectID uuid.UUID, req *dto.AssignMemberRequest) (*dto.ProjectMemberResponse, error) {
	// Verify project exists
	var project models.Project
	if err := s.db.First(&project, "id = ?", projectID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound("Project")
		}
		return nil, apperrors.ErrDatabaseError(err)
	}

	// Verify member exists
	member, err := s.memberRepo.GetByID(req.MemberID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound("Member")
		}
		return nil, apperrors.ErrDatabaseError(err)
	}

	// Check if member is already assigned
	existing, _ := s.memberRepo.GetProjectMember(projectID, req.MemberID)
	if existing != nil {
		return nil, apperrors.ErrConflict("Member is already assigned to this project")
	}

	// Set default values
	allocationRate := 1.0
	if req.AllocationRate != nil {
		allocationRate = *req.AllocationRate
	}

	hourlyRateSnapshot := member.HourlyRate
	if req.HourlyRateSnapshot != nil {
		hourlyRateSnapshot = *req.HourlyRateSnapshot
	}

	projectMember := &models.ProjectMember{
		ProjectID:          projectID,
		MemberID:           req.MemberID,
		AllocationRate:     allocationRate,
		HourlyRateSnapshot: &hourlyRateSnapshot,
	}

	if err := s.memberRepo.AssignToProject(projectMember); err != nil {
		return nil, apperrors.ErrDatabaseError(err)
	}

	// Reload with member data
	pm, err := s.memberRepo.GetProjectMember(projectID, req.MemberID)
	if err != nil {
		return nil, apperrors.ErrDatabaseError(err)
	}

	return s.toProjectMemberResponse(pm), nil
}

// RemoveMemberFromProject removes a member from a project
func (s *MemberService) RemoveMemberFromProject(projectID, memberID uuid.UUID) error {
	// Verify assignment exists
	_, err := s.memberRepo.GetProjectMember(projectID, memberID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.ErrNotFound("Project member assignment")
		}
		return apperrors.ErrDatabaseError(err)
	}

	if err := s.memberRepo.RemoveFromProject(projectID, memberID); err != nil {
		return apperrors.ErrDatabaseError(err)
	}

	return nil
}

// toMemberResponse converts a Member model to MemberResponse DTO
func (s *MemberService) toMemberResponse(member *models.Member) *dto.MemberResponse {
	return &dto.MemberResponse{
		ID:         member.ID,
		UserID:     member.UserID,
		Name:       member.Name,
		Email:      member.Email,
		Role:       member.Role,
		HourlyRate: member.HourlyRate,
		Department: member.Department,
		CreatedAt:  member.CreatedAt,
		UpdatedAt:  member.UpdatedAt,
	}
}

// toProjectMemberResponse converts a ProjectMember model to ProjectMemberResponse DTO
func (s *MemberService) toProjectMemberResponse(pm *models.ProjectMember) *dto.ProjectMemberResponse {
	response := &dto.ProjectMemberResponse{
		ID:                 pm.ID,
		ProjectID:          pm.ProjectID,
		MemberID:           pm.MemberID,
		JoinedAt:           pm.JoinedAt.Format("2006-01-02"),
		AllocationRate:     pm.AllocationRate,
		HourlyRateSnapshot: pm.HourlyRateSnapshot,
	}

	if pm.LeftAt != nil {
		leftAt := pm.LeftAt.Format("2006-01-02")
		response.LeftAt = &leftAt
	}

	if pm.Member.ID != uuid.Nil {
		response.Member = &dto.MemberResponse{
			ID:         pm.Member.ID,
			Name:       pm.Member.Name,
			Email:      pm.Member.Email,
			Role:       pm.Member.Role,
			HourlyRate: pm.Member.HourlyRate,
			Department: pm.Member.Department,
		}
	}

	return response
}
