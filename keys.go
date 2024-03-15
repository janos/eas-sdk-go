// Copyright (c) 2024, Janoš Guljaš <janos@resenje.org>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eas

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/crypto"
)

func HexParsePrivateKey(h string) (*ecdsa.PrivateKey, error) {
	return crypto.HexToECDSA(h)
}
