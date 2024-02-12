package main

import (
	"path/filepath"

	"fmt"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var CLI = &cobra.Command{
	Use:   "secret",
	Short: "Secret is an API key and other secrets manager",
}

var encodingKey string

// Sets the file to the users home directory instead of using a relative file of the project
func secretsPath() string {
	home, _ := homedir.Dir()
	return filepath.Join(home, ".secrets")
}

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Sets a secret in your secret storage",
	Run: func(cmd *cobra.Command, args []string) {
		v := NewVault(encodingKey, secretsPath())
		key, value := args[0], args[1]
		err := v.Set(key, value)
		if err != nil {
			panic(err)
		}
		fmt.Println("Value set successfully!")
	},
}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets a secret in your secret storage",
	Run: func(cmd *cobra.Command, args []string) {
		v := NewVault(encodingKey, secretsPath())
		key := args[0]
		value, err := v.Get(key)
		if err != nil {
			fmt.Println("no value set")
			return
		}
		fmt.Printf("%s = %s\n", key, value)
	},
}

func init() {
	CLI.PersistentFlags().StringVarP(&encodingKey, "key", "k", "", "the key to use when encoding and decoding secrets")
	CLI.AddCommand(setCmd)
	CLI.AddCommand(getCmd)
}
