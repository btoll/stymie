// Copyright © 2017 Benjamin Toll <ben@benjamintoll.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package cmd

import (
	"fmt"

	//	diceware "github.com/btoll/diceware-go/lib"
	"github.com/spf13/cobra"
)

func (k *Key) getFields() {
	for {
		var s string

		fmt.Print("URL: ")
		_, err := fmt.Scanf("%s", &s)
		CheckError(err)
		k.Fields["url"] = s

		for {
			fmt.Print("Username: ")
			if _, err := fmt.Scanf("%s", &s); err != nil {
				fmt.Println("Cannot be blank!!")
			} else {
				k.Fields["username"] = s
				break
			}
		}

		for {
			fmt.Print("Password: ")
			if _, err := fmt.Scanf("%s", &s); err != nil {
				fmt.Println("Cannot be blank!!")
			} else {
				k.Fields["password"] = s
				break
			}
		}

		//		fmt.Println(diceware.GetPassphrase(6))
		return
	}
}

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new key",
	//	Long: `A longer description that spans multiple lines and likely contains examples
	//and usage of using your command. For example:
	//
	//Cobra is a CLI library for Go that empowers applications.
	//This application is a tool to generate the needed files
	//to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("[stymie] No key name provided, aborting.")
			return
		}

		stymie := &Stymie{}
		stymie.GetFileContents()

		newkey := args[0]

		if _, ok := stymie.Keys[newkey]; !ok {
			k := &Key{
				Fields: make(map[string]string),
			}

			k.getFields()

			// Add the new key => struct.
			stymie.Keys[newkey] = k

			stymie.PutFileContents()

			fmt.Println("[stymie] Successfully created key.")
		} else {
			fmt.Println("[stymie] Key already exists, exiting.")
		}
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
}
