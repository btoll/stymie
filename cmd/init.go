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
	"os"

	"github.com/spf13/cobra"
)

type GPGConfig struct {
	Armor     bool
	Sign      bool
	Recipient string
}

type StymieConfig struct {
	Dir string
	*GPGConfig
}

type Task struct {
	Msg string
	Run func(done chan bool)
}

func getStymieConfig(c *StymieConfig) {
	for {
		var s string

		fmt.Print("Enter the full path of the directory to install .stymie.d [~/.stymie.d]: ")
		fmt.Scanf("%s", &s)
		if s == "" {
			//            inputs[2] = os.Getenv("HOME") + ".stymie.d"
			c.Dir = "stymie.d"
		} else {
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

func makeDir(dir string, done chan bool) {
	os.Mkdir(dir, 0700)
	done <- true
}

func createConfigFile(c *StymieConfig, done chan bool) {
	// Create config file.
	f, err := os.Create(c.Dir + "/c")
	defer f.Close()

	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	b, err := json.Marshal(c.GPGConfig)

	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	f.Write(b)

	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	done <- true
	//		        return util.encrypt(JSON.stringify(gpgOptions, null, 4))
	//		        .then(writeFile(`${stymieDir}/c`))
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
		c := StymieConfig{
			Dir:       os.Getenv("HOME"),
			GPGConfig: &GPGConfig{},
		}

		getStymieConfig(&c)

		done := make(chan bool)

		tasks := []Task{
			{
				Msg: "Creating project directory " + c.Dir,
				Run: func(done chan bool) {
					makeDir(c.Dir, done)
				},
			},
			{
				Msg: "Creating stymie config file",
				Run: func(done chan bool) {
					createConfigFile(&c, done)
					done <- true
				},
			},
		}

		for _, task := range tasks {
			fmt.Println(task.Msg)
			go task.Run(done)
			<-done
		}
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
}
