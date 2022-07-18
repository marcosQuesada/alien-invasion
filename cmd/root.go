package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	cfgFile       string
	totalAliens   int
	maxIterations int
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "alien invasion",
	Short: "alien invasion root command",
	Long:  `alien invasion root command`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "planet.csv", "planet X config file")
	rootCmd.PersistentFlags().IntVar(&totalAliens, "aliens", 2, "total aliens on map")
	rootCmd.PersistentFlags().IntVar(&maxIterations, "max-iterations", 10000, "max alien iterations")
}
