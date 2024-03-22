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

	type record struct {
		Key        eas.UID
		Value      []byte   `abi:"val"`
		ABigNumber *big.Int `abi:"abn"`
	}

	type nestedSliceWithTuples struct {
		Message string `abi:"msg"`
		Records []record
	}

	type simpleTuple struct {
		T1 [2]string
		T2 *big.Int
		T3 eas.UID
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

	t.Run("all supported type in a tuple", newSchemaTest("address F1, string F2, bool F3, bytes32 F4, bytes32 F5, bytes F6, uint8 F7, uint16 F8, uint32 F9, uint64 FA, uint256 FB, (string[2] T1, uint256 T2, bytes32 T3) FC, address[] S1, string[] S2, bool[] S3, bytes32[] S4, bytes32[] S5, bytes[] S6, bytes S7, uint16[] S8, uint32[] S9, uint64[] SA, uint256[] SB, (string[2] T1, uint256 T2, bytes32 T3)[] SC, address[2] A1, string[2] A2, bool[2] A3, bytes32[2] A4, bytes32[2] A5, bytes[2] A6, uint8[2] A7, uint16[2] A8, uint32[2] A9, uint64[2] AA, uint256[2] AB, (string[2] T1, uint256 T2, bytes32 T3)[2] AC", struct {
		F1 common.Address
		F2 string
		F3 bool
		F4 [32]byte
		F5 eas.UID
		F6 []byte
		F7 uint8
		F8 uint16
		F9 uint32
		FA uint64
		FB *big.Int
		FC simpleTuple
		S1 []common.Address
		S2 []string
		S3 []bool
		S4 [][32]byte
		S5 []eas.UID
		S6 [][]byte
		S7 []uint8
		S8 []uint16
		S9 []uint32
		SA []uint64
		SB []*big.Int
		SC []simpleTuple
		A1 [2]common.Address
		A2 [2]string
		A3 [2]bool
		A4 [2][32]byte
		A5 [2]eas.UID
		A6 [2][]byte
		A7 [2]uint8
		A8 [2]uint16
		A9 [2]uint32
		AA [2]uint64
		AB [2]*big.Int
		AC [2]simpleTuple
	}{
		F1: common.Address{1, 2, 3, 4, 5},
		F2: "ethereum attestation service",
		F3: true,
		F4: [32]byte{10, 20, 30},
		F5: eas.UID{21, 22, 23, 24, 25},
		F6: []byte{0, 0, 2, 100, 25},
		F7: 7,
		F8: 8,
		F9: 9,
		FA: 10,
		FB: big.NewInt(11),
		FC: simpleTuple{
			T1: [2]string{"one", "two"},
			T2: big.NewInt(3),
			T3: eas.UID{1, 2, 3, 4},
		},
		S1: []common.Address{{1, 2, 3, 4, 5}, {9, 7, 8}},
		S2: []string{"ethereum attestation service"},
		S3: []bool{true, false, true},
		S4: [][32]byte{{10, 20, 30}, {0, 0, 5, 0, 1}},
		S5: []eas.UID{{21, 22, 23, 24, 25}},
		S6: [][]byte{{0, 0, 2, 100, 25}, {100, 7}},
		S7: []uint8{7, 8},
		S8: []uint16{8, 9, 10},
		S9: []uint32{9, 10, 11, 12},
		SA: []uint64{10, 11},
		SB: []*big.Int{big.NewInt(11), big.NewInt(12)},
		SC: []simpleTuple{
			{
				T1: [2]string{"one", "two"},
				T2: big.NewInt(3),
				T3: eas.UID{1, 2, 3, 4},
			},
			{
				T1: [2]string{"three", "four"},
				T2: big.NewInt(44),
				T3: eas.UID{0, 6, 8, 2, 6},
			},
			{
				T1: [2]string{"five", "six"},
				T2: big.NewInt(50),
				T3: eas.UID{8, 4, 39, 7},
			},
		},
		A1: [2]common.Address{{1, 2, 3, 4, 5}, {9, 7, 8}},
		A2: [2]string{"ethereum attestation service", "something"},
		A3: [2]bool{true, false},
		A4: [2][32]byte{{10, 20, 30}, {0, 0, 5, 0, 1}},
		A5: [2]eas.UID{{21, 22, 23, 24, 25}},
		A6: [2][]byte{{0, 0, 2, 100, 25}, {100, 7}},
		A7: [2]uint8{7, 8},
		A8: [2]uint16{8, 9},
		A9: [2]uint32{9, 10},
		AA: [2]uint64{10, 11},
		AB: [2]*big.Int{big.NewInt(11), big.NewInt(12)},
		AC: [2]simpleTuple{
			{
				T1: [2]string{"seven", "eight"},
				T2: big.NewInt(32),
				T3: eas.UID{1, 2, 3, 4},
			},
			{
				T1: [2]string{"nine", "ten"},
				T2: big.NewInt(99),
				T3: eas.UID{9, 3, 2},
			},
		},
	}))

	t.Run("nested slice of tuples",
		newSchemaTest("string msg, (bytes32 Key, bytes val, uint256 abn)[] Records", nestedSliceWithTuples{
			Message: "Hello, 世界",
			Records: []record{
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
