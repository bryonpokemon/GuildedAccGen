package guilded

import "net/http"

type GuildeadClient struct {
	Client    *http.Client
	Token     string
	Username  string
	Email     string
	Password  string
	ClienID   string
	DeviceID  string
	EmailBase string
}

type RegisterPayload struct {
	ExtraInfo ExtraInfo `json:"extraInfo,omitempty"`
	Name      string    `json:"name,omitempty"`
	Email     string    `json:"email,omitempty"`
	Password  string    `json:"password,omitempty"`
	FullName  string    `json:"fullName,omitempty"`
}
type ExtraInfo struct {
	Platform string `json:"platform,omitempty"`
}
