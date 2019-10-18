// Go Substrate RPC Client (GSRPC) provides APIs and types around Polkadot and any Substrate-based chain RPC calls
// Copyright (C) 2019  Centrifuge GmbH
//
// This file is part of Go Substrate RPC Client (GSRPC).
//
// GSRPC is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// GSRPC is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package types_test

import (
	"fmt"
	"testing"

	. "github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/stretchr/testify/assert"
)

var examplePhaseApp = Phase{
	IsApplyExtrinsic: true,
	AsApplyExtrinsic: 42,
}

var examplePhaseFin = Phase{
	IsFinalization: true,
}

var exampleEventApp = EventSystemExtrinsicSuccess{
	Phase:  examplePhaseApp,
	Topics: []Hash{{1, 2}},
}

var exampleEventFin = EventSystemExtrinsicSuccess{
	Phase:  examplePhaseFin,
	Topics: []Hash{{1, 2}},
}

var exampleEventFinEnc = []byte{0x1, 0x4, 0x1, 0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0} //nolint:lll

func TestEventSystemExtrinsicSuccess_EncodedLength(t *testing.T) {
	assertEncodedLength(t, []encodedLengthAssert{
		{exampleEventApp, 38},
		{exampleEventFin, 34},
	})
}

func TestEventSystemExtrinsicSuccess_Encode(t *testing.T) {
	encoded, err := EncodeToBytes(exampleEventFin)
	assert.NoError(t, err)
	assert.Equal(t, exampleEventFinEnc, encoded)
}

func TestEventSystemExtrinsicSuccess_Decode(t *testing.T) {
	decoded := EventSystemExtrinsicSuccess{}
	err := DecodeFromBytes(exampleEventFinEnc, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, exampleEventFin, decoded)
}

func TestEventSystemExtrinsicSuccess_Hash(t *testing.T) {
	assertHash(t, []hashAssert{
		{exampleEventFin, MustHexDecodeString(
			"0xfb1a0568e74c9e2ed9ec6a7cca8b680a24ca442e5cf391ca6d863e3b35a4c962")},
	})
}

func ExampleEventRecordsRaw_Decode() {
	e := EventRecordsRaw(MustHexDecodeString("0x100000000000000000000100000000000000020000000302d43593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d8eaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a48266d00000000000000000000000000000010a5d4e80000000000000000000000000002000000000000")) //nolint:lll

	events := EventRecords{}
	err := e.Decode(ExamplaryMetadataV8, &events)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Got %v System_ExtrinsicSuccess events\n", len(events.System_ExtrinsicSuccess))
	fmt.Printf("Got %v System_ExtrinsicFailed events\n", len(events.System_ExtrinsicFailed))
	fmt.Printf("Got %v Indices_NewAccountIndex events\n", len(events.Indices_NewAccountIndex))
	fmt.Printf("Got %v Balances_Transfer events\n", len(events.Balances_Transfer))
	t := events.Balances_Transfer[0]
	fmt.Printf("Transfer: %v tokens from %#x to\n%#x with a fee of %v", t.Value, t.From, t.To, t.Fees)

	// Output: Got 1 System_ExtrinsicSuccess events
	// Got 1 System_ExtrinsicFailed events
	// Got 1 Indices_NewAccountIndex events
	// Got 1 Balances_Transfer events
	// Transfer: 109 tokens from 0x3593c715fdd31c61141abd04a99fd6822c8558854ccde39a5684e7a56da27d8e to
	// 0xaf04151687736326c9fea17e25fc5287613693c912909cb226aa4794f26a4826 with a fee of 3906250000
}

func TestDispatchError(t *testing.T) {
	assertRoundtrip(t, DispatchError{HasModule: true, Module: 0xf1, Error: 0xa2})
	assertRoundtrip(t, DispatchError{HasModule: false, Error: 0xa2})
}
