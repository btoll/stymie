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
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all saved keys",
	//	Long: `A longer description that spans multiple lines and likely contains examples
	//and usage of using your command. For example:
	//
	//Cobra is a CLI library for Go that empowers applications.
	//This application is a tool to generate the needed files
	//to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		stymie := &Stymie{}
		keyfile := GetKeyFile()

		b, err := ioutil.ReadFile(keyfile)
		CheckError(err)

		// TODO: Error checking.
		decrypted := stymie.Decrypt(b)

		// Fill the `stymie` struct with the decrypted json.
		err = json.Unmarshal(decrypted, stymie)
		CheckError(err)

		fmt.Println("Saved keys:\n")

		for key := range stymie.Keys {
			fmt.Println(key)
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
