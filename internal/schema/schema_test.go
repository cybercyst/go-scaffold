package schema

import (
	"testing"

	"github.com/spf13/afero"
)

func TestValidateSchemaShouldWorkWithBackstageLikeSchema(t *testing.T) {
	fs := afero.NewMemMapFs()
	afero.WriteFile(fs, "schema.yaml", []byte(`
  - title: 'Let''s get started, please name your web service'
    required:
      - project_name
    properties:
      project_name:
        title: Name
        type: string
        description: Unique name of the component
        'ui:autofocus': true
      owner:
        title: Owner
        type: string
        description: Then name of the person/team responsible for this web service
        'ui:autofocus': true
  - title: Let's configure your project's Docker image
    required:
      - image_repository
    properties:
      image_registry:
        title: Image Registry
        type: string
        description: >-
          The registry you would like to store your image in, default is Docker
          Hub
        default: docker.io
        'ui:autofocus': true
      image_repository:
        title: Repository Name
        type: string
        description: The repository you want to push your Docker image to
        'ui:autofocus': true
      image_tag:
        title: Initial Image Tag
        type: string
        description: 'Initial version of the Docker image tag, e.g. 0.0.1'
        'ui:autofocus': true
        default: 0.0.1
      image_port:
        title: Container Port
        type: integer
        description: >-
          The port on your host machine that your running container will be
          mapped to, default is 4000
        default: 4000
        'ui:autofocus': true
`), 0644)

	schema, err := LoadSchema(fs)
	if err != nil {
		t.Fatalf("got unexpected error while parsing schema: %s", err)
	}

	input := map[string]interface{}{
		"project_name":     "my-project",
		"image_repository": "my-image-repo-name",
	}

	err = ValidateInput(schema, &input)
	if err != nil {
		t.Error("Got error when validating test input", err)
	}
}

func TestValidateSchemaShouldThrowNoErrorWhenInputMatchesSchema(t *testing.T) {
	fs := afero.NewMemMapFs()
	afero.WriteFile(fs, "schema.yaml", []byte(`
name: My Template
type: object
schema:
  project_name:
    type: string
required:
  - project_name
`), 0644)

	schema, err := LoadSchema(fs)
	if err != nil {
		t.Fatalf("got unexpected error while parsing schema: %s", err)
	}

	input := map[string]interface{}{
		"project_name": "my-project",
	}

	err = ValidateInput(schema, &input)
	if err != nil {
		t.Error("Got error when validating test input", err)
	}
}

func TestValidateSchemaShouldThrowErrorWhenInputDoesntMatchSchema(t *testing.T) {
	fs := afero.NewMemMapFs()
	afero.WriteFile(fs, "schema.yaml", []byte(`
name: My Template
type: object
schema:
  project_name:
    type: string
required:
  - project_name
`), 0644)

	schema, err := LoadSchema(fs)
	if err != nil {
		t.Error("got unexpected error while parsing schema")
	}

	input := map[string]interface{}{
		"invalid_key": "No One Cares About This Value",
	}

	err = ValidateInput(schema, &input)
	if err == nil {
		t.Error("Did not get an expected error when passing bad user input")
	}
}
