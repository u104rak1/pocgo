package seed

import (
	"context"

	"github.com/ucho456job/pocgo/internal/infrastructure/postgres/model"
	"github.com/uptrace/bun"
)

const (
	johnDoeID   = "01J9R83RCMVQ1FJK60P0BS23T3"
	janeSmithID = "01J9R844GCZZK02ZW76J5Q32M8"
)

func SaveUser(db *bun.DB) error {
	data := []model.UserModel{
		{ID: johnDoeID, Name: "John Doe", Email: "john@example.com"},
		{ID: janeSmithID, Name: "Jane Smith", Email: "jane@example.com"},
	}
	if _, err := db.NewInsert().Model(&data).Exec(context.Background()); err != nil {
		return err
	}
	return nil
}
