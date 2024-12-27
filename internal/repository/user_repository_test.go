package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/PakornBank/learn-go/internal/model"
	"github.com/PakornBank/learn-go/internal/testutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestNewUserRepository(t *testing.T) {
	_, gormDB, _ := testutil.DbMock(t)
	userRepo := NewUserRepository(gormDB)

	assert.Equal(t, gormDB, userRepo.db)
}

func TestUserRepository_Create(t *testing.T) {
	sqlDB, gormDB, mock := testutil.DbMock(t)
	defer sqlDB.Close()
	userRepo := NewUserRepository(gormDB)

	tests := []struct {
		name    string
		user    *model.User
		mockFn  func(sqlmock.Sqlmock, *model.User)
		wantErr bool
	}{
		{
			name: "successful creation",
			user: &model.User{
				ID:           uuid.New(),
				Email:        "test1@example.com",
				PasswordHash: "hashedpassword1",
				FullName:     "User Name1",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			},
			mockFn: func(mock sqlmock.Sqlmock, user *model.User) {
				mock.ExpectBegin()
				rows := sqlmock.NewRows(
					[]string{"id", "created_at", "updated_at"}).
					AddRow(user.ID, user.CreatedAt, user.UpdatedAt)
				mock.ExpectQuery(`INSERT INTO "users"`).
					WithArgs(user.Email, user.PasswordHash, user.FullName,
						user.ID, user.CreatedAt, user.UpdatedAt).
					WillReturnRows(rows)
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "failed creation",
			user: &model.User{
				ID:           uuid.New(),
				Email:        "test2@example.com",
				PasswordHash: "hashedpassword2",
				FullName:     "User Name2",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			},
			mockFn: func(mock sqlmock.Sqlmock, user *model.User) {
				mock.ExpectBegin()
				mock.ExpectQuery(`INSERT INTO "users"`).
					WithArgs(user.Email, user.PasswordHash, user.FullName,
						user.ID, user.CreatedAt, user.UpdatedAt).
					WillReturnError(sql.ErrConnDone)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(mock, tt.user)
			err := userRepo.Create(context.Background(), tt.user)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, sql.ErrConnDone, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Nil(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_FindByEmail(t *testing.T) {
	sqlDB, gormDB, mock := testutil.DbMock(t)
	defer sqlDB.Close()
	userRepo := NewUserRepository(gormDB)
	testUUID := uuid.New()
	now := time.Now()

	tests := []struct {
		name    string
		email   string
		mockFn  func(sqlmock.Sqlmock, string)
		want    *model.User
		wantErr bool
	}{
		{
			name:  "user found",
			email: "test1@example.com",
			mockFn: func(mock sqlmock.Sqlmock, email string) {
				rows := sqlmock.NewRows([]string{"id", "email",
					"password_hash", "full_name", "created_at", "updated_at"}).
					AddRow(testUUID, email, "hashedpassword1", "User Name1", now, now)
				mock.ExpectQuery(`SELECT .* FROM "users" WHERE email = \$1 (.+) LIMIT \$2`).WithArgs(email, 1).WillReturnRows(rows)
			},
			want: &model.User{
				ID:           testUUID,
				Email:        "test1@example.com",
				PasswordHash: "hashedpassword1",
				FullName:     "User Name1",
				CreatedAt:    now,
				UpdatedAt:    now,
			},
			wantErr: false,
		},
		{
			name:  "user not found",
			email: "test2@example.com",
			mockFn: func(mock sqlmock.Sqlmock, email string) {
				rows := sqlmock.NewRows([]string{"id", "email",
					"password_hash", "full_name", "created_at", "updated_at"})
				mock.ExpectQuery(`SELECT .* FROM "users" WHERE email = \$1 (.+) LIMIT \$2`).WithArgs(email, 1).WillReturnRows(rows)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(mock, tt.email)
			got, err := userRepo.FindByEmail(context.Background(), tt.email)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, gorm.ErrRecordNotFound, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.Nil(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_FindById(t *testing.T) {
	sqlDB, gormDB, mock := testutil.DbMock(t)
	defer sqlDB.Close()
	userRepo := NewUserRepository(gormDB)
	testUUID := uuid.New()
	now := time.Now()

	tests := []struct {
		name    string
		id      uuid.UUID
		mockFn  func(sqlmock.Sqlmock)
		want    *model.User
		wantErr bool
	}{
		{
			name: "user found",
			id:   testUUID,
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "email",
					"password_hash", "full_name", "created_at", "updated_at"}).
					AddRow(testUUID, "test1@example.com", "hashedpassword1", "User Name1", now, now)
				mock.ExpectQuery(`SELECT .* FROM "users" WHERE id = \$1 (.+) LIMIT \$2`).WithArgs(testUUID, 1).WillReturnRows(rows)
			},
			want: &model.User{
				ID:           testUUID,
				Email:        "test1@example.com",
				PasswordHash: "hashedpassword1",
				FullName:     "User Name1",
				CreatedAt:    now,
				UpdatedAt:    now,
			},
			wantErr: false,
		},
		{
			name: "user not found",
			id:   testUUID,
			mockFn: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "email",
					"password_hash", "full_name", "created_at", "updated_at"})
				mock.ExpectQuery(`SELECT .* FROM "users" WHERE id = \$1 (.+) LIMIT \$2`).WithArgs(testUUID, 1).WillReturnRows(rows)
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFn(mock)
			got, err := userRepo.FindById(context.Background(), testUUID.String())
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, gorm.ErrRecordNotFound, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.Nil(t, mock.ExpectationsWereMet())
		})
	}
}
