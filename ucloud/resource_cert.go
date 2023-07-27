package ucloud

import (
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/ucloud/ucloud-sdk-go/services/ucdn"
	"github.com/ucloud/ucloud-sdk-go/ucloud/request"
	"golang.org/x/net/context"
)

var (
	_ resource.Resource              = &ucloudCertResource{}
	_ resource.ResourceWithConfigure = &ucloudCertResource{}
)

type ucloudCertResourceModel struct {
	CertName   types.String `tfsdk:"cert_name"`
	UserCert   types.String `tfsdk:"user_cert"`
	PrivateKey types.String `tfsdk:"private_key"`
	CaCert     types.String `tfsdk:"ca_cert"`
}

type ucloudCertResource struct {
	client *ucdn.UCDNClient
}

func NewUcloudCertResource() resource.Resource {
	return &ucloudCertResource{}
}

func (r *ucloudCertResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cert"
}

func (r *ucloudCertResource) Schema(_ context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "The resource provides the cdn certificate",
		Attributes: map[string]schema.Attribute{
			"cert_name": &schema.StringAttribute{
				Description: "The name of certificate",
				Required:    true,
			},
			"user_cert": &schema.StringAttribute{
				Description: "Certificate,e.g.cert.pem",
				Required:    true,
			},
			"private_key": &schema.StringAttribute{
				Description: "Private key,e.g.,privkey.pem",
				Required:    true,
				Sensitive:   true,
			},
			"ca_cert": &schema.StringAttribute{
				Description: "CA of the certificate,e.g.,chain.pem",
				Optional:    true,
			},
		},
	}
}

func (r *ucloudCertResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(ucloudClients).cdnClient
}

func (r *ucloudCertResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var model ucloudCertResourceModel

	diags := req.Plan.Get(ctx, &model)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addCertificateRequest := &ucdn.AddCertificateRequest{
		CommonBase: request.CommonBase{
			ProjectId: &r.client.GetConfig().ProjectId,
		},
		CertName:   model.CertName.ValueStringPointer(),
		UserCert:   model.UserCert.ValueStringPointer(),
		PrivateKey: model.PrivateKey.ValueStringPointer(),
		CaCert:     model.CaCert.ValueStringPointer(),
	}
	addCertificateResponse, err := r.client.AddCertificate(addCertificateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"[API ERROR] Failed to Add Cert",
			err.Error(),
		)
		return
	}
	if addCertificateResponse.RetCode != 0 {
		resp.Diagnostics.AddError(
			"[API ERROR] Failed to Add Cert",
			addCertificateResponse.Message,
		)
		return
	}
	resp.State.Set(ctx, &model)
}

func (r *ucloudCertResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

func (r *ucloudCertResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan ucloudCertResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if state.CertName != plan.CertName {
		resp.Diagnostics.AddError(
			"[API ERROR] Failed to Update Cert",
			"not allowed to modify cert_name")
		return
	}
	resp.Diagnostics.AddWarning("[API WARNING]", "not implemented to update cert")
	resp.State.Set(ctx, req.State.Raw)
}

func (r *ucloudCertResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var model ucloudCertResourceModel

	diags := req.State.Get(ctx, &model)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	deleteCertificateRequest := ucdn.DeleteCertificateRequest{
		CommonBase: request.CommonBase{
			ProjectId: &r.client.GetConfig().ProjectId,
		},
		CertName: model.CertName.ValueStringPointer(),
	}
	deleteCertificateResponse, err := r.client.DeleteCertificate(&deleteCertificateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"[API ERROR] Failed to Del Cert",
			err.Error(),
		)
		return
	}
	if deleteCertificateResponse.RetCode != 0 {
		resp.Diagnostics.AddError(
			"[API ERROR] Failed to Del Cert",
			deleteCertificateResponse.Message,
		)
		return
	}
}
