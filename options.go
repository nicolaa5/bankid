package bankid

import (
	"encoding/base64"
	"fmt"
	"unicode/utf8"
)


type Option func(RequestBody) (RequestBody, error)

func NewRequest[T RequestBody](opts ...Option) (T, error) {
    var request T

    for _, opt := range opts {
        val, err := opt(request)
		switch err.(type) {
		case nil:
			request = val.(T)
			continue
		case RequiredInputError:
			return request, err
		default:
			fmt.Printf("Warning: %v\n", err)
		}
    }

    return request, nil
}

func WithEndUserIP(endUserIP string) Option {
    return func(rb RequestBody) (RequestBody, error) {
		if endUserIP == "" {
			return nil, RequiredInputError{ Message: fmt.Sprintf("EndUserIP is missing but required by BankID for request: %T", rb)}
		}

		switch v := rb.(type) {
		case AuthRequest:
			v.EndUserIP = endUserIP
			return v, nil
		case SignRequest:
			v.EndUserIP = endUserIP
			return v, nil
		}

		return nil, fmt.Errorf("unkown type: %T", rb)
    }
}

func WithUserVisibleData(userVisibleData string) Option {
    return func(rb RequestBody) (RequestBody, error) {
		switch v := (rb).(type) {
		case AuthRequest:
			if userVisibleData == "" || !utf8.ValidString(userVisibleData) {
				return nil, fmt.Errorf("Optional input is not set: UserVisibleData")
			}

			encodedData := base64.StdEncoding.EncodeToString([]byte(userVisibleData))
			v.UserVisibleData = encodedData
			return v, nil

		case SignRequest:
			// required for sign requests
			if userVisibleData == "" || !utf8.ValidString(userVisibleData) {
				return nil, RequiredInputError{ Message: fmt.Sprintf("UserVisibleData is missing or invalid, it's required by BankID for request: %T", v)}
			}

			encodedData := base64.StdEncoding.EncodeToString([]byte(userVisibleData))
			v.UserVisibleData = encodedData
			return v, nil
		}
		
		return nil, fmt.Errorf("unkown type: %T", rb)
    }
}

func WithUserNonVisibleData(userNonVisibleData string) Option {
	return func(rb RequestBody) (RequestBody, error) {

		if userNonVisibleData == "" || !utf8.ValidString(userNonVisibleData) {
			return nil, fmt.Errorf("Optional input is not set: UserNonVisibleData")
		}

		encodedData := base64.StdEncoding.EncodeToString([]byte(userNonVisibleData))

		switch v := (rb).(type) {
		case AuthRequest:
			v.UserVisibleData = encodedData
			return v, nil

		case SignRequest:
			v.UserVisibleData = encodedData
			return v, nil
		}

		return nil, fmt.Errorf("unkown type: %T", rb)
	}
}