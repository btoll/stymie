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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/btoll/diceware"
	sillypass "github.com/btoll/sillypass-go"
)

// https://talks.golang.org/2012/10things.slide#4
type GPGConfig struct {
	Armor     bool   `json:"armor"`
	Sign      bool   `json:"sign"`
	Recipient string `json:"recipient"`
}

type Key struct {
	Fields map[string]string `json:"fields"`
}

type PassConfig struct {
	Diceware  int `json:"diceware"`  // Number of words in a Diceware passphrase.
	Sillypass int `json:"sillypass"` // Number of characters in a Sillypass password.
}

type Stymie struct {
	Dir string     `json:"dir"`
	GPG *GPGConfig `json:"gpg"`
	// TODO: Why can't PassConfig be a pointer? 20180603
	PassConfig PassConfig      `json:"passConfig"`
	Keys       map[string]*Key `json:"keys"`
}

/* ----------------------------------------------------------- */
// Private
/* ----------------------------------------------------------- */
func addNewFields(k *Key) *Key {
	var s, n, v string

	fmt.Print("Create another field? [y/N]: ")
	fmt.Scanf("%s", &s)
	switch s {
	case "y":
		fallthrough
	case "Y":
		for {
			fmt.Print("Name: ")

			if _, err := fmt.Scanf("%s", &n); err != nil {
				fmt.Println("Cannot be blank!")
			} else {
				fmt.Print("Value: ")

				if _, err := fmt.Scanf("%s", &v); err != nil {
					fmt.Println("Cannot be blank!")
				} else {
					k.Fields[n] = v
					break
				}
			}
		}

		addNewFields(k)
	}

	return k
}

func formatError(err error) error {
	return fmt.Errorf("[ERROR] %v\n", err)
}

func getKeyFile() string {
	return GetStymieDir() + "/k"
}

func spawnGPG(cmd string, b []byte) []byte {
	gpgCmd := exec.Command("bash", "-c", cmd)
	gpgIn, err := gpgCmd.StdinPipe()
	formatError(err)

	gpgOut, err := gpgCmd.StdoutPipe()
	formatError(err)

	gpgCmd.Start()
	gpgIn.Write(b)
	gpgIn.Close()

	gpgBytes, err := ioutil.ReadAll(gpgOut)
	formatError(err)

	gpgCmd.Wait()

	return gpgBytes
}

/* ----------------------------------------------------------- */
// Public
/* ----------------------------------------------------------- */
func GetKeyFields(passConfig PassConfig) *Key {
	k := &Key{
		Fields: make(map[string]string),
	}

	for {
		var s string

		fmt.Print("URL: ")
		// Note that we don't care if there's an error here!
		fmt.Scanf("%s", &s)

		k.Fields["url"] = s

		for {
			fmt.Print("Username: ")
			if _, err := fmt.Scanf("%s", &s); err != nil {
				fmt.Println("Cannot be blank!")
			} else {
				k.Fields["username"] = s
				break
			}
		}

		fmt.Println("Password generation method:")
		k.Fields["password"] = k.getPassword(passConfig)
		k = addNewFields(k)
		break
	}

	return k
}

func GetStymieDir() string {
	return os.Getenv("HOME") + "/.stymie.d"
}

/* ----------------------------------------------------------- */
// Key methods
/* ----------------------------------------------------------- */
func (k *Key) generatePassphrase(fn func() string) string {
	var t string
	s := fn()
	fmt.Println(s)

	fmt.Print("Accept? [Y/n]: ")
	fmt.Scanf("%s", &t)
	switch t {
	case "n":
		fallthrough
	case "N":
		return k.generatePassphrase(fn)
	default:
		// Remove spaces (nop for Sillypass).
		return strings.Replace(s, " ", "", -1)
	}

	return ""
}

func (k *Key) getPassword(passConfig PassConfig) string {
	var s string

	fmt.Println("\t(1) Diceware (passphrase)")
	fmt.Println("\t(2) Sillypass (mixed-case, alphanumeric, random characters)")
	fmt.Println("\t(3) I'll generate it myself")
	fmt.Println("\t(4) Skip")
	fmt.Print("Select [1]: ")
	fmt.Scanf("%s", &s)
	switch s {
	case "2":
		return k.generatePassphrase(func() string {
			return sillypass.Generate(passConfig.Sillypass)
		})
	case "3":
		for {
			fmt.Print("Custom password: ")
			if _, err := fmt.Scanf("%s", &s); err != nil {
				fmt.Println("Cannot be blank!")
			} else {
				return s
				break
			}
		}
	case "4":
		break
	default:
		return k.generatePassphrase(func() string {
			return diceware.Generate(passConfig.Diceware)
		})
		//			k.generatePassphrase(diceware.Generate)
	}

	return ""
}

func (k *Key) getUpdatedFields(passConfig PassConfig) *Key {
	newkey := &Key{
		Fields: make(map[string]string),
	}

	for key, value := range k.Fields {
		var newvalue string

		if key == "password" {
			var s string
			fmt.Printf("Edit %s (%s):\n", key, value)
			if s = k.getPassword(passConfig); s == "" {
				s = value
			}
			newkey.Fields[key] = s
		} else {
			fmt.Printf("Edit %s (%s): ", key, value)
			// Usually, an error here means that nothing was entered (just a newline, e.g. [Enter]).
			if _, err := fmt.Scanf("%s", &newvalue); err != nil {
				newkey.Fields[key] = value
			} else {
				newkey.Fields[key] = newvalue
			}
		}
	}

	return addNewFields(newkey)
}

/* ----------------------------------------------------------- */
// Stymie methods
/* ----------------------------------------------------------- */
func (c *Stymie) Decrypt(b []byte) []byte {
	// Gather the args from the GPGConfig struct to send to the `gpg` binary.
	return spawnGPG("gpg -d", b)
}

func (c *Stymie) Encrypt(b []byte) []byte {
	// Gather the args from the GPG struct to send to the `gpg` binary.
	args := []string{"-r", c.GPG.Recipient}

	if c.GPG.Armor {
		args = append(args, "-a")
	}

	if c.GPG.Sign {
		args = append(args, "-s")
	}

	cmd := fmt.Sprintf("gpg %s -e", strings.Join(args, " "))
	//	fmt.Println("cmd", cmd)

	return spawnGPG(cmd, b)
}

func (c *Stymie) GetFileContents() error {
	// Maybe pass filename is as func param?
	keyfile := getKeyFile()

	b, err := ioutil.ReadFile(keyfile)
	if err != nil {
		return formatError(err)
	}

	// TODO: Error checking.
	decrypted := c.Decrypt(b)

	// Fill the `stymie` struct with the decrypted json.
	err = json.Unmarshal(decrypted, c)
	formatError(err)

	return nil
}

func (c *Stymie) getConfig() {
	for {
		var s string
		var i int

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

		fmt.Print("How many words should Diceware use to generate a passphrase? [6]: ")
		fmt.Scanf("%s", &s)
		i, _ = strconv.Atoi(s)
		// Equals zero if there was an error, such as a user entering characters that couldn't be converted to an integer.
		if i == 0 {
			c.PassConfig.Diceware = 6
		} else {
			c.PassConfig.Diceware = i
		}

		fmt.Print("How many characters should Sillypass use to generate a password? [12]: ")
		fmt.Scanf("%s", &s)
		// Equals zero if there was an error, such as a user entering characters that couldn't be converted to an integer.
		i, _ = strconv.Atoi(s)
		if i == 0 {
			c.PassConfig.Sillypass = 12
		} else {
			c.PassConfig.Sillypass = i
		}

		return
	}
}

func (c *Stymie) makeConfigFile() error {
	f, err := os.Create(c.Dir + "/k")
	defer f.Close()
	if err != nil {
		return formatError(err)
	}

	gpgConfig, err := json.Marshal(c.GPG)
	if err != nil {
		return formatError(err)
	}

	passConfig, err := json.Marshal(c.PassConfig)
	if err != nil {
		return formatError(err)
	}

	// Stuff the gpgConfig into the json.
	d := fmt.Sprintf("{ \"dir\": \"%s\", \"passConfig\": %s, \"gpg\": %s, \"keys\": {} }", c.Dir, string(passConfig), string(gpgConfig))

	f.Write(c.Encrypt([]byte(d)))

	if err != nil {
		return formatError(err)
	}

	return nil
}

func (c *Stymie) makeDir() {
	os.Mkdir(c.Dir, 0700)
}

func (c *Stymie) PutFileContents() {
	// Back to json (maybe combine this with the actual encryption?).
	b, err := json.Marshal(c)
	formatError(err)

	// Pretty-print the json.
	var out bytes.Buffer
	json.Indent(&out, b, "", "\t")

	// TODO: Error checking.
	encrypted := c.Encrypt(out.Bytes())

	err = ioutil.WriteFile(getKeyFile(), encrypted, 0700)
	formatError(err)
}

// Implement `Stringer` interface.
//func (gpg *GPGConfig) String() string {
//	args := []string{"-r", gpg.Recipient}
//
//	if gpg.Armor {
//		args = append(args, "-a")
//	}
//
//	if gpg.Sign {
//		args = append(args, "-s")
//	}
//
//	return strings.Join(args, " ")
//}
