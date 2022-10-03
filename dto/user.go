package dto

type User struct {
	Name      string `json:"Name"`
	Email     string `json:"Email"`
	SlackHook string `json:"SlackHook"`
	Password  string `json:"Password"`
}

type LoginInput struct {
	Email    string `json:"Email"`
	Password string `json:"Password"`
}
