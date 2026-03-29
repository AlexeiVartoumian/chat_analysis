package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{ // main entrypoint of commandtree

	Use:   "zero", // how to use the commnad
	Short: "zero is  a cli tool for basic math ops",
	Long:  "zero is a cli tool for blablabla",
	Run: func(cmd *cobra.Command, args []string) {

	}, // the function to execute
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "forambula error '%s' ", err)
		os.Exit(1)
	}
}
