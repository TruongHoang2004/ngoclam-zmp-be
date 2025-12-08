package zalo

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

type UserPhoneNumber struct {
	Number string `json:"number"`
}
