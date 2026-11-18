/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"netcontrol/internal/scanner"
	"os"

	"github.com/spf13/cobra"
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "scans a network",
	Long:  `scans a network`,
	Args:  cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			if err := scanner.ScanNetworkInterfaces(); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		} else {
			if err := scanner.ScanNetwork(args[0]); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
