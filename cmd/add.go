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

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new key",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			exit("No key name provided, aborting.")
		}
		s, err := stymie.GetStymie[*plugin.GPG]()
		if err != nil {
			exit(fmt.Sprintf("%s", err))
		}
		decryptedKeys := stymie.Keys{}
		err = json.Unmarshal(s.Keys, &decryptedKeys)
		newkey := args[0]

		if _, ok := decryptedKeys[newkey]; !ok {
			// Add the new key => struct.
			decryptedKeys[newkey] = s.GetKeyFields()
			keys, err := json.Marshal(decryptedKeys)
			if err != nil {
				exit(fmt.Sprintf("%s", err))
			}
			encryptedKeys, err := s.Encrypt(keys)
			if err != nil {
				exit(fmt.Sprintf("%s", err))
			}
			s.Keys = encryptedKeys
			err = s.PutFileContents()
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
