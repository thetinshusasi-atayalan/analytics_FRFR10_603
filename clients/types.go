package clients

import (
	"fmt"

	"atayalan.com/analytics/customValidators"
)

type ClientType string

const (
	FiveG ClientType = "5G"
	Wlan  ClientType = "Wlan"
)

type DataUsage struct {
	BytesIn  int64 `json:"bytesIn"  validate:"required|isdefault"`
	BytesOut int64 `json:"bytesOut"  validate:"required|isdefault"`
}

type Usage struct {
	ClientId        string     `json:"clientId" validate:"required"`
	ClientType      ClientType `json:"clientType" validate:"required,oneof=5G Wlan"`
	Dnn             string     `json:"dnn"  validate:"required"`
	BytesToClient   int64      `json:"bytesToClient"  validate:"required|isdefault"`
	BytesFromClient int64      `json:"bytesFromClient"  validate:"required|isdefault"`
	IP              string     `json:"ip"  validate:"required,ipv4"`
	UPFIP           string     `json:"upfIp"  validate:"required,ipv4"`
}
type ClientUsage struct {
	ClientId   string     `json:"clientId" validate:"required"`
	ClientType ClientType `json:"clientType" validate:"required"`
	DataUsage
}
type DnnUsage struct {
	DataUsage
	Dnn string `json:"dnn"  validate:"required"`
}

func (config *Usage) ValidateStruct() []*customValidators.ErrorResponse {
	errors := make([]*customValidators.ErrorResponse, 0)

	err := customValidators.CustomValidate.Struct(config)
	if err != nil {
		fmt.Printf("ValidateStruct Usage %v\n", err)
		errors = append(errors, customValidators.ConstructValidationError(err)...)

	}

	if len(errors) == 0 {
		return nil
	}

	return errors
}
