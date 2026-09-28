package capability

import (
	"strings"
	"testing"
)

func TestInspectDeclarationShape(t *testing.T) {
	shape, err := InspectDeclaration("package jev\nnamespace example\nentity Receipt\nproperty status string\nactivity observe")
	if err != nil {
		t.Fatal(err)
	}
	if shape.Package != "jev" || shape.Namespace != "example" || len(shape.Entities) != 1 || len(shape.Activities) != 1 {
		t.Fatalf("unexpected shape: %#v", shape)
	}
	if err := shape.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestInspectDeclarationRejectsOrphanProperty(t *testing.T) {
	if _, err := InspectDeclaration("package jev\nproperty status string"); err == nil {
		t.Fatal("orphan property should be rejected")
	}
}

func TestInspectDeclarationRejectsTampering(t *testing.T) {
	shape, err := InspectDeclaration("package jev\nactivity observe")
	if err != nil {
		t.Fatal(err)
	}
	shape.ShapeDigest = strings.Repeat("0", 64)
	if err := shape.Validate(); err == nil {
		t.Fatal("tampered declaration shape should be rejected")
	}
}
