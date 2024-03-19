// Copyright (c) 2024, Janoš Guljaš <janos@resenje.org>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eas_test

import (
	"context"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"resenje.org/eas"
)

func TestEASContract_Revoke(t *testing.T) {
	client := newClient(t)
	ctx := context.Background()

	schemaUID := newSchema(t, client, "string message")

	attestationUID := newAttestation(t, client, schemaUID, &eas.AttestOptions{Revocable: true}, "Hello!")

	a, err := client.EAS.GetAttestation(ctx, attestationUID)
	assertNilError(t, err)

	assertEqual(t, "revocation time", a.RevocationTime, eas.ZeroTime)

	tx, wait, err := client.EAS.Revoke(ctx, schemaUID, attestationUID, nil)
	assertNilError(t, err)

	client.backend.Commit()

	_, isPending, err := client.backend.Client().(ethereum.TransactionReader).TransactionByHash(ctx, tx.Hash())
	assertNilError(t, err)

	assertEqual(t, "is pending", isPending, false)

	r, err := wait(ctx)
	assertNilError(t, err)

	assertEqual(t, "attestation uid", r.UID, attestationUID)
	assertEqual(t, "schema uid", r.Schema, schemaUID)

	a, err = client.EAS.GetAttestation(ctx, attestationUID)
	assertNilError(t, err)

	if time.Since(a.RevocationTime) > time.Minute {
		t.Errorf("too old revocation time %v", a.RevocationTime)
	}
}

func TestEASContract_MultiRevoke(t *testing.T) {
	client := newClient(t)
	ctx := context.Background()

	schemaUID := newSchema(t, client, "string message")

	var attestationUIDs []eas.UID
	for i := 0; i < 10; i++ {
		attestationUIDs = append(attestationUIDs, newAttestation(t, client, schemaUID, &eas.AttestOptions{Revocable: true}, "Hello!"))
	}

	for _, uid := range attestationUIDs {
		a, err := client.EAS.GetAttestation(ctx, uid)
		assertNilError(t, err)

		assertEqual(t, "revocation time", a.RevocationTime, eas.ZeroTime)
	}

	_, wait, err := client.EAS.MultiRevoke(ctx, schemaUID, []eas.UID{attestationUIDs[2], attestationUIDs[3], attestationUIDs[5]})
	assertNilError(t, err)

	client.backend.Commit()

	_, err = wait(ctx)
	assertNilError(t, err)

	for i, uid := range attestationUIDs {
		a, err := client.EAS.GetAttestation(ctx, uid)
		assertNilError(t, err)

		if i == 2 || i == 3 || i == 5 {
			if time.Since(a.RevocationTime) > time.Minute {
				t.Errorf("too old revocation time %v", a.RevocationTime)
			}
		} else {
			assertEqual(t, "revocation time", a.RevocationTime, eas.ZeroTime)
		}
	}
}

func TestEASContract_FilterRevoked(t *testing.T) {
	client := newClient(t)
	ctx := context.Background()

	schemaUID := newSchema(t, client, "string message")

	var attestationUIDs []eas.UID
	for i := 0; i < 10; i++ {
		attestationUIDs = append(attestationUIDs, newAttestation(t, client, schemaUID, &eas.AttestOptions{Revocable: true}, "Hello!"))
	}

	for _, i := range []int{1, 3, 4, 6, 7} {
		_, wait, err := client.EAS.Revoke(ctx, schemaUID, attestationUIDs[i], nil)
		assertNilError(t, err)

		client.backend.Commit()

		_, err = wait(ctx)
		assertNilError(t, err)
	}

	t.Run("all", func(t *testing.T) {
		it, err := client.EAS.FilterRevoked(ctx, 0, nil, nil, nil, nil)
		assertNilError(t, err)
		defer it.Close()

		count := 0

		for it.Next() {
			r := it.Value()

			a, err := client.EAS.GetAttestation(ctx, r.UID)
			assertNilError(t, err)

			if time.Since(a.RevocationTime) > time.Minute {
				t.Errorf("too old revocation time %v", a.RevocationTime)
			}

			count++
		}
		assertNilError(t, it.Error())

		assertEqual(t, "count", count, 5)
	})
}

func TestEASContract_WatchRevoked(t *testing.T) {
	client := newClient(t)
	ctx := context.Background()

	schemaUID := newSchema(t, client, "string message")

	var attestationUIDs []eas.UID
	for i := 0; i < 10; i++ {
		attestationUIDs = append(attestationUIDs, newAttestation(t, client, schemaUID, &eas.AttestOptions{Revocable: true}, "Hello!"))
	}

	sink := make(chan *eas.EASRevoked)
	sub, err := client.EAS.WatchRevoked(ctx, nil, sink, nil, nil, nil)
	assertNilError(t, err)

	count := 0

	go func() {
		defer sub.Unsubscribe()

		for _, i := range []int{1, 3, 4, 6, 7} {
			_, wait, err := client.EAS.Revoke(ctx, schemaUID, attestationUIDs[i], nil)
			assertNilError(t, err)

			client.backend.Commit()

			_, err = wait(ctx)
			assertNilError(t, err)

			select {
			case <-time.After(100 * time.Millisecond):
			case <-ctx.Done():
			}
		}
	}()

loop:
	for {
		select {
		case r := <-sink:
			a, err := client.EAS.GetAttestation(ctx, r.UID)
			assertNilError(t, err)

			if time.Since(a.RevocationTime) > time.Minute {
				t.Errorf("too old revocation time %v", a.RevocationTime)
			}

			count++
		case err, ok := <-sub.Err():
			if !ok {
				break loop
			}
			if err != nil {
				t.Fatal(err)
			}
		case <-ctx.Done():
			if err != nil {
				t.Fatal(err)
			}
		}
	}

	assertEqual(t, "count", count, 5)
}
