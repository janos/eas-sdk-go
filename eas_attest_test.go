// Copyright (c) 2024, Janoš Guljaš <janos@resenje.org>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eas_test

import (
	"context"
	"testing"

	"resenje.org/eas"
)

func TestEASContract_Attest(t *testing.T) {
	client := newClient(t)
	ctx := context.Background()

	schemaUID := newSchema(t, client, "string message")

	_, wait, err := client.EAS.Attest(ctx, schemaUID, &eas.AttestOptions{Revocable: true}, "Hello!")
	assertNilError(t, err)

	client.backend.Commit()

	r, err := wait(ctx)
	assertNilError(t, err)
	assertEqual(t, "schema uid", r.Schema, schemaUID)
	assertEqual(t, "attester", r.Attester, client.account)
}

func TestEASContract_MultiAttest(t *testing.T) {
	client := newClient(t)
	ctx := context.Background()

	schemaUID := newSchema(t, client, "string message")

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

		fields, err := a.Fields("string message")
		assertNilError(t, err)

		assertEqual(t, "message", fields[0].Name, "message")
		assertEqual(t, "message", fields[0].Type, "string")
		assertEqual(t, "message", fields[0].Value, schemas[i][0])

		var message string
		err = a.ScanValues(&message)
		assertNilError(t, err)
		assertEqual(t, "message", message, schemas[i][0].(string))

		count++
	}

	assertEqual(t, "count", count, len(schemas))
}

func newAttestation(t testing.TB, client *Client, schemaUID eas.UID, o *eas.AttestOptions, values ...any) eas.UID {
	t.Helper()
	ctx := context.Background()

	_, wait, err := client.EAS.Attest(ctx, schemaUID, o, values...)
	assertNilError(t, err)

	client.backend.Commit()

	r, err := wait(ctx)
	assertNilError(t, err)

	return r.UID
}
