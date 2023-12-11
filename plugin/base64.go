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
	"encoding/base64"
)

type Base64 struct {
	Name string `json:"name,noempty"`
}

func (b *Base64) Configure() error {
	b.Name = "base64"
	return nil
}

func (b *Base64) Decrypt(chars []byte) ([]byte, error) {
	return base64.StdEncoding.DecodeString(string(chars))
}

func (b *Base64) Encrypt(chars []byte) ([]byte, error) {
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(chars)))
	base64.StdEncoding.Encode(dst, chars)
	return dst, nil
}
