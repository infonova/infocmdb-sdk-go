package main

import (
	"github.com/spf13/cobra"
)

// ciCmd represents the artifactory command
var ciListCmd = &cobra.Command{
	Use:   "list",
	Short: "interact with the ci list api",
	//Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		//i, err := v1.NewCMDB("localhost.yml")


	},
}

func init() {
	ciCmd.AddCommand(ciListCmd)
}



