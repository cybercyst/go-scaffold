package template

import (
	"reflect"
	"strings"
	"testing"

	"github.com/spf13/afero"
)

func TestReadYamlConfig(t *testing.T) {
	fs := afero.NewMemMapFs()
	afero.WriteFile(fs, "template.yaml", []byte(strings.TrimSpace(`
title: My Test Template
version: v1.0.0
`)), 0644)

	got, err := LoadConfig(fs)
	if err != nil {
		t.Fatalf("unexpected error thrown while reading template config: %s", err)
	}

	want := &TemplateConfig{
		Title:   "My Test Template",
		Version: "v1.0.0",
	}

	if reflect.DeepEqual(got, want) == false {
		t.Fatalf("want %+v and got %+v", want, got)
	}
}

func TestReadYmlConfig(t *testing.T) {
	fs := afero.NewMemMapFs()
	afero.WriteFile(fs, "template.yml", []byte(strings.TrimSpace(`
title: My Test Template YML
version: v1.0.0
`)), 0644)

	got, err := LoadConfig(fs)
	if err != nil {
		t.Fatalf("unexpected error thrown while reading template config: %s", err)
	}

	want := &TemplateConfig{
		Title:   "My Test Template YML",
		Version: "v1.0.0",
	}

	if reflect.DeepEqual(got, want) == false {
		t.Fatalf("want %+v and got %+v", want, got)
	}
}

func TestReadJsonConfig(t *testing.T) {
	fs := afero.NewMemMapFs()
	afero.WriteFile(fs, "template.json", []byte(strings.TrimSpace(`
{"title": "My Test Template JSON", "version": "v1.0.0"}
`)), 0644)

	got, err := LoadConfig(fs)
	if err != nil {
		t.Fatalf("unexpected error thrown while reading template config: %s", err)
	}

	want := &TemplateConfig{
		Title:   "My Test Template JSON",
		Version: "v1.0.0",
	}

	if reflect.DeepEqual(got, want) == false {
		t.Fatalf("want %+v and got %+v", want, got)
	}
}

func TestNoTemplateConfig(t *testing.T) {
	fs := afero.NewMemMapFs()

	_, err := LoadConfig(fs)
	if err == nil {
		t.Fatal("failed to get expected error when missing template config")
	}
}

func TestTemplateConfigMissingRequiredFields(t *testing.T) {
	fs := afero.NewMemMapFs()
	afero.WriteFile(fs, "template.yaml", []byte(""), 0644)

	_, err := LoadConfig(fs)
	if err == nil {
		t.Fatal("failed to get expected error when config is missing required properties")
	}
}
