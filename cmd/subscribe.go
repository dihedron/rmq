// Copyright © 2018 Andrea Funtò - released under the MIT License

package cmd

import (
	"github.com/dihedron/rmq/rmq"
	"github.com/spf13/cobra"
)

// subscibeCmd represents the subscibe command
var subscribeCmd = &cobra.Command{
	Use:     "subscribe",
	Short:   "Subscribe to a queue",
	Aliases: []string{"s"},
	Long: `subscribe allows to subscribe to a queue; the application will wait 
for messages until either the connection is dropped, or CTRL-C is hit by the user.`,
	Run: rmq.Subscribe,
}

func init() {
	rootCmd.AddCommand(subscribeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// subscribeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// subscribeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
