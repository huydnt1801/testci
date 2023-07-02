package driver

type Driver struct {
	ID          int    `json:"id"`
	PhoneNumber string `json:"phoneNumber"`
	FullName    string `json:"fullName"`
	License     string `json:"license"`
	Password    string `json:"-"`
}
