package zalo

// UserInfo represents the user information returned by Zalo API.
// Fields are mapped to JSON tags as per expected response.
// Note: The exact response structure depends on Zalo API documentation,
// but we include common fields here.
type UserInfo struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Picture struct {
		Data struct {
			URL string `json:"url"`
		} `json:"data"`
	} `json:"picture"`
	// Add other fields as needed based on actual API response
}
