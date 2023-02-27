package template

import (
	"testing"

	"github.com/cybercyst/go-scaffold/internal/utils"
	"github.com/spf13/afero"

	"github.com/stretchr/testify/assert"
)

func TestDownloadShouldFetchAllSteps(t *testing.T) {
	fs := utils.CreateTestFs(map[string]string{
		"/template1/template.yaml": `
title: Template One
version: v0.0.1
steps:
  - source: ./template
  - source: /template2
		`,
		"/template1/template/file_one.txt": `
{{ var_one }}
		`,
		"/template2/template.yaml": `
title: Template Two
version: v0.0.2
steps:
  - source: ./template
		`,
		"/template2/template/file_two.txt": `
{{ var_two }}
		`,
	})

	templateFs := afero.NewBasePathFs(fs, "/template1")
	isDir, _ := afero.IsDir(templateFs, "./template")
	assert.True(t, isDir)

	got, err := NewTemplate(fs, "/template1")
	assert.NoError(t, err)

	assert.NotNil(t, got)
	assert.Len(t, got.Templates, 2)
}

func TestDownloadShouldThrowErrorWithRecursiveTemplateDependency(t *testing.T) {
	// Template 1 needs Template 2 which needs Template 1 (This is a circular dependency and shouldn't be allowed)
	fs := utils.CreateTestFs(map[string]string{
		"/template1/template.yaml": `
title: Template One
version: v0.0.1
steps:
  - source: ./template
  - source: /template2
		`,
		"/template1/template/file_one.txt": `
{{ var_one }}
		`,
		"/template2/template.yaml": `
title: Template Two
version: v0.0.2
steps:
  - source: ./template
  - source: /template1
		`,
		"/template2/template/file_two.txt": `
{{ var_two }}
		`,
	})

	templateFs := afero.NewBasePathFs(fs, "/template1")
	isDir, _ := afero.IsDir(templateFs, "./template")
	assert.True(t, isDir)

	got, err := NewTemplate(fs, "/template1")
	assert.ErrorContains(t, err, "using template at /template1 would cause a circular dependency")
	assert.Nil(t, got)
}
