package go_cookiecutter

import (
	"testing"

	"github.com/spf13/afero"
)

func TestGenerateTemplateFiles(t *testing.T) {
	templateFs := afero.NewMemMapFs()
	afero.WriteFile(templateFs, "{{ project_name }}.md", []byte(`
This file should be rendered with variables replaced.
project_name is {{ project_name }}
`), 0644)
	outputFs := afero.NewMemMapFs()
	input := map[string]interface{}{
		"project_name": "My Project",
	}

	err := generateTemplateFiles(templateFs, outputFs, &input)
	if err != nil {
		t.Fatalf("unexpected error thrown while generating template: %s", err)
	}

	_, err = outputFs.Stat("My Project.md")
	if err != nil {
		t.Fatalf("expected file 'My Project.md' not generated: %s", err)
	}

	got, err := afero.ReadFile(outputFs, "My Project.md")
	if err != nil {
		t.Fatalf("unexpected error while reading generated file 'My Project.md': %s", err)
	}

	want := `
This file should be rendered with variables replaced.
project_name is My Project
`
	if string(got) != want {
		t.Fatalf("expected %s but got %s", want, got)
	}
}
