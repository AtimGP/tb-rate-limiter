package models

type LimiterRequest struct {
	Key		string
	Limit	int
	Window	string
}

type LimiterResponse struct {
	Allowed		bool
	Remaining	int
	ResetAfter	string
}