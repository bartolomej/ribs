package data

import (
	"os"
	"path"
	"reflect"
	"testing"
)

func TestExtractFrontMatter(t *testing.T) {
	data := `---
aliases: phpfusion, php-fusion-cms, php-fusion-themes, php-fusion-infusions
created_by: PHP-Fusion Inc
---
PHP-Fusion is an all in one integrated and scalable platform a lightweight open source content management system (CMS) written in PHP that will fit any purpose when it comes to website productions, whether you are creating community portals or personal sites.
`

	fm := extractFrontMatter(data)

	if fm[0] == '-' {
		t.Errorf("First char is '-': %s", fm)
	}
	if fm[len(fm)-1] == '-' {
		t.Errorf("Last char is '-': %s", fm)
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
3D modeling uses specialized software to create a digital model of a physical object.`

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

func TestCloneExploreRepo(t *testing.T) {
	MakeDirIfNotExists(CacheRoot())

	if !fileExists(CacheRoot()) {
		t.Errorf("Cache directory %s not created", CacheRoot())
	}

	err := cloneExploreRepo(CacheRoot())
	if err != nil {
		t.Errorf("Command error is not nil: %s", err.Error())
	}

	if !fileExists(path.Join(CacheRoot(), "explore")) {
		t.Errorf("Repository directory not downloaded")
	}
}

func TestExtractExploreData(t *testing.T) {
	ExtractExploreData()

	if !fileExists(path.Join(DataPath(), "topics.json")) {
		t.Errorf("'topics.json' file not created under %s", DataPath())
	}
	if !fileExists(path.Join(DataPath(), "collections.json")) {
		t.Errorf("'collections.json' file not created under %s", DataPath())
	}
	if fileExists(path.Join(CacheRoot(), "explore")) {
		t.Errorf("'explore' directory not removed under %s", CacheRoot())
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
