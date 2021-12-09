package sumoapp

import (
	"fmt"
	"strings"

	"github.com/r3labs/diff"
	"gopkg.in/yaml.v2"
)

const (
	ColorRed   string = "\033[31m"
	ColorGreen string = "\033[32m"
	ColorReset string = "\033[0m"
)

func taggedDiff(tag string, a, b interface{}) (diff.Changelog, error) {
	// It turns out the Differ object stores the state of previous Diff() call
	//  so every new diff requires a new differ
	d, err := diff.NewDiffer(diff.DisableStructValues())
	if err != nil {
		return diff.Changelog{}, err
	}

	changes, err := d.Diff(a, b)
	if err != nil {
		return diff.Changelog{}, err
	}
	result := make([]diff.Change, len(changes))
	for i, c := range changes {
		c2 := c
		c2.Path = append([]string{tag}, c2.Path...)
		result[i] = c2
	}
	return result, nil
}

func displayDiff(cs *changeSet) {
	allChanges := []diff.Changelog{
		cs.ChangelogVar,
		cs.ChangelogPanel,
		cs.ChangelogSavedSearches,
		cs.ChangelogDashboard,
		cs.ChangelogFolder,
	}

	var changelogs diff.Changelog
	for _, c := range allChanges {
		changelogs = append(changelogs, c...)
	}

	fmt.Println("Found", len(changelogs), "changes")
	for _, change := range changelogs {
		fmt.Println("")
		fmt.Println("TYPE: ", change.Type)
		fmt.Println("PATH: ", change.Path)
		if change.Type == "update" {
			fmt.Println(writeDeleted(change.From))
			fmt.Println(writeCreated(change.To), ColorReset)
		} else if change.Type == "create" {
			fmt.Println(writeCreated(change.To), ColorReset)
		} else if change.Type == "delete" {
			fmt.Println(writeDeleted(change.From), ColorReset)
		}
	}
}

func writeDeleted(object interface{}) string {
	return writeYamlObject(object, ColorRed)
}

func writeCreated(object interface{}) string {
	return writeYamlObject(object, ColorGreen)
}

func writeYamlObject(object interface{}, linePrefix string) string {
	p, err := yaml.Marshal(&object)
	if err != nil {
		msg := fmt.Errorf("Unable to marshall: %w", err)
		return fmt.Sprint(msg)
	}
	pstring := strings.Trim(string(p), "\n")
	return linePrefix + strings.ReplaceAll(pstring, "\n", "\n"+linePrefix)
}
