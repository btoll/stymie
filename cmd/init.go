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
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func (c *Stymie) getConfig() {
	for {
		var s string

		fmt.Print("Enter the full path of the directory to install .stymie.d [~/.stymie.d]: ")
		fmt.Scanf("%s", &s)
		if s != "" {
			c.Dir = s
		}

		for {
			fmt.Print("Enter the email address or key ID of your public key: ")
			if _, err := fmt.Scanf("%s", &s); err != nil {
				fmt.Println("Cannot be blank!!")
			} else {
				c.GPG.Recipient = s
				break
			}
		}

		fmt.Print("Should GPG/PGP encrypt the password files as binary? [Y/n]: ")
		fmt.Scanf("%s", &s)
		switch s {
		case "n":
			fallthrough
		case "N":
			c.GPG.Armor = true
		default:
			c.GPG.Armor = false
		}

		fmt.Print("Should GPG/PGP also sign the password files? (Recommended) [Y/n]: ")
		fmt.Scanf("%s", &s)
		switch s {
		case "n":
			fallthrough
		case "N":
			c.GPG.Sign = false
		default:
			c.GPG.Sign = true
		}

		return
	}
}

func (c *Stymie) makeConfigFile() error {
	f, err := os.Create(c.Dir + "/k")
	defer f.Close()
	if err != nil {
		return FormatError(err)
	}

	b, err := json.Marshal(c.GPG)

	if err != nil {
		return FormatError(err)
	}

	// Stuff the gpgConfig into the json.
	d := fmt.Sprintf("{ \"dir\": \"%s\", \"gpg\": %s, \"keys\": {} }", c.Dir, string(b))

	f.Write(c.Encrypt([]byte(d)))

	if err != nil {
		return FormatError(err)
	}

	return nil
}

func (c *Stymie) makeDir() {
	os.Mkdir(c.Dir, 0700)
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initializes the `stymie` store",
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
