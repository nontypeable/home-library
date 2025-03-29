package repository

import (
	"context"
	"database/sql"
	"home-library/internal/services/user/entities"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewRepository(sqlxDB)

	t.Run("create new user with valid data", func(t *testing.T) {
		userID := uuid.New()
		now := time.Now()
		user := &entities.User{
			UserID:      userID,
			FirstName:   "Evgeny",
			LastName:    "Koveshnikov",
			Email:       "evgeny@example.com",
			PhoneNumber: "+79001234567",
			Password:    "hashedPassword",
			UserType:    "user",
			IsActive:    true,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		mock.ExpectExec("INSERT INTO users").
			WithArgs(
				userID,
				user.FirstName,
				user.LastName,
				user.Email,
				user.PhoneNumber,
				user.Password,
				user.UserType,
				user.IsActive,
				user.CreatedAt,
				user.UpdatedAt,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))

		id, err := repo.CreateUser(context.Background(), user)

		assert.NoError(t, err)
		assert.Equal(t, userID, id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("create admin user", func(t *testing.T) {
		userID := uuid.New()
		now := time.Now()
		user := &entities.User{
			UserID:      userID,
			FirstName:   "Admin",
			LastName:    "Adminov",
			Email:       "admin@example.com",
			PhoneNumber: "+79001234568",
			Password:    "hashedPassword",
			UserType:    "admin",
			IsActive:    true,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		mock.ExpectExec("INSERT INTO users").
			WithArgs(
				userID,
				user.FirstName,
				user.LastName,
				user.Email,
				user.PhoneNumber,
				user.Password,
				user.UserType,
				user.IsActive,
				user.CreatedAt,
				user.UpdatedAt,
			).
			WillReturnResult(sqlmock.NewResult(1, 1))

		id, err := repo.CreateUser(context.Background(), user)

		assert.NoError(t, err)
		assert.Equal(t, userID, id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("attempt to create user with duplicate email", func(t *testing.T) {
		userID := uuid.New()
		now := time.Now()
		user := &entities.User{
			UserID:      userID,
			FirstName:   "Sergey",
			LastName:    "Sergeev",
			Email:       "evgeny@example.com", // same email as in first test
			PhoneNumber: "+79001234569",
			Password:    "hashedPassword",
			UserType:    "user",
			IsActive:    true,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		mock.ExpectExec("INSERT INTO users").
			WithArgs(
				userID,
				user.FirstName,
				user.LastName,
				user.Email,
				user.PhoneNumber,
				user.Password,
				user.UserType,
				user.IsActive,
				user.CreatedAt,
				user.UpdatedAt,
			).
			WillReturnError(&pq.Error{
				Code:    "23505",
				Message: "duplicate key value violates unique constraint \"users_email_key\"",
			})

		id, err := repo.CreateUser(context.Background(), user)

		assert.Error(t, err)
		assert.Equal(t, uuid.Nil, id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("attempt to create user with duplicate phone number", func(t *testing.T) {
		userID := uuid.New()
		now := time.Now()
		user := &entities.User{
			UserID:      userID,
			FirstName:   "Alexey",
			LastName:    "Alexeev",
			Email:       "alexey@example.com",
			PhoneNumber: "+79001234567", // same phone number as in first test
			Password:    "hashedPassword",
			UserType:    "user",
			IsActive:    true,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		mock.ExpectExec("INSERT INTO users").
			WithArgs(
				userID,
				user.FirstName,
				user.LastName,
				user.Email,
				user.PhoneNumber,
				user.Password,
				user.UserType,
				user.IsActive,
				user.CreatedAt,
				user.UpdatedAt,
			).
			WillReturnError(&pq.Error{
				Code:    "23505",
				Message: "duplicate key value violates unique constraint \"users_phone_number_key\"",
			})

		id, err := repo.CreateUser(context.Background(), user)

		assert.Error(t, err)
		assert.Equal(t, uuid.Nil, id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("attempt to create user with empty email", func(t *testing.T) {
		userID := uuid.New()
		now := time.Now()
		user := &entities.User{
			UserID:      userID,
			FirstName:   "Dmitry",
			LastName:    "Dmitriev",
			Email:       "",
			PhoneNumber: "+79001234570",
			Password:    "hashedPassword",
			UserType:    "user",
			IsActive:    true,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		mock.ExpectExec("INSERT INTO users").
			WithArgs(
				userID,
				user.FirstName,
				user.LastName,
				user.Email,
				user.PhoneNumber,
				user.Password,
				user.UserType,
				user.IsActive,
				user.CreatedAt,
				user.UpdatedAt,
			).
			WillReturnError(&pq.Error{
				Code:    "23514",
				Message: "new row for relation \"users\" violates check constraint \"users_email_check\"",
			})

		id, err := repo.CreateUser(context.Background(), user)

		assert.Error(t, err)
		assert.Equal(t, uuid.Nil, id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("attempt to create user with empty password", func(t *testing.T) {
		userID := uuid.New()
		now := time.Now()
		user := &entities.User{
			UserID:      userID,
			FirstName:   "Konstantin",
			LastName:    "Konstantinov",
			Email:       "konstantin@example.com",
			PhoneNumber: "+79001234572",
			Password:    "",
			UserType:    "user",
			IsActive:    true,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		mock.ExpectExec("INSERT INTO users").
			WithArgs(
				userID,
				user.FirstName,
				user.LastName,
				user.Email,
				user.PhoneNumber,
				user.Password,
				user.UserType,
				user.IsActive,
				user.CreatedAt,
				user.UpdatedAt,
			).
			WillReturnError(&pq.Error{
				Code:    "23514",
				Message: "new row for relation \"users\" violates check constraint \"users_password_check\"",
			})

		id, err := repo.CreateUser(context.Background(), user)

		assert.Error(t, err)
		assert.Equal(t, uuid.Nil, id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("attempt to create user with invalid user type", func(t *testing.T) {
		userID := uuid.New()
		now := time.Now()
		user := &entities.User{
			UserID:      userID,
			FirstName:   "Maxim",
			LastName:    "Maximov",
			Email:       "maxim@example.com",
			PhoneNumber: "+79001234573",
			Password:    "hashedPassword",
			UserType:    "invalid_type",
			IsActive:    true,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		mock.ExpectExec("INSERT INTO users").
			WithArgs(
				userID,
				user.FirstName,
				user.LastName,
				user.Email,
				user.PhoneNumber,
				user.Password,
				user.UserType,
				user.IsActive,
				user.CreatedAt,
				user.UpdatedAt,
			).
			WillReturnError(&pq.Error{
				Code:    "23514",
				Message: "new row for relation \"users\" violates check constraint \"users_user_type_check\"",
			})

		id, err := repo.CreateUser(context.Background(), user)

		assert.Error(t, err)
		assert.Equal(t, uuid.Nil, id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("attempt to create user with empty first name", func(t *testing.T) {
		userID := uuid.New()
		now := time.Now()
		user := &entities.User{
			UserID:      userID,
			FirstName:   "",
			LastName:    "Petrov",
			Email:       "petr@example.com",
			PhoneNumber: "+79001234574",
			Password:    "hashedPassword",
			UserType:    "user",
			IsActive:    true,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		mock.ExpectExec("INSERT INTO users").
			WithArgs(
				userID,
				user.FirstName,
				user.LastName,
				user.Email,
				user.PhoneNumber,
				user.Password,
				user.UserType,
				user.IsActive,
				user.CreatedAt,
				user.UpdatedAt,
			).
			WillReturnError(&pq.Error{
				Code:    "23514",
				Message: "new row for relation \"users\" violates check constraint \"users_first_name_check\"",
			})

		id, err := repo.CreateUser(context.Background(), user)

		assert.Error(t, err)
		assert.Equal(t, uuid.Nil, id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("attempt to create user with empty last name", func(t *testing.T) {
		userID := uuid.New()
		now := time.Now()
		user := &entities.User{
			UserID:      userID,
			FirstName:   "Ivan",
			LastName:    "",
			Email:       "ivan2@example.com",
			PhoneNumber: "+79001234575",
			Password:    "hashedPassword",
			UserType:    "user",
			IsActive:    true,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		mock.ExpectExec("INSERT INTO users").
			WithArgs(
				userID,
				user.FirstName,
				user.LastName,
				user.Email,
				user.PhoneNumber,
				user.Password,
				user.UserType,
				user.IsActive,
				user.CreatedAt,
				user.UpdatedAt,
			).
			WillReturnError(&pq.Error{
				Code:    "23514",
				Message: "new row for relation \"users\" violates check constraint \"users_last_name_check\"",
			})

		id, err := repo.CreateUser(context.Background(), user)

		assert.Error(t, err)
		assert.Equal(t, uuid.Nil, id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("attempt to create user with empty phone number", func(t *testing.T) {
		userID := uuid.New()
		now := time.Now()
		user := &entities.User{
			UserID:      userID,
			FirstName:   "Nikolay",
			LastName:    "Nikolaev",
			Email:       "nikolay@example.com",
			PhoneNumber: "",
			Password:    "hashedPassword",
			UserType:    "user",
			IsActive:    true,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		mock.ExpectExec("INSERT INTO users").
			WithArgs(
				userID,
				user.FirstName,
				user.LastName,
				user.Email,
				user.PhoneNumber,
				user.Password,
				user.UserType,
				user.IsActive,
				user.CreatedAt,
				user.UpdatedAt,
			).
			WillReturnError(&pq.Error{
				Code:    "23514",
				Message: "new row for relation \"users\" violates check constraint \"users_phone_number_check\"",
			})

		id, err := repo.CreateUser(context.Background(), user)

		assert.Error(t, err)
		assert.Equal(t, uuid.Nil, id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database connection error", func(t *testing.T) {
		userID := uuid.New()
		now := time.Now()
		user := &entities.User{
			UserID:      userID,
			FirstName:   "Oleg",
			LastName:    "Olegov",
			Email:       "oleg@example.com",
			PhoneNumber: "+79001234576",
			Password:    "hashedPassword",
			UserType:    "user",
			IsActive:    true,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		mock.ExpectExec("INSERT INTO users").
			WithArgs(
				userID,
				user.FirstName,
				user.LastName,
				user.Email,
				user.PhoneNumber,
				user.Password,
				user.UserType,
				user.IsActive,
				user.CreatedAt,
				user.UpdatedAt,
			).
			WillReturnError(&pq.Error{
				Code:    "08006",
				Message: "connection to database failed",
			})

		id, err := repo.CreateUser(context.Background(), user)

		assert.Error(t, err)
		assert.Equal(t, uuid.Nil, id)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetUserByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewRepository(sqlxDB)

	t.Run("successfully get user by email", func(t *testing.T) {
		userID := uuid.New()
		now := time.Now()
		expectedUser := &entities.User{
			UserID:      userID,
			FirstName:   "Evgeny",
			LastName:    "Koveshnikov",
			Email:       "evgeny@example.com",
			PhoneNumber: "+79001234567",
			Password:    "hashedPassword",
			UserType:    "user",
			IsActive:    true,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		rows := sqlmock.NewRows([]string{
			"user_id", "first_name", "last_name", "email", "phone_number",
			"password", "user_type", "is_active", "created_at", "updated_at",
		}).AddRow(
			expectedUser.UserID,
			expectedUser.FirstName,
			expectedUser.LastName,
			expectedUser.Email,
			expectedUser.PhoneNumber,
			expectedUser.Password,
			expectedUser.UserType,
			expectedUser.IsActive,
			expectedUser.CreatedAt,
			expectedUser.UpdatedAt,
		)

		mock.ExpectQuery("SELECT \\* FROM users WHERE email = \\$1 AND deleted_at IS NULL").
			WithArgs(expectedUser.Email).
			WillReturnRows(rows)

		user, err := repo.GetUserByEmail(context.Background(), expectedUser.Email)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("user not found", func(t *testing.T) {
		email := "nonexistent@example.com"

		mock.ExpectQuery("SELECT \\* FROM users WHERE email = \\$1 AND deleted_at IS NULL").
			WithArgs(email).
			WillReturnError(sql.ErrNoRows)

		user, err := repo.GetUserByEmail(context.Background(), email)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("user is deleted", func(t *testing.T) {
		email := "deleted@example.com"

		mock.ExpectQuery("SELECT \\* FROM users WHERE email = \\$1 AND deleted_at IS NULL").
			WithArgs(email).
			WillReturnError(sql.ErrNoRows)

		user, err := repo.GetUserByEmail(context.Background(), email)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database connection error", func(t *testing.T) {
		email := "test@example.com"

		mock.ExpectQuery("SELECT \\* FROM users WHERE email = \\$1 AND deleted_at IS NULL").
			WithArgs(email).
			WillReturnError(&pq.Error{
				Code:    "08006",
				Message: "connection to database failed",
			})

		user, err := repo.GetUserByEmail(context.Background(), email)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty email", func(t *testing.T) {
		email := ""

		mock.ExpectQuery("SELECT \\* FROM users WHERE email = \\$1 AND deleted_at IS NULL").
			WithArgs(email).
			WillReturnError(&pq.Error{
				Code:    "23514",
				Message: "new row for relation \"users\" violates check constraint \"users_email_check\"",
			})

		user, err := repo.GetUserByEmail(context.Background(), email)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestIsUserExist(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := NewRepository(sqlxDB)

	t.Run("user exists by email", func(t *testing.T) {
		email := "evgeny@example.com"
		phoneNumber := "+79001234567"

		rows := sqlmock.NewRows([]string{"count"}).AddRow(1)

		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users WHERE email = \\$1 OR phone_number = \\$2 AND deleted_at IS NULL").
			WithArgs(email, phoneNumber).
			WillReturnRows(rows)

		exists, err := repo.IsUserExist(context.Background(), email, phoneNumber)

		assert.NoError(t, err)
		assert.True(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("user exists by phone number", func(t *testing.T) {
		email := "nonexistent@example.com"
		phoneNumber := "+79001234567"

		rows := sqlmock.NewRows([]string{"count"}).AddRow(1)

		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users WHERE email = \\$1 OR phone_number = \\$2 AND deleted_at IS NULL").
			WithArgs(email, phoneNumber).
			WillReturnRows(rows)

		exists, err := repo.IsUserExist(context.Background(), email, phoneNumber)

		assert.NoError(t, err)
		assert.True(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("user exists by both email and phone number", func(t *testing.T) {
		email := "evgeny@example.com"
		phoneNumber := "+79001234567"

		rows := sqlmock.NewRows([]string{"count"}).AddRow(1)

		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users WHERE email = \\$1 OR phone_number = \\$2 AND deleted_at IS NULL").
			WithArgs(email, phoneNumber).
			WillReturnRows(rows)

		exists, err := repo.IsUserExist(context.Background(), email, phoneNumber)

		assert.NoError(t, err)
		assert.True(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("user does not exist", func(t *testing.T) {
		email := "nonexistent@example.com"
		phoneNumber := "+79001234568"

		rows := sqlmock.NewRows([]string{"count"}).AddRow(0)

		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users WHERE email = \\$1 OR phone_number = \\$2 AND deleted_at IS NULL").
			WithArgs(email, phoneNumber).
			WillReturnRows(rows)

		exists, err := repo.IsUserExist(context.Background(), email, phoneNumber)

		assert.NoError(t, err)
		assert.False(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("user is deleted", func(t *testing.T) {
		email := "deleted@example.com"
		phoneNumber := "+79001234569"

		rows := sqlmock.NewRows([]string{"count"}).AddRow(0)

		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users WHERE email = \\$1 OR phone_number = \\$2 AND deleted_at IS NULL").
			WithArgs(email, phoneNumber).
			WillReturnRows(rows)

		exists, err := repo.IsUserExist(context.Background(), email, phoneNumber)

		assert.NoError(t, err)
		assert.False(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database connection error", func(t *testing.T) {
		email := "test@example.com"
		phoneNumber := "+79001234570"

		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users WHERE email = \\$1 OR phone_number = \\$2 AND deleted_at IS NULL").
			WithArgs(email, phoneNumber).
			WillReturnError(&pq.Error{
				Code:    "08006",
				Message: "connection to database failed",
			})

		exists, err := repo.IsUserExist(context.Background(), email, phoneNumber)

		assert.Error(t, err)
		assert.False(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("empty parameters", func(t *testing.T) {
		email := ""
		phoneNumber := ""

		rows := sqlmock.NewRows([]string{"count"}).AddRow(0)

		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM users WHERE email = \\$1 OR phone_number = \\$2 AND deleted_at IS NULL").
			WithArgs(email, phoneNumber).
			WillReturnRows(rows)

		exists, err := repo.IsUserExist(context.Background(), email, phoneNumber)

		assert.NoError(t, err)
		assert.False(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
