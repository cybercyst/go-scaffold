package generate

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

	gotCreatedFiles, err := GenerateTemplateFiles(templateFs, outputFs, &input)
	if err != nil {
		t.Fatalf("unexpected error thrown while generating template: %s", err)
	}

	wantCreatedFiles := []string{"My Project.md"}
	if len(gotCreatedFiles) != len(wantCreatedFiles) {
		t.Fatalf("received wrong number of created files as response")
	}
	if gotCreatedFiles[0] != wantCreatedFiles[0] {
		t.Fatalf("received wrong created files as response")
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
