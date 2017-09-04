package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

type GPGConfig struct {
	Armor     bool   `json:"armor"`
	Sign      bool   `json:"sign"`
	Recipient string `json:"recipient"`
}

type Stymie struct {
	Dir  string                 `json:"dir"`
	GPG  *GPGConfig             `json:"gpg"`
	Keys map[string]interface{} `json:"keys"` // TODO
}

func spawnGPG(cmd string, b []byte) []byte {
	gpgCmd := exec.Command("bash", "-c", cmd)
	gpgIn, err := gpgCmd.StdinPipe()
	CheckError(err)

	gpgOut, err := gpgCmd.StdoutPipe()
	CheckError(err)

	gpgCmd.Start()
	gpgIn.Write(b)
	gpgIn.Close()

	gpgBytes, err := ioutil.ReadAll(gpgOut)
	CheckError(err)

	gpgCmd.Wait()

	return gpgBytes
}

func CheckError(err error) {
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
}

func (c *Stymie) Decrypt(b []byte) []byte {
	// Gather the args from the GPGConfig struct to send to the `gpg` binary.
	return spawnGPG("gpg -d", b)
}

func (c *Stymie) Encrypt(b []byte) []byte {
	// Gather the args from the GPG struct to send to the `gpg` binary.
	cmd := fmt.Sprintf("gpg %s -e", c.GPG)
	return spawnGPG(cmd, b)
}

func GetStymieDir() string {
	return os.Getenv("HOME") + "/.stymie.d"
}
