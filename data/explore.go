package data

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

const (
	exploreRepoUrl = "https://github.com/github/explore"
)

type Topic struct {
	Aliases          []string
	DisplayName      string
	ShortDescription string
	Topic            string
}

type Collection struct {
	Items       []string `yaml:"items"`
	DisplayName string   `yaml:"display_name"`
}

func ExploreData() {
	MakeDirIfNotExists(CacheRoot())
	defer os.Remove(path.Join(CacheRoot(), "explore"))
	cloneErr := cloneExploreRepo(CacheRoot())
	if cloneErr != nil {
		log.Fatalf("Failed to clone %s repo: %s", exploreRepoUrl, cloneErr.Error())
	}

	// parse topics directory
	var topics []Topic
	topicsPath := path.Join(CacheRoot(), "explore", "topics")
	topicFiles, topicErr := ioutil.ReadDir(topicsPath)
	check(topicErr)

	for _, f := range topicFiles {
		tData, tErr := ioutil.ReadFile(path.Join(topicsPath, f.Name(), "index.md"))
		check(tErr)
		topics = append(topics, parseTopicData(string(tData)))
	}

	tJson, tErr := json.Marshal(topics)
	check(tErr)
	tWriteErr := ioutil.WriteFile(path.Join(DataPath(), "topics.json"), tJson, 0775)
	check(tWriteErr)

	// parse collections directory
	var collections []Collection
	collPath := path.Join(CacheRoot(), "explore", "collections")
	collFiles, collErr := ioutil.ReadDir(collPath)
	check(collErr)

	for _, f := range collFiles {
		cData, tErr := ioutil.ReadFile(path.Join(collPath, f.Name(), "index.md"))
		check(tErr)
		collections = append(collections, parseCollectionData(string(cData)))
	}

	cJson, cErr := json.Marshal(collections)
	check(cErr)
	cWriteErr := ioutil.WriteFile(path.Join(DataPath(), "collections.json"), cJson, 0775)
	check(cWriteErr)
}

func cloneExploreRepo(targetDir string) error {
	if _, err := os.Stat(path.Join(targetDir, "explore")); !os.IsNotExist(err) {
		// repository is already cloned
		return nil
	}
	cmd := exec.Command("git", "clone", exploreRepoUrl)
	cmd.Dir = targetDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func parseTopicData(s string) Topic {
	// intermediate struct form required
	// aliases field does not match yaml array syntax
	type TopicYaml struct {
		Aliases          string `yaml:"aliases"`
		DisplayName      string `yaml:"display_name"`
		ShortDescription string `yaml:"short_description"`
		Topic            string `yaml:"topic"`
	}
	t := TopicYaml{}
	fm := extractFrontMatter(s)
	err := yaml.Unmarshal([]byte(fm), &t)
	check(err)
	// convert intermediate form to plain Topic struct
	return Topic{
		Aliases:          strings.Split(t.Aliases, ", "),
		DisplayName:      t.DisplayName,
		ShortDescription: t.ShortDescription,
		Topic:            t.Topic,
	}
}

func parseCollectionData(s string) Collection {
	fm := extractFrontMatter(s)
	c := Collection{}
	err := yaml.Unmarshal([]byte(fm), &c)
	check(err)
	return c
}

func extractFrontMatter(s string) string {
	i0 := strings.IndexAny(s, "---") + 4
	i1 := strings.LastIndexAny(s, "---") - 3
	return s[i0:i1]
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
