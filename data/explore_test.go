package data

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"testing"
)

func TestCloneExploreRepo(t *testing.T) {
	MakeDirIfNotExists(CacheRoot())

	if !fileExists(CacheRoot()) {
		t.Error(fmt.Sprintf("Cache directory %s not created", CacheRoot()))
	}

	err := cloneExploreRepo(CacheRoot())
	if err != nil {
		t.Error(fmt.Sprintf("Command error is not nil: %s", err.Error()))
	}

	if !fileExists(path.Join(CacheRoot(), "explore")) {
		t.Error(fmt.Sprintf("Repository directory not downloaded"))
	}
}

func TestParseTopicsDir(t *testing.T) {
	data := `---
aliases: 3d-printing, 3d-graphics, 3d-models
display_name: 3D
short_description: 3D modeling is the process of virtually developing the surface
  and structure of a 3D object.
topic: 3d
wikipedia_url: https://en.wikipedia.org/wiki/3D_modeling
---
3D modeling uses specialized software to create a digital model of a physical object. It is an aspect of 3D computer graphics, used for video games, 3D printing, and VR, among other applications.`

	actual := parseTopicData(data)

	expected := Topic{
		Aliases:          []string{"3d-printing", "3d-graphics", "3d-models"},
		DisplayName:      "3D",
		ShortDescription: "3D modeling is the process of virtually developing the surface and structure of a 3D object.",
		Topic:            "3d",
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected topics struct not equal to actual struct: %s", actual)
	}
}

func TestParseCollectionDir(t *testing.T) {
	data := `---
items:
 - tensorflow/models
 - Theano/Theano
display_name: Model Zoos of machine and deep learning technologies
created_by: alanbraz
---
Model Zoo is a common way that open source frameworks and companies organize their machine learning and deep learning models.
`

	actual := parseCollectionData(data)

	expected := Collection{
		Items:       []string{"tensorflow/models", "Theano/Theano"},
		DisplayName: "Model Zoos of machine and deep learning technologies",
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected collection struct not equal to actual struct: %s", actual)
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
