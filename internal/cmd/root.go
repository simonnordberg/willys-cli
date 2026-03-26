package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/simonnordberg/willys-cli/internal/env"
	"github.com/simonnordberg/willys-cli/internal/tui"
	"github.com/simonnordberg/willys-cli/internal/willys"
	"github.com/spf13/cobra"
)

var batch string

func Execute() {
	root := &cobra.Command{
		Use:   "willys",
		Short: "Willys.se grocery store CLI",
		Long:  "Search products, browse categories, and manage your shopping cart at Willys.se.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			_ = env.Load(".env")
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if batch != "" {
				return RunBatch(batch)
			}
			c, err := GetClient()
			if err != nil {
				return err
			}
			return tui.Run(c)
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	root.Flags().StringVarP(&batch, "batch", "i", "", "CSV file with batch operations")

	root.AddCommand(loginCmd(), logoutCmd(), statusCmd())
	root.AddCommand(searchCmd(), categoriesCmd(), browseCmd())
	root.AddCommand(cartCmd())
	root.AddCommand(ordersCmd())
	root.AddCommand(diffCmd())

	if err := root.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func GetCreds() (string, string, error) {
	username := strings.Trim(os.Getenv("WILLYS_USERNAME"), `"`)
	password := strings.Trim(os.Getenv("WILLYS_PASSWORD"), `"`)
	if username == "" || password == "" {
		return "", "", fmt.Errorf("credentials required: set WILLYS_USERNAME and WILLYS_PASSWORD env vars or use a .env file")
	}
	return username, password, nil
}

func GetClient() (*willys.Client, error) {
	c := willys.NewClient()
	if c.IsLoggedIn() {
		return c, nil
	}
	username, password, err := GetCreds()
	if err != nil {
		return nil, err
	}
	if _, err := c.Login(username, password); err != nil {
		return nil, err
	}
	return c, nil
}
