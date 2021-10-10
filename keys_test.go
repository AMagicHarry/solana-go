// Copyright 2021 github.com/gagliardetto
// This file has been modified by github.com/gagliardetto
//
// Copyright 2020 dfuse Platform Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package solana

import (
	"encoding/hex"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPublicKeyFromBytes(t *testing.T) {
	tests := []struct {
		name     string
		inHex    string
		expected PublicKey
	}{
		{
			"empty",
			"",
			MustPublicKeyFromBase58("11111111111111111111111111111111"),
		},
		{
			"smaller than required",
			"010203040506",
			MustPublicKeyFromBase58("4wBqpZM9k69W87zdYXT2bRtLViWqTiJV3i2Kn9q7S6j"),
		},
		{
			"equal to 32 bytes",
			"0102030405060102030405060102030405060102030405060102030405060101",
			MustPublicKeyFromBase58("4wBqpZM9msxygzsdeLPq6Zw3LoiAxJk3GjtKPpqkcsi"),
		},
		{
			"longer than required",
			"0102030405060102030405060102030405060102030405060102030405060101FFFFFFFFFF",
			MustPublicKeyFromBase58("4wBqpZM9msxygzsdeLPq6Zw3LoiAxJk3GjtKPpqkcsi"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			bytes, err := hex.DecodeString(test.inHex)
			require.NoError(t, err)

			actual := PublicKeyFromBytes(bytes)
			assert.Equal(t, test.expected, actual, "%s != %s", test.expected, actual)
		})
	}
}

func TestPublicKeyFromBase58(t *testing.T) {
	tests := []struct {
		name        string
		in          string
		expected    PublicKey
		expectedErr error
	}{
		{
			"hand crafted",
			"SerumkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA",
			MustPublicKeyFromBase58("SerumkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA"),
			nil,
		},
		{
			"hand crafted error",
			"SerkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA",
			zeroPublicKey,
			errors.New("invalid length, expected 32, got 30"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := PublicKeyFromBase58(test.in)
			if test.expectedErr == nil {
				require.NoError(t, err)
				assert.Equal(t, test.expected, actual)
			} else {
				assert.Equal(t, test.expectedErr, err)
			}
		})
	}
}

func TestPrivateKeyFromSolanaKeygenFile(t *testing.T) {
	tests := []struct {
		inFile      string
		expected    PrivateKey
		expectedPub PublicKey
		expectedErr error
	}{
		{
			"testdata/standard.solana-keygen.json",
			MustPrivateKeyFromBase58("66cDvko73yAf8LYvFMM3r8vF5vJtkk7JKMgEKwkmBC86oHdq41C7i1a2vS3zE1yCcdLLk6VUatUb32ZzVjSBXtRs"),
			MustPublicKeyFromBase58("F8UvVsKnzWyp2nF8aDcqvQ2GVcRpqT91WDsAtvBKCMt9"),
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.inFile, func(t *testing.T) {
			actual, err := PrivateKeyFromSolanaKeygenFile(test.inFile)
			if test.expectedErr == nil {
				require.NoError(t, err)
				assert.Equal(t, test.expected, actual)
				assert.Equal(t, test.expectedPub, actual.PublicKey(), "%s != %s", test.expectedPub, actual.PublicKey())

			} else {
				assert.Equal(t, test.expectedErr, err)
			}
		})
	}
}

func TestPublicKey_MarshalText(t *testing.T) {
	keyString := "4wBqpZM9k69W87zdYXT2bRtLViWqTiJV3i2Kn9q7S6j"
	keyParsed := MustPublicKeyFromBase58(keyString)

	var key PublicKey
	err := key.UnmarshalText([]byte(keyString))
	require.NoError(t, err)

	assert.True(t, keyParsed.Equals(key))

	keyText, err := key.MarshalText()
	require.NoError(t, err)
	assert.Equal(t, []byte(keyString), keyText)

	type IdentityToSlotsBlocks map[PublicKey][2]int64

	var payload IdentityToSlotsBlocks
	data := `{"` + keyString + `":[3,4]}`
	err = json.Unmarshal([]byte(data), &payload)
	require.NoError(t, err)

	assert.Equal(t,
		IdentityToSlotsBlocks{
			keyParsed: [2]int64{3, 4},
		},
		payload,
	)
}

func TestPublicKeySlice(t *testing.T) {
	slice := make(PublicKeySlice, 0)
	require.False(t, slice.Has(BPFLoaderProgramID))

	slice.Append(BPFLoaderProgramID)
	require.True(t, slice.Has(BPFLoaderProgramID))
	require.Len(t, slice, 1)

	slice.UniqueAppend(BPFLoaderProgramID)
	require.Len(t, slice, 1)
	slice.Append(ConfigProgramID)
	require.Len(t, slice, 2)
	require.True(t, slice.Has(ConfigProgramID))
}

func TestIsNativeProgramID(t *testing.T) {
	require.True(t, isNativeProgramID(ConfigProgramID))
}

func TestCreateWithSeed(t *testing.T) {
	{
		got, err := CreateWithSeed(PublicKey{}, "limber chicken: 4/45", PublicKey{})
		require.NoError(t, err)
		require.True(t, got.Equals(MustPublicKeyFromBase58("9h1HyLCW5dZnBVap8C5egQ9Z6pHyjsh5MNy83iPqqRuq")))
	}
}

func TestCreateProgramAddress(t *testing.T) {
	program_id := MustPublicKeyFromBase58("BPFLoaderUpgradeab1e11111111111111111111111")
	public_key := MustPublicKeyFromBase58("SeedPubey1111111111111111111111111111111111")

	{
		got, err := CreateProgramAddress([][]byte{
			{},
			{1},
		},
			program_id,
		)
		require.NoError(t, err)
		require.True(t, got.Equals(MustPublicKeyFromBase58("BwqrghZA2htAcqq8dzP1WDAhTXYTYWj7CHxF5j7TDBAe")))
	}

	{
		got, err := CreateProgramAddress([][]byte{
			[]byte("☉"),
			{0},
		},
			program_id,
		)
		require.NoError(t, err)
		require.True(t, got.Equals(MustPublicKeyFromBase58("13yWmRpaTR4r5nAktwLqMpRNr28tnVUZw26rTvPSSB19")))
	}

	{
		got, err := CreateProgramAddress([][]byte{
			[]byte("Talking"),
			[]byte("Squirrels"),
		},
			program_id,
		)
		require.NoError(t, err)
		require.True(t, got.Equals(MustPublicKeyFromBase58("2fnQrngrQT4SeLcdToJAD96phoEjNL2man2kfRLCASVk")))
	}

	{
		got, err := CreateProgramAddress([][]byte{
			public_key[:],
			{1},
		},
			program_id,
		)
		require.NoError(t, err)
		require.True(t, got.Equals(MustPublicKeyFromBase58("976ymqVnfE32QFe6NfGDctSvVa36LWnvYxhU6G2232YL")))
	}
}

// https://github.com/solana-labs/solana/blob/216983c50e0a618facc39aa07472ba6d23f1b33a/sdk/program/src/pubkey.rs#L590
func TestFindProgramAddress(t *testing.T) {
	for i := 0; i < 1_000; i++ {

		program_id := NewWallet().PrivateKey.PublicKey()
		address, bump_seed, err := FindProgramAddress(
			[][]byte{
				[]byte("Lil'"),
				[]byte("Bits"),
			},
			program_id,
		)
		require.NoError(t, err)

		got, err := CreateProgramAddress(
			[][]byte{
				[]byte("Lil'"),
				[]byte("Bits"),
				[]byte{bump_seed},
			},
			program_id,
		)
		require.NoError(t, err)
		require.Equal(t, address, got)
	}
}

func TestFindTokenMetadataAddress(t *testing.T) {
	// Zuuper Grapes (TOILET)
	// https://solscan.io/token/77K8mr457qxUSSNSfi4sSj5euP8DyuJJWHAUQVW8QCp3
	mint := MustPublicKeyFromBase58("77K8mr457qxUSSNSfi4sSj5euP8DyuJJWHAUQVW8QCp3")
	metadataPDA, bumpSeed, err := FindTokenMetadataAddress(mint)
	require.NoError(t, err)
	// https://solscan.io/account/GfihrEYCPrvUyrMyMQPdhGEStxa9nKEK2Wfn9iK4AZq2
	assert.Equal(t, metadataPDA, MustPublicKeyFromBase58("GfihrEYCPrvUyrMyMQPdhGEStxa9nKEK2Wfn9iK4AZq2"))
	assert.Equal(t, bumpSeed, uint8(0xfd))
}
