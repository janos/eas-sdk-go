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
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient/simulated"

	"resenje.org/eas"
	"resenje.org/eas/internal/deployment"
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

	sim := simulated.NewBackend(types.GenesisAlloc{
		accountAddress: {
			Balance: balance,
		},
	})
	t.Cleanup(func() {
		if err := sim.Close(); err != nil {
			t.Error(err)
		}
	})

	backend := sim.Client()

	// deploy schema registry contract
	_, _, _, wait, err := deployment.DeploySchemaRegistry(ctx, backend, privateKey)
	assertNilError(t, err)

	sim.Commit()

	schemaRegistryAddress, err := wait(ctx)
	assertNilError(t, err)

	// deploy eas contract
	_, _, _, wait, err = deployment.DeployEAS(ctx, backend, privateKey, schemaRegistryAddress)
	assertNilError(t, err)

	sim.Commit()

	easAddress, err := wait(ctx)
	assertNilError(t, err)

	// construct client
	c, err := eas.NewClient(ctx, "", privateKey, easAddress, &eas.Options{
		Backend: backend,
	})
	assertNilError(t, err)

	return &Client{
		Client:  c,
		account: accountAddress,
		backend: sim,
	}
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
