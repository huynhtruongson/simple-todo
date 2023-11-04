package common

type SuccessResponse struct {
	Data interface{} `json:"data"`
}

func NewSimpleSuccessResponse(data interface{}) *SuccessResponse {
	return &SuccessResponse{
		Data: data,
	}
}
