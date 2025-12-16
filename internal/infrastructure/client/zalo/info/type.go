package info

// UserInfo represents the user information returned by Zalo API.
type UserInfo struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Picture struct {
		Data struct {
			URL string `json:"url"`
		} `json:"data"`
	} `json:"picture"`
	Birthday string `json:"birthday"`
	Gender   string `json:"gender"`
}

type UserPhoneNumberData struct {
	Number string `json:"number"`
}

type UserPhoneNumberResponse struct {
	Data    UserPhoneNumberData `json:"data"`
	Error   int                 `json:"error"`
	Message string              `json:"message"`
}
