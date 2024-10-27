// Copyright © 2024 Benjamin Toll <ben@benjamintoll.com>
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
	"encoding/json"
	"fmt"

	"github.com/btoll/stymie/plugin"
	"github.com/btoll/stymie/stymie"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get the fields of the given key",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			exit("No key name provided, aborting.")
		}
		keyname := args[0]
		s, err := stymie.GetStymie[*plugin.GPG]()
		if err != nil {
			exit(fmt.Sprintf("%s", err))
		}
		decryptedKeys := stymie.Keys{}
		err = json.Unmarshal(s.Keys, &decryptedKeys)
		if err != nil {
			exit(fmt.Sprintf("%s", err))
		}
		if _, ok := decryptedKeys[keyname]; ok {
			fmt.Printf("[stymie] Printing key `%s`\n\n", keyname)
			for key, value := range decryptedKeys[keyname].Fields {
				fmt.Printf("%s: %s\n", key, value)
			}
		} else {
			exit("Key doesn't exist, exiting.")
		}
	},
}

func init() {
	RootCmd.AddCommand(getCmd)
}
