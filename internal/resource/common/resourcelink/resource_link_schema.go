package resourcelink

import (
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func Schema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description: "The ID of the resource.",
			Required:    true,
		},
		"location": schema.StringAttribute{
			Description: "A read-only URL that references the resource. If the resource is not currently URL-accessible, this property will be null.",
			Computed:    true,
			Optional:    false,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
	}
}