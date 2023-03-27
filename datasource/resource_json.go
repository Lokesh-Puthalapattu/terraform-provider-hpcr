// Copyright 2022 IBM Corp.
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
// limitations under the License.package datasource

package datasource

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ibm-hyper-protect/terraform-provider-hpcr/common"
	"github.com/ibm-hyper-protect/terraform-provider-hpcr/fp"
	E "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/either"
	F "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/function"
	J "github.com/ibm-hyper-protect/terraform-provider-hpcr/fp/json"
)

var (
	jsonBytes = F.Flow3(
		getJsonE,
		common.MapRefAnyE,
		E.Chain(J.Stringify[any]),
	)
)

func ResourceJSON() *schema.Resource {
	return &schema.Resource{
		Create: jsonUnencrypted.F1,
		Read:   jsonUnencrypted.F2,
		Delete: jsonUnencrypted.F3,
		Schema: map[string]*schema.Schema{
			common.KeyJSON:     &schemaJsonIn,
			common.KeyRendered: &schemaRenderedOut,
			common.KeySha256:   &schemaSha256Out,
		},
		Description: "Generates an base64 encoded token from the JSON serialization of the input",
	}
}

func ResourceJSONEncrypted() *schema.Resource {
	return &schema.Resource{
		Create: jsonEncrypted.F1,
		Read:   jsonEncrypted.F2,
		Delete: jsonEncrypted.F3,
		Schema: map[string]*schema.Schema{
			// input parameters
			common.KeyJSON: &schemaJsonIn,
			common.KeyCert: &schemaCertIn,
			// output parameters
			common.KeyRendered: &schemaRenderedOut,
			common.KeySha256:   &schemaSha256Out,
			common.KeyChecksum: &schemaChecksumOut,
		},
		Description: "Generates an base64 encoded token from the JSON serialization of the input",
	}
}

func resourceEncJson(ctx *Context) func(d fp.ResourceData) ResourceDataE {

	// get the update method depending on the context
	update := updateEncryptedResource(ctx)
	hashWithCert := createHashWithCert(ctx)

	return func(d fp.ResourceData) ResourceDataE {
		// marshal input text
		jsonE := jsonBytes(d)

		return F.Pipe2(
			jsonE,
			E.Chain(hashWithCert(d)),
			E.Chain(F.Flow3(
				checksumMatchO(d),
				update(d)(jsonE),
				getResourceData(d),
			),
			),
		)
	}
}

func resourceJson(ctx *Context) func(d fp.ResourceData) ResourceDataE {

	return func(d fp.ResourceData) ResourceDataE {
		// marshal input text
		jsonE := jsonBytes(d)

		return F.Pipe2(
			jsonE,
			createHashE,
			E.Chain(F.Flow3(
				checksumMatchO(d),
				updateBase64Resource(d)(jsonE),
				getResourceData(d),
			),
			),
		)
	}
}

var (
	jsonUnencrypted = resourceLifeCycle(resourceJson)
	jsonEncrypted   = resourceLifeCycle(resourceEncJson)
)
