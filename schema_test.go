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
	"resenje.org/eas"
)

func TestSchema(t *testing.T) {
	type Record struct {
		Key        eas.UID
		Value      []byte   `abi:"val"`
		ABigNumber *big.Int `abi:"abn"`
	}

	type NestedSliceWithTuples struct {
		Message string `abi:"msg"`
		Records []Record
	}

	t.Run("address", newSchemaTest("address", common.Address{1, 2, 3}))
	t.Run("string", newSchemaTest("string", "Message"))
	t.Run("bool", newSchemaTest("bool", true))
	t.Run("bytes32", newSchemaTest("bytes32", [32]byte{4, 5}))
	t.Run("uid", newSchemaTest("bytes32", eas.UID{6, 7}))
	t.Run("bytes", newSchemaTest("bytes", []byte{25, 25, 45}))
	t.Run("uint8", newSchemaTest("uint8", uint8(2)))
	t.Run("uint16", newSchemaTest("uint16", uint16(2)))
	t.Run("uint32", newSchemaTest("uint32", uint32(2)))
	t.Run("uint64", newSchemaTest("uint64", uint64(2)))
	t.Run("uint256", newSchemaTest("uint256", big.NewInt(44)))

	t.Run("nested slice of tuples",
		newSchemaTest[NestedSliceWithTuples]("string msg, (bytes32 Key, bytes val, uint256 abn)[] Records", NestedSliceWithTuples{
			Message: "Hello, 世界",
			Records: []Record{
				{
					Key:        eas.HexDecodeUID("0xb4d0ab81afc3474119212c28a8303ae693510a13d4024dae15eae99a59e2aa7c"),
					Value:      []byte{1, 2, 3, 4, 5, 6, 7},
					ABigNumber: big.NewInt(100),
				},
				{
					Key:        eas.HexDecodeUID("0x2d65177371a8c15a75480b2dedd127d52d9ba9f2aa78bf52886e102a6a76333d"),
					Value:      []byte{0, 0, 1, 0, 0},
					ABigNumber: big.NewInt(520),
				},
			},
		}),
	)

	t.Run("slice", newSchemaTest("string[]", []string{"1", "2"}))

	t.Run("array", newSchemaTest("string[2]", [2]string{"4", "5"}))

	t.Run("nested slices and arrays",
		newSchemaTest("string[2][][3]", [3][][2]string{
			{},
			{
				{"", "ab"},
				{"cd", "efg"},
				{"a", ""},
			},
			{
				{"one", "two"},
			},
		}),
	)
}

func newSchemaTest[T any](wantSchema string, wantAttestation T) func(*testing.T) {
	return func(t *testing.T) {
		client := newClient(t)
		ctx := context.Background()

		var schemaType T
		got, err := eas.NewSchema(schemaType)
		assertNilError(t, err)
		assertEqual(t, "schema", got, wantSchema)

		_, waitRegister, err := client.SchemaRegistry.Register(ctx, eas.MustNewSchema(schemaType), common.Address{}, false)
		assertNilError(t, err)

		client.backend.Commit()

		schemaRegistration, err := waitRegister(ctx)
		assertNilError(t, err)

		schema, err := client.SchemaRegistry.GetSchema(ctx, schemaRegistration.UID)
		assertNilError(t, err)

		assertEqual(t, "schema", schema.Schema, wantSchema)

		_, waitAttest, err := client.EAS.Attest(ctx, schema.UID, nil, wantAttestation)
		assertNilError(t, err)

		client.backend.Commit()

		attestConfirmation, err := waitAttest(ctx)
		assertNilError(t, err)

		a, err := client.EAS.GetAttestation(ctx, attestConfirmation.UID)
		assertNilError(t, err)

		var gotAttestation T
		err = a.ScanValues(&gotAttestation)
		assertNilError(t, err)
		assertEqual(t, "attestation", gotAttestation, wantAttestation)

		gotAttestationReflection := reflect.New(reflect.TypeOf(wantAttestation)).Interface()
		err = a.ScanValues(gotAttestationReflection)
		assertNilError(t, err)
		var wantAttestationReflection any = &wantAttestation
		assertEqual(t, "attestation reflection", gotAttestationReflection, wantAttestationReflection)

		var gotAttestationInterface any = new(T)
		var wantAttestationInterface any = &wantAttestation
		err = a.ScanValues(gotAttestationInterface)
		assertNilError(t, err)
		assertEqual(t, "attestation interface", gotAttestationInterface, wantAttestationInterface)
	}
}
