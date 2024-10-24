// Copyright © 2024 Benjamin Toll <ben@benjamintoll.com>
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

import "fmt"

type Encrypter interface {
	Base64 | GPG | Plaintext
	Configure() error
	//	GetPlugin(v string) (Encrypter, error)
	Encrypt([]byte) ([]byte, error)
	Decrypt([]byte) ([]byte, error)
}

type Plugin[T Encrypter] struct {
	Encrypter T
}

func (p *Plugin[T]) f() *Plugin[T] {
	fmt.Printf("%+v\n", p)
	//	fmt.Println(p.Encrypter.Name)
	return p
}

func New(s string) *GPG {
	//func New[T Encrypter](s string) T {
	//	var plugin T
	switch s {
	//	case "base64":
	//		p := Plugin[Base64]{
	//			Encrypter: Base64{
	//				Name: "base64",
	//			},
	//		}
	//		return p
	case "gpg":
		return &GPG{}

		//		plugin = Plugin[GPG]{
		//			Encrypter: GPG{
		//				Name: "gpg",
		//			},
		//		}
		//		return plugin
		//	case "plaintext":
		//		p := Plugin[Plaintext]{
		//			Encrypter: Plaintext{
		//				Name: "plaintext",
		//			},
		//		}
		//		return p
		//	default:
		//		return nil
	default:
		return nil
	}
}

//func Configure[E Encrypter](E) error {
//	return e.Configure()
//}

//func GetPlugin[E Encrypter](E) (Encrypter, error) {
//	_, err := stymie.GetConfigFileContents()
//	if err != nil {
//		return nil, err
//	}
// Fill the `stymie` struct with the decrypted json.
//	err = json.Unmarshal(b, s)
//	fmt.Println("s", s)
//	if err != nil {
//		return nil, formatError(err)
//	}

//	return nil, nil
//}

//func Decrypt(b []byte) ([]byte, error) {
//	return e.Decrypt(b)
//}
//
//func Encrypt(b []byte) ([]byte, error) {
//	return e.Encrypt(b)
//}
