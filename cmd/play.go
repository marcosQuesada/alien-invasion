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
		fmt.Println("play called")
		writer := game.NewChannelWriter()
		g := game.NewGame(writer)
		go g.Run()

		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)

		defer log.Println("Application close down")
		for {
			select {
			case <-signals:
				return
			case m, ok := <-writer.Chan():
				if !ok {
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
