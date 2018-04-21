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

	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Delete a key",
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

		toRemove := args[0]

		stymie := &Stymie{}
		if err := stymie.GetFileContents(); err != nil {
			fmt.Print(err)
			return
		}

		if _, ok := stymie.Keys[toRemove]; ok {
			var s string
			fmt.Print("Are you sure you wish to delete the key [y/N]: ")
			fmt.Scanf("%s", &s)
			switch s {
			case "":
				fallthrough
			case "n":
				fallthrough
			case "N":
				fmt.Println("[stymie] Operation aborted.")
			default:
				delete(stymie.Keys, toRemove)
				stymie.PutFileContents()
				fmt.Println("[stymie] Successfully removed key.")
			}
		} else {
			fmt.Println("[stymie] Key doesn't exist, exiting.")
		}
	},
}

func init() {
	RootCmd.AddCommand(rmCmd)
}
