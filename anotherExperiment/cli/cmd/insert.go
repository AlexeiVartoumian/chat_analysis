package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var InsertCmd = &cobra.Command{
	Use: "insert",
	//Aliases: append([]string{"addition"}), // can call command in other way
	Short: "insert into db",
	Long:  "ha",
	Args:  cobra.ExactArgs(2), // enforces two args
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("inserting csv %s into table %s.\n\n", args[0], args[1])
		CsvFile(args[0], args[1])
	},
}

func init() {
	rootCmd.AddCommand(InsertCmd)
}
