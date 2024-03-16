// Copyright (c) 2024, Janoš Guljaš <janos@resenje.org>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eas

import (
	"crypto/ecdsa"
	"fmt"
	"io"
	"io/fs"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
)

func HexParsePrivateKey(h string) (*ecdsa.PrivateKey, error) {
	return crypto.HexToECDSA(h)
}

func LoadEthereumKeyFile(fs fs.FS, filename, auth string) (*ecdsa.PrivateKey, error) {
	f, err := fs.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	const maxFileLength = 1024

	data, err := io.ReadAll(io.LimitReader(f, maxFileLength))
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	k, err := keystore.DecryptKey(data, auth)
	if err != nil {
		return nil, fmt.Errorf("decrypt key: %w", err)
	}

	return k.PrivateKey, nil
}
