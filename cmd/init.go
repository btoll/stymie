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

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes the `stymie` database",
	//	Long: `A longer description that spans multiple lines and likely contains examples
	//and usage of using your command. For example:
	//
	//Cobra is a CLI library for Go that empowers applications.
	//This application is a tool to generate the needed files
	//to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		stymie := &Stymie{
			Dir:  GetStymieDir(),
			GPG:  &GPGConfig{},
			Keys: nil, // TODO
		}

		stymie.getConfig()

		fmt.Printf("Creating project directory %s\n", stymie.Dir)
		stymie.makeDir()

		fmt.Println("Creating stymie config file")
		stymie.makeConfigFile()

		fmt.Println("Installation complete!")
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
}
