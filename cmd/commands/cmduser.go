package commands

import (
	"encoding/base64"
	"fmt"

	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
	"github.com/ya-breeze/geekbudgetbe/pkg/auth"
)

func CmdUser() *cobra.Command {
	res := &cobra.Command{
		Use:   "user",
		Short: "Control users",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	res.AddCommand(NewUserAdd())

	return res
}

func NewUserAdd() *cobra.Command {
	res := &cobra.Command{
		Use:   "add",
		Short: "Add a new user",
		RunE: func(cmd *cobra.Command, args []string) error {
			var username string
			fmt.Print("Enter username: ")
			fmt.Scanln(&username)

			fmt.Print("Enter password: ")
			password, err := gopass.GetPasswd()
			if err != nil {
				return fmt.Errorf("Error reading password: %w", err)
			}

			hashed, err := auth.HashPassword(password)
			if err != nil {
				return fmt.Errorf("Error hashing password: %w", err)
			}
			fmt.Printf("Use '%s:%s' to add a user\n", username, base64.StdEncoding.EncodeToString(hashed))

			return nil
		},
	}

	return res
}
