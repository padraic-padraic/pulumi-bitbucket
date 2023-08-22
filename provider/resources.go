// Copyright 2016-2018, Pulumi Corporation.
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

package bitbucket

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ettle/strcase"
	"github.com/padraic-padraic/pulumi-bitbucket/provider/pkg/version"
	"github.com/padraic-padraic/terraform-provider-bitbucket/v2/bitbucket"
	"github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfbridge"
	shim "github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfshim"
	shimv2 "github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfshim/sdk-v2"
	"github.com/pulumi/pulumi/sdk/v3/go/common/resource"
	"github.com/pulumi/pulumi/sdk/v3/go/common/tokens"
	"github.com/pulumi/pulumi/sdk/v3/go/common/util/contract"
)

// all of the token components used below.
const (
	// This variable controls the default name of the package in the package
	mainMod = "index" // the bitbucket module
)

func convertName(name string) string {
	idx := strings.Index(name, "_")
	contract.Assertf(idx > 0 && idx < len(name)-1, "Invalid snake case name %s", name)
	name = name[idx+1:]
	contract.Assertf(len(name) > 0, "Invalid snake case name %s", name)
	return strcase.ToPascal(name)
}

func makeDataSource(mod string, name string) tokens.ModuleMember {
	name = convertName(name)
	return tfbridge.MakeDataSource("bitbucket", mod, "get"+name)
}

func makeResource(mod string, res string) tokens.Type {
	return tfbridge.MakeResource("bitbucket", mod, convertName(res))
}

// preConfigureCallback is called before the providerConfigure function of the underlying provider.
// It should validate that the provider can be configured, and provide actionable errors in the case
// it cannot be. Configuration variables can be read from `vars` using the `stringValue` function -
// for example `stringValue(vars, "accessKey")`.
func preConfigureCallback(vars resource.PropertyMap, c shim.ResourceConfig) error {
	return nil
}

// Provider returns additional overlaid schema and metadata associated with the provider..
func Provider() tfbridge.ProviderInfo {
	// Instantiate the Terraform provider
	p := shimv2.NewProvider(bitbucket.Provider())
	// Create a Pulumi provider mapping
	prov := tfbridge.ProviderInfo{
		P:    p,
		Name: "bitbucket",
		// DisplayName is a way to be able to change the casing of the provider
		// name when being displayed on the Pulumi registry
		DisplayName: "Bitbucket",
		// The default publisher for all packages is Pulumi.
		// Change this to your personal name (or a company name) that you
		// would like to be shown in the Pulumi Registry if this package is published
		// there.
		Publisher: "padraic-padraic",
		// LogoURL is optional but useful to help identify your package in the Pulumi Registry
		// if this package is published there.
		//
		// You may host a logo on a domain you control or add an SVG logo for your package
		// in your repository and use the raw content URL for that file as your logo URL.
		LogoURL: "https://raw.githubusercontent.com/padraic-padraic/pulumi-bitbucket/main/docs/bitbucket.png",
		// PluginDownloadURL is an optional URL used to download the Provider
		// for use in Pulumi programs
		// e.g https://github.com/org/pulumi-provider-name/releases/
		PluginDownloadURL: "github://api.github.com/padraic-padraic/pulumi-bitbucket",
		Description:       "A Pulumi package for creating and managing Bitbucket resources",
		// category/cloud tag helps with categorizing the package in the Pulumi Registry.
		// For all available categories, see `Keywords` in
		// https://www.pulumi.com/docs/guides/pulumi-packages/schema/#package.
		Keywords: []string{
			"pulumi",
			"bitbucket",
			"category/versioncontrol",
		},
		License:    "Apache-2.0",
		Homepage:   "https://github.com/padraic-padraic/pulumi-bitbucket",
		Repository: "https://github.com/padraic-padraic/pulumi-bitbucket",
		// The GitHub Org for the provider - defaults to `terraform-providers`. Note that this
		// should match the TF provider module's require directive, not any replace directives.
		Version:   version.Version,
		GitHubOrg: "padraic-padraic",
		Config: map[string]*tfbridge.SchemaInfo{
			// Add any required configuration here, or remove the example below if
			// no additional points are required.
			// "region": {
			// 	Type: tfbridge.MakeType("region", "Region"),
			// 	Default: &tfbridge.DefaultInfo{
			// 		EnvVars: []string{"AWS_REGION", "AWS_DEFAULT_REGION"},
			// 	},
			// },
			"username": {
				Default: &tfbridge.DefaultInfo{
					EnvVars: []string{"BITBUCKET_USERNAME"},
				},
				MarkAsOptional: tfbridge.True(),
			},
			"password": {
				Default: &tfbridge.DefaultInfo{
					EnvVars: []string{"BITBUCKET_PASSWORD"},
				},
				MarkAsOptional: tfbridge.True(),
			},
			"oauth_client_id": {
				Default: &tfbridge.DefaultInfo{
					EnvVars: []string{"BITBUCKET_OAUTH_CLIENT_ID"},
				},
				MarkAsOptional: tfbridge.True(),
			},
			"oauth_client_secret": {
				Default: &tfbridge.DefaultInfo{
					EnvVars: []string{"BITBUCKET_OAUTH_CLIENT_SECRET"},
				},
				MarkAsOptional: tfbridge.True(),
			},
			"oauth_token": {
				Default: &tfbridge.DefaultInfo{
					EnvVars: []string{"BITBUCKET_OAUTH_TOKEN"},
				},
				MarkAsOptional: tfbridge.True(),
			},
		},
		PreConfigureCallback: preConfigureCallback,
		Resources: map[string]*tfbridge.ResourceInfo{
			// Map each resource in the Terraform provider to a Pulumi type. Two examples
			// are below - the single line form is the common case. The multi-line form is
			// needed only if you wish to override types or other default options.
			//
			// "aws_iam_role": {Tok: makeResource(mainMod(mainMod, "aws_iam_role")}
			//
			// "aws_acm_certificate": {
			// 	Tok: Tok: makeResource(mainMod(mainMod, "aws_acm_certificate"),
			// 	Fields: map[string]*tfbridge.SchemaInfo{
			// 		"tags": {Type: tfbridge.MakeType("bitbucket", "Tags")},
			// 	},
			// },
			"bitbucket_branch_restriction": {
				Tok: makeResource(mainMod, "bitbucket_branch_restriction"),
			},
			"bitbucket_branching_model": {
				Tok: makeResource(mainMod, "bitbucket_branching_model"),
			},
			"bitbucket_commit_file": {
				Tok: makeResource(mainMod, "bitbucket_commit_file"),
			},
			"bitbucket_default_reviewers": {
				Tok: makeResource(mainMod, "bitbucket_default_reviewers"),
			},
			"bitbucket_deploy_key": {
				Tok: makeResource(mainMod, "bitbucket_deploy_key"),
			},
			"bitbucket_deployment": {
				Tok: makeResource(mainMod, "bitbucket_deployment"),
			},
			"bitbucket_deployment_variable": {
				Tok: makeResource(mainMod, "bitbucket_deployment_variable"),
			},
			"bitbucket_forked_repository": {
				Tok: makeResource(mainMod, "bitbucket_forked_repository"),
			},
			"bitbucket_group": {
				Tok: makeResource(mainMod, "bitbucket_group"),
			},
			"bitbucket_group_membership": {
				Tok: makeResource(mainMod, "bitbucket_group_membership"),
			},
			"bitbucket_hook": {
				Tok: makeResource(mainMod, "bitbucket_hook"),
			},
			"bitbucket_pipeline_schedule": {
				Tok: makeResource(mainMod, "bitbucket_pipeline_schedule"),
			},
			"bitbucket_pipeline_ssh_key": {
				Tok: makeResource(mainMod, "bitbucket_pipeline_ssh_key"),
			},
			"bitbucket_pipeline_ssh_known_host": {
				Tok: makeResource(mainMod, "bitbucket_pipeline_ssh_known_host"),
			},
			"bitbucket_project": {
				Tok: makeResource(mainMod, "bitbucket_project"),
			},
			"bitbucket_project_branching_model": {
				Tok: makeResource(mainMod, "bitbucket_project_branching_model"),
			},
			"bitbucket_project_default_reviewers": {
				Tok: makeResource(mainMod, "bitbucket_project_default_reviewers"),
			},
			"bitbucket_repository": {
				Tok: makeResource(mainMod, "bitbucket_repository"),
			},
			"bitbucket_repository_group_permission": {
				Tok: makeResource(mainMod, "bitbucket_repository_group_permission"),
			},
			"bitbucket_repository_user_permission": {
				Tok: makeResource(mainMod, "bitbucket_repository_user_permission"),
			},
			"bitbucket_repository_variable": {
				Tok: makeResource(mainMod, "bitbucket_repository_variable"),
			},
			"bitbucket_ssh_key": {
				Tok: makeResource(mainMod, "bitbucket_ssh_key"),
			},
			"bitbucket_workspace_hook": {
				Tok: makeResource(mainMod, "bitbucket_workspace_hook"),
			},
			"bitbucket_workspace_variable": {
				Tok: makeResource(mainMod, "bitbucket_workspace_variable"),
			},
		},
		DataSources: map[string]*tfbridge.DataSourceInfo{
			// Map each resource in the Terraform provider to a Pulumi function. An example
			// is below.
			// "aws_ami": {Tok: makeDataSource(mainMod, "aws_ami")},
			"bitbucket_current_user": {
				Tok: makeDataSource(mainMod, "bitbucket_current_user"),
			},
			"bitbucket_deployment": {
				Tok: makeDataSource(mainMod, "bitbucket_deployment"),
			},
			"bitbucket_group": {
				Tok: makeDataSource(mainMod, "bitbucket_group"),
			},
			"bitbucket_group_members": {
				Tok: makeDataSource(mainMod, "bitbucket_group_members"),
			},
			"bitbucket_groups": {
				Tok: makeDataSource(mainMod, "bitbucket_groups"),
			},
			"bitbucket_hook_types": {
				Tok: makeDataSource(mainMod, "bitbucket_hook_types"),
			},
			"bitbucket_ip_ranges": {
				Tok: makeDataSource(mainMod, "bitbucket_ip_ranges"),
			},
			"bitbucket_pipeline_oidc_config": {
				Tok: makeDataSource(mainMod, "bitbucket_pipeline_oidc_config"),
			},
			"bitbucket_pipeline_oidc_config_keys": {
				Tok: makeDataSource(mainMod, "bitbucket_pipeline_oidc_config_keys"),
			},
			"bitbucket_user": {
				Tok: makeDataSource(mainMod, "bitbucket_user"),
			},
			"bitbucket_workspace": {
				Tok: makeDataSource(mainMod, "bitbucket_workspace"),
			},
			"bitbucket_workspace_members": {
				Tok: makeDataSource(mainMod, "bitbucket_workspace_members"),
			},
		},
		JavaScript: &tfbridge.JavaScriptInfo{
			PackageName: "@pulumiverse/bitbucket",

			// List any npm dependencies and their versions
			Dependencies: map[string]string{
				"@pulumi/pulumi": "^3.0.0",
			},
			DevDependencies: map[string]string{
				"@types/node": "^10.0.0", // so we can access strongly typed node definitions.
				"@types/mime": "^2.0.0",
			},
			// See the documentation for tfbridge.OverlayInfo for how to lay out this
			// section, or refer to the AWS provider. Delete this section if there are
			// no overlay files.
			//Overlay: &tfbridge.OverlayInfo{},
		},
		Python: &tfbridge.PythonInfo{
			PackageName: "pulumiverse_bitbucket",

			// List any Python dependencies and their version ranges
			Requires: map[string]string{
				"pulumi": ">=3.0.0,<4.0.0",
			},
		},
		Golang: &tfbridge.GolangInfo{
			ImportBasePath: filepath.Join(
				fmt.Sprintf("github.com/padraic-padraic/pulumi-%[1]s/sdk/", "bitbucket"),
				tfbridge.GetModuleMajorVersion(version.Version),
				"go",
				"bitbucket",
			),
			GenerateResourceContainerTypes: true,
		},
		CSharp: &tfbridge.CSharpInfo{
			RootNamespace: "Pulumiverse",

			PackageReferences: map[string]string{
				"Pulumi": "3.*",
			},
		},
		Java: &tfbridge.JavaInfo{
			BasePackage: "com.padraic-padraic",
		},
	}

	prov.SetAutonaming(255, "-")

	return prov
}
