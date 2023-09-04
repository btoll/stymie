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

package plugin

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
)

// https://talks.golang.org/2012/10things.slide#4
type GPG struct {
	Armor     bool   `json:"armor"`
	Sign      bool   `json:"sign"`
	Recipient string `json:"recipient"`
}

func (g *GPG) Configure() error {
	var s string
	fmt.Print("Enter the email address or key ID of your public key: ")
	if _, err := fmt.Scanf("%s", &s); err != nil {
		return errors.New("Public key ID cannot be blank.")
	} else {
		g.Recipient = s
	}

	fmt.Print("Should GPG/PGP encrypt the password files as binary? [Y/n]: ")
	fmt.Scanf("%s", &s)
	switch s {
	case "n":
		fallthrough
	case "N":
		g.Armor = true
	default:
		g.Armor = false
	}

	fmt.Print("Should GPG/PGP also sign the password files? (Recommended) [Y/n]: ")
	fmt.Scanf("%s", &s)
	switch s {
	case "n":
		fallthrough
	case "N":
		g.Sign = false
	default:
		g.Sign = true
	}
	return nil
}

func (g *GPG) Decrypt(b []byte) ([]byte, error) {
	return spawnGPG("gpg -d", b)
}

func (g *GPG) Encrypt(b []byte) ([]byte, error) {
	args := []string{"-r", g.Recipient}

	if g.Armor {
		args = append(args, "-a")
	}

	if g.Sign {
		args = append(args, "-s")
	}

	cmd := fmt.Sprintf("gpg %s -e", strings.Join(args, " "))
	//	fmt.Println("cmd", cmd)

	return spawnGPG(cmd, b)
}

func spawnGPG(cmd string, b []byte) ([]byte, error) {
	gpgCmd := exec.Command("bash", "-c", cmd)
	gpgIn, err := gpgCmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	gpgOut, err := gpgCmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	gpgCmd.Start()
	gpgIn.Write(b)
	gpgIn.Close()

	gpgBytes, err := ioutil.ReadAll(gpgOut)
	if err != nil {
		return nil, err
	}

	gpgCmd.Wait()
	return gpgBytes, nil
}
