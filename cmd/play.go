package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/marcosQuesada/alien-invasion/internal/game"
	"github.com/spf13/cobra"
)

// playCmd represents the play command
var playCmd = &cobra.Command{
	Use:   "play",
	Short: "Alien Invasion Play command",
	Long:  `Alien Invasion Play command run Game Play`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Alien Invasion started reading from file %s total aliens %d \n", cfgFile, totalAliens)

		m, err := game.LoaDefinitionsFromFile(cfgFile)
		if err != nil {
			log.Fatalf("unable to load file %s error %v", cfgFile, err)
		}

		rnd := game.NewRandomProvider()
		runner := game.NewEngine(m, rnd)
		writer := game.NewChannelWriter()
		done := make(chan struct{})
		g := game.NewRunner(runner, writer, done)

		defer runner.Dump(writer)

		for i := 0; i < totalAliens; i++ {
			a := game.NewAlien(fmt.Sprintf("Alien-%d", i), maxIterations)
			if exit := runner.AssignRandomPosition(a); exit {
				return
			}
		}

		go g.Run()

		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)

		for {
			select {
			case <-done:
				return
			case <-signals:
				return
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(playCmd)

}
