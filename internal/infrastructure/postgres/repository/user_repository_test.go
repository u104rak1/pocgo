package repository_test

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	userDomain "github.com/u104rak1/pocgo/internal/domain/user"
	"github.com/u104rak1/pocgo/internal/infrastructure/postgres/repository"
)

func TestUserRepository_Save(t *testing.T) {
	repo, mock, ctx, db := PrepareTestRepository(t, repository.NewUserRepository)
	user, err := userDomain.New("sato taro", "sato@example.com")
	assert.NoError(t, err)

	expectQuery := fmt.Sprintf(`
		INSERT INTO "users" AS "user" ("id", "name", "email", "deleted_at")
		VALUES ('%s', '%s', '%s', DEFAULT)
		ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name, email = EXCLUDED.email
		RETURNING "deleted_at"
	`, user.IDString(), user.Name(), user.Email())

	tests := []struct {
		caseName string
		prepare  func()
		wantErr  bool
	}{
		{
			caseName: "Positive: ユーザー保存が成功する",
			prepare: func() {
				mock.ExpectQuery(regexp.QuoteMeta(expectQuery)).WillReturnRows(sqlmock.NewRows([]string{"deleted_at"}))
			},
			wantErr: false,
		},
		{
			caseName: "Negative: SQLエラーで失敗する",
			prepare: func() {
				mock.ExpectQuery(regexp.QuoteMeta(expectQuery)).WillReturnError(assert.AnError)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			tt.prepare()
			err := repo.Save(ctx, user)

			if tt.wantErr {
				assert.Error(t, assert.AnError)
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
		wantErr  bool
	}{
		{
			caseName: "Positive: トランザクション内でユーザー保存が成功する",
			prepare: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(expectQuery)).WillReturnRows(sqlmock.NewRows([]string{"deleted_at"}))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			caseName: "Negative: SQLエラーでロールバックされる",
			prepare: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(expectQuery)).WillReturnError(assert.AnError)
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range testsWithTx {
		t.Run(tt.caseName, func(t *testing.T) {
			tt.prepare()
			err := unitOfWork.RunInTx(ctx, func(ctx context.Context) error {
				return repo.Save(ctx, user)
			})

			if tt.wantErr {
				assert.Error(t, err)
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
	user, err := userDomain.New("sato taro", "sato@example.com")
	assert.NoError(t, err)

	expectQuery := fmt.Sprintf(`
		SELECT "user"."id", "user"."name", "user"."email", "user"."deleted_at"
		FROM "users" AS "user"
		WHERE (id = '%s') AND "user"."deleted_at" IS NULL
	`, user.IDString())

	tests := []struct {
		caseName string
		prepare  func()
		wantUser *userDomain.User
		wantErr  bool
	}{
		{
			caseName: "Positive: IDでユーザー取得が成功する",
			prepare: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "deleted_at"}).
					AddRow(user.IDString(), user.Name(), user.Email(), nil)
				mock.ExpectQuery(regexp.QuoteMeta(expectQuery)).WillReturnRows(rows)
			},
			wantUser: user,
			wantErr:  false,
		},
		{
			caseName: "Positive: ユーザーが見つからない場合、nilを返す",
			prepare: func() {
				mock.ExpectQuery(regexp.QuoteMeta(expectQuery)).WillReturnError(sql.ErrNoRows)
			},
			wantUser: nil,
			wantErr:  false,
		},
		{
			caseName: "Negative: SQLエラーで失敗する",
			prepare: func() {
				mock.ExpectQuery(regexp.QuoteMeta(expectQuery)).WillReturnError(assert.AnError)
			},
			wantUser: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			tt.prepare()
			foundUser, err := repo.FindByID(ctx, user.ID())

			if tt.wantErr {
				assert.Error(t, err)
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
	user, err := userDomain.New("sato taro", "sato@example.com")
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
		wantErr  bool
	}{
		{
			caseName: "Positive: メールアドレスでユーザー取得が成功する",
			prepare: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "email", "deleted_at"}).
					AddRow(user.IDString(), user.Name(), user.Email(), nil)
				mock.ExpectQuery(regexp.QuoteMeta(expectQuery)).WillReturnRows(rows)
			},
			wantUser: user,
			wantErr:  false,
		},
		{
			caseName: "Positive: ユーザーが見つからない場合、nilを返す",
			prepare: func() {
				mock.ExpectQuery(regexp.QuoteMeta(expectQuery)).WillReturnError(sql.ErrNoRows)
			},
			wantUser: nil,
			wantErr:  false,
		},
		{
			caseName: "Negative: SQLエラーで失敗する",
			prepare: func() {
				mock.ExpectQuery(regexp.QuoteMeta(expectQuery)).WillReturnError(assert.AnError)
			},
			wantUser: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			tt.prepare()
			foundUser, err := repo.FindByEmail(ctx, user.Email())

			if tt.wantErr {
				assert.Error(t, err)
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

func TestUserRepository_ExistsByID(t *testing.T) {
	repo, mock, ctx, _ := PrepareTestRepository(t, repository.NewUserRepository)
	user, err := userDomain.New("sato taro", "sato@example.com")
	assert.NoError(t, err)

	expectQuery := fmt.Sprintf(`
		SELECT EXISTS
			(SELECT "user"."id", "user"."name", "user"."email", "user"."deleted_at"
			FROM "users" AS "user"
			WHERE (id = '%s') AND "user"."deleted_at" IS NULL)
	`, user.IDString())

	tests := []struct {
		caseName   string
		prepare    func()
		wantExists bool
		wantErr    bool
	}{
		{
			caseName: "Positive: IDでユーザーの存在確認が成功する",
			prepare: func() {
				mock.ExpectQuery(regexp.QuoteMeta(expectQuery)).WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow(1))
			},
			wantExists: true,
			wantErr:    false,
		},
		{
			caseName: "Positive: ユーザーが存在しない場合、falseを返す",
			prepare: func() {
				mock.ExpectQuery(regexp.QuoteMeta(expectQuery)).WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow(0))
			},
			wantExists: false,
			wantErr:    false,
		},
		{
			caseName: "Negative: SQLエラーで失敗する",
			prepare: func() {
				mock.ExpectQuery(regexp.QuoteMeta(expectQuery)).WillReturnError(assert.AnError)
			},
			wantExists: false,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			tt.prepare()
			exists, err := repo.ExistsByID(ctx, user.ID())

			if tt.wantErr {
				assert.Error(t, err)
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
		wantErr    bool
	}{
		{
			caseName: "Positive: メールアドレスでユーザーの存在確認が成功する",
			prepare: func() {
				mock.ExpectQuery(regexp.QuoteMeta(expectQuery)).WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow(1))
			},
			wantExists: true,
			wantErr:    false,
		},
		{
			caseName: "Positive: ユーザーが存在しない場合、falseを返す",
			prepare: func() {
				mock.ExpectQuery(regexp.QuoteMeta(expectQuery)).WillReturnRows(sqlmock.NewRows([]string{"1"}).AddRow(0))
			},
			wantExists: false,
			wantErr:    false,
		},
		{
			caseName: "Negative: SQLエラーで失敗する",
			prepare: func() {
				mock.ExpectQuery(regexp.QuoteMeta(expectQuery)).WillReturnError(assert.AnError)
			},
			wantExists: false,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			tt.prepare()
			exists, err := repo.ExistsByEmail(ctx, email)

			if tt.wantErr {
				assert.Error(t, err)
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
