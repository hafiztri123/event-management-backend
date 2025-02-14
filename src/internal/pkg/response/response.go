package response

import "time"



type SuccessResponse struct {
	Timestamp time.Time `json:"timestamp"`
	Data interface{} `json:"data"`
}

type ErrorResponse struct {
	Timestamp time.Time `json:"timestamp"`
	Message string `json:"string"`
}




