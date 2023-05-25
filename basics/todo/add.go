package todo

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a task to your task list.",
	Run: func(cmd *cobra.Command, args []string) {
		task := strings.Join(args, " ")
		_, err := CreateTask(task)
		if err != nil {
			fmt.Println("Uh oh. Something went wrong:", err)
			os.Exit(1)
		}
		fmt.Printf("Added \"%s\" to your task list.\n", task)
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
}
