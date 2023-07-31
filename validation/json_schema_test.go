// Copyright 2022 IBM Corp.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package validation

import (
	"testing"

	_ "embed"

	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	I "github.com/IBM/fp-go/identity"
	"github.com/stretchr/testify/assert"
)

//go:embed samples/simple.yml
var TrivialContract string

func TestDiagContract(t *testing.T) {
	// construct the diagnostics
	diags := DiagContract(TrivialContract, nil)
	assert.Empty(t, diags)
}

func TestJsonSchema(t *testing.T) {
	schemaE := GetContractSchema()
	assert.True(t, E.IsRight(schemaE))
}

func TestValidate(t *testing.T) {
	// validator function
	validatorE := F.Pipe1(
		GetContractSchema(),
		E.Map[error](F.Flow2(
			validate[RawMap],
			ValidateYAML[RawMap],
		)),
	)
	// validate the data
	resE := F.Pipe1(
		validatorE,
		E.Chain(I.Ap[E.Either[error, RawMap]](TrivialContract)),
	)
	// validate that we have some data
	assert.True(t, E.IsRight(resE))
}
