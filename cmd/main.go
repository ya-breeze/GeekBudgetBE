package main

import (
	"github.com/google/uuid"

	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/database/models"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

func main() {
	db, err := database.OpenSqlite()
	if err != nil {
		panic("failed to connect database")
	}
	if err := database.Migrate(db); err != nil {
		panic("failed to migrate database")
	}

	userId := "123e4567-e89b-12d3-a456-426614174000"
	account := models.Account{
		ID:     uuid.New(),
		UserId: userId,
		AccountNoId: goserver.AccountNoId{
			Name:        "Test Account",
			Type:        "CHECKING",
			Description: "Test Account Description",
		},
	}
	db.Create(&account)

	result, err := db.Model(&models.Account{}).Where("user_id = ?", userId).Rows()
	if err != nil {
		panic(err)
	}
	defer result.Close()

	for result.Next() {
		var account models.Account
		if err := db.ScanRows(result, &account); err != nil {
			panic(err)
		}
		println("Account:")
		println(account.ID.String())
		println(account.UserId)
		println(account.Name)
		println(account.Type)
		println(account.Description)
		println()
	}
}
