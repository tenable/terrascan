/*
    Copyright (C) 2022 Tenable, Inc.

	Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
	Package cft provides the Template type that models a CloudFormation template.

 	The sub-packages of cft contains various tools for working with templates
*/

package cftv1

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Tags is a mapping from YAML short tags to full instrincic function names
var Tags = map[string]string{
	"!And":           "Fn::And",
	"!Base64":        "Fn::Base64",
	"!Cidr":          "Fn::Cidr",
	"!Equals":        "Fn::Equals",
	"!FindInMap":     "Fn::FindInMap",
	"!GetAZs":        "Fn::GetAZs",
	"!GetAtt":        "Fn::GetAtt",
	"!If":            "Fn::If",
	"!ImportValue":   "Fn::ImportValue",
	"!Join":          "Fn::Join",
	"!Not":           "Fn::Not",
	"!Or":            "Fn::Or",
	"!Select":        "Fn::Select",
	"!Split":         "Fn::Split",
	"!Sub":           "Fn::Sub",
	"!Ref":           "Ref",
	"!Condition":     "Condition",
	"!Rain::Embed":   "Rain::Embed",
	"!Rain::Include": "Rain::Include",
	"!Rain::Env":     "Rain::Env",
	"!Rain::S3Http":  "Rain::S3Http",
	"!Rain::S3":      "Rain::S3",
	"!Rain::Module":  "Rain::Module",
}

// Template represents a CloudFormation template. The Template type
// is minimal for now but will likely grow new features as needed by rain.
type Template struct {
	*yaml.Node
}

// TODO - We really need a convenient Template data structure
// that lets us easily access elements.
// t.Resources["MyResource"].Properties["MyProp"]
//
// Add a Model attribute to the struct and an Init function to populate it.
// t.Model.Resources

// Map returns the template as a map[string]interface{}
func (t Template) Map() (map[string]interface{}, error) {
	var out map[string]interface{}
	err := t.Decode(&out)
	if err != nil {
		return nil, fmt.Errorf("error converting template to map: %s", err)
	}
	return out, nil
}

// AppendStateMap appends a "State" section to the template
func AppendStateMap(state Template) *yaml.Node {
	state.Node.Content[0].Content = append(state.Node.Content[0].Content,
		&yaml.Node{Kind: yaml.ScalarNode, Value: "State"})
	stateMap := &yaml.Node{Kind: yaml.MappingNode, Content: make([]*yaml.Node, 0)}
	state.Node.Content[0].Content = append(state.Node.Content[0].Content, stateMap)
	return stateMap
}

// File reads the file and returns string
func (a *CFTV1) File(fileName string) (Template, error) {
	source, err := os.ReadFile(fileName)
	if err != nil {
		return Template{}, fmt.Errorf("unable to read file: %s", err)
	}

	return String(string(source))
}

// String returns a cft.Template parsed from a string
func String(input string) (Template, error) {
	var node yaml.Node
	err := yaml.Unmarshal([]byte(input), &node)
	if err != nil {
		return Template{}, fmt.Errorf("invalid YAML: %s", err)
	}

	return Node(&node)
}

// Node returns a cft.Template parse from a *yaml.Node
func Node(node *yaml.Node) (Template, error) {
	err := TransformNode(node)
	return Template{Node: node}, err
}

// TransformNode takes a *yaml.Node and convert tag-style names into map-style,
// and converts other scalars into a canonical format
func TransformNode(n *yaml.Node) error {
	// Fix badly-parsed numbers
	if n.ShortTag() == "!!float" && n.Value[0] == '0' {
		n.Tag = "!!str"
	}

	// Fix badly-parsed timestamps which are often used for versions in cloudformation
	if n.ShortTag() == "!!timestamp" {
		n.Tag = "!!str"
	}

	// Convert tag-style intrinsics into map-style
	for tag, funcName := range Tags {
		if n.ShortTag() == tag {
			body := Clone(n)

			// Fix empty Fn values (should never be null)
			if body.Tag == "!!null" {
				body.Tag = "!!str"
			} else {
				body.Tag = ""
			}

			// Wrap in a map
			*n = yaml.Node{
				Kind: yaml.MappingNode,
				Tag:  "!!map",
				Content: []*yaml.Node{
					{
						Kind:  yaml.ScalarNode,
						Style: 0,
						Tag:   "!!str",
						Value: funcName,
					},
					body,
				},
			}
			break
		}
	}
	// Convert GetAtts
	if n.Kind == yaml.MappingNode && len(n.Content) == 2 {
		if n.Content[0].Value == "Fn::GetAtt" && n.Content[1].Kind == yaml.ScalarNode {
			err := parseGetAtt(n.Content[1])
			if err != nil {
				return err
			}
		}
	}
	for _, child := range n.Content {
		err := TransformNode(child)
		if err != nil {
			return err
		}
	}
	return nil
}

// Convert string GetAtt into array format so that it it's easier to compare
func parseGetAtt(n *yaml.Node) error {
	parts := strings.SplitN(n.Value, ".", 2)

	if len(parts) != 2 {
		return errors.New("GetAtt requires two parameters")
	}

	*n = yaml.Node{
		Kind: yaml.SequenceNode,

		HeadComment: n.HeadComment,
		LineComment: n.LineComment,
		FootComment: n.FootComment,

		Content: []*yaml.Node{
			{
				Kind:  yaml.ScalarNode,
				Style: 0,
				Tag:   "!!str",
				Value: parts[0],
			},
			{
				Kind:  yaml.ScalarNode,
				Style: 0,
				Tag:   "!!str",
				Value: parts[1],
			},
		},
	}
	return nil
}

// Clone returns a copy of the provided node
func Clone(node *yaml.Node) *yaml.Node {
	if node == nil {
		return nil
	}

	out := &yaml.Node{
		Kind:        node.Kind,
		Style:       node.Style,
		Tag:         node.Tag,
		Value:       node.Value,
		Anchor:      node.Anchor,
		Alias:       Clone(node.Alias),
		Content:     make([]*yaml.Node, len(node.Content)),
		HeadComment: node.HeadComment,
		LineComment: node.LineComment,
		FootComment: node.FootComment,
		Line:        node.Line,
		Column:      node.Column,
	}

	for i, child := range node.Content {
		out.Content[i] = Clone(child)
	}
	return out
}
