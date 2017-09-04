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
- error handling
- 'export HISTIGNORE="stymie *:$HISTIGNORE"\n'
- remove .stymie.d dir on error
*/

package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

type GPGConfig struct {
	Armor     bool   `json:"armor"`
	Sign      bool   `json:"sign"`
	Recipient string `json:"recipient"`
}

type Stymie struct {
	Dir string
	*GPGConfig
}

type Task struct {
	Run func()
}

// Implement `Stringer` interface.
func (gpg *GPGConfig) String() string {
	args := []string{"-r", gpg.Recipient}

	if gpg.Armor {
		args = append(args, "-a")
	}

	if gpg.Sign {
		args = append(args, "-s")
	}

	return strings.Join(args, " ")
}

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
				c.Recipient = s
				break
			}
		}

		fmt.Print("Should GPG/PGP encrypt the password files as binary? [Y/n]: ")
		fmt.Scanf("%s", &s)
		switch s {
		case "n":
			fallthrough
		case "N":
			c.Armor = true
		default:
			c.Armor = false
		}

		fmt.Print("Should GPG/PGP also sign the password files? (Recommended) [Y/n]: ")
		fmt.Scanf("%s", &s)
		switch s {
		case "n":
			fallthrough
		case "N":
			c.Sign = false
		default:
			c.Sign = true
		}

		return
	}
}

func (c *Stymie) makeConfigFile() {
	f, err := os.Create(c.Dir + "/k")
	defer f.Close()

	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	fmt.Println("config", c)

	b, err := json.Marshal(c.GPGConfig)

	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	// Stuff the gpgConfig into the json.
	json := fmt.Sprintf("{ \"gpg\": %s, \"keys\": {} }", string(b))

	f.Write(c.encrypt(json))

	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
}

func (c *Stymie) makeDir() {
	os.Mkdir(c.Dir, 0700)
}

func (c *Stymie) encrypt(s string) []byte {
	// Gather the args from the GPGConfig struct to send to the `gpg` binary.
	cmd := fmt.Sprintf("gpg %s -e", c)
	gpgCmd := exec.Command("bash", "-c", cmd)
	gpgIn, _ := gpgCmd.StdinPipe()
	gpgOut, _ := gpgCmd.StdoutPipe()

	gpgCmd.Start()
	gpgIn.Write([]byte(s))
	gpgIn.Close()

	gpgBytes, _ := ioutil.ReadAll(gpgOut)
	gpgCmd.Wait()

	return gpgBytes
}

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		stymie := &Stymie{
			Dir:       os.Getenv("HOME") + "/.stymie.d",
			GPGConfig: &GPGConfig{},
		}

		stymie.getConfig()

		fmt.Printf("Creating project directory %s\n", stymie.Dir)
		stymie.makeDir()

		fmt.Println("Creating stymie config file")
		stymie.makeConfigFile()

		fmt.Println("Installation complete")
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
}
