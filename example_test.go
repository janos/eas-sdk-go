// Copyright (c) 2024, Janoš Guljaš <janos@resenje.org>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package eas_test

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"resenje.org/eas"
)

var (
	endpointSepolia        = "https://ethereum-sepolia-rpc.publicnode.com/"
	contractAddressSepolia = common.HexToAddress("0xC2679fBD37d54388Ce493F1DB75320D236e1815e")
)

func ExampleEASContract_GetAttestation() {
	ctx := context.Background()

	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}

	c, err := eas.NewClient(ctx, endpointSepolia, privateKey, contractAddressSepolia, nil)
	if err != nil {
		log.Fatal(err)
	}

	attestationUID := eas.HexDecodeUID("0xac812932f5cee90a457d57a9fbd7b142b21ba99b809f982bbf86947f295281ff")

	a, err := c.EAS.GetAttestation(ctx, attestationUID)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Attestation UID", a.UID)
	log.Println("Attestation Time", a.Time)

	var schemaUID eas.UID
	var name string

	if err := a.ScanValues(&schemaUID, &name); err != nil {
		log.Fatal(err)
	}

	log.Println("Attestation")
	log.Println("Schema UID:", schemaUID)
	log.Println("Name:", name)
}

type RoadTrip struct {
	ID           uint64      `abi:"id"`
	VIN          string      `abi:"vin"` // Vehicle Identification Number
	VehicleOwner string      `abi:"vehicleOwner"`
	Passengers   []Passenger `abi:"passengers"`
}

type Passenger struct {
	Name     string `abi:"name"`
	CanDrive bool   `abi:"canDrive"`
}

func (p Passenger) CanDriveString() string {
	if p.CanDrive {
		return "can drive"
	}
	return "cannot drive"
}

// Attest a road trip with a structured schema.
func Example_structuredAttestation() {
	ctx := context.Background()

	// Use a fake key here. Use your own funded key to be able to send transactions.
	privateKey, err := eas.HexParsePrivateKey("a896e1f28a6453e8db4794f11ea185befd04c4e4f06790e37e8d1cc90a611948")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Wallet address:", crypto.PubkeyToAddress(privateKey.PublicKey))

	// Construct a client that will interact with EAS contracts.
	c, err := eas.NewClient(ctx, endpointSepolia, privateKey, contractAddressSepolia, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Create the Schema on chain.
	tx, waitRegistration, err := c.SchemaRegistry.Register(ctx, eas.MustNewSchema(RoadTrip{}), common.Address{}, true)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Waiting schema registration transaction:", tx.Hash())
	schemaRegistration, err := waitRegistration(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Just check the schema definition.
	schema, err := c.SchemaRegistry.GetSchema(ctx, schemaRegistration.UID)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Schema UID:", schema.UID)
	log.Println("Schema:", schema.Schema)

	// Attest a road trip on chain.
	tx, waitAttestation, err := c.EAS.Attest(ctx,
		schema.UID,
		&eas.AttestOptions{Revocable: true},
		RoadTrip{
			ID:           4242,
			VIN:          "1FA6P8CF5L5100421",
			VehicleOwner: "Richard Hammond",
			Passengers: []Passenger{
				{
					Name:     "James May",
					CanDrive: true,
				},
				{
					Name:     "Jeremy Clarkson",
					CanDrive: true,
				},
				{
					Name:     "The Stig",
					CanDrive: false,
				},
			},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Waiting attest transaction:", tx.Hash())
	attestConfirmation, err := waitAttestation(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Get the attestation to verify it.
	a, err := c.EAS.GetAttestation(ctx, attestConfirmation.UID)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Attestation UID", a.UID)
	log.Println("Attestation Time", a.Time)

	var roadTrip RoadTrip
	if err := a.ScanValues(&roadTrip); err != nil {
		log.Fatal(err)
	}

	log.Println("Road trip:", roadTrip.ID)
	log.Println("Vehicle Identification Number:", roadTrip.VIN)
	log.Println("Vehicle owner:", roadTrip.VehicleOwner)
	for i, p := range roadTrip.Passengers {
		log.Printf("Passenger %v: %s (%s)", i, p.Name, p.CanDriveString())
	}
}
