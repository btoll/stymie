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
	"sort"

	"github.com/btoll/stymie/stymie"
	"github.com/spf13/cobra"
)

type ByKey []string

func (a ByKey) Len() int           { return len(a) }
func (a ByKey) Less(i, j int) bool { return a[i] < a[j] }
func (a ByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type Interface interface {
	// Len is the number of elements in the collection.
	Len() int
	// Less reports whether the element with
	// index i should sort before the element with index j.
	Less(i, j int) bool
	// Swap swaps the elements with indexes i and j.
	Swap(i, j int)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all saved keys",
	Run: func(cmd *cobra.Command, args []string) {
		s, err := stymie.GetStymie()
		if err != nil {
			exit(fmt.Sprintf("%s", err))
		}
		decryptedKeys := stymie.Keys{}
		err = json.Unmarshal(s.Keys, &decryptedKeys)
		if err != nil {
			exit(fmt.Sprintf("%s", err))
		}
		if len(decryptedKeys) == 0 {
			fmt.Println("[stymie] No installed keys.")
		} else {
			fmt.Println("[stymie] Saved keys:")
			keys := make(ByKey, len(decryptedKeys))
			j := 0
			for key := range decryptedKeys {
				keys[j] = key
				j = j + 1
			}
			sort.Sort(keys)
			for _, key := range keys {
				fmt.Println(key)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
