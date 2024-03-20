// Copyright (c) 2024, Janoš Guljaš <janos@resenje.org>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eas_test

import (
	"testing"

	"resenje.org/eas"
)

func TestNewSchema(t *testing.T) {

	type Record struct {
		Key   eas.UID
		Value []byte `json:"val"`
	}

	type Schema struct {
		Message string
		Records []Record
	}

	got, err := eas.NewSchema(Schema{})
	assertNilError(t, err)
	assertEqual(t, "schema", got, "string Message, (bytes32 Key, bytes val)[] Record")
}
