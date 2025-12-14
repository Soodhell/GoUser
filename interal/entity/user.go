package entity

type User struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Password  string `json:"password"`
	PathImage string `json:"path_image"`
	NameImage string `json:"name_image"`
	Roles     int    `json:"roles"`
}
