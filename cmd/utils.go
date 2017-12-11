package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
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

type Stymie struct {
	Dir  string          `json:"dir"`
	GPG  *GPGConfig      `json:"gpg"`
	Keys map[string]*Key `json:"keys"`
}

func spawnGPG(cmd string, b []byte) []byte {
	gpgCmd := exec.Command("bash", "-c", cmd)
	gpgIn, err := gpgCmd.StdinPipe()
	FormatError(err)

	gpgOut, err := gpgCmd.StdoutPipe()
	FormatError(err)

	gpgCmd.Start()
	gpgIn.Write(b)
	gpgIn.Close()

	gpgBytes, err := ioutil.ReadAll(gpgOut)
	FormatError(err)

	gpgCmd.Wait()

	return gpgBytes
}

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

func FormatError(err error) error {
	return fmt.Errorf("[ERROR] %v\n", err)
}

func (c *Stymie) GetFileContents() error {
	// Maybe pass filename is as func param?
	keyfile := GetKeyFile()

	b, err := ioutil.ReadFile(keyfile)
	if err != nil {
		return FormatError(err)
	}

	// TODO: Error checking.
	decrypted := c.Decrypt(b)

	// Fill the `stymie` struct with the decrypted json.
	err = json.Unmarshal(decrypted, c)
	FormatError(err)

	return nil
}

func GetKeyFile() string {
	return GetStymieDir() + "/k"
}

func GetStymieDir() string {
	return os.Getenv("HOME") + "/.stymie.d"
}

func (c *Stymie) PutFileContents() {
	// Back to json (maybe combine this with the actual encryption?).
	b, err := json.Marshal(c)
	FormatError(err)

	// Pretty-print the json.
	var out bytes.Buffer
	json.Indent(&out, b, "", "\t")

	// TODO: Error checking.
	encrypted := c.Encrypt(out.Bytes())

	err = ioutil.WriteFile(GetKeyFile(), encrypted, 0700)
	FormatError(err)
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
