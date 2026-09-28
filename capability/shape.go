package capability

import (
	"fmt"
	"strings"

	"github.com/kimjooyoon/jev-gooo/envelope"
)

type PropertyShape struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type EntityShape struct {
	Name       string         `json:"name"`
	Properties []PropertyShape `json:"properties"`
}

type DeclarationShape struct {
	DeclarationDigest string       `json:"declaration_digest"`
	Package          string       `json:"package"`
	Namespace        string       `json:"namespace,omitempty"`
	Entities         []EntityShape `json:"entities"`
	Activities       []string      `json:"activities"`
	ShapeDigest      string       `json:"shape_digest"`
}

func InspectDeclaration(source string) (DeclarationShape, error) {
	declaration, err := envelope.BindDeclaration(source)
	if err != nil {
		return DeclarationShape{}, err
	}
	shape := DeclarationShape{
		DeclarationDigest: declaration.Digest,
		Entities:          []EntityShape{},
		Activities:        []string{},
	}
	currentEntity := -1
	for lineNumber, raw := range strings.Split(source, string([]byte{10})) {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "//") || strings.HasPrefix(line, "#") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		switch fields[0] {
		case "package":
			if len(fields) < 2 {
				return DeclarationShape{}, fmt.Errorf("line %d: package name is required", lineNumber+1)
			}
			shape.Package = fields[1]
		case "namespace":
			if len(fields) < 2 {
				return DeclarationShape{}, fmt.Errorf("line %d: namespace name is required", lineNumber+1)
			}
			shape.Namespace = fields[1]
		case "entity":
			if len(fields) < 2 {
				return DeclarationShape{}, fmt.Errorf("line %d: entity name is required", lineNumber+1)
			}
			shape.Entities = append(shape.Entities, EntityShape{Name: fields[1], Properties: []PropertyShape{}})
			currentEntity = len(shape.Entities) - 1
		case "property":
			if currentEntity < 0 {
				return DeclarationShape{}, fmt.Errorf("line %d: property must belong to an entity", lineNumber+1)
			}
			if len(fields) < 3 {
				return DeclarationShape{}, fmt.Errorf("line %d: property name and type are required", lineNumber+1)
			}
			shape.Entities[currentEntity].Properties = append(shape.Entities[currentEntity].Properties, PropertyShape{
				Name: fields[1],
				Type: strings.Join(fields[2:], " "),
			})
		case "activity":
			if len(fields) < 2 {
				return DeclarationShape{}, fmt.Errorf("line %d: activity name is required", lineNumber+1)
			}
			shape.Activities = append(shape.Activities, fields[1])
		}
	}
	if shape.Package == "" {
		return DeclarationShape{}, fmt.Errorf("gooo declaration package is required")
	}
	shape.ShapeDigest = shape.digest()
	return shape, nil
}

func (s DeclarationShape) Validate() error {
	if !validDigest(s.DeclarationDigest) || !validDigest(s.ShapeDigest) || strings.TrimSpace(s.Package) == "" {
		return fmt.Errorf("declaration shape is incomplete")
	}
	for _, entity := range s.Entities {
		if strings.TrimSpace(entity.Name) == "" {
			return fmt.Errorf("declaration shape entity name is required")
		}
		for _, property := range entity.Properties {
			if strings.TrimSpace(property.Name) == "" || strings.TrimSpace(property.Type) == "" {
				return fmt.Errorf("declaration shape property is incomplete")
			}
		}
	}
	for _, activity := range s.Activities {
		if strings.TrimSpace(activity) == "" {
			return fmt.Errorf("declaration shape activity name is required")
		}
	}
	if s.digest() != s.ShapeDigest {
		return fmt.Errorf("declaration shape digest does not match")
	}
	return nil
}

func (s DeclarationShape) digest() string {
	parts := []string{"gooo-declaration-shape", s.DeclarationDigest, s.Package, s.Namespace}
	for _, entity := range s.Entities {
		parts = append(parts, "entity", entity.Name)
		for _, property := range entity.Properties {
			parts = append(parts, "property", property.Name, property.Type)
		}
	}
	for _, activity := range s.Activities {
		parts = append(parts, "activity", activity)
	}
	return digest(parts...)
}
