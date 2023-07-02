package user

import (
	"github.com/huydnt1801/chuyende/internal/ent"
	"github.com/mitchellh/mapstructure"
)

func DecodeUser(input *ent.User) (*User, error) {
	u := &User{}
	err := mapstructure.Decode(input, u)
	if err != nil {
		return nil, err
	}
	u.Password = input.Password
	return u, nil
}

func MustDecodeUser(input *ent.User) *User {
	out, err := DecodeUser(input)
	if err != nil {
		panic(err)
	}
	return out
}
