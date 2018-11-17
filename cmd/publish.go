// Copyright © 2018 Andrea Funtò - released under the MIT License

package cmd

import (
	"github.com/dihedron/rmq/rmq"
	"github.com/spf13/cobra"
)

// publishCmd represents the publish command
var publishCmd = &cobra.Command{
	Use:     "publish",
	Aliases: []string{"p"},
	Short:   "A brief description of your command",
	Long: `publish allows to send messages to a queue; the application will read 
the message contents from STDIN up to a double newlineand then send it as 
plain/text; the application will then ask for more message until either the 
connection is dropped, or CTRL-C is hit by the user.`,
	Run: rmq.Publish,
}

func init() {
	rootCmd.AddCommand(publishCmd)
}
