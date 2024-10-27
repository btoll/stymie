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

/*
TODO
- 'export HISTIGNORE="stymie *:$HISTIGNORE"\n'
- allow .stymie.d to be installed anywhere and even change the name
- if not accepting default stymie location, get full pathname to put in Dir field
- remove .stymie.d dir on error (when initializing)
*/

package cmd

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"

	"github.com/btoll/stymie/plugin"
	"github.com/btoll/stymie/stymie"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes the `stymie` database",
	Run: func(cmd *cobra.Command, args []string) {
		files, err := ioutil.ReadDir("./plugin")
		if err != nil {
			exit(fmt.Sprintf("%s", err))
		}
		var pluginNames []string
		fmt.Println("Please choose a plugin:")
		for _, file := range files {
			f := file.Name()
			pname := f[:len(f)-len(filepath.Ext(f))]
			if pname != "plugin" {
				fmt.Printf("\t(%d) %s\n", len(pluginNames)+1, pname)
				pluginNames = append(pluginNames, pname)
			}
		}
		var str string
		fmt.Print("Select: ")
		fmt.Scanf("%s", &str)
		num, err := strconv.Atoi(str)
		// Adjust for zero-based.
		num -= 1
		if err != nil || num < 0 || num >= len(pluginNames) {
			exit("You must make a valid selection.")
		}
		pluginChoice := pluginNames[num]
		fmt.Printf("Installing the %s plugin.\n", pluginChoice)
		pluginFactory := plugin.Factory[plugin.GPG]{}
		p := pluginFactory.Create()
		stymie := stymie.New[*plugin.GPG](p)
		err = stymie.Init()
		if err != nil {
			exit(fmt.Sprintf("%s", err))
		}
		fmt.Println("Installation complete!")
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
}
