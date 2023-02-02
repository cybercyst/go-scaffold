package template

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestTemplateConfigShouldUnarshalValidConfig(t *testing.T) {
	file := []byte(strings.TrimSpace(`
title: My Template
version: v1.0.0
description: |
  My really great template that does a lot of magical
  things and is really beyond belief
tags:
  - template
  - tag_one
  - tag_two
schema:
  title: Please input your project name
  type: object
  required:
    - project_name
  properties:
    project_name:
      title: Name
      type: string
`))

	var got TemplateConfig
	err := yaml.Unmarshal(file, &got)
	if err != nil {
		t.Errorf("unexpected error while unmarshalling config: %s", err)
	}

	want := "My Template"
	if got.Title != want {
		t.Errorf("expected value %s and got %s", want, got)
	}

	want = "v1.0.0"
	if got.Version != want {
		t.Errorf("expected value %s and got %s", want, got)
	}

	// optional fields are fine
	want = ""
	if got.Description != want {
		t.Errorf("expected value %s and got %s", want, got)
	}
}
