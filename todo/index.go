package todo

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "todo",
	Short: "This is a todo list in your terminal",
}
