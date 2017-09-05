package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

type GPGConfig struct {
	Armor     bool   `json:"armor"`
	Sign      bool   `json:"sign"`
	Recipient string `json:"recipient"`
}

type Key struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Stymie struct {
	Dir  string          `json:"dir"`
	GPG  *GPGConfig      `json:"gpg"`
	Keys map[string]*Key `json:"keys"`
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

func GetKeyFile() string {
	return GetStymieDir() + "/k"
}

func GetStymieDir() string {
	return os.Getenv("HOME") + "/.stymie.d"
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
