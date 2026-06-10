package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// LoadYAMLWithComments loads a YAML file preserving comments.
func LoadYAMLWithComments(path string) (*yaml.Node, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var doc yaml.Node
	if err := yaml.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return &doc, nil
}

// SaveYAMLWithComments saves a YAML node to file preserving comments.
func SaveYAMLWithComments(path string, doc *yaml.Node) error {
	data, err := yaml.Marshal(doc)
	if err != nil {
		return fmt.Errorf("failed to marshal YAML: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// FindMappingValue finds a value node in a mapping by key.
func FindMappingValue(node *yaml.Node, key string) *yaml.Node {
	if node == nil || node.Kind != yaml.MappingNode {
		return nil
	}

	for i := 0; i < len(node.Content)-1; i += 2 {
		keyNode := node.Content[i]
		if keyNode.Value == key {
			return node.Content[i+1]
		}
	}

	return nil
}

// SetMappingValue sets a value in a mapping node.
// If the key exists, it updates the value.
// If the key doesn't exist, it adds a new key-value pair.
func SetMappingValue(node *yaml.Node, key string, value interface{}) error {
	if node == nil || node.Kind != yaml.MappingNode {
		return fmt.Errorf("node is not a mapping")
	}

	// Find existing key
	for i := 0; i < len(node.Content)-1; i += 2 {
		keyNode := node.Content[i]
		if keyNode.Value == key {
			// Update existing value
			valueNode := node.Content[i+1]
			return setNodeValue(valueNode, value)
		}
	}

	// Key not found, add new key-value pair
	keyNode := &yaml.Node{
		Kind:  yaml.ScalarNode,
		Value: key,
		Tag:   "!!str",
	}

	valueNode := &yaml.Node{}
	if err := setNodeValue(valueNode, value); err != nil {
		return err
	}

	node.Content = append(node.Content, keyNode, valueNode)
	return nil
}

// setNodeValue sets the value of a scalar node.
func setNodeValue(node *yaml.Node, value interface{}) error {
	switch v := value.(type) {
	case string:
		node.Kind = yaml.ScalarNode
		node.Value = v
		node.Tag = "!!str"
		node.Style = 0
	case int:
		node.Kind = yaml.ScalarNode
		node.Value = fmt.Sprintf("%d", v)
		node.Tag = "!!int"
		node.Style = 0
	case int64:
		node.Kind = yaml.ScalarNode
		node.Value = fmt.Sprintf("%d", v)
		node.Tag = "!!int"
		node.Style = 0
	case bool:
		node.Kind = yaml.ScalarNode
		if v {
			node.Value = "true"
		} else {
			node.Value = "false"
		}
		node.Tag = "!!bool"
		node.Style = 0
	case float64:
		node.Kind = yaml.ScalarNode
		node.Value = fmt.Sprintf("%g", v)
		node.Tag = "!!float"
		node.Style = 0
	default:
		return fmt.Errorf("unsupported type: %T", value)
	}
	return nil
}

// GetStringValue gets a string value from a mapping node.
func GetStringValue(node *yaml.Node, key string) string {
	valueNode := FindMappingValue(node, key)
	if valueNode == nil {
		return ""
	}
	return valueNode.Value
}

// GetIntValue gets an int value from a mapping node.
func GetIntValue(node *yaml.Node, key string) int {
	valueNode := FindMappingValue(node, key)
	if valueNode == nil {
		return 0
	}
	var v int
	fmt.Sscanf(valueNode.Value, "%d", &v)
	return v
}

// GetBoolValue gets a bool value from a mapping node.
func GetBoolValue(node *yaml.Node, key string) bool {
	valueNode := FindMappingValue(node, key)
	if valueNode == nil {
		return false
	}
	return valueNode.Value == "true"
}
