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

// https://stackoverflow.com/questions/70394814/create-new-object-of-typed-value-via-go-go-1-18-generics
// https://gotipplay.golang.org/p/IJErmO1mrJh

package plugin

type Plugin interface {
	*Base64 | *GPG | *Plaintext
	Configure() error
	Encrypt([]byte) ([]byte, error)
	Decrypt([]byte) ([]byte, error)
}

type Factory[T any] struct{}

func (f Factory[T]) Create() *T {
	var a T
	return &a
}

//func New[T Plugin](b []byte) T {
//	t := new(T)
//	err := json.Unmarshal(b, t)
//	if err != nil {
//		log.Fatal(err)
//	}
//	return *t
//}
