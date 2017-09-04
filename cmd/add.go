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

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new item",
	//	Long: `A longer description that spans multiple lines and likely contains examples
	//and usage of using your command. For example:
	//
	//Cobra is a CLI library for Go that empowers applications.
	//This application is a tool to generate the needed files
	//to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		stymie := &Stymie{}

		b, err := ioutil.ReadFile(GetStymieDir() + "/k")
		CheckError(err)

		// TODO: Error checking.
		bjson := stymie.Decrypt(b)

		// Fill the `stymie` struct with the decrypted json.
		json.Unmarshal(bjson, stymie)

		fmt.Println("stymie.Dir", stymie.Dir)
		fmt.Println("stymie.GPG", stymie.GPG)
	},
}

func init() {
	RootCmd.AddCommand(addCmd)
}
