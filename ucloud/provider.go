package ucloud

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/ucloud/ucloud-sdk-go/services/ucdn"
	"github.com/ucloud/ucloud-sdk-go/ucloud"
	"github.com/ucloud/ucloud-sdk-go/ucloud/auth"
)

const ApiEndpoint = "https://api.ucloud.cn"

type ucloudProviderModel struct {
	Region     types.String `tfsdk:"region"`
	Zone       types.String `tfsdk:"zone"`
	ProjectId  types.String `tfsdk:"project_id"`
	PublicKey  types.String `tfsdk:"public_key"`
	PrivateKey types.String `tfsdk:"private_key"`
}

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &ucloudProvider{}
)

// Wrapper of Ucloud client
type ucloudClients struct {
	cdnClient *ucdn.UCDNClient
}

type ucloudProvider struct{}

func New() provider.Provider {
	return &ucloudProvider{}
}

// Metadata returns the provider type name.
func (p *ucloudProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "st-ucloud"
}

// Schema defines the provider-level schema for configuration data.
func (p *ucloudProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The Ucloud provider is used to interact with the many resources supported by Ucloud. " +
			"The provider needs to be configured with the proper credentials before it can be used.",
		Attributes: map[string]schema.Attribute{
			"region": schema.StringAttribute{
				Description: "Ucloud region",
				Optional:    true,
			},
			"zone": schema.StringAttribute{
				Description: "Ucloud zone",
				Optional:    true,
			},
			"project_id": schema.StringAttribute{
				Description: "Project id should not be empty if public_key/private_key belongs to sub-account",
				Optional:    true,
			},
			"public_key": schema.StringAttribute{
				Description: "Public key for Ucloud API. May also be provided via UCLOUD_PUBLIC_KEY environment variable",
				Optional:    true,
			},
			"private_key": schema.StringAttribute{
				Description: "Secret key for Ucloud API. May also be provided via UCLOUD_SECRET_KEY environment variable",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

// Configure prepares a UCloud API client for data sources and resources.
func (p *ucloudProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var model ucloudProviderModel
	diags := req.Config.Get(ctx, &model)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.
	if model.Region.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("region"),
			"Unknown Region",
			"The provider cannot create the UCloud API client as there is an unknown configuration value for the"+
				"region. Set the value statically in the configuration, or use the UCLOUD_REGION environment variable.",
		)
	}

	if model.Zone.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("zone"),
			"Unknown Zone",
			"The provider cannot create the UCloud API client as there is an unknown configuration value for the"+
				"zone. Set the value statically in the configuration, or use the UCLOUD_ZONE environment variable.",
		)
	}

	if model.ProjectId.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("project_id"),
			"Unknown ProjectId",
			"The provider cannot create the UCloud API client as there is an unknown configuration value for the"+
				"project_id. Set the value statically in the configuration, or use the UCLOUD_PROJECT_ID environment variable.",
		)
	}

	if model.PublicKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("public_key"),
			"Unknown PublicKey",
			"The provider cannot create the UCloud API client as there is an unknown configuration value for the"+
				"public_key. Set the value statically in the configuration, or use the UCLOUD_PUBLIC_KEY environment variable.",
		)
	}

	if model.PrivateKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("private_key"),
			"Unknown PrivateKey",
			"The provider cannot create the UCloud API client as there is an unknown configuration value for the"+
				"private_key. Set the value statically in the configuration, or use the UCLOUD_PRIVATE_KEY environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	var (
		region,
		zone,
		projectId,
		publicKey,
		privateKey string
	)

	// Default values to environment variables, but override
	// with Terraform configuration value if set.
	if !model.Region.IsNull() {
		region = model.Region.ValueString()
	} else {
		region = os.Getenv("UCLOUD_REGION")
	}

	if !model.Zone.IsNull() {
		zone = model.Zone.ValueString()
	} else {
		zone = os.Getenv("UCLOUD_ZONE")
	}

	if !model.ProjectId.IsNull() {
		projectId = model.ProjectId.ValueString()
	} else {
		projectId = os.Getenv("UCLOUD_PROJECT_ID")
	}

	if !model.PublicKey.IsNull() {
		publicKey = model.PublicKey.ValueString()
	} else {
		publicKey = os.Getenv("UCLOUD_PUBLIC_KEY")
	}

	if !model.PrivateKey.IsNull() {
		privateKey = model.PrivateKey.ValueString()
	} else {
		privateKey = os.Getenv("UCLOUD_PRIVATE_KEY")
	}

	// If any of the expected configuration are missing, return
	// errors with provider-specific guidance.
	if region == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("region"),
			"Missing Region",
			"The provider cannot create the UCloud API client as there is a "+
				"missing or empty value for the Region. Set the "+
				"value in the configuration or use the UCLOUD_REGION"+
				"environment variable. If either is already set, ensure the value "+
				"is not empty.",
		)
	}

	if zone == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("zone"),
			"Missing Zone",
			"The provider cannot create the UCloud API client as there is a "+
				"missing or empty value for the Zone. Set the "+
				"value in the configuration or use the UCLOUD_ZONE"+
				"environment variable. If either is already set, ensure the value "+
				"is not empty.",
		)
	}

	if projectId == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("project_id"),
			"Missing ProjectId",
			"The provider cannot create the UCloud API client as there is a "+
				"missing or empty value for the project_id. Set the "+
				"value in the configuration or use the UCLOUD_PROJECT_ID"+
				"environment variable. If either is already set, ensure the value "+
				"is not empty.",
		)
	}

	if publicKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("public_key"),
			"Missing PublicKey",
			"The provider cannot create the UCloud API client as there is a "+
				"missing or empty value for the public_key. Set the "+
				"value in the configuration or use the UCLOUD_PUBLIC_KEY"+
				"environment variable. If either is already set, ensure the value "+
				"is not empty.",
		)
	}

	if privateKey == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("private_key"),
			"Missing PrivateKey",
			"The provider cannot create the UCloud API client as there is a "+
				"missing or empty value for the private_key. Set the "+
				"value in the configuration or use the UCLOUD_PRIVATE_KEY"+
				"environment variable. If either is already set, ensure the value "+
				"is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	cfg := ucloud.Config{
		BaseUrl:   ApiEndpoint,
		Region:    region,
		Zone:      zone,
		ProjectId: projectId,
	}
	keys := auth.Credential{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}
	client := ucdn.NewClient(&cfg, &keys)

	// UCloud clients wrapper
	ucloudClients := ucloudClients{
		cdnClient: client,
	}

	resp.DataSourceData = ucloudClients
	resp.ResourceData = ucloudClients
}

func (p *ucloudProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewCertDataSource,
	}
}

func (p *ucloudProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewSslCertificateResource,
		NewCdnDomainResource,
		NewCdnDomainSslResource,
	}
}
