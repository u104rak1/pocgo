package repository_test

import (
	"database/sql"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	authDomain "github.com/u104rak1/pocgo/internal/domain/authentication"
	idVO "github.com/u104rak1/pocgo/internal/domain/value_object/id"
	"github.com/u104rak1/pocgo/internal/infrastructure/postgres/repository"
)

func TestAuthenticationRepository_Save(t *testing.T) {
	repo, mock, ctx, _ := PrepareTestRepository(t, repository.NewAuthenticationRepository)
	userID := idVO.NewUserIDForTest("user_id_1")
	auth, err := authDomain.New(userID, "password123")
	assert.NoError(t, err)

	expectQuery := fmt.Sprintf(`
		INSERT INTO "authentications" AS "authentication" ("user_id", "password_hash", "deleted_at")
		VALUES ('%s', '%s', DEFAULT)
		ON CONFLICT (user_id) DO UPDATE SET password_hash = EXCLUDED.password_hash
		RETURNING "deleted_at"
	`, auth.UserIDString(), auth.PasswordHash())

	tests := []struct {
		caseName string
		prepare  func()
		wantErr  bool
	}{
		{
			caseName: "Positive: 認証情報の保存が成功する",
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
			err := repo.Save(ctx, auth)

			if tt.wantErr {
				assert.Error(t, assert.AnError)
			} else {
				assert.NoError(t, err)
			}
			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestAuthenticationRepository_FindByUserID(t *testing.T) {
	repo, mock, ctx, _ := PrepareTestRepository(t, repository.NewAuthenticationRepository)
	userID := idVO.NewUserIDForTest("user_id_1")
	auth, err := authDomain.New(userID, "password123")
	assert.NoError(t, err)

	expectQuery := fmt.Sprintf(`
		SELECT "authentication"."user_id", "authentication"."password_hash", "authentication"."deleted_at"
		FROM "authentications" AS "authentication"
		WHERE (user_id = '%s') AND "authentication"."deleted_at" IS NULL
	`, auth.UserIDString())

	tests := []struct {
		caseName string
		prepare  func()
		wantAuth *authDomain.Authentication
		wantErr  bool
	}{
		{
			caseName: "Positive: ユーザーIDで認証情報取得が成功する",
			prepare: func() {
				rows := sqlmock.NewRows([]string{"user_id", "password_hash", "deleted_at"}).
					AddRow(auth.UserIDString(), auth.PasswordHash(), nil)
				mock.ExpectQuery(regexp.QuoteMeta(expectQuery)).WillReturnRows(rows)
			},
			wantAuth: auth,
			wantErr:  false,
		},
		{
			caseName: "Positive: 認証情報が見つからない場合、nilを返す",
			prepare: func() {
				mock.ExpectQuery(regexp.QuoteMeta(expectQuery)).WillReturnError(sql.ErrNoRows)
			},
			wantAuth: nil,
			wantErr:  false,
		},
		{
			caseName: "Negative: SQLエラーで失敗する",
			prepare: func() {
				mock.ExpectQuery(regexp.QuoteMeta(expectQuery)).WillReturnError(assert.AnError)
			},
			wantAuth: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			tt.prepare()
			foundAuth, err := repo.FindByUserID(ctx, userID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, foundAuth)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantAuth, foundAuth)
			}
			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestAuthenticationRepository_ExistsByUserID(t *testing.T) {
	repo, mock, ctx, _ := PrepareTestRepository(t, repository.NewAuthenticationRepository)
	userID := idVO.NewUserIDForTest("user_id_1")

	expectQuery := fmt.Sprintf(`
		SELECT EXISTS
			(SELECT "authentication"."user_id", "authentication"."password_hash", "authentication"."deleted_at"
			FROM "authentications" AS "authentication"
			WHERE (user_id = '%s') AND "authentication"."deleted_at" IS NULL)
	`, userID.String())

	tests := []struct {
		caseName   string
		prepare    func()
		wantExists bool
		wantErr    bool
	}{
		{
			caseName: "Positive: ユーザーIDで認証情報の存在確認が成功する",
			prepare: func() {
				mock.ExpectQuery(regexp.QuoteMeta(expectQuery)).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
			},
			wantExists: true,
			wantErr:    false,
		},
		{
			caseName: "Positive: 認証情報が存在しない場合、falseを返す",
			prepare: func() {
				mock.ExpectQuery(regexp.QuoteMeta(expectQuery)).WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))
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
			exists, err := repo.ExistsByUserID(ctx, userID)

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
