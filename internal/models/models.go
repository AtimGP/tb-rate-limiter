package models

type LimiterRequest struct {
	Key		string	`json:"key"`
	Limit	int		`json:"limit"`
	Window	string	`json:"window"`
}

type LimiterResponse struct {
	Allowed		bool	`json:"allowed"`
	Remaining	int		`json:"remaining"`
	ResetAfter	string	`json:"reset_after"`
}