package data

import (
	"testing"
)

func TestFetchOpenIssues(t *testing.T) {
	actual, err := FetchOpenIssues("bartolomej", "scng-api")
	if err != nil {
		t.Errorf("Fetch error: %s", err.Error())
	}
	t.Log(actual)
}

func TestFetchRepo(t *testing.T) {
	actual, err := FetchRepo("bartolomej", "scng-api")
	if err != nil {
		t.Errorf("Fetch error: %s", err.Error())
	}
	t.Log(actual)
}
