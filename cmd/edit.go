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

func (k *Key) getUpdatedFields() *Key {
	newkey := &Key{
		Fields: make(map[string]string),
	}

	for key, value := range k.Fields {
		var newvalue string

		fmt.Printf("Edit %s (%s): ", key, value)

		// Usually, an error here means that nothing was entered (just a newline, e.g. [Enter]).
		if _, err := fmt.Scanf("%s", &newvalue); err != nil {
			newkey.Fields[key] = value
		} else {
			newkey.Fields[key] = newvalue
		}
	}

	return newkey
}

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit a key",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("[stymie] No key name provided, aborting.")
			return
		}

		stymie := &Stymie{}
		if err := stymie.GetFileContents(); err != nil {
			fmt.Print(err)
			return
		}

		keyname := args[0]

		if _, ok := stymie.Keys[keyname]; ok {
			key := stymie.Keys[keyname]
			stymie.Keys[keyname] = key.getUpdatedFields()
			stymie.PutFileContents()
			fmt.Printf("\n[stymie] Updated key `%s`\n", keyname)
		} else {
			fmt.Println("[stymie] Key doesn't exist, exiting.")
		}
	},
}

func init() {
	RootCmd.AddCommand(editCmd)
}
