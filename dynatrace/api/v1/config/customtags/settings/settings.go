/**
* @license
* Copyright 2020 Dynatrace LLC
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package customtags

import (
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/api/v1/config/common"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/terraform/hcl"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Settings struct {
	EntitySelector  string            `json:"-"`    // Specifies the entities where you want to update tags.
	Tags            common.TagFilters `json:"tags"` // A list of tags to be added to monitored entities.
	MatchedEntities int64             `json:"matchedEntitiesCount"`
	CurrentState    string            `json:"-"`
	ExportName      string            `json:"-"` // If set it will get used preferrably before the EntitySelector as the Name of this setting
}

func (me *Settings) Schema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"entity_selector": {
			Type:        schema.TypeString,
			Description: "Specifies the entities where you want to update tags",
			Required:    true,
			ForceNew:    true,
		},
		"tags": {
			Type:        schema.TypeList,
			Description: "Specifies the entities where you want to update tags",
			MaxItems:    1,
			Required:    true,
			Elem:        &schema.Resource{Schema: new(common.TagFilters).Schema()},
		},
		"matched_entities": {
			Type:        schema.TypeInt,
			Description: "The number of monitored entities where the tags have been added.",
			Optional:    true,
			Computed:    true,
		},
		"current_state": {
			Type:        schema.TypeString,
			Description: "For internal use: current state of tags in JSON format",
			Optional:    true,
			Computed:    true,
		},
	}
}

func (me *Settings) Name() string {
	if len(me.ExportName) > 0 {
		return me.ExportName
	}
	return me.EntitySelector
}

func (me *Settings) MarshalHCL(properties hcl.Properties) error {
	return properties.EncodeAll(map[string]interface{}{
		"entity_selector": me.EntitySelector,
		"tags":            me.Tags,
	})
}

func (me *Settings) UnmarshalHCL(decoder hcl.Decoder) error {
	return decoder.DecodeAll(map[string]interface{}{
		"entity_selector": &me.EntitySelector,
		"tags":            &me.Tags,
	})
}
