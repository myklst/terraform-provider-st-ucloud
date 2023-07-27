package ucloud

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/ucloud/ucloud-sdk-go/services/ucdn"
	"github.com/ucloud/ucloud-sdk-go/ucloud"
	"github.com/ucloud/ucloud-sdk-go/ucloud/auth"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &ucloudProvider{}
)

type ucloudProviderModel struct {
	APIUrl     types.String `tfsdk:"api_url"`
	PrivateKey types.String `tfsdk:"private_key"`
	PublicKey  types.String `tfsdk:"public_key"`
	ProjectId  types.String `tfsdk:"project_id"`
	Region     types.String `tfsdk:"region"`
	Zone       types.String `tfsdk:"zone"`
}

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
			"api_url": schema.StringAttribute{
				Description: "",
				Required:    true,
			},
			"public_key": schema.StringAttribute{
				Description: "Public key for Ucloud API. May also be provided via UCLOUD_PUBLIC_KEY environment variable",
				Required:    true,
			},
			"private_key": schema.StringAttribute{
				Description: "Secret key for Ucloud API. May also be provided via UCLOUD_SECRET_KEY environment variable",
				Required:    true,
				Sensitive:   true,
			},
			"project_id": schema.StringAttribute{
				Description: "Project id should not be empty if public_key/private_key belongs to sub-account",
				Optional:    true,
			},
			"region": schema.StringAttribute{
				Description: "Ucloud region",
				Required:    true,
			},
			"zone": schema.StringAttribute{
				Description: "Ucloud zone",
				Required:    true,
			},
		},
	}
}

// Configure prepares a Ucloud API client for data sources and resources.
func (p *ucloudProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var model ucloudProviderModel

	cred := auth.NewCredential()
	diags := req.Config.Get(ctx, &model)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	cred.PublicKey = model.PublicKey.ValueString()
	cred.PrivateKey = model.PrivateKey.ValueString()

	cfg := ucloud.NewConfig()
	cfg.BaseUrl = model.APIUrl.ValueString()
	cfg.ProjectId = model.ProjectId.ValueString()
	cfg.Zone = model.Zone.ValueString()
	cfg.Region = model.Zone.ValueString()

	ucdnClient := ucdn.NewClient(&cfg, &cred)

	// Ucloud clients wrapper
	ucloudClients := ucloudClients{
		cdnClient: ucdnClient,
	}

	resp.DataSourceData = ucloudClients
	resp.ResourceData = ucloudClients
}

func (p *ucloudProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewUcloudCertDataSource,
	}
}

func (p *ucloudProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewUcloudCertResource,
		NewUcloudCdnDomainResource,
	}
}
