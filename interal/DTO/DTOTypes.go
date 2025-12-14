package DTO

type ResponseError struct {
	Error string `json:"error"`
}

type ResponseToken struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
}

type ResponseSuccess struct {
	Success bool   `json:"success"`
	Mail    string `json:"mail"`
	Roles   int    `json:"roles"`
}

type ReturnUser struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Image string `json:"image"`
	Roles int    `json:"roles"`
}

type ErrorUser struct {
	Code int
	Msg  string
}

type FileReturn struct {
	AvatarPath string
	AvatarName string
}

func (e *ErrorUser) Error() string {
	return e.Msg
}
