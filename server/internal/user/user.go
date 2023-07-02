package user

type User struct {
	ID          int    `json:"id"`
	PhoneNumber string `json:"phoneNumber"`
	ImageURL    string `json:"imageUrl"`
	FullName    string `json:"fullName"`
	Confirmed   bool   `json:"confirmed"`
	Password    string `json:"-"`
}

type UserUpdate struct {
	ID        int
	FullName  string
	Password  string
	Confirmed *bool
	ImageURL  string
}

type UserParams struct {
	PhoneNumber string
	ID          int
}
