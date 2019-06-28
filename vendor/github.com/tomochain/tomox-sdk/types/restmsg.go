package types

import (
	"fmt"
)

type RestResponse struct {
	Status string           `json:"status"`
	Data   interface{}      `json:"data"`
	Meta   RestResponseMeta `json:"meta"`
}

func (r *RestResponse) String() string {
	return fmt.Sprintf("Status: %v / Data: %v / Meta: %v", r.Status, r.Data, r.Meta.String())
}

type RestResponseMeta struct {
	Total int `json:"total"`
}

func (m *RestResponseMeta) String() string {
	return fmt.Sprintf("%v", m.Total)
}
