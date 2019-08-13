package common

type KeyModels map[string]interface{}

type Request struct {
	Username string `json:"username,omitemtpy"`
	Password string `json:"password,omitempty"`
	Mail     string `json:"mail,omitempty"`
	Role     string `json:"role,omitempty"`
}
