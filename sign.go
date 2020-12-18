// Copyright 2018 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package feed

import (
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/btcsuite/btcd/btcec"
)

const signatureLength = 65

// Hash represents the 32 byte Keccak256 hash of arbitrary data.
type Hash [HashLength]byte

// Signature is an alias for a static byte array with the size of a signature
type Signature [signatureLength]byte

// Signer signs feed update payloads
type Signer interface {
	Sign(Hash) (Signature, error)
	Address() Address
}

// GenericSigner implements the Signer interface
// It is the vanilla signer that probably should be used in most cases
type GenericSigner struct {
	PrivKey *ecdsa.PrivateKey
	address Address
}

// NewGenericSigner builds a signer that will sign everything with the provided private key
func NewGenericSigner(privKey *ecdsa.PrivateKey) *GenericSigner {
	addr, err := NewEthereumAddress(privKey.PublicKey)
	if err != nil {
		panic(err)
	}

	s := &GenericSigner{
		PrivKey: privKey,
	}
	copy(s.address[:], addr)
	return s
}

// addEthereumPrefix adds the ethereum prefix to the data.
func addEthereumPrefix(data []byte) []byte {
	return []byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data))
}

// hashWithEthereumPrefix returns the hash that should be signed for the given data.
func hashWithEthereumPrefix(data []byte) ([]byte, error) {
	return LegacyKeccak256(addEthereumPrefix(data))
}

// Sign signs the supplied data
// It wraps the ethereum crypto.Sign() method
func (s *GenericSigner) Sign(data Hash) (signature Signature, err error) {
	//hash, err := hashWithEthereumPrefix(data[:])
	//if err != nil {
	//return
	//}

	sig, err := s.sign(data[:], true)
	if err != nil {
		return
	}
	copy(signature[:], sig)
	return
}

// sign the provided hash and convert it to the ethereum (r,s,v) format.
func (s *GenericSigner) sign(sighash []byte, isCompressedKey bool) ([]byte, error) {
	signature, err := btcec.SignCompact(btcec.S256(), (*btcec.PrivateKey)(s.PrivKey), sighash, false)
	if err != nil {
		return nil, err
	}

	// Convert to Ethereum signature format with 'recovery id' v at the end.
	fmt.Println(hex.EncodeToString(signature))
	v := signature[0]
	copy(signature, signature[1:])
	signature[64] = v
	fmt.Println(hex.EncodeToString(signature))

	return signature, nil
}

// Address returns the public key of the signer's private key
func (s *GenericSigner) Address() Address {
	return s.address
}

// getUserAddr extracts the address of the feed update signer
func getUserAddr(digest Hash, signature Signature) (Address, error) {
	pub, err := Recover(signature[:], digest[:])
	if err != nil {
		return Address{}, err
	}
	a, _ := NewEthereumAddress(*pub)
	var aa Address
	copy(aa[:], a)
	return aa, nil
}

// Recover verifies signature with the data base provided.
// It is using `btcec.RecoverCompact` function.
func Recover(signature, data []byte) (*ecdsa.PublicKey, error) {
	if len(signature) != 65 {
		return nil, errors.New("invalid signature length")
	}
	// Convert to btcec input format with 'recovery id' v at the beginning.
	btcsig := make([]byte, 65)
	btcsig[0] = signature[64]
	copy(btcsig[1:], signature)

	//hash, err := hashWithEthereumPrefix(data)
	//if err != nil {
	//return nil, err
	//}

	p, _, err := btcec.RecoverCompact(btcec.S256(), btcsig, data)
	return (*ecdsa.PublicKey)(p), err
}
