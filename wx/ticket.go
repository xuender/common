package wx

// Ticket JS-SDK使用授权
type Ticket struct {
	CommonError
	ExpiresIn int    `json:"expires_in"`
	Ticket    string `json:"ticket"`
}
