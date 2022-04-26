package linenumber

import (
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type lineObj map[string]interface{}

func GetLineConfig(data []byte) (lineObj, error) {
	docNode := yaml.Node{}
	if len(data) == 0 {
		return nil, nil
	}

	err := yaml.Unmarshal(data, &docNode)
	if err != nil {
		return nil, err
	}

	if len(docNode.Content) == 0 {
		return nil, err
	}
	return ParseNodes(*docNode.Content[0]), nil
}

func ParseNodes(nodes yaml.Node) lineObj {
	obj := make(map[string]interface{})
	for i := 0; i < len(nodes.Content)-1; i = i + 2 {
		if nodes.Content[i+1].Kind == yaml.ScalarNode {
			obj[nodes.Content[i].Value] = nodes.Content[i+1].Line
			continue
		}
		if nodes.Content[i+1].Kind == yaml.SequenceNode {
			obj[nodes.Content[i].Value] = []lineObj{}
			for _, cont := range nodes.Content[i+1].Content {
				val := obj[nodes.Content[i].Value].([]lineObj)
				if len(cont.Content) == 0 {
					obj[nodes.Content[i].Value] = append(val,
						lineObj{
							cont.Value: cont.Line,
						},
					)
					continue
				}
				obj[nodes.Content[i].Value] = append(val, ParseNodes(*cont))
			}
			continue
		}
		if nodes.Content[i+1].Kind == yaml.MappingNode {
			obj[nodes.Content[i].Value] = ParseNodes(*nodes.Content[i+1])
		}
	}
	return obj
}

func GetAttributeLineNo(linesData interface{}, traversalPath string) int {
	lines := linesData.(lineObj)
	TraverseEelement, _ := getTraverseElements(traversalPath)
	for _, elem := range TraverseEelement {
		if elem.Position == nil {
			val := lines[elem.Name]
			if val != nil {
				lines = val.(lineObj)
				continue
			}
			break
		}
		linesArray := lines[elem.Name].([]lineObj)
		lines = linesArray[*elem.Position]
	}

	for _, line := range lines {
		if val, ok := line.(int); ok {
			return val
		}
	}
	return -1
}

type TraverseEelement struct {
	Name     string
	Position *int64
}

func getTraverseElements(traversalPath string) ([]*TraverseEelement, error) {
	traverseElems := make([]*TraverseEelement, 0)
	if len(traversalPath) == 0 {
		return nil, fmt.Errorf("traversal information not available")
	}
	pathElems := strings.Split(traversalPath, ".")
	for i := range pathElems {
		elems := strings.Split(pathElems[i], "[")
		t := TraverseEelement{
			Name: elems[0],
		}
		if len(elems) > 1 {
			x := elems[1][:len(elems[1])-1]
			j, err := strconv.ParseInt(x, 10, strconv.IntSize)
			if err != nil {
				return nil, fmt.Errorf("incorrect value for index in %s, error: %v", pathElems[i], err)
			}
			t.Position = &j
		}
		traverseElems = append(traverseElems, &t)
	}
	return traverseElems, nil
}
