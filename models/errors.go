package models

import (
	"encoding/json"
)

type JSONSerializable interface {
	ToJSON() []byte
}

// ErrorResponse model
// swagger:model ErrorResponse
type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// ErrorResponseWithRequestStatus model
// swagger:model ErrorResponseWithRequestStatus
type ErrorResponseWithRequestStatus struct {
	Status        int           `json:"status"`
	RequestStatus RequestStatus `json:"requestStatus"`
}

func (r ErrorResponse) ToJSON() []byte {
	jsonM, _ := json.Marshal(r)
	return jsonM
}

func (r ErrorResponseWithRequestStatus) ToJSON() []byte {
	jsonM, _ := json.Marshal(r)
	return jsonM
}
