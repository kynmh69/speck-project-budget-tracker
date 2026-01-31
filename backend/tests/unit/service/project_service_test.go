package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/your-org/project-budget-tracker/backend/internal/dto"
	apperrors "github.com/your-org/project-budget-tracker/backend/internal/errors"
	"github.com/your-org/project-budget-tracker/backend/internal/models"
	"github.com/your-org/project-budget-tracker/backend/internal/service"
)

// MockProjectRepository はテスト用のモックリポジトリ
type MockProjectRepository struct {
	mock.Mock
}

func (m *MockProjectRepository) Create(ctx context.Context, project *models.Project) error {
	args := m.Called(ctx, project)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	project.ID = 1
	return nil
}

func (m *MockProjectRepository) GetByID(ctx context.Context, id uint) (*models.Project, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Project), args.Error(1)
}

func (m *MockProjectRepository) GetByIDWithDetails(ctx context.Context, id uint) (*models.Project, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Project), args.Error(1)
}

func (m *MockProjectRepository) List(ctx context.Context, ownerID uint, params dto.ProjectListParams) (*dto.ProjectListResult, error) {
	args := m.Called(ctx, ownerID, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ProjectListResult), args.Error(1)
}

func (m *MockProjectRepository) Update(ctx context.Context, project *models.Project) error {
	args := m.Called(ctx, project)
	return args.Error(0)
}

func (m *MockProjectRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockProjectRepository) IsOwner(ctx context.Context, projectID, userID uint) (bool, error) {
	args := m.Called(ctx, projectID, userID)
	return args.Bool(0), args.Error(1)
}

func (m *MockProjectRepository) GetProjectStats(ctx context.Context, projectID uint) (*dto.ProjectStats, error) {
	args := m.Called(ctx, projectID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ProjectStats), args.Error(1)
}

func TestProjectService_CreateProject(t *testing.T) {
	tests := []struct {
		name      string
		userID    uint
		req       dto.CreateProjectRequest
		setupMock func(*MockProjectRepository)
		wantErr   bool
		errType   error
	}{
		{
			name:   "正常: プロジェクト作成成功",
			userID: 1,
			req: dto.CreateProjectRequest{
				Name:        "テストプロジェクト",
				Description: "テスト用のプロジェクトです",
				Status:      "planning",
			},
			setupMock: func(repo *MockProjectRepository) {
				repo.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "正常: 全フィールド指定",
			userID: 1,
			req: dto.CreateProjectRequest{
				Name:         "フルスペックプロジェクト",
				Description:  "全フィールド入力",
				Status:       "in_progress",
				StartDate:    stringPtr("2024-01-01"),
				EndDate:      stringPtr("2024-12-31"),
				BudgetAmount: float64Ptr(1000000),
			},
			setupMock: func(repo *MockProjectRepository) {
				repo.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name:   "異常: リポジトリエラー",
			userID: 1,
			req: dto.CreateProjectRequest{
				Name:   "エラープロジェクト",
				Status: "planning",
			},
			setupMock: func(repo *MockProjectRepository) {
				repo.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockProjectRepository)
			tt.setupMock(mockRepo)

			svc := service.NewProjectService(mockRepo)
			ctx := context.Background()

			result, err := svc.CreateProject(ctx, tt.userID, tt.req)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.req.Name, result.Name)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestProjectService_GetProject(t *testing.T) {
	now := time.Now()
	testProject := &models.Project{
		ID:          1,
		OwnerID:     1,
		Name:        "テストプロジェクト",
		Description: "説明文",
		Status:      "in_progress",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	tests := []struct {
		name      string
		projectID string
		userID    uint
		setupMock func(*MockProjectRepository)
		wantErr   bool
		errType   error
	}{
		{
			name:      "正常: プロジェクト取得成功",
			projectID: "1",
			userID:    1,
			setupMock: func(repo *MockProjectRepository) {
				repo.On("GetByIDWithDetails", mock.Anything, uint(1)).Return(testProject, nil)
				repo.On("IsOwner", mock.Anything, uint(1), uint(1)).Return(true, nil)
				repo.On("GetProjectStats", mock.Anything, uint(1)).Return(&dto.ProjectStats{
					TotalTasks:     10,
					CompletedTasks: 5,
				}, nil)
			},
			wantErr: false,
		},
		{
			name:      "異常: 無効なID形式",
			projectID: "invalid",
			userID:    1,
			setupMock: func(repo *MockProjectRepository) {},
			wantErr:   true,
		},
		{
			name:      "異常: プロジェクトが存在しない",
			projectID: "999",
			userID:    1,
			setupMock: func(repo *MockProjectRepository) {
				repo.On("GetByIDWithDetails", mock.Anything, uint(999)).Return(nil, apperrors.ErrNotFound)
			},
			wantErr: true,
			errType: apperrors.ErrNotFound,
		},
		{
			name:      "異常: 権限がない",
			projectID: "1",
			userID:    2,
			setupMock: func(repo *MockProjectRepository) {
				repo.On("GetByIDWithDetails", mock.Anything, uint(1)).Return(testProject, nil)
				repo.On("IsOwner", mock.Anything, uint(1), uint(2)).Return(false, nil)
			},
			wantErr: true,
			errType: apperrors.ErrForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockProjectRepository)
			tt.setupMock(mockRepo)

			svc := service.NewProjectService(mockRepo)
			ctx := context.Background()

			result, err := svc.GetProject(ctx, tt.projectID, tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errType != nil {
					assert.True(t, errors.Is(err, tt.errType))
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestProjectService_ListProjects(t *testing.T) {
	tests := []struct {
		name      string
		userID    uint
		params    dto.ProjectListParams
		setupMock func(*MockProjectRepository)
		wantErr   bool
		wantCount int
	}{
		{
			name:   "正常: プロジェクト一覧取得",
			userID: 1,
			params: dto.ProjectListParams{
				Page:    1,
				PerPage: 10,
			},
			setupMock: func(repo *MockProjectRepository) {
				repo.On("List", mock.Anything, uint(1), mock.Anything).Return(&dto.ProjectListResult{
					Projects:   []models.Project{{ID: 1, Name: "Project 1"}, {ID: 2, Name: "Project 2"}},
					TotalCount: 2,
					Page:       1,
					PerPage:    10,
					TotalPages: 1,
				}, nil)
			},
			wantErr:   false,
			wantCount: 2,
		},
		{
			name:   "正常: ステータスでフィルタ",
			userID: 1,
			params: dto.ProjectListParams{
				Page:    1,
				PerPage: 10,
				Status:  "in_progress",
			},
			setupMock: func(repo *MockProjectRepository) {
				repo.On("List", mock.Anything, uint(1), mock.Anything).Return(&dto.ProjectListResult{
					Projects:   []models.Project{{ID: 1, Name: "Active Project", Status: "in_progress"}},
					TotalCount: 1,
					Page:       1,
					PerPage:    10,
					TotalPages: 1,
				}, nil)
			},
			wantErr:   false,
			wantCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockProjectRepository)
			tt.setupMock(mockRepo)

			svc := service.NewProjectService(mockRepo)
			ctx := context.Background()

			result, err := svc.ListProjects(ctx, tt.userID, tt.params)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Len(t, result.Projects, tt.wantCount)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestProjectService_UpdateProject(t *testing.T) {
	now := time.Now()
	testProject := &models.Project{
		ID:          1,
		OwnerID:     1,
		Name:        "更新前",
		Description: "説明文",
		Status:      "planning",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	tests := []struct {
		name      string
		projectID string
		userID    uint
		req       dto.UpdateProjectRequest
		setupMock func(*MockProjectRepository)
		wantErr   bool
	}{
		{
			name:      "正常: プロジェクト更新成功",
			projectID: "1",
			userID:    1,
			req: dto.UpdateProjectRequest{
				Name:   stringPtr("更新後"),
				Status: stringPtr("in_progress"),
			},
			setupMock: func(repo *MockProjectRepository) {
				repo.On("GetByID", mock.Anything, uint(1)).Return(testProject, nil)
				repo.On("IsOwner", mock.Anything, uint(1), uint(1)).Return(true, nil)
				repo.On("Update", mock.Anything, mock.Anything).Return(nil)
				repo.On("GetByIDWithDetails", mock.Anything, uint(1)).Return(testProject, nil)
				repo.On("GetProjectStats", mock.Anything, uint(1)).Return(&dto.ProjectStats{}, nil)
			},
			wantErr: false,
		},
		{
			name:      "異常: 権限がない",
			projectID: "1",
			userID:    2,
			req: dto.UpdateProjectRequest{
				Name: stringPtr("更新後"),
			},
			setupMock: func(repo *MockProjectRepository) {
				repo.On("GetByID", mock.Anything, uint(1)).Return(testProject, nil)
				repo.On("IsOwner", mock.Anything, uint(1), uint(2)).Return(false, nil)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockProjectRepository)
			tt.setupMock(mockRepo)

			svc := service.NewProjectService(mockRepo)
			ctx := context.Background()

			result, err := svc.UpdateProject(ctx, tt.projectID, tt.userID, tt.req)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestProjectService_DeleteProject(t *testing.T) {
	now := time.Now()
	testProject := &models.Project{
		ID:        1,
		OwnerID:   1,
		Name:      "削除対象",
		CreatedAt: now,
		UpdatedAt: now,
	}

	tests := []struct {
		name      string
		projectID string
		userID    uint
		setupMock func(*MockProjectRepository)
		wantErr   bool
	}{
		{
			name:      "正常: プロジェクト削除成功",
			projectID: "1",
			userID:    1,
			setupMock: func(repo *MockProjectRepository) {
				repo.On("GetByID", mock.Anything, uint(1)).Return(testProject, nil)
				repo.On("IsOwner", mock.Anything, uint(1), uint(1)).Return(true, nil)
				repo.On("Delete", mock.Anything, uint(1)).Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "異常: 権限がない",
			projectID: "1",
			userID:    2,
			setupMock: func(repo *MockProjectRepository) {
				repo.On("GetByID", mock.Anything, uint(1)).Return(testProject, nil)
				repo.On("IsOwner", mock.Anything, uint(1), uint(2)).Return(false, nil)
			},
			wantErr: true,
		},
		{
			name:      "異常: プロジェクトが存在しない",
			projectID: "999",
			userID:    1,
			setupMock: func(repo *MockProjectRepository) {
				repo.On("GetByID", mock.Anything, uint(999)).Return(nil, apperrors.ErrNotFound)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockProjectRepository)
			tt.setupMock(mockRepo)

			svc := service.NewProjectService(mockRepo)
			ctx := context.Background()

			err := svc.DeleteProject(ctx, tt.projectID, tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func stringPtr(s string) *string {
	return &s
}

func float64Ptr(f float64) *float64 {
	return &f
}
