package utils

import (
	"fmt"
	"strings"

	"atayalan.com/analytics/customValidators"
)

type ErrWithParam struct {
	Prefix string
	Msg    string
	Param  []InvalidParam
}

func (e ErrWithParam) Error() string {
	var s string
	if e.Msg != "" {
		return e.Msg
	}

	for _, v := range e.Param {
		if s == "" {
			s = v.String()
		} else {
			s += fmt.Sprintf(", %s", v.String())
		}
	}
	return strings.TrimSpace(fmt.Sprintf("%s: %s", e.Prefix, s))
}

func (e ErrWithParam) Details(format string) ErrorDetails {
	return ErrorDetails{
		Detail:        fmt.Sprintf(format, e.Error()),
		InvalidParams: e.Param,
	}
}

func NewErrWithParam(msg string, field string, reason string) error {
	return &ErrWithParam{
		Msg: msg,
		Param: []InvalidParam{
			{
				Param:  field,
				Reason: reason,
			},
		},
	}
}

func NewErrWithParams(msg string, prm []InvalidParam) error {
	return &ErrWithParam{
		Msg:   msg,
		Param: prm,
	}
}

func NewErrWithParamsWithPrefix(prefix string, prm []InvalidParam) error {
	return &ErrWithParam{
		Prefix: prefix,
		Param:  prm,
	}
}

type ErrorDetails struct {
	Detail        string         `json:"detail"`
	InvalidParams []InvalidParam `json:"invalidParams,omitempty"`
}

type InvalidParam struct {
	Param  string `json:"param"`
	Reason string `json:"reason"`
}

func (p *InvalidParam) String() string {
	return fmt.Sprintf("%s %s", p.Param, p.Reason)
}

func NewValidateErrorDetails(format string, msg string, ve []*customValidators.ErrorResponse) ErrorDetails {
	res := ErrorDetails{
		Detail:        fmt.Sprintf(format, msg),
		InvalidParams: make([]InvalidParam, 0),
	}

	for _, e := range ve {
		res.InvalidParams = append(res.InvalidParams, InvalidParam{
			Param:  e.FailedField,
			Reason: strings.TrimRight(fmt.Sprintf("must be %s %s", e.Tag, e.Value), " "),
		})
	}

	return res
}

func GetErrorDetails(format string, e error) ErrorDetails {
	if err, ok := e.(*ErrWithParam); ok {
		return err.Details(format)
	}

	return ErrorDetails{
		Detail: fmt.Sprintf(format, e.Error()),
	}
}
