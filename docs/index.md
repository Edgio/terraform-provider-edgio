---
layout: "Edgio"
page_title: "Provider: Edgio"
description: |-
  Learn about the Edgio Terraform provider.
---
# Edgecast Terraform Provider

Edgio allows you to securely deliver content to your customers. Use the Edgio Terraform provider to manage and provision your Property, Envirovment and CDN configuration.

## Prerequisites
Use of the Edgio Terraform provider requires:
* An Edgio customer account in good standing.
* An Edgio API client with the following scopes enabled:
    * `app.accounts`
    * `app.config`

 [Learn more.](guides/authentication.md)

## Basic Setup
We will now cover how to set up the Edgio Terraform provider using the [Terraform CLI](https://learn.hashicorp.com/tutorials/terraform/install-cli) version 0.14 or later. Perform the following steps:
1. Create and define a **.tf** configuration file for the Edgio Terraform provider. 
    1. Click `Use Provider` to display a template for the Edgio Terraform provider.
    1. Copy and paste the above template into your **.tf** configuration file.
    1. Define each desired variable within a variable block. 

        !> Do not define sensitive information (e.g., API credentials) within a variable block. One method for securely defining variables is to set variable values from within a **.tfvars** file (e.g., terraform.tfvars) that is excluded from source control. 

        For example, you can define the `client_secret` variable within a variable block.

            variable "client_secret" {
                description = "Client secret for the Edgio API client." 
                type = string
                sensitive = true
            }
    1. From within the `edgio` provider block, define `client_id` and `client_secret`. 

        For example, you can set the **client_secret** argument to the variable defined in the previous step:
            
            provider "edgio" {
                client_secret = var.client_secret
                client_id = var.client_id
            }
    1. Save your changes.
1. Create a **terraform.tfvars** file and set your variables within it.

    !> Exclude this file from source control, since it contains sensitive information. 

    1. Set each desired variable. 

        For example, you can set the `client_secret` variable for .

            client_secret = "AjWhnhjMUjasNywG8rrfANJ3tAxt1POk"
    1. Save your changes.
    1. Configure your repository to prevent this file from being included in source control by adding it to your **.gitignore** file.
1. Create and define a **.tf** configuration file for the desired resource(s).
    1. Create a **.tf** configuration file. You may add all of the desired resources to this configuration file or create resource-specific configuration files.
    1. From the documentation for the desired resource, copy and paste the desired example resource into your **.tf** configuration file.
    1. Modify the resource block as needed.
    1. Save your changes.
1. Initialize Terraform and install the latest version of the Edgio provider by running the following command from Terminal:

        $ terraform init
1. Create a plan for your resources by running the following command:

        $ terraform plan
1. Apply your plan by running the following command:

        $ terraform apply
    Review the plan and then type **yes** to apply it.
    Terraform will use the Edgio provider to update your configuration to match the configuration defined within your plan.

## Resources
Learn how to get started with Terraform:
* [Use the Command Line Interface](https://learn.hashicorp.com/collections/terraform/cli)
* [Reuse Configuration with Modules](https://learn.hashicorp.com/collections/terraform/modules)
* [Write Terraform Configuration](https://learn.hashicorp.com/collections/terraform/configuration-language)
* [Glossary](https://www.terraform.io/docs/glossary)

Learn about basic Edgio concepts:
* [Getting Started (Create Edgio Account)](https://docs.edg.io/applications/v7/getting_started#create-account)
* [Getting Started (Create Edgio Property)](https://docs.edg.io/applications/v7/getting_started#create-property)
* [Security](https://docs.edg.io/applications/v7/security)
* [Sites](https://docs.edg.io/applications/v7/sites_frameworks)
* [Edge Functions](https://docs.edg.io/applications/v7/edge_functions)