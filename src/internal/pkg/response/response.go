package response

import "time"




type Response struct {
	Timestamp time.Time `json:"timestamp"`
	Message string `json:"message,omitempty"`
	Data interface{} `json:"data,omitempty"`
}




