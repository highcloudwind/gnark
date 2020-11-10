// Copyright 2020 ConsenSys AG
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by gnark/internal/generators DO NOT EDIT

package groth16

import (
	"encoding/binary"

	curve "github.com/consensys/gurvy/bw761"

	"io"

	"github.com/fxamacker/cbor/v2"
)

// WriteTo writes binary encoding of the Proof elements to writer
// points are stored in compressed form Ar | Krs | Bs
// use WriteRawTo(...) to encode the proof without point compression
func (proof *Proof) WriteTo(w io.Writer) (n int64, err error) {
	return proof.writeTo(w, false)
}

// WriteRawTo writes binary encoding of the Proof elements to writer
// points are stored in uncompressed form Ar | Krs | Bs
// use WriteTo(...) to encode the proof with point compression
func (proof *Proof) WriteRawTo(w io.Writer) (n int64, err error) {
	return proof.writeTo(w, true)
}

func (proof *Proof) writeTo(w io.Writer, raw bool) (int64, error) {
	var enc *curve.Encoder
	if raw {
		enc = curve.NewEncoder(w, curve.RawEncoding())
	} else {
		enc = curve.NewEncoder(w)
	}

	if err := enc.Encode(&proof.Ar); err != nil {
		return enc.BytesWritten(), err
	}
	if err := enc.Encode(&proof.Krs); err != nil {
		return enc.BytesWritten(), err
	}
	if err := enc.Encode(&proof.Bs); err != nil {
		return enc.BytesWritten(), err
	}
	return enc.BytesWritten(), nil
}

// ReadFrom attempts to decode a Proof from reader
// Proof must be encoded through WriteTo (compressed) or WriteRawTo (uncompressed)
// note that we don't check that the points are on the curve or in the correct subgroup at this point
func (proof *Proof) ReadFrom(r io.Reader) (n int64, err error) {

	dec := curve.NewDecoder(r)

	if err := dec.Decode(&proof.Ar); err != nil {
		return dec.BytesRead(), err
	}
	if err := dec.Decode(&proof.Krs); err != nil {
		return dec.BytesRead(), err
	}
	if err := dec.Decode(&proof.Bs); err != nil {
		return dec.BytesRead(), err
	}

	return dec.BytesRead(), nil
}

// WriteTo writes binary encoding of the key elements to writer
// points are compressed
// use WriteRawTo(...) to encode the key without point compression
func (vk *VerifyingKey) WriteTo(w io.Writer) (n int64, err error) {
	return vk.writeTo(w, false)
}

// WriteRawTo writes binary encoding of the key elements to writer
// points are not compressed
// use WriteTo(...) to encode the key with point compression
func (vk *VerifyingKey) WriteRawTo(w io.Writer) (n int64, err error) {
	return vk.writeTo(w, true)
}

func (vk *VerifyingKey) writeTo(w io.Writer, raw bool) (n int64, err error) {
	var written int

	// encode public input names
	var pBytes []byte
	pBytes, err = cbor.Marshal(vk.PublicInputs)
	if err != nil {
		return
	}
	err = binary.Write(w, binary.BigEndian, uint64(len(pBytes)))
	if err != nil {
		return
	}
	n += 8
	written, err = w.Write(pBytes)
	n += int64(written)
	if err != nil {
		return
	}

	// write vk.E
	buf := vk.E.Bytes()
	written, err = w.Write(buf[:])
	n += int64(written)
	if err != nil {
		return
	}

	var enc *curve.Encoder
	if raw {
		enc = curve.NewEncoder(w, curve.RawEncoding())
	} else {
		enc = curve.NewEncoder(w)
	}

	err = enc.Encode(&vk.G2.GammaNeg)
	n += enc.BytesWritten()
	if err != nil {
		return
	}

	err = enc.Encode(&vk.G2.DeltaNeg)
	n += enc.BytesWritten()
	if err != nil {
		return
	}

	err = enc.Encode(vk.G1.K)
	n += enc.BytesWritten()
	return
}

// ReadFrom attempts to decode a VerifyingKey from reader
// VerifyingKey must be encoded through WriteTo (compressed) or WriteRawTo (uncompressed)
// note that we don't check that the points are on the curve or in the correct subgroup at this point
// TODO while Proof points correctness is checkd in the Verifier, here may be a good place to check key
func (vk *VerifyingKey) ReadFrom(r io.Reader) (n int64, err error) {

	var read int
	var buf [curve.SizeOfGT]byte

	read, err = io.ReadFull(r, buf[:8])
	n += int64(read)
	if err != nil {
		return
	}
	lPublicInputs := binary.BigEndian.Uint64(buf[:8])

	bPublicInputs := make([]byte, lPublicInputs)
	read, err = io.ReadFull(r, bPublicInputs)
	n += int64(read)
	if err != nil {
		return
	}
	err = cbor.Unmarshal(bPublicInputs, &vk.PublicInputs)
	if err != nil {
		return
	}

	// read vk.E

	read, err = r.Read(buf[:])
	n += int64(read)
	if err != nil {
		return
	}
	err = vk.E.SetBytes(buf[:])
	if err != nil {
		return
	}

	dec := curve.NewDecoder(r)

	err = dec.Decode(&vk.G2.GammaNeg)
	n += dec.BytesRead()
	if err != nil {
		return
	}

	err = dec.Decode(&vk.G2.DeltaNeg)
	n += dec.BytesRead()
	if err != nil {
		return
	}

	err = dec.Decode(&vk.G1.K)
	n += dec.BytesRead()

	return
}

// WriteTo writes binary encoding of the key elements to writer
// points are compressed
// use WriteRawTo(...) to encode the key without point compression
func (pk *ProvingKey) WriteTo(w io.Writer) (n int64, err error) {
	return pk.writeTo(w, false)
}

// WriteRawTo writes binary encoding of the key elements to writer
// points are not compressed
// use WriteTo(...) to encode the key with point compression
func (pk *ProvingKey) WriteRawTo(w io.Writer) (n int64, err error) {
	return pk.writeTo(w, true)
}

func (pk *ProvingKey) writeTo(w io.Writer, raw bool) (int64, error) {
	n, err := pk.Domain.WriteTo(w)
	if err != nil {
		return n, err
	}

	var enc *curve.Encoder
	if raw {
		enc = curve.NewEncoder(w, curve.RawEncoding())
	} else {
		enc = curve.NewEncoder(w)
	}

	toEncode := []interface{}{
		&pk.G1.Alpha,
		&pk.G1.Beta,
		&pk.G1.Delta,
		pk.G1.A,
		pk.G1.B,
		pk.G1.Z,
		pk.G1.K,
		&pk.G2.Beta,
		&pk.G2.Delta,
		pk.G2.B,
	}

	for _, v := range toEncode {
		if err := enc.Encode(v); err != nil {
			return n + enc.BytesWritten(), err
		}
	}

	return n + enc.BytesWritten(), nil

}

// ReadFrom attempts to decode a ProvingKey from reader
// ProvingKey must be encoded through WriteTo (compressed) or WriteRawTo (uncompressed)
// note that we don't check that the points are on the curve or in the correct subgroup at this point
// TODO while Proof points correctness is checkd in the Verifier, here may be a good place to check key
func (pk *ProvingKey) ReadFrom(r io.Reader) (int64, error) {

	n, err := pk.Domain.ReadFrom(r)
	if err != nil {
		return n, err
	}

	dec := curve.NewDecoder(r)

	toDecode := []interface{}{
		&pk.G1.Alpha,
		&pk.G1.Beta,
		&pk.G1.Delta,
		&pk.G1.A,
		&pk.G1.B,
		&pk.G1.Z,
		&pk.G1.K,
		&pk.G2.Beta,
		&pk.G2.Delta,
		&pk.G2.B,
	}

	for _, v := range toDecode {
		if err := dec.Decode(v); err != nil {
			return n + dec.BytesRead(), err
		}
	}

	return n + dec.BytesRead(), nil
}
