package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/your-org/project-budget-tracker/backend/internal/models"
)

// MemberRepository handles database operations for members
type MemberRepository struct {
	db *gorm.DB
}

// NewMemberRepository creates a new MemberRepository
func NewMemberRepository(db *gorm.DB) *MemberRepository {
	return &MemberRepository{db: db}
}

// Create creates a new member
func (r *MemberRepository) Create(member *models.Member) error {
	return r.db.Create(member).Error
}

// GetByID retrieves a member by ID
func (r *MemberRepository) GetByID(id uuid.UUID) (*models.Member, error) {
	var member models.Member
	if err := r.db.First(&member, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &member, nil
}

// GetByEmail retrieves a member by email
func (r *MemberRepository) GetByEmail(email string) (*models.Member, error) {
	var member models.Member
	if err := r.db.First(&member, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &member, nil
}

// MemberListParams represents parameters for listing members
type MemberListParams struct {
	Page       int
	PerPage    int
	Search     string
	Department string
}

// List retrieves members with pagination and filtering
func (r *MemberRepository) List(params MemberListParams) ([]models.Member, int64, error) {
	var members []models.Member
	var total int64

	query := r.db.Model(&models.Member{})

	// Apply search filter
	if params.Search != "" {
		searchPattern := "%" + params.Search + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ?", searchPattern, searchPattern)
	}

	// Apply department filter
	if params.Department != "" {
		query = query.Where("department = ?", params.Department)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (params.Page - 1) * params.PerPage
	if err := query.Order("name ASC").Offset(offset).Limit(params.PerPage).Find(&members).Error; err != nil {
		return nil, 0, err
	}

	return members, total, nil
}

// Update updates a member
func (r *MemberRepository) Update(member *models.Member) error {
	return r.db.Save(member).Error
}

// Delete soft deletes a member
func (r *MemberRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Member{}, "id = ?", id).Error
}

// GetByProjectID retrieves members assigned to a project
func (r *MemberRepository) GetByProjectID(projectID uuid.UUID) ([]models.Member, error) {
	var members []models.Member
	if err := r.db.
		Joins("JOIN project_members ON project_members.member_id = members.id").
		Where("project_members.project_id = ? AND project_members.left_at IS NULL", projectID).
		Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

// AssignToProject assigns a member to a project
func (r *MemberRepository) AssignToProject(projectMember *models.ProjectMember) error {
	return r.db.Create(projectMember).Error
}

// RemoveFromProject removes a member from a project (sets left_at)
func (r *MemberRepository) RemoveFromProject(projectID, memberID uuid.UUID) error {
	return r.db.Model(&models.ProjectMember{}).
		Where("project_id = ? AND member_id = ? AND left_at IS NULL", projectID, memberID).
		Update("left_at", gorm.Expr("CURRENT_DATE")).Error
}

// GetProjectMember retrieves a project member assignment
func (r *MemberRepository) GetProjectMember(projectID, memberID uuid.UUID) (*models.ProjectMember, error) {
	var pm models.ProjectMember
	if err := r.db.
		Preload("Member").
		Where("project_id = ? AND member_id = ? AND left_at IS NULL", projectID, memberID).
		First(&pm).Error; err != nil {
		return nil, err
	}
	return &pm, nil
}

// GetProjectMembers retrieves all active project member assignments
func (r *MemberRepository) GetProjectMembers(projectID uuid.UUID) ([]models.ProjectMember, error) {
	var projectMembers []models.ProjectMember
	if err := r.db.
		Preload("Member").
		Where("project_id = ? AND left_at IS NULL", projectID).
		Find(&projectMembers).Error; err != nil {
		return nil, err
	}
	return projectMembers, nil
}
