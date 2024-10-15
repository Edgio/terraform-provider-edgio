// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"flag"
	"log"
	"os"

	"terraform-provider-edgio/internal/edgio_provider"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	var debugVar bool

	// Output current  directory
	s, _ := os.Getwd()
	log.Println("Current directory: ", s)

	flag.BoolVar(&debugVar, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		// NOTE: This is not a typical Terraform Registry provider address,
		// such as registry.terraform.io/hashicorp/hashicups. This specific
		// provider address is used in these tutorials in conjunction with a
		// specific Terraform CLI configuration for manual development testing
		// of this provider.
		Address: "hashicorp.com/edu/edgio",
		Debug:   debugVar,
	}

	providerFunc := func() provider.Provider {
		return edgio_provider.New()
	}

	err := providerserver.Serve(context.Background(), providerFunc, opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
