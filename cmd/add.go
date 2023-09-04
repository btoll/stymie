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

	"github.com/btoll/libstymie"
	"github.com/btoll/stymie/plugin"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new key",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("[stymie] No key name provided, aborting.")
			return
		}

		stymie := libstymie.New(&plugin.GPG{})
		if err := stymie.GetFileContents(); err != nil {
			fmt.Print(err)
			return
		}

		newkey := args[0]

		if _, ok := stymie.Keys[newkey]; !ok {
			// Add the new key => struct.
			stymie.Keys[newkey] = stymie.GetKeyFields()
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
