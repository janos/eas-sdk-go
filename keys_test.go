// Copyright (c) 2024, Janoš Guljaš <janos@resenje.org>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eas_test

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"resenje.org/eas"
)

func TestHexParsePrivateKey(t *testing.T) {
	pk, err := eas.HexParsePrivateKey("933c798b990a6be3fb91ae2fd3b6593f61d6d478548091205ee948b1de9c9f19")
	assertNilError(t, err)

	got := crypto.FromECDSA(pk)
	want := []byte{147, 60, 121, 139, 153, 10, 107, 227, 251, 145, 174, 47, 211, 182, 89, 63, 97, 214, 212, 120, 84, 128, 145, 32, 94, 233, 72, 177, 222, 156, 159, 25}

	if !bytes.Equal(got, want) {
		t.Error("invalid public key")
	}
}

func TestLoadEthereumKeyFile(t *testing.T) {
	pk, err := eas.LoadEthereumKeyFile(os.DirFS("./testdata"), "UTC--2024-03-16T20-25-58.090Z--0x08752c431c3e38b12e94a0f195b166590e764831", "12345")
	if err != nil {
		log.Fatal(err)
	}

	got := crypto.FromECDSA(pk)

	want := []byte{0, 221, 217, 75, 162, 191, 163, 150, 69, 242, 230, 214, 174, 16, 231, 254, 222, 224, 0, 39, 194, 195, 237, 208, 30, 45, 38, 149, 107, 69, 106, 10}

	if !bytes.Equal(got, want) {
		t.Error("invalid public key")
	}
}
