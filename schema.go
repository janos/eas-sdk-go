// Copyright (c) 2024, Janoš Guljaš <janos@resenje.org>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eas

import (
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"resenje.org/taint"
)

func NewSchema(args ...any) (string, error) {
	schemas := make([]string, 0, len(args))
	for _, arg := range args {
		_, s, _, err := getABINewTypeArguments(arg, "", "", nil)
		if err != nil {
			return "", err
		}
		schemas = append(schemas, s)
	}
	if len(schemas) == 1 {
		return stripWrappedBrackets(schemas[0]), nil
	}
	return strings.Join(schemas, ", "), nil
}

func MustNewSchema(args ...any) string {
	s, err := NewSchema(args...)
	if err != nil {
		panic(err)
	}
	return s
}

func stripWrappedBrackets(s string) string {
	l := len(s)
	if l < 2 {
		return s
	}
	if s[0] != '(' {
		return s
	}
	if s[l-1] != ')' {
		return s
	}

	if !hasRedundantParentheses(s) {
		return s
	}

	return s[1 : l-1]
}

func hasRedundantParentheses(s string) bool {
	stack := make([]rune, 0)
	for _, char := range s {
		switch char {
		case '(':
			stack = append(stack, char)
		case ')':
			if len(stack) == 0 {
				// mismatched closing bracket
				return false
			}
			stack = stack[:len(stack)-1]
		}
	}

	// If the first and last characters are redundant parentheses, the stack
	// will be empty.
	return len(stack) == 0
}

func encodeAttestationValues(values []any) ([]byte, error) {
	var args abi.Arguments

	for i, v := range values {
		t, err := getABIType(v)
		if err != nil {
			return nil, fmt.Errorf("abi type for argument %v: %w", i, err)
		}
		args = append(args, abi.Argument{
			Type: t,
		})
	}

	data, err := args.Pack(values...)
	if err != nil {
		return nil, fmt.Errorf("pack abi: %w", err)
	}
	return data, nil
}

func scanAttestationValues(data []byte, args ...any) error {
	var abiArgs abi.Arguments

	for i, arg := range args {
		t, err := getABIType(reflect.Indirect(reflect.ValueOf(arg)).Interface())
		if err != nil {
			return fmt.Errorf("abi type for argument %v: %w", i, err)
		}
		abiArgs = append(abiArgs, abi.Argument{
			Type: t,
		})
	}

	values, err := abiArgs.Unpack(data)
	if err != nil {
		return fmt.Errorf("unpack abi: %w", err)
	}

	if len(args) != len(values) {
		return errors.New("unable to unpack all fields")
	}

	for i, arg := range args {
		if err := taint.InjectWithTag(values[i], arg, "abi"); err != nil {
			return fmt.Errorf("inject value for argument %v: %w", i, err)
		}
	}

	return nil
}

func getABIType(v any) (abi.Type, error) {
	ty, internalType, components, err := getABINewTypeArguments(v, "", "", nil)
	if err != nil {
		return abi.Type{}, err
	}
	return abi.NewType(ty, internalType, components)
}

func getABINewTypeArguments(v any, ty, internalType string, components []abi.ArgumentMarshaling) (string, string, []abi.ArgumentMarshaling, error) {
	switch v.(type) {
	case common.Address:
		return "address" + ty, "address" + internalType, components, nil
	case string:
		return "string" + ty, "string" + internalType, components, nil
	case bool:
		return "bool" + ty, "bool" + internalType, components, nil
	case [32]byte, UID:
		return "bytes32" + ty, "bytes32" + internalType, components, nil
	case []byte:
		return "bytes" + ty, "bytes" + internalType, components, nil
	case uint8:
		return "uint8" + ty, "uint8" + internalType, components, nil
	case uint16:
		return "uint16" + ty, "uint16" + internalType, components, nil
	case uint32:
		return "uint32" + ty, "uint32" + internalType, components, nil
	case uint64:
		return "uint64" + ty, "uint64" + internalType, components, nil
	case *big.Int:
		return "uint256" + ty, "uint256" + internalType, components, nil
	default:
		t := reflect.TypeOf(v)
		switch t.Kind() {
		case reflect.Struct:
			v := reflect.ValueOf(v)
			numFields := v.NumField()
			internalTypeParts := make([]string, 0, numFields)
			components := make([]abi.ArgumentMarshaling, 0, numFields)
			for i := 0; i < numFields; i++ {
				fv := v.Type().Field(i)
				var nv any
				switch v.Field(i).Kind() {
				case reflect.Slice:
					nv = reflect.MakeSlice(fv.Type, 0, 0).Interface()
				default:
					nv = reflect.Indirect(reflect.New(fv.Type)).Interface()
				}
				fty, fit, fc, err := getABINewTypeArguments(nv, "", "", nil)
				if err != nil {
					return "", "", nil, fmt.Errorf("unsupported type %T", v)
				}
				name := abiArgumentNameFromTag(fv.Tag)
				if name == "" {
					name = fv.Name
				}
				components = append(components, abi.ArgumentMarshaling{
					Name:         name,
					Type:         fty,
					InternalType: fit,
					Components:   fc,
				})
				internalTypeParts = append(internalTypeParts, fit+" "+name)
			}
			return "tuple" + ty, "(" + strings.Join(internalTypeParts, ", ") + ")" + internalType, components, nil
		case reflect.Slice:
			e := reflect.Indirect(reflect.New(t.Elem())).Interface()
			return getABINewTypeArguments(e, "[]"+ty, "[]"+internalType, components)
		case reflect.Array:
			len := reflect.ValueOf(v).Len()
			e := reflect.Indirect(reflect.New(t.Elem())).Interface()
			prefix := fmt.Sprintf("[%v]", len)
			return getABINewTypeArguments(e, prefix+ty, prefix+internalType, components)
		case reflect.Pointer:
			return getABINewTypeArguments(reflect.Indirect(reflect.New(t.Elem())).Interface(), ty, internalType, components)
		default:
			return "", "", nil, fmt.Errorf("unsupported type %T", v)
		}
	}
}

func abiArgumentNameFromTag(structTag reflect.StructTag) (keyName string) {
	tag := structTag.Get("abi")
	if tag == "" {
		return ""
	}
	return strings.Split(tag, ",")[0]
}
