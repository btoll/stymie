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

/*
TODO
- 'export HISTIGNORE="stymie *:$HISTIGNORE"\n'
- allow .stymie.d to be installed anywhere and even change the name
- if not accepting default stymie location, get full pathname to put in Dir field
- remove .stymie.d dir on error
*/

package cmd

import (
	"fmt"

	"github.com/btoll/libstymie"
	"github.com/btoll/stymie/plugin"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes the `stymie` database",
	Run: func(cmd *cobra.Command, args []string) {
		stymie := libstymie.New(&plugin.GPG{})
		err := stymie.Init()
		if err != nil {
			exit(fmt.Sprintf("%s", err))
		}

		//fmt.Printf("Created project directory %s.\n", stymie.Dir)
		//fmt.Println("Created stymie config file.")
		fmt.Println("Installation complete!")
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
}
