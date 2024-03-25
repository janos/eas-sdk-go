// Copyright (c) 2024, Janoš Guljaš <janos@resenje.org>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eas_test

import (
	"context"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient/simulated"

	"resenje.org/eas"
	"resenje.org/eas/eastest"
)

type Client struct {
	*eas.Client
	account common.Address
	backend *simulated.Backend
}

func newClient(t testing.TB) *Client {
	t.Helper()

	ctx := context.Background()

	privateKey, err := crypto.GenerateKey()
	assertNilError(t, err)

	balance := new(big.Int)
	balance.SetString("100000000000000000000", 10)

	accountAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	sim, easAddress := eastest.NewSimulatedBackend(t, map[common.Address]*big.Int{
		accountAddress: balance,
	})

	// construct client
	c, err := eas.NewClient(ctx, "", privateKey, easAddress, &eas.Options{
		Backend: sim.Client(),
	})
	assertNilError(t, err)

	return &Client{
		Client:  c,
		account: accountAddress,
		backend: sim,
	}
}

func TestClient_Address(t *testing.T) {
	c := newClient(t)

	var zeroAddress common.Address
	if c.account == zeroAddress {
		t.Error("zero address account")
	}

	assertEqual(t, "address", c.account, c.Address())
}

func assertEqual[T any](t testing.TB, name string, got, want T) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %s %+v, want %+v", name, got, want)
	}
}

func assertNilError(t testing.TB, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("got error %[1]T %[1]q", err)
	}
}
