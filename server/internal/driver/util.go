package driver

import (
	"github.com/huydnt1801/chuyende/internal/ent"
	"github.com/mitchellh/mapstructure"
)

func DecodeDriver(input *ent.VehicleDriver) (*Driver, error) {
	u := &Driver{}
	err := mapstructure.Decode(input, u)
	if err != nil {
		return nil, err
	}
	u.Password = input.Password
	return u, nil
}

func MustDecodeDriver(input *ent.VehicleDriver) *Driver {
	out, err := DecodeDriver(input)
	if err != nil {
		panic(err)
	}
	return out
}
