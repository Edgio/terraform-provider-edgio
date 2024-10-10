package resources

import (
	"context"
	"terraform-provider-edgio/internal/edgio_api"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-edgio/internal/edgio_provider/models"
	"terraform-provider-edgio/internal/edgio_provider/utility"
)

type CDNConfigurationResource struct {
	client edgio_api.EdgioClientInterface
}

func NewCDNConfigurationResource(client edgio_api.EdgioClientInterface) resource.Resource {
	return &CDNConfigurationResource{
		client: client,
	}
}

func (r *CDNConfigurationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "edgio_cdn_configuration"
}

func (r *CDNConfigurationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"environment_id": schema.StringAttribute{
				Required: true,
			},
			"configuration_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"rules": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					JSONEqualityModifier{},
				},
			},
			"origins": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required: true,
						},
						"type": schema.StringAttribute{
							Optional: true,
							Computed: true,
						},
						"hosts": schema.ListNestedAttribute{
							Optional: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"weight": schema.Int64Attribute{
										Optional: true,
										Computed: true,
									},
									"dns_max_ttl": schema.Int64Attribute{
										Optional: true,
										Computed: true,
									},
									"dns_preference": schema.StringAttribute{
										Optional: true,
										Computed: true,
									},
									"max_hard_pool": schema.Int64Attribute{
										Optional: true,
										Computed: true,
									},
									"dns_min_ttl": schema.Int64Attribute{
										Optional: true,
										Computed: true,
									},
									"location": schema.ListNestedAttribute{
										Optional: true,
										Computed: true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"port": schema.Int64Attribute{
													Optional: true,
													Computed: true,
												},
												"hostname": schema.StringAttribute{
													Optional: true,
													Computed: true,
												},
											},
										},
									},
									"max_pool": schema.Int64Attribute{
										Optional: true,
										Computed: true,
									},
									"balancer": schema.StringAttribute{
										Optional: true,
										Computed: true,
									},
									"scheme": schema.StringAttribute{
										Optional: true,
										Computed: true,
									},
									"override_host_header": schema.StringAttribute{
										Optional: true,
										Computed: true,
									},
									"sni_hint_and_strict_san_check": schema.StringAttribute{
										Optional: true,
										Computed: true,
									},
									"use_sni": schema.BoolAttribute{
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"balancer": schema.StringAttribute{
							Optional: true,
							Computed: true,
						},
						"override_host_header": schema.StringAttribute{
							Optional: true,
							Computed: true,
						},
						"shields": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"global": schema.StringAttribute{
									Optional: true,
									Computed: true,
								},
								"apac": schema.StringAttribute{
									Optional: true,
									Computed: true,
								},
								"emea": schema.StringAttribute{
									Optional: true,
									Computed: true,
								},
								"us_west": schema.StringAttribute{
									Optional: true,
									Computed: true,
								},
								"us_east": schema.StringAttribute{
									Optional: true,
									Computed: true,
								},
							},
						},
						"pci_certified_shields": schema.BoolAttribute{
							Optional: true,
							Computed: true,
						},
						"tls_verify": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"use_sni": schema.BoolAttribute{
									Optional: true,
									Computed: true,
								},
								"sni_hint_and_strict_san_check": schema.StringAttribute{
									Optional: true,
									Computed: true,
								},
								"allow_self_signed_certs": schema.BoolAttribute{
									Optional: true,
									Computed: true,
								},
								"pinned_certs": schema.ListAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Computed:    true,
								},
							},
						},
						"retry": schema.SingleNestedAttribute{
							Optional: true,
							Attributes: map[string]schema.Attribute{
								"status_codes": schema.ListAttribute{
									ElementType: types.Int64Type,
									Optional:    true,
									Computed:    true,
								},
								"ignore_retry_after_header": schema.BoolAttribute{
									Optional: true,
									Computed: true,
								},
								"after_seconds": schema.Int64Attribute{
									Optional: true,
									Computed: true,
								},
								"max_requests": schema.Int64Attribute{
									Optional: true,
									Computed: true,
								},
								"max_wait_seconds": schema.Int64Attribute{
									Optional: true,
									Computed: true,
								},
							},
						},
					},
				},
			},
			"hostnames": schema.ListNestedAttribute{
				Required: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"hostname": schema.StringAttribute{
							Required: true,
						},
						"default_origin_name": schema.StringAttribute{
							Optional: true,
							Computed: true,
						},
						"report_code": schema.Int64Attribute{
							Optional: true,
							Computed: true,
						},
						"tls": schema.SingleNestedAttribute{
							Optional: true,
							Computed: true,
							Attributes: map[string]schema.Attribute{
								"npn": schema.BoolAttribute{
									Optional: true,
									Computed: true,
								},
								"alpn": schema.BoolAttribute{
									Optional: true,
									Computed: true,
								},
								"protocols": schema.StringAttribute{
									Optional: true,
									Computed: true,
								},
								"use_sigalgs": schema.BoolAttribute{
									Optional: true,
									Computed: true,
								},
								"sni": schema.BoolAttribute{
									Optional: true,
									Computed: true,
								},
								"sni_strict": schema.BoolAttribute{
									Optional: true,
									Computed: true,
								},
								"sni_host_match": schema.BoolAttribute{
									Optional: true,
									Computed: true,
								},
								"client_renegotiation": schema.BoolAttribute{
									Optional: true,
									Computed: true,
								},
								"options": schema.StringAttribute{
									Optional: true,
									Computed: true,
								},
								"cipher_list": schema.StringAttribute{
									Optional: true,
									Computed: true,
								},
								"named_curve": schema.StringAttribute{
									Optional: true,
									Computed: true,
								},
								"oscp": schema.BoolAttribute{
									Optional: true,
									Computed: true,
								},
								"pem": schema.StringAttribute{
									Optional: true,
									Computed: true,
								},
								"ca": schema.StringAttribute{
									Optional: true,
									Computed: true,
								},
							},
						},
						"directory": schema.StringAttribute{
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"experiments": schema.ListAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"edge_functions_sources": schema.MapAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"edge_function_init_script": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
		},
	}
}

func (r *CDNConfigurationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.CDNConfigurationModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	cdnConfig := utility.ConvertCdnConfigToNative(&plan)
	cfg, err := r.client.UploadCdnConfiguration(&cdnConfig)

	if err != nil {
		resp.Diagnostics.AddError("Error creating CDN configuration", err.Error())
		return
	}

	status, err := r.client.GetCDNConfiguration(cfg.ConfigurationID)

	if err != nil {
		resp.Diagnostics.AddError("Error reading CDN configuration", err.Error())
		return
	}

	state := utility.ConvertNativeToCdnConfig(status)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *CDNConfigurationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.CDNConfigurationModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	cdnConfig, err := r.client.GetCDNConfiguration(state.ConfigurationID.ValueString())

	if err != nil {
		resp.Diagnostics.AddError("Error reading CDN configuration", err.Error())
		return
	}

	state = utility.ConvertNativeToCdnConfig(cdnConfig)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *CDNConfigurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.CDNConfigurationModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	cdnConfig := utility.ConvertCdnConfigToNative(&plan)

	cfg, err := r.client.UploadCdnConfiguration(&cdnConfig)

	if err != nil {
		resp.Diagnostics.AddError("Error creating CDN configuration", err.Error())
		return
	}

	status, err := r.client.GetCDNConfiguration(cfg.ConfigurationID)

	if err != nil {
		resp.Diagnostics.AddError("Error reading CDN configuration", err.Error())
		return
	}

	state := utility.ConvertNativeToCdnConfig(status)
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *CDNConfigurationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddWarning("Operation Not Supported", "This resource does not support deletion.")
}
