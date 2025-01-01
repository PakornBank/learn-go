package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/PakornBank/learn-go/internal/model"
	"github.com/PakornBank/learn-go/internal/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTest(t *testing.T) (*sql.DB, *gorm.DB, sqlmock.Sqlmock, *UserRepository) {
	sqlDB, gormDB, mock := testutil.DbMock(t)
	userRepo := NewUserRepository(gormDB)
	return sqlDB, gormDB, mock, userRepo
}

func TestNewUserRepository(t *testing.T) {
	_, gormDB, _, _ := setupTest(t)
	userRepo := NewUserRepository(gormDB)
	assert.Equal(t, gormDB, userRepo.db)
}

func TestUserRepository_Create(t *testing.T) {
	sqlDB, _, mock, userRepo := setupTest(t)
	defer sqlDB.Close()

	mockUser := testutil.NewMockUser()

	tests := []struct {
		name    string
		user    *model.User
		mockFn  func(sqlmock.Sqlmock)
		wantErr bool
		errType error
	}{
		{
			name: "successful creation",
			user: &model.User{
				Email:        mockUser.Email,
				PasswordHash: mockUser.PasswordHash,
				FullName:     mockUser.FullName,
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
					AddRow(mockUser.ID, mockUser.CreatedAt, mockUser.UpdatedAt)
				mock.ExpectQuery(`INSERT INTO "users"`).
					WithArgs(mockUser.Email, mockUser.PasswordHash, mockUser.FullName).
					WillReturnRows(rows)
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "database error",
			user: &model.User{
				Email:        mockUser.Email,
				PasswordHash: mockUser.PasswordHash,
				FullName:     mockUser.FullName,
			},
			mockFn: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "users"`).
					WithArgs(mockUser.Email, mockUser.PasswordHash, mockUser.FullName).
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			wantErr: true,
			errType: sql.ErrConnDone,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(mock)
			err := userRepo.Create(context.Background(), tt.user)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errType, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, tt.user.ID)
				assert.NotZero(t, tt.user.CreatedAt)
				assert.NotZero(t, tt.user.UpdatedAt)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_FindByEmail(t *testing.T) {
	sqlDB, _, mock, userRepo := setupTest(t)
	defer sqlDB.Close()

	mockUser := testutil.NewMockUser()

	tests := []struct {
		name    string
		mockFn  func(sqlmock.Sqlmock)
		want    *model.User
		wantErr bool
		errType error
	}{
		{
			name: "user found",
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "full_name", "created_at", "updated_at"}).
					AddRow(mockUser.ID, mockUser.Email, mockUser.PasswordHash, mockUser.FullName, mockUser.CreatedAt, mockUser.UpdatedAt)
				mock.ExpectQuery(`SELECT .* FROM "users" WHERE email = \$1 (.+) LIMIT \$2`).
					WithArgs(mockUser.Email, 1).
					WillReturnRows(rows)
			},
			want:    &mockUser,
			wantErr: false,
		},
		{
			name: "user not found",
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "full_name", "created_at", "updated_at"})
				mock.ExpectQuery(`SELECT .* FROM "users" WHERE email = \$1 (.+) LIMIT \$2`).
					WithArgs(mockUser.Email, 1).
					WillReturnRows(rows)
			},
			want:    nil,
			wantErr: true,
			errType: gorm.ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(mock)
			got, err := userRepo.FindByEmail(context.Background(), mockUser.Email)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errType, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_FindByID(t *testing.T) {
	sqlDB, _, mock, userRepo := setupTest(t)
	defer sqlDB.Close()

	mockUser := testutil.NewMockUser()

	tests := []struct {
		name    string
		id      uuid.UUID
		mockFn  func(sqlmock.Sqlmock)
		want    *model.User
		wantErr bool
		errType error
	}{
		{
			name: "user found",
			id:   mockUser.ID,
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "full_name", "created_at", "updated_at"}).
					AddRow(mockUser.ID, mockUser.Email, mockUser.PasswordHash, mockUser.FullName, mockUser.CreatedAt, mockUser.UpdatedAt)
				mock.ExpectQuery(`SELECT .* FROM "users" WHERE id = \$1 (.+) LIMIT \$2`).
					WithArgs(mockUser.ID, 1).
					WillReturnRows(rows)
			},
			want:    &mockUser,
			wantErr: false,
		},
		{
			name: "user not found",
			id:   mockUser.ID,
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "full_name", "created_at", "updated_at"})
				mock.ExpectQuery(`SELECT .* FROM "users" WHERE id = \$1 (.+) LIMIT \$2`).
					WithArgs(mockUser.ID, 1).
					WillReturnRows(rows)
			},
			want:    nil,
			wantErr: true,
			errType: gorm.ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(mock)
			got, err := userRepo.FindByID(context.Background(), tt.id.String())

			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errType, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
