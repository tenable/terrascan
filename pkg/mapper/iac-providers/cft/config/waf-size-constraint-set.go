/*
    Copyright (C) 2021 Accurics, Inc.

	Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/

package config

import "github.com/awslabs/goformation/v5/cloudformation/waf"

type FieldToMatchBlock struct {
	Data string `json:"data,omitempty"`
	Type string `json:"type"`
}

type SizeConstraintSetBlock struct {
	ComparisonOperator string              `json:"comparison_operator"`
	Size               int                 `json:"size"`
	TextTransformation string              `json:"text_transformation"`
	FieldToMatch       []FieldToMatchBlock `json:"field_to_match"`
}

type WafSizeConstraintSetConfig struct {
	Config
	Name              string                   `json:"name"`
	SizeConstraintSet []SizeConstraintSetBlock `json:"size_constraints,omitempty"`
}

func GetWafSizeConstraintSetConfig(w *waf.SizeConstraintSet) []AWSResourceConfig {
	sizeConstraintSet := setSizeConstraintSet(w)

	cf := WafSizeConstraintSetConfig{
		Config: Config{
			Name: w.Name,
		},
		Name:              w.Name,
		SizeConstraintSet: sizeConstraintSet,
	}
	return []AWSResourceConfig{{
		Resource: cf,
		Metadata: w.AWSCloudFormationMetadata,
	}}
}

func setSizeConstraintSet(w *waf.SizeConstraintSet) []SizeConstraintSetBlock {
	sizeConstraintSet := make([]SizeConstraintSetBlock, len(w.SizeConstraints))

	for i := 0; i < len(w.SizeConstraints); i++ {
		sizeConstraintSet[i].Size = w.SizeConstraints[i].Size
		sizeConstraintSet[i].TextTransformation = w.SizeConstraints[i].TextTransformation
		sizeConstraintSet[i].ComparisonOperator = w.SizeConstraints[i].ComparisonOperator
		sizeConstraintSet[i].FieldToMatch = setFieldToMatch(w, i)
	}

	return sizeConstraintSet
}

func setFieldToMatch(w *waf.SizeConstraintSet, index int) []FieldToMatchBlock {
	fieldToMatchBlock := make([]FieldToMatchBlock, 1)
	fieldToMatchBlock[0].Type = w.SizeConstraints[index].FieldToMatch.Type

	if w.SizeConstraints[index].FieldToMatch.Data != "" {
		fieldToMatchBlock[0].Data = w.SizeConstraints[index].FieldToMatch.Data
	}

	return fieldToMatchBlock
}
