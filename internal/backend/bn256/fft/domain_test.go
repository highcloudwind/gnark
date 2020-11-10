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

package fft

import (
	"bytes"
	"reflect"
	"testing"
)

func TestDomainSerialization(t *testing.T) {
	domain := NewDomain(1 << 6)
	var reconstructed Domain

	var buf bytes.Buffer
	written, err := domain.WriteTo(&buf)
	if err != nil {
		t.Fatal(err)
	}
	var read int64
	read, err = reconstructed.ReadFrom(&buf)
	if err != nil {
		t.Fatal(err)
	}

	if written != read {
		t.Fatal("didn't read as many bytes as we wrote")
	}
	if !reflect.DeepEqual(domain, &reconstructed) {
		t.Fatal("Domain.SetBytes(Bytes()) failed")
	}
}
