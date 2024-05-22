package bankid

import (
	"encoding/base64"
	"fmt"
)

type ProcessOption func(RequestBody) (RequestBody, error)

// Modify input data based on BankID requirements or leave the input unchanged
func process[T RequestBody](opts ...ProcessOption) (T, error) {
	var request T

	for _, opt := range opts {
		val, err := opt(request)
		switch err.(type) {
		case nil:
			request = val.(T)
			continue
		default:
			fmt.Printf("Warning: %v\n", err)
		}
	}

	return request, nil
}

// Ensures UserVisibleData will be base64 encoded if defined
func processUserVisibleData(userVisibleData string) ProcessOption {
	return func(rb RequestBody) (RequestBody, error) {

		// accept empty string
		if userVisibleData == "" {
			return rb, nil
		}

		// in case the input is not base64 encoded we do it ad-hoc instead of failing the request
		if _, err := base64.StdEncoding.DecodeString(userVisibleData); err != nil {
			userVisibleData = base64.StdEncoding.EncodeToString([]byte(userVisibleData))
		}

		switch v := (rb).(type) {
		case AuthRequest:
			v.UserVisibleData = userVisibleData
			return v, nil

		case SignRequest:
			v.UserVisibleData = userVisibleData
			return v, nil

		case PhoneAuthRequest:
			v.UserVisibleData = userVisibleData
			return v, nil

		case PhoneSignRequest:
			v.UserVisibleData = userVisibleData
			return v, nil
		}

		return rb, nil
	}
}

// Ensures UserNonVisibleData will be base64 encoded if defined
func processUserNonVisibleData(userNonVisibleData string) ProcessOption {
	return func(rb RequestBody) (RequestBody, error) {

		// accept empty string
		if userNonVisibleData == "" {
			return rb, nil
		}

		// in case the input is not base64 encoded we do it ad-hoc instead of failing the request
		if _, err := base64.StdEncoding.DecodeString(userNonVisibleData); err != nil {
			userNonVisibleData = base64.StdEncoding.EncodeToString([]byte(userNonVisibleData))
		}

		switch v := (rb).(type) {
		case AuthRequest:
			v.UserNonVisibleData = userNonVisibleData
			return v, nil

		case SignRequest:
			v.UserNonVisibleData = userNonVisibleData
			return v, nil

		case PhoneAuthRequest:
			v.UserNonVisibleData = userNonVisibleData
			return v, nil

		case PhoneSignRequest:
			v.UserNonVisibleData = userNonVisibleData
			return v, nil
		}

		return rb, nil
	}
}

// Ensures UserVisibleDataFormat is set to default markdown if input is empty
func processUserVisibleDataFormat(userVisibleDataFormat string) ProcessOption {
	return func(rb RequestBody) (RequestBody, error) {

		if userVisibleDataFormat == "" {
			// set to default formatting style
			userVisibleDataFormat = "simpleMarkdownV1"
		}

		switch v := (rb).(type) {
		case AuthRequest:
			v.UserVisibleDataFormat = userVisibleDataFormat
			return v, nil

		case SignRequest:
			v.UserVisibleDataFormat = userVisibleDataFormat
			return v, nil

		case PhoneAuthRequest:
			v.UserVisibleDataFormat = userVisibleDataFormat
			return v, nil

		case PhoneSignRequest:
			v.UserVisibleDataFormat = userVisibleDataFormat
			return v, nil
		}

		return rb, nil
	}
}