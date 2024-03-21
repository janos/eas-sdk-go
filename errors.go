package eas

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/rpc"
)

type Error struct {
	StatusCode int
	Code       int
	Message    string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%v %s code=%v %s", e.StatusCode, http.StatusText(e.StatusCode), e.Code, e.Message)
}

type ContractError struct {
	Name      string
	Arguments []ContractErrorArgument
}

func (e *ContractError) Error() string {
	al := len(e.Arguments)
	if al == 0 {
		return e.Name
	}
	b := new(strings.Builder)
	b.WriteString(e.Name)
	b.WriteString(" (")
	for i, a := range e.Arguments {
		if a.Name != "" {
			b.WriteString(a.Name)
		} else {
			b.WriteString("arg")
			b.WriteString(strconv.Itoa(i))
		}
		if a.Type != "" {
			b.WriteString("(")
			b.WriteString(a.Type)
			b.WriteString(")")
		}
		b.WriteString("=")
		b.WriteString(fmt.Sprintf("%v", a.Value))
		if i != al-1 {
			b.WriteString(", ")
		}
	}
	b.WriteString(")")
	return b.String()
}

type ContractErrorArgument struct {
	Name  string
	Type  string
	Value any
}

func unpackError(err error, abi *abi.ABI) error {
	{
		var e rpc.DataError
		if errors.As(err, &e) {
			data := e.ErrorData()
			switch data := data.(type) {
			case string:
				b, decodeErr := hex.DecodeString(strings.TrimPrefix(data, "0x"))
				if decodeErr != nil {
					return fmt.Errorf("%w: %w", err, decodeErr)
				}
				contractErr := unpackErrorData(abi, b)
				if contractErr != nil {
					return fmt.Errorf("%w: %w", err, contractErr)
				}
				return err
			case []byte:
				contractErr := unpackErrorData(abi, data)
				if contractErr != nil {
					return fmt.Errorf("%w: %w", err, contractErr)
				}
				return err
			}
		}
	}

	{
		var e rpc.HTTPError
		if errors.As(err, &e) {
			var v struct {
				Error struct {
					Code    int    `json:"code"`
					Message string `json:"message"`
				} `json:"error"`
			}
			if err := json.Unmarshal(e.Body, &v); err != nil {
				return &Error{
					StatusCode: e.StatusCode,
					Message:    string(e.Body),
				}
			}
			if v.Error.Message == "" {
				return &Error{
					StatusCode: e.StatusCode,
					Message:    string(e.Body),
				}
			}
			return &Error{
				StatusCode: e.StatusCode,
				Code:       v.Error.Code,
				Message:    v.Error.Message,
			}
		}
	}

	return err
}

func unpackErrorData(abi *abi.ABI, data []byte) *ContractError {
	abiError, err := abi.ErrorByID([4]byte(data[:4]))
	if abiError == nil || err != nil {
		return nil
	}

	v, err := abiError.Unpack(data)
	if err == nil {
		values, ok := v.([]any)
		if !ok {
			return nil
		}
		output := ContractError{
			Name:      abiError.Name,
			Arguments: make([]ContractErrorArgument, 0, len(abiError.Inputs)),
		}
		for i, input := range abiError.Inputs {
			output.Arguments = append(output.Arguments, ContractErrorArgument{
				Name:  input.Name,
				Type:  input.Type.String(),
				Value: values[i],
			})
		}
		return &output
	}
	return nil
}
