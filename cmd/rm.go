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

	"github.com/btoll/stymie/libstymie"
	"github.com/btoll/stymie/plugin"
	"github.com/spf13/cobra"
)

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Delete a key",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			exit("No key name provided, aborting.")
		}

		toRemove := args[0]

		stymie := libstymie.New(&plugin.GPG{})
		err := stymie.GetFileContents()
		if err != nil {
			exit(fmt.Sprintf("%s", err))
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
				fmt.Println("[stymie] Key not deleted.")
			case "y":
				fallthrough
			case "Y":
				delete(stymie.Keys, toRemove)
				err := stymie.PutFileContents()
				if err != nil {
					exit(fmt.Sprintf("%s", err))
				}
				fmt.Println("[stymie] Successfully removed key.")
			default:
				fmt.Println("[stymie] Key not deleted.")
			}
		} else {
			fmt.Println("[stymie] Key doesn't exist, exiting.")
		}
	},
}

func init() {
	RootCmd.AddCommand(rmCmd)
}
