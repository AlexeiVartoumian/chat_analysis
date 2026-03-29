package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:     "add",
	Aliases: append([]string{"addition"}), // can call command in other way
	Short:   "Add 2 numbers",
	Long:    "bla blabla",
	Args:    cobra.ExactArgs(2), // enforces two args
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Addition of %s and %s = %s.\n\n", args[0], args[1], Add(args[0], args[1]))
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
