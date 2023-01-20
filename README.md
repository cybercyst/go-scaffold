# go-cookiecutter

Go-cookiecutter is a cookiecutter-like CLI & library to generate templated artifacts

## Installation

```bash
$ task install
```

## Usage

### Creating a template

You can see [same-repository](https://github.com/cybercyst/sample-template) for an example template

```
$ tree sample-template

├── schema.yaml
└── template
    └── {{ project_name }}
        └── {{ project_name }}.txt

3 directories, 2 files
```

At the root of the template we expect two things

- schema.yaml
- template

#### `schema.yaml`

```
$ cat sample-template\schema.yaml

name: Test Template
type: object
schema:
  project_name:
    type: string
required:
  - project_name
```

`schema.yaml` uses (json-schema)[https://json-schema.org/] to validate the variables that will be interpolated in the template files

#### `template`

`template` is a directory containing files that will be parsed with (pongo2)[https://github.com/flosch/pongo2].

NOTE: Both the contents and the filenames of files will be rendered with the variables provided

### Create template input

Create a YAML file to store the variables we want to have replaced in our generated artifact

```
$ cat input.yaml

project_name: My Project
```

### Generating the template

```
$ go-cookiecutter /path/to/sample-template --input-file /path/to/input.yaml --output-directory /output/path

$ tree /output/path

path
└── My Project
    └── My Project.txt

2 directories, 1 file

$ cat /output/path/My\ Project/My\ Project.txt

My Project

```

