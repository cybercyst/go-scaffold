package generate

import (
	"testing"

	"github.com/spf13/afero"
)

func TestGenerateYamlFilter(t *testing.T) {
	templateFs := afero.NewMemMapFs()
	afero.WriteFile(templateFs, "template_test", []byte(`{{ yaml | yaml }}`), 0644)
	outputFs := afero.NewMemMapFs()
	input := map[string]interface{}{
		"yaml": map[string]interface{}{
			"strVal":  "a",
			"numVal":  123,
			"boolVal": true,
			"arrayVal": []string{
				"one",
				"two",
				"banana",
			},
			"objVal": map[string]interface{}{
				"nestedVal": "hi there",
			},
		},
	}

	_, err := GenerateTemplateFiles(templateFs, outputFs, &input)
	if err != nil {
		t.Fatalf("unexpected error thrown while generating template: %s", err)
	}

	got, err := afero.ReadFile(outputFs, "template_test")
	if err != nil {
		t.Fatalf("unexpected error while reading generated file 'template_test': %s", err)
	}

	// interface values are sorted alphabetically
	want := `arrayVal:
    - one
    - two
    - banana
boolVal: true
numVal: 123
objVal:
    nestedVal: hi there
strVal: a
`
	if string(got) != want {
		t.Fatalf("expected %s but got %s", want, got)
	}
}

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

func TestExecuteTemplateShouldHaveSlugifyFilter(t *testing.T) {
	template := "{{ my_variable | slugify }}"
	input := map[string]interface{}{
		"my_variable": "My WeIrDly Capped Value",
	}
	got, err := executeTemplate([]byte(template), &input)
	if err != nil {
		t.Fatalf("unexpected error while executing template: %s", err)
	}
	want := "my-weirdly-capped-value"

	if got != want {
		t.Fatalf("expected %s to be %s after executing template", got, want)
	}
}
