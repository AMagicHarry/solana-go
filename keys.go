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
	"crypto"
	"crypto/ed25519"
	crypto_rand "crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io/ioutil"
	"math"

	"filippo.io/edwards25519"
	"github.com/mr-tron/base58"
)

type PrivateKey []byte

func MustPrivateKeyFromBase58(in string) PrivateKey {
	out, err := PrivateKeyFromBase58(in)
	if err != nil {
		panic(err)
	}
	return out
}

func PrivateKeyFromBase58(privkey string) (PrivateKey, error) {
	res, err := base58.Decode(privkey)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func PrivateKeyFromSolanaKeygenFile(file string) (PrivateKey, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("read keygen file: %w", err)
	}

	var values []byte
	err = json.Unmarshal(content, &values)
	if err != nil {
		return nil, fmt.Errorf("decode keygen file: %w", err)
	}

	return PrivateKey([]byte(values)), nil
}

func (k PrivateKey) String() string {
	return base58.Encode(k)
}

func NewRandomPrivateKey() (PrivateKey, error) {
	pub, priv, err := ed25519.GenerateKey(crypto_rand.Reader)
	if err != nil {
		return nil, err
	}
	var publicKey PublicKey
	copy(publicKey[:], pub)
	return PrivateKey(priv), nil
}

func (k PrivateKey) Sign(payload []byte) (Signature, error) {
	p := ed25519.PrivateKey(k)
	signData, err := p.Sign(crypto_rand.Reader, payload, crypto.Hash(0))
	if err != nil {
		return Signature{}, err
	}

	var signature Signature
	copy(signature[:], signData)

	return signature, err
}

func (k PrivateKey) PublicKey() PublicKey {
	p := ed25519.PrivateKey(k)
	pub := p.Public().(ed25519.PublicKey)

	var publicKey PublicKey
	copy(publicKey[:], pub)

	return publicKey
}

type PublicKey [PublicKeyLength]byte

func PublicKeyFromBytes(in []byte) (out PublicKey) {
	byteCount := len(in)
	if byteCount == 0 {
		return
	}

	max := PublicKeyLength
	if byteCount < max {
		max = byteCount
	}

	copy(out[:], in[0:max])
	return
}

func MustPublicKeyFromBase58(in string) PublicKey {
	out, err := PublicKeyFromBase58(in)
	if err != nil {
		panic(err)
	}
	return out
}

func PublicKeyFromBase58(in string) (out PublicKey, err error) {
	val, err := base58.Decode(in)
	if err != nil {
		return out, fmt.Errorf("decode: %w", err)
	}

	if len(val) != PublicKeyLength {
		return out, fmt.Errorf("invalid length, expected %v, got %d", PublicKeyLength, len(val))
	}

	copy(out[:], val)
	return
}

func (p PublicKey) MarshalText() ([]byte, error) {
	return []byte(base58.Encode(p[:])), nil
}

func (p *PublicKey) UnmarshalText(data []byte) (err error) {
	*p, err = PublicKeyFromBase58(string(data))
	if err != nil {
		return fmt.Errorf("invalid public key %q: %w", data, err)
	}
	return
}

func (p PublicKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(base58.Encode(p[:]))
}

func (p *PublicKey) UnmarshalJSON(data []byte) (err error) {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	*p, err = PublicKeyFromBase58(s)
	if err != nil {
		return fmt.Errorf("invalid public key %q: %w", s, err)
	}
	return
}

func (p PublicKey) Equals(pb PublicKey) bool {
	return p == pb
}

// ToPointer returns a pointer to the pubkey.
func (p PublicKey) ToPointer() *PublicKey {
	return &p
}

func (p PublicKey) Bytes() []byte {
	return []byte(p[:])
}

var zeroPublicKey = PublicKey{}

// IsZero returns whether the public key is zero.
// NOTE: the System Program public key is also zero.
func (p PublicKey) IsZero() bool {
	return p == zeroPublicKey
}

func (p PublicKey) String() string {
	return base58.Encode(p[:])
}

type PublicKeySlice []PublicKey

// UniqueAppend appends the provided pubkey only if it is not
// already present in the slice.
// Returns true when the provided pubkey wasn't already present.
func (slice *PublicKeySlice) UniqueAppend(pubkey PublicKey) bool {
	if !slice.Has(pubkey) {
		slice.Append(pubkey)
		return true
	}
	return false
}

func (slice *PublicKeySlice) Append(pubkey PublicKey) {
	*slice = append(*slice, pubkey)
}

func (slice PublicKeySlice) Has(pubkey PublicKey) bool {
	for _, key := range slice {
		if key.Equals(pubkey) {
			return true
		}
	}
	return false
}

var nativeProgramIDs = PublicKeySlice{
	BPFLoaderProgramID,
	BPFLoaderDeprecatedProgramID,
	FeatureProgramID,
	ConfigProgramID,
	StakeProgramID,
	VoteProgramID,
	Secp256k1ProgramID,
	SystemProgramID,
	SysVarClockPubkey,
	SysVarEpochSchedulePubkey,
	SysVarFeesPubkey,
	SysVarInstructionsPubkey,
	SysVarRecentBlockHashesPubkey,
	SysVarRentPubkey,
	SysVarRewardsPubkey,
	SysVarSlotHashesPubkey,
	SysVarSlotHistoryPubkey,
	SysVarStakeHistoryPubkey,
}

// https://github.com/solana-labs/solana/blob/216983c50e0a618facc39aa07472ba6d23f1b33a/sdk/program/src/pubkey.rs#L372
func isNativeProgramID(key PublicKey) bool {
	return nativeProgramIDs.Has(key)
}

const (
	/// Number of bytes in a pubkey.
	PublicKeyLength = 32
	// Maximum length of derived pubkey seed.
	MaxSeedLength = 32
	// Maximum number of seeds.
	MaxSeeds = 16
	// // Maximum string length of a base58 encoded pubkey.
	// MaxBase58Length = 44
)

// Ported from https://github.com/solana-labs/solana/blob/216983c50e0a618facc39aa07472ba6d23f1b33a/sdk/program/src/pubkey.rs#L159
func CreateWithSeed(base PublicKey, seed string, owner PublicKey) (PublicKey, error) {
	if len(seed) > MaxSeedLength {
		return PublicKey{}, errors.New("Max seed length exceeded")
	}

	// let owner = owner.as_ref();
	// if owner.len() >= PDA_MARKER.len() {
	//     let slice = &owner[owner.len() - PDA_MARKER.len()..];
	//     if slice == PDA_MARKER {
	//         return Err(PubkeyError::IllegalOwner);
	//     }
	// }

	b := make([]byte, 0, 64+len(seed))
	b = append(b, base[:]...)
	b = append(b, seed[:]...)
	b = append(b, owner[:]...)
	hash := sha256.Sum256(b)
	return PublicKeyFromBytes(hash[:]), nil
}

const PDA_MARKER = "ProgramDerivedAddress"

// Create a program address.
// Ported from https://github.com/solana-labs/solana/blob/216983c50e0a618facc39aa07472ba6d23f1b33a/sdk/program/src/pubkey.rs#L204
func CreateProgramAddress(seeds [][]byte, programID PublicKey) (PublicKey, error) {
	if len(seeds) > MaxSeeds {
		return PublicKey{}, errors.New("Max seed length exceeded")
	}

	for _, seed := range seeds {
		if len(seed) > MaxSeedLength {
			return PublicKey{}, errors.New("Max seed length exceeded")
		}
	}

	if isNativeProgramID(programID) {
		return PublicKey{}, fmt.Errorf("illegal owner: %s is a native program", programID)
	}

	buf := []byte{}
	for _, seed := range seeds {
		buf = append(buf, seed...)
	}

	buf = append(buf, programID[:]...)
	buf = append(buf, []byte(PDA_MARKER)...)
	hash := sha256.Sum256(buf)

	_, err := new(edwards25519.Point).SetBytes(hash[:])
	isOnCurve := err == nil
	if isOnCurve {
		return PublicKey{}, errors.New("invalid seeds; address must fall off the curve")
	}

	return PublicKeyFromBytes(hash[:]), nil
}

// Find a valid program address and its corresponding bump seed.
func FindProgramAddress(seed [][]byte, programID PublicKey) (PublicKey, uint8, error) {
	var address PublicKey
	var err error
	bumpSeed := uint8(math.MaxUint8)
	for bumpSeed != 0 {
		address, err = CreateProgramAddress(append(seed, []byte{byte(bumpSeed)}), programID)
		if err == nil {
			return address, bumpSeed, nil
		}
		bumpSeed--
	}
	return PublicKey{}, bumpSeed, errors.New("unable to find a valid program address")
}

func FindAssociatedTokenAddress(
	wallet PublicKey,
	mint PublicKey,
) (PublicKey, uint8, error) {
	return findAssociatedTokenAddressAndBumpSeed(
		wallet,
		mint,
		SPLAssociatedTokenAccountProgramID,
	)
}

func findAssociatedTokenAddressAndBumpSeed(
	walletAddress PublicKey,
	splTokenMintAddress PublicKey,
	programID PublicKey,
) (PublicKey, uint8, error) {
	return FindProgramAddress([][]byte{
		walletAddress[:],
		TokenProgramID[:],
		splTokenMintAddress[:],
	},
		programID,
	)
}

// FindTokenMetadataAddress returns the token metadata program-derived address given a SPL token mint address.
func FindTokenMetadataAddress(mint PublicKey) (PublicKey, uint8, error) {
	seed := [][]byte{
		[]byte("metadata"),
		TokenMetadataProgramID[:],
		mint[:],
	}
	return FindProgramAddress(seed, TokenMetadataProgramID)
}
