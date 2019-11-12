package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

// ciCmd represents the artifactory command
var ciCmd = &cobra.Command{
	Use:   "ci",
	Short: "interact with the ci api",
	//Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ci called")
	},
}

func init() {
	rootCmd.AddCommand(ciCmd)

}



