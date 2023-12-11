package plugin

import (
	"github.com/btoll/stymie/libstymie"
)

func New(s string) libstymie.Plugin {
	var p libstymie.Plugin
	if s == "base64" {
		p = &Base64{}
	}
	if s == "gpg" {
		p = &GPG{}
	}
	if s == "plaintext" {
		p = &Plaintext{}
	}
	return p
}
