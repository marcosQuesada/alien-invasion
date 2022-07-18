package cmd

import (
	"context"
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
		g := game.NewRunner(runner, writer)
		defer m.Dump(writer)

		for i := 0; i < totalAliens; i++ {
			a := game.NewAlien(fmt.Sprintf("Alien-%d", i), maxIterations)
			if exit := runner.AssignRandomPosition(a); exit {
				return
			} // @TODO: Move to Populate!
		}

		ctx, cancel := context.WithCancel(context.Background())
		go g.Run(ctx)

		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)

		for {
			select {
			case <-signals:
				cancel()
				return
			case m, ok := <-writer.Chan():
				if !ok {
					cancel()
					return
				}
				fmt.Printf("%s \n", m)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(playCmd)

}
