/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/spf13/cobra"
)

type PomodoroCycle struct {
	cycle      int32
	work_time  int32
	while_time int32
	duration   time.Duration
}

var durationMapping = map[string]time.Duration{"s": time.Second, "m": time.Minute, "h": time.Hour}

func (pomodoroCycle *PomodoroCycle) Start(done chan int) {
	for i := 0; int32(i) < pomodoroCycle.cycle; i++ {
		timer := time.NewTimer(time.Duration(pomodoroCycle.work_time) * pomodoroCycle.duration)
		<-timer.C

		beeep.Notify("Take a break", "Work time over take a break", "./icon1.png")

		timer = time.NewTimer(time.Duration(pomodoroCycle.while_time) * pomodoroCycle.duration)
		<-timer.C

		beeep.Alert("Go back to work", "Time to go back to work", "./icon1.png")

	}
	done <- 1
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "zentimer",
	Short: "Productivity timer application",
	Long: `This is a pomodoro timer app. Use --work-time and --while-time flags for corresponding
		work and while away time you need reminder for.
	`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: func(cmd *cobra.Command, args []string) error {

		work_time, _ := cmd.Flags().GetInt32("work-time")
		while_time, _ := cmd.Flags().GetInt32("while-time")
		cycles, _ := cmd.Flags().GetInt32("cycles")
		d, _ := cmd.Flags().GetString("duration")

		duration, present := durationMapping[d]

		if !present {
			return fmt.Errorf("duration can only be  s,m or h for seconds minutes or hours")
		}

		pomodoroCycle := PomodoroCycle{
			work_time:  work_time,
			while_time: while_time,
			cycle:      cycles,
			duration:   duration,
		}

		done := make(chan int)

		go pomodoroCycle.Start(done)

		<-done
		return nil
	},
	Args: cobra.NoArgs,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.zentimer.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().Int32("work-time", 45, "Time for which you want to work in minutes")
	rootCmd.Flags().Int32("while-time", 15, "Time for which you want to rest in minutes")
	rootCmd.Flags().Int32("cycles", 2, "Number of work cycles")
	rootCmd.Flags().String("duration", "m", "Number of work cycles")

}
