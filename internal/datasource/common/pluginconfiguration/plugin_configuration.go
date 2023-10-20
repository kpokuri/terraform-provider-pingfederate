package pluginconfiguration

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	fieldAttrTypes = map[string]attr.Type{
		"name":            basetypes.StringType{},
		"value":           basetypes.StringType{},
		"encrypted_value": basetypes.StringType{},
		"inherited":       basetypes.BoolType{},
	}

	rowAttrTypes = map[string]attr.Type{
		"fields":      basetypes.ListType{ElemType: basetypes.ObjectType{AttrTypes: fieldAttrTypes}},
		"default_row": basetypes.BoolType{},
	}

	tableAttrTypes = map[string]attr.Type{
		"name":      basetypes.StringType{},
		"rows":      basetypes.ListType{ElemType: basetypes.ObjectType{AttrTypes: rowAttrTypes}},
		"inherited": basetypes.BoolType{},
	}

	configurationAttrTypes = map[string]attr.Type{
		"fields": basetypes.ListType{ElemType: types.ObjectType{AttrTypes: fieldAttrTypes}},
		"tables": basetypes.ListType{ElemType: types.ObjectType{AttrTypes: tableAttrTypes}},
	}
)

func AttrType() map[string]attr.Type {
	return configurationAttrTypes
}

func ToDataSourceSchema() schema.SingleNestedAttribute {
	fieldsListDefault, _ := types.ListValue(types.ObjectType{AttrTypes: fieldAttrTypes}, []attr.Value{})
	return schema.SingleNestedAttribute{
		Description: "Plugin instance configuration.",
		Required:    false,
		Optional:    false,
		Computed:    true,
		Attributes: map[string]schema.Attribute{
			"tables": schema.ListNestedAttribute{
				Description: "List of configuration tables.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The name of the table.",
							Required:    false,
							Optional:    false,
							Computed:    true,
						},
						"rows": schema.ListNestedAttribute{
							Description: "List of table rows.",
							Required:    false,
							Optional:    false,
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"fields": schema.ListNestedAttribute{
										Description: "The configuration fields in the row.",
										Required:    false,
										Optional:    false,
										Computed:    true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"name": schema.StringAttribute{
													Description: "The name of the configuration field.",
													Required:    false,
													Optional:    false,
													Computed:    true,
												},
												"value": schema.StringAttribute{
													Description: "The value for the configuration field. For encrypted or hashed fields, GETs will not return this attribute.",
													Required:    false,
													Optional:    false,
													Computed:    true,
												},
												"encrypted_value": schema.StringAttribute{
													Description: "For encrypted or hashed fields, this attribute contains the encrypted representation of the field's value, if a value is defined.",
													Required:    false,
													Optional:    false,
													Computed:    true,
												},
												"inherited": schema.BoolAttribute{
													Description: "Whether this field is inherited from its parent instance. If true, the value/encrypted value properties become read-only. The default value is false.",
													Required:    false,
													Optional:    false,
													Computed:    true,
												},
											},
										},
									},
									"default_row": schema.BoolAttribute{
										Description: "Whether this row is the default.",
										Required:    false,
										Optional:    false,
										Computed:    true,
									},
								},
							},
						},
						"inherited": schema.BoolAttribute{
							Description: "Whether this table is inherited from its parent instance. If true, the rows become read-only. The default value is false.",
							Required:    false,
							Optional:    false,
							Computed:    true,
						},
					},
				},
			},
			"fields": schema.ListNestedAttribute{
				Description: "List of configuration fields.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				Default:     listdefault.StaticValue(fieldsListDefault),
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "The name of the configuration field.",
							Required:    false,
							Optional:    false,
							Computed:    true,
						},
						"value": schema.StringAttribute{
							Description: "The value for the configuration field. For encrypted or hashed fields, GETs will not return this attribute.",
							Required:    false,
							Optional:    false,
							Computed:    true,
						},
						"encrypted_value": schema.StringAttribute{
							Description: "For encrypted or hashed fields, this attribute contains the encrypted representation of the field's value, if a value is defined.",
							Required:    false,
							Optional:    false,
							Computed:    true,
						},
						"inherited": schema.BoolAttribute{
							Description: "Whether this field is inherited from its parent instance. If true, the value/encrypted value properties become read-only. The default value is false.",
							Required:    false,
							Optional:    false,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}
