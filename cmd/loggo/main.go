package main

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "loggo",
	Short: "Loggo - a lightweight and easy-to-setup log aggregation system",
}

func main() {
	// All error handling and printing is handled in subcommands and cobra itself
	_ = rootCmd.Execute()
}
