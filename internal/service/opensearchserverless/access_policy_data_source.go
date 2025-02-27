package opensearchserverless

import (
	"context"

	awstypes "github.com/aws/aws-sdk-go-v2/service/opensearchserverless/types"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/enum"
	"github.com/hashicorp/terraform-provider-aws/internal/framework"
	"github.com/hashicorp/terraform-provider-aws/internal/framework/flex"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @FrameworkDataSource(name="Access Policy")
func newDataSourceAccessPolicy(context.Context) (datasource.DataSourceWithConfigure, error) {
	return &dataSourceAccessPolicy{}, nil
}

const (
	DSNameAccessPolicy = "Access Policy Data Source"
)

type dataSourceAccessPolicy struct {
	framework.DataSourceWithConfigure
}

func (d *dataSourceAccessPolicy) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) { // nosemgrep:ci.meta-in-func-name
	resp.TypeName = "aws_opensearchserverless_access_policy"
}

func (d *dataSourceAccessPolicy) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"description": schema.StringAttribute{
				Computed: true,
			},
			"id": framework.IDAttribute(),
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(3, 32),
				},
			},
			"policy": schema.StringAttribute{
				Computed: true,
			},
			"policy_version": schema.StringAttribute{
				Computed: true,
			},
			"type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					enum.FrameworkValidate[awstypes.AccessPolicyType](),
				},
			},
		},
	}
}
func (d *dataSourceAccessPolicy) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	conn := d.Meta().OpenSearchServerlessClient(ctx)

	var data dataSourceAccessPolicyData
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	out, err := FindAccessPolicyByNameAndType(ctx, conn, data.Name.ValueString(), data.Type.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			create.ProblemStandardMessage(names.OpenSearchServerless, create.ErrActionReading, DSNameAccessPolicy, data.Name.String(), err),
			err.Error(),
		)
		return
	}

	data.ID = flex.StringToFramework(ctx, out.Name)
	data.Description = flex.StringToFramework(ctx, out.Description)
	data.Name = flex.StringToFramework(ctx, out.Name)
	data.Type = flex.StringValueToFramework(ctx, out.Type)
	data.PolicyVersion = flex.StringToFramework(ctx, out.PolicyVersion)

	policyBytes, err := out.Policy.MarshalSmithyDocument()

	if err != nil {
		resp.Diagnostics.AddError(
			create.ProblemStandardMessage(names.OpenSearchServerless, create.ErrActionReading, DSNameAccessPolicy, data.Name.String(), err),
			err.Error(),
		)
	}

	pb := string(policyBytes)
	data.Policy = flex.StringToFramework(ctx, &pb)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

type dataSourceAccessPolicyData struct {
	Description   types.String `tfsdk:"description"`
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Policy        types.String `tfsdk:"policy"`
	PolicyVersion types.String `tfsdk:"policy_version"`
	Type          types.String `tfsdk:"type"`
}
