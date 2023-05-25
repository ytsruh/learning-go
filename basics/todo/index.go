package todo

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "todo",
	Short: "This is a todo list in your terminal",
}

func RunTodo() {
	// Bind a users home directory and setup the database there
	home, _ := homedir.Dir()
	dbPath := filepath.Join(home, "tasks.db")
	must(InitDB(dbPath))
	must(RootCmd.Execute())
}

func must(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
