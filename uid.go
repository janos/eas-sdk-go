// Copyright (c) 2024, Janoš Guljaš <janos@resenje.org>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eas

import (
	"bytes"
	"encoding/hex"

	"github.com/ethereum/go-ethereum/common"
)

type UID [32]byte

func HexDecodeUID(s string) UID {
	return [32]byte(common.HexToHash(s))
}

func (u UID) String() string {
	return "0x" + hex.EncodeToString(u[:])
}

func (u UID) IsZero() bool {
	return u == zeroUID
}

func (u UID) MarshalText() ([]byte, error) {
	result := make([]byte, len(u)*2+2)
	copy(result, `0x`)
	hex.Encode(result[2:], u[:])
	return result, nil
}

func (u *UID) UnmarshalText(b []byte) error {
	b = bytes.TrimPrefix(b, []byte("0x"))
	_, err := hex.Decode(u[:], b)
	return err
}

var zeroUID = [32]byte{}
