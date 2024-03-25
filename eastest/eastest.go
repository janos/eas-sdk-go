// Copyright (c) 2024, Janoš Guljaš <janos@resenje.org>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eastest

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient/simulated"

	"resenje.org/eas/internal/deployment"
)

// NewSimulatedBackend returns a simulated backend that has accounts with
// provided balances set on it and EAS contracts deployed on the returned
// address.
func NewSimulatedBackend(t testing.TB, accounts map[common.Address]*big.Int) (*simulated.Backend, common.Address) {
	t.Helper()

	ctx := context.Background()

	privateKey, err := crypto.GenerateKey()
	assertNilError(t, err)

	balance := new(big.Int)
	balance.SetString("100000000000000000000", 10)

	accountAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	alloc := types.GenesisAlloc{
		accountAddress: {
			Balance: balance,
		},
	}

	for address, balance := range accounts {
		alloc[address] = types.Account{
			Balance: balance,
		}
	}

	sim := simulated.NewBackend(alloc)
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

	return sim, easAddress
}

func assertNilError(t testing.TB, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("got error %[1]T %[1]q", err)
	}
}
