package response

// Request Response
type ReqRes struct {
	Type string `json:"type"`
	Msg  string `json:"msg"`
}

// Create an error response
func NewReqResFromError(err error) *ReqRes {
	return &ReqRes{
		Type: "error",
		Msg:  err.Error(),
	}
}
