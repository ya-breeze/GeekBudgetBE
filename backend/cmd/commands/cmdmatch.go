//nolint:forbidigo // it's okay to use fmt in this file
package commands

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/common"
)

func CmdMatch(log *slog.Logger) *cobra.Command {
	var username, matcherID, transactionID string

	res := &cobra.Command{
		Use:   "match",
		Short: "Match transaction against matcher",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, logger, err := createConfigAndLogger(cmd)
			if err != nil {
				return err
			}

			storage := database.NewStorage(logger, cfg)
			if err = storage.Open(); err != nil {
				return fmt.Errorf("failed to open storage: %w", err)
			}

			userID, err := storage.GetUserID(username)
			if err != nil {
				return fmt.Errorf("failed to get user ID by username %q: %w", username, err)
			}
			matcher, err := storage.GetMatcherRuntime(userID, matcherID)
			if err != nil {
				return fmt.Errorf("failed to get matcher %q: %w", matcherID, err)
			}
			transaction, err := storage.GetTransaction(userID, transactionID)
			if err != nil {
				return fmt.Errorf("failed to get transaction %q: %w", transactionID, err)
			}

			status := common.Match(&matcher, &transaction)
			statusStr := "unknown"
			switch status {
			case common.MatchResultSuccess:
				statusStr = "success"
			case common.MatchResultWrongDescription:
				statusStr = fmt.Sprintf("wrong description. Regexp %q doesn't match %q",
					matcher.DescriptionRegexp.String(), transaction.Description)
			case common.MatchResultWrongPartnerAccount:
				statusStr = fmt.Sprintf("wrong partner account. Regexp %q doesn't match %q",
					matcher.PartnerAccountRegexp.String(), transaction.PartnerAccount)
			}

			fmt.Println("Match result:", statusStr)
			return nil
		},
		Args: cobra.NoArgs,
	}

	res.Flags().StringVarP(&username, "username", "u", "", "username")
	res.Flags().StringVarP(&matcherID, "matcher-id", "m", "", "matcher ID")
	res.Flags().StringVarP(&transactionID, "transaction-id", "t", "", "transaction ID")

	_ = res.MarkFlagRequired("username")
	_ = res.MarkFlagRequired("matcher-id")
	_ = res.MarkFlagRequired("transaction-id")

	return res
}

// func NewUserAdd(log *slog.Logger) *cobra.Command {
// 	res := &cobra.Command{
// 		Use:   "add",
// 		Short: "Add a new user",
// 		RunE: func(_ *cobra.Command, _ []string) error {
// 			var username string
// 			fmt.Print("Enter username: ")
// 			_, err := fmt.Scanln(&username)
// 			if err != nil {
// 				return fmt.Errorf("error reading username: %w", err)
// 			}

// 			fmt.Print("Enter password: ")
// 			password, err := gopass.GetPasswd()
// 			if err != nil {
// 				return fmt.Errorf("error reading password: %w", err)
// 			}

// 			hashed, err := auth.HashPassword(password)
// 			if err != nil {
// 				return fmt.Errorf("error hashing password: %w", err)
// 			}
// 			fmt.Printf("Use '%s:%s' to add a user\n", username, base64.StdEncoding.EncodeToString(hashed))

// 			return nil
// 		},
// 	}

// 	return res
// }
