package v1

type UserCreateRequest struct {
	Username string `json:"username,omitempty"`
}

type UserLoginRequest struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type UserLoginResponse struct {
	Username   string `json:"username,omitempty"`
	Privileges string `json:"privileges,omitempty"`
}
