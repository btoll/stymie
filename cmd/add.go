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
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new key",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			exit("No key name provided, aborting.")
		}

		//		stymie := libstymie.New(&plugin.GPG{})
		stymie := libstymie.Stymie{}
		if err := stymie.GetFileContents(); err != nil {
			exit(fmt.Sprintf("%s", err))
		}

		newkey := args[0]

		if _, ok := stymie.Keys[newkey]; !ok {
			// Add the new key => struct.
			stymie.Keys[newkey] = stymie.GetKeyFields()
			err := stymie.PutFileContents()
			if err != nil {
				exit(fmt.Sprintf("%s", err))
			}
			fmt.Println("[stymie] Successfully created key.")
		} else {
			exit("Key already exists, exiting.")
		}
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
}
