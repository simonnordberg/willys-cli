package cmd

import (
	"fmt"

	"github.com/simonnordberg/willys-cli/internal/willys"
	"github.com/spf13/cobra"
)

func loginCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "Log in and save session",
		RunE: func(cmd *cobra.Command, args []string) error {
			username, password, err := GetCreds()
			if err != nil {
				return err
			}
			c := willys.NewClient()
			cust, err := c.Login(username, password)
			if err != nil {
				return err
			}
			fmt.Printf("Logged in as %s %s (%s)\n", cust.FirstName, cust.LastName, cust.Email)
			return nil
		},
	}
}

func logoutCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Clear saved session",
		RunE: func(cmd *cobra.Command, args []string) error {
			willys.ClearSession()
			fmt.Println("Session cleared.")
			return nil
		},
	}
}

func statusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Check login status",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := willys.NewClient()
			if c.IsLoggedIn() {
				cust, err := c.GetCustomer()
				if err != nil {
					return err
				}
				fmt.Printf("Logged in as %s %s (%s)\n", cust.FirstName, cust.LastName, cust.Email)
			} else {
				fmt.Println("Not logged in.")
			}
			return nil
		},
	}
}
