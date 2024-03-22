// Copyright (c) 2024, Janoš Guljaš <janos@resenje.org>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eas_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"resenje.org/eas"
)

func TestEASContract_Attest(t *testing.T) {
	client := newClient(t)
	ctx := context.Background()

	schemaUID := registerSchema(t, client, "string message")

	_, wait, err := client.EAS.Attest(ctx, schemaUID, &eas.AttestOptions{Revocable: true}, "Hello!")
	assertNilError(t, err)

	client.backend.Commit()

	r, err := wait(ctx)
	assertNilError(t, err)
	assertEqual(t, "schema uid", r.Schema, schemaUID)
	assertEqual(t, "attester", r.Attester, client.account)
}

func TestEASContract_GetAttestation(t *testing.T) {
	client := newClient(t)
	ctx := context.Background()

	schemaUID := registerSchema(t, client, "string message")

	_, wait, err := client.EAS.Attest(ctx, schemaUID, &eas.AttestOptions{Revocable: true}, "Hello!")
	assertNilError(t, err)

	client.backend.Commit()

	r, err := wait(ctx)
	assertNilError(t, err)

	a, err := client.EAS.GetAttestation(ctx, r.UID)
	assertNilError(t, err)

	assertEqual(t, "schema uid", a.Schema, schemaUID)
	assertEqual(t, "attester", a.Attester, client.account)
}

func TestEASContract_GetAttestation_structured(t *testing.T) {
	client := newClient(t)
	ctx := context.Background()

	type KV struct {
		Key   string
		Value string `abi:"somethingStrange"`
	}

	type Schema struct {
		ID      uint64
		Map     []KV
		Comment string
	}

	schemaUID := registerSchema(t, client, eas.MustNewSchema(Schema{}))

	attestationValues := Schema{
		ID: 3,
		Map: []KV{
			{"k1", "v1"},
			{"k2", "v2"},
		},
		Comment: "Hey",
	}

	_, wait, err := client.EAS.Attest(ctx, schemaUID, nil, attestationValues)
	assertNilError(t, err)

	client.backend.Commit()

	r, err := wait(ctx)
	assertNilError(t, err)

	a, err := client.EAS.GetAttestation(ctx, r.UID)
	assertNilError(t, err)

	assertEqual(t, "schema uid", a.Schema, schemaUID)
	assertEqual(t, "attester", a.Attester, client.account)

	var validationValues Schema
	err = a.ScanValues(&validationValues)

	assertNilError(t, err)
	assertEqual(t, "data", validationValues, attestationValues)
}

func TestEASContract_MultiAttest(t *testing.T) {
	client := newClient(t)
	ctx := context.Background()

	schemaUID := registerSchema(t, client, "string message")

	schemas := [][]any{
		{"one"},
		{"two"},
		{"three"},
		{"four"},
		{"five"},
	}

	_, wait, err := client.EAS.MultiAttest(ctx, schemaUID, &eas.AttestOptions{Revocable: true}, schemas...)
	assertNilError(t, err)

	client.backend.Commit()

	r, err := wait(ctx)
	assertNilError(t, err)

	count := 0
	for i, e := range r {
		a, err := client.EAS.GetAttestation(ctx, e.UID)
		assertNilError(t, err)

		assertNilError(t, err)
		assertEqual(t, "schema uid", e.Schema, schemaUID)
		assertEqual(t, "attester", e.Attester, client.account)

		var message string
		err = a.ScanValues(&message)
		assertNilError(t, err)
		assertEqual(t, "message", message, schemas[i][0].(string))

		count++
	}

	assertEqual(t, "count", count, len(schemas))
}

func TestEASContract_FilterAttested(t *testing.T) {
	client := newClient(t)
	ctx := context.Background()

	schemaUID := registerSchema(t, client, "string message")

	for i := 0; i < 10; i++ {
		attest(t, client, schemaUID, &eas.AttestOptions{Revocable: true}, fmt.Sprintf("Hello %v!", i))
	}

	t.Run("all", func(t *testing.T) {
		it, err := client.EAS.FilterAttested(ctx, 0, nil, nil, nil, nil)
		assertNilError(t, err)
		defer it.Close()

		count := 0

		for it.Next() {
			r := it.Value()

			a, err := client.EAS.GetAttestation(ctx, r.UID)
			assertNilError(t, err)

			var message string
			err = a.ScanValues(&message)
			assertNilError(t, err)

			assertEqual(t, "message", message, fmt.Sprintf("Hello %v!", count))

			count++
		}
		assertNilError(t, it.Error())

		assertEqual(t, "count", count, 10)
	})
}

func TestEASContract_WatchAttested(t *testing.T) {
	client := newClient(t)
	ctx := context.Background()

	schemaUID := registerSchema(t, client, "string message")

	sink := make(chan *eas.EASAttested)
	sub, err := client.EAS.WatchAttested(ctx, nil, sink, nil, nil, nil)
	assertNilError(t, err)

	count := 0

	go func() {
		defer sub.Unsubscribe()

		for i := 0; i < 10; i++ {
			attest(t, client, schemaUID, &eas.AttestOptions{Revocable: true}, fmt.Sprintf("Hello %v!", i))

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

			var message string
			err = a.ScanValues(&message)
			assertNilError(t, err)

			assertEqual(t, "message", message, fmt.Sprintf("Hello %v!", count))

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

	assertEqual(t, "count", count, 10)
}

func attest(t testing.TB, client *Client, schemaUID eas.UID, o *eas.AttestOptions, values ...any) eas.UID {
	t.Helper()

	ctx := context.Background()

	_, wait, err := client.EAS.Attest(ctx, schemaUID, o, values...)
	assertNilError(t, err)

	client.backend.Commit()

	r, err := wait(ctx)
	assertNilError(t, err)

	return r.UID
}
