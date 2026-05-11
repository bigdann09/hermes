/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/bigdann09/notifications/internal/infrastructure/database/seeders"
	"github.com/spf13/cobra"
)

// seedCmd represents the seed command
var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed data to the database",
	Long:  `Hermes CLI Tool for external calls`,
	Run: func(cmd *cobra.Command, args []string) {
		seeders.Seed()
	},
}

func init() {
	rootCmd.AddCommand(seedCmd)
	seedCmd.PersistentFlags().String("all", "a", "Seed all tables")
}
