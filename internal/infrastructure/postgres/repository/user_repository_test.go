package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	userDomain "github.com/ucho456job/pocgo/internal/domain/user"
	"github.com/ucho456job/pocgo/internal/infrastructure/postgres/repository"
	"github.com/ucho456job/pocgo/pkg/ulid"
)

func TestUserRepository_Save(t *testing.T) {
	repo, mock, ctx, db := PrepareTestRepository(t, repository.NewUserRepository)
	userID := ulid.GenerateStaticULID("user")
	user, err := userDomain.New(userID, "sato taro", "sato@example.com")
	assert.NoError(t, err)

	expectQuery := fmt.Sprintf(`
		INSERT INTO "users" AS "user" ("id", "name", "email", "deleted_at")
		VALUES ('%s', '%s', '%s', DEFAULT)
		ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, email = EXCLUDED.email
		RETURNING "deleted_at"
	`, user.ID(), user.Name(), user.Email())

	ErrDB := errors.New("database error")

	tests := []struct {
		caseName string
		prepare  func()
		wantErr  error
	}{
		{
			caseName: "Successfully saves user.",
			prepare: func() {
				mock.ExpectQuery(regexp.QuoteMeta(expectQuery)).WillReturnRows(sqlmock.NewRows([]string{"deleted_at"}))
			},
			wantErr: nil,
		},
		{
			caseName: "Error occurs in database during save operation.",
			prepare: func() {
				mock.ExpectQuery(regexp.QuoteMeta(expectQuery)).WillReturnError(ErrDB)
			},
			wantErr: ErrDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			tt.prepare()
			err := repo.Save(ctx, user)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, ErrDB)
			} else {
				assert.NoError(t, err)
			}
			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}

	unitOfWork := repository.NewUnitOfWork(db)
	testsWithTx := []struct {
		caseName string
		prepare  func()
		wantErr  error
	}{
		{
			caseName: "Successfully saves user with transaction.",
			prepare: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(expectQuery)).WillReturnRows(sqlmock.NewRows([]string{"deleted_at"}))
				mock.ExpectCommit()
			},
			wantErr: nil,
		},
		{
			caseName: "Error occurs in database, transaction is rolled back.",
			prepare: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(expectQuery)).WillReturnError(ErrDB)
				mock.ExpectRollback()
			},
			wantErr: ErrDB,
		},
	}

	for _, tt := range testsWithTx {
		t.Run(tt.caseName, func(t *testing.T) {
			tt.prepare()
			err := unitOfWork.RunInTx(ctx, func(ctx context.Context) error {
				return repo.Save(ctx, user)
			})

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestUserRepository_FindByID(t *testing.T) {
	repo, mock, ctx, _ := PrepareTestRepository(t, repository.NewUserRepository)
	userID := ulid.GenerateStaticULID("user")
	user, err := userDomain.New(userID, "sato taro", "sato@example.com")
	assert.NoError(t, err)

	expectQuery := fmt.Sprintf(`
		SELECT "user"."id", "user"."name", "user"."email", "user"."deleted_at"
		FROM "users" AS "user"
		WHERE (id = '%s') AND "user"."deleted_at" IS NULL
	`, userID)

	tests := []struct {
		caseName string
		prepare  func()
		wantUser *userDomain.User
		wantErr  error
	}{
		{
			caseName: "Successfully finds user by ID.",
			prepare: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "deleted_at"}).
					AddRow(userID, user.Name(), user.Email(), nil)
				mock.ExpectQuery(regexp.QuoteMeta(expectQuery)).WillReturnRows(rows)
			},
			wantUser: user,
			wantErr:  nil,
		},
		{
			caseName: "Returns ErrUserNotFound when no user is found.",
			prepare: func() {
				mock.ExpectQuery(regexp.QuoteMeta(expectQuery)).WillReturnError(sql.ErrNoRows)
			},
			wantUser: nil,
			wantErr:  userDomain.ErrUserNotFound,
		},
		{
			caseName: "Returns database error when unknown error occurs.",
			prepare: func() {
				mock.ExpectQuery(regexp.QuoteMeta(expectQuery)).WillReturnError(ErrDB)
			},
			wantUser: nil,
			wantErr:  ErrDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			tt.prepare()
			foundUser, err := repo.FindByID(ctx, userID)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, foundUser)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantUser, foundUser)
			}
			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestUserRepository_FindByEmail(t *testing.T) {
	repo, mock, ctx, _ := PrepareTestRepository(t, repository.NewUserRepository)
	userID := ulid.GenerateStaticULID("user")
	user, err := userDomain.New(userID, "sato taro", "sato@example.com")
	assert.NoError(t, err)

	expectQuery := fmt.Sprintf(`
		SELECT "user"."id", "user"."name", "user"."email", "user"."deleted_at"
		FROM "users" AS "user"
		WHERE (email = '%s') AND "user"."deleted_at" IS NULL
	`, user.Email())

	tests := []struct {
		caseName string
		prepare  func()
		wantUser *userDomain.User
		wantErr  error
	}{
		{
			caseName: "Successfully finds user by email.",
			prepare: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "deleted_at"}).
					AddRow(userID, user.Name(), user.Email(), nil)
				mock.ExpectQuery(regexp.QuoteMeta(expectQuery)).WillReturnRows(rows)
			},
			wantUser: user,
			wantErr:  nil,
		},
		{
			caseName: "Returns ErrUserNotFound when no user is found.",
			prepare: func() {
				mock.ExpectQuery(regexp.QuoteMeta(expectQuery)).WillReturnError(sql.ErrNoRows)
			},
			wantUser: nil,
			wantErr:  userDomain.ErrUserNotFound,
		},
		{
			caseName: "Returns database error when unknown error occurs.",
			prepare: func() {
				mock.ExpectQuery(regexp.QuoteMeta(expectQuery)).WillReturnError(ErrDB)
			},
			wantUser: nil,
			wantErr:  ErrDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			tt.prepare()
			foundUser, err := repo.FindByEmail(ctx, user.Email())

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, foundUser)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantUser, foundUser)
			}
			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestUserRepository_ExistsByEmail(t *testing.T) {
	repo, mock, ctx, _ := PrepareTestRepository(t, repository.NewUserRepository)
	email := "sato@example.com"

	expectQuery := fmt.Sprintf(`
		SELECT EXISTS
			(SELECT "user"."id", "user"."name", "user"."email", "user"."deleted_at"
			FROM "users" AS "user"
			WHERE (email = '%s') AND "user"."deleted_at" IS NULL)
	`, email)

	tests := []struct {
		caseName   string
		prepare    func()
		wantExists bool
		wantErr    error
	}{
		{
			caseName: "Email exists in the database.",
			prepare: func() {
				mock.ExpectQuery(regexp.QuoteMeta(expectQuery)).WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow(1))
			},
			wantExists: true,
			wantErr:    nil,
		},
		{
			caseName: "Email does not exist in the database.",
			prepare: func() {
				mock.ExpectQuery(regexp.QuoteMeta(expectQuery)).WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow(0))
			},
			wantExists: false,
			wantErr:    nil,
		},
		{
			caseName: "Error occurs in database when checking email existence.",
			prepare: func() {
				mock.ExpectQuery(regexp.QuoteMeta(expectQuery)).WillReturnError(ErrDB)
			},
			wantExists: false,
			wantErr:    ErrDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			tt.prepare()
			exists, err := repo.ExistsByEmail(ctx, email)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.False(t, exists)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantExists, exists)
			}
			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

// TODO
// func TestUserRepository_Delete(t *testing.T) {
// 	repo, mock, ctx, _ := PrepareTestRepository(t, repository.NewUserRepository)
// 	userID := ulid.GenerateStaticULID("user")

// 	expectQuery := fmt.Sprintf(`
// 		UPDATE "users" AS "user"
// 		SET "deleted_at" = '%s'
// 		WHERE "user"."deleted_at" IS NULL AND ("id" = '%s')
// 		`, sqlmock.AnyArg(), userID)
// 	// expectQuery := `
// 	// UPDATE "users" AS "user"
// 	// SET "deleted_at" = ?
// 	// WHERE "user"."deleted_at" IS NULL AND ("id" = ?)
// 	// `

// 	tests := []struct {
// 		caseName string
// 		prepare  func()
// 		wantErr  error
// 	}{
// 		{
// 			caseName: "Happy path: delete user by ID.",
// 			prepare: func() {
// 				mock.ExpectQuery(regexp.QuoteMeta(expectQuery)).WillReturnRows(sqlmock.NewRows([]string{"deleted_at"}))
// 			},
// 			wantErr: nil,
// 		},
// 		// {
// 		// 	caseName: "Error case: database error during delete.",
// 		// 	prepare: func() {
// 		// 		mock.ExpectExec(regexp.QuoteMeta(expectQuery)).
// 		// 			WithArgs(deletedAt, userID).
// 		// 			WillReturnError(ErrDB)
// 		// 	},
// 		// 	wantErr: ErrDB,
// 		// },
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.caseName, func(t *testing.T) {
// 			tt.prepare()
// 			err := repo.Delete(ctx, userID)

// 			if tt.wantErr != nil {
// 				assert.ErrorIs(t, err, tt.wantErr)
// 			} else {
// 				assert.NoError(t, err)
// 			}
// 			err = mock.ExpectationsWereMet()
// 			assert.NoError(t, err)
// 		})
// 	}
// }
