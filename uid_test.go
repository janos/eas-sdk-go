// Copyright (c) 2024, Janoš Guljaš <janos@resenje.org>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eas_test

import (
	"encoding/json"
	"testing"

	"resenje.org/eas"
)

func TestUID(t *testing.T) {
	in := "0x6fa753be53bb614388f8a8116b139df78dac74311d6f91d6b561ac8517032170"
	got := eas.HexDecodeUID(in)

	want := [32]byte{111, 167, 83, 190, 83, 187, 97, 67, 136, 248, 168, 17, 107, 19, 157, 247, 141, 172, 116, 49, 29, 111, 145, 214, 181, 97, 172, 133, 23, 3, 33, 112}
	if got != want {
		t.Errorf("got uid %v, want %v", [32]byte(got), want)
	}

	assertEqual(t, "string", got.String(), in)
	assertEqual(t, "is zero", got.IsZero(), false)
}

func TestUID_invalid(t *testing.T) {
	in := "invalid hex string"
	got := eas.HexDecodeUID(in)

	want := [32]byte{}
	if got != want {
		t.Errorf("got uid %v, want %v", [32]byte(got), want)
	}

	assertEqual(t, "string", got.String(), "0x0000000000000000000000000000000000000000000000000000000000000000")
	assertEqual(t, "is zero", got.IsZero(), true)
}

func TestUID_marshalJSON(t *testing.T) {
	want := "0x6fa753be53bb614388f8a8116b139df78dac74311d6f91d6b561ac8517032170"
	u := eas.HexDecodeUID(want)

	b, err := json.Marshal(u)
	assertNilError(t, err)

	got := string(b)
	assertEqual(t, "encoded", got, `"`+want+`"`)
}

func TestUID_unmarshalJSON(t *testing.T) {
	var got eas.UID
	err := json.Unmarshal([]byte(`"0x6fa753be53bb614388f8a8116b139df78dac74311d6f91d6b561ac8517032170"`), &got)
	assertNilError(t, err)

	want := [32]byte{111, 167, 83, 190, 83, 187, 97, 67, 136, 248, 168, 17, 107, 19, 157, 247, 141, 172, 116, 49, 29, 111, 145, 214, 181, 97, 172, 133, 23, 3, 33, 112}
	if got != want {
		t.Errorf("got uid %v, want %v", [32]byte(got), want)
	}
}
