package configuration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/emits-io/core"
)

// FILE constant for configuration file name
const FILE = "emits.json"

// Configuration contains all options used to establish processing of FILE
type Configuration struct {
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Description string   `json:"description"`
	Author      string   `json:"author"`
	License     string   `json:"license"`
	Script      []Script `json:"script"`
	Task        []Task   `json:"task"`
	File        []File   `json:"file"`
}

// Script contains all the options used to establish a script on Configuration
type Script struct {
	Name string   `json:"name"`
	Task []string `json:"task"`
}

// Task contains all the options used to establish a task on Configuration
type Task struct {
	Name string `json:"name"`
	Path Path   `json:"path"`
}

// Path contains all the options used to establish a path on Task
type Path struct {
	Include []string `json:"include"`
	Exclude []string `json:"exclude"`
}

// File contains all the options used to establish a file on Configuration
type File struct {
	Type   []string `json:"type"`
	Parse  Parse    `json:"parse"`
	Modify Modify   `json:"modify"`
}

// Modify contains all the options used to establish a modify on File
type Modify struct {
	Plugin []Plugin `json:"plugin"`
	Regex  []Regex  `json:"regex"`
}

// Parse contains all the options used to establish a parse on File
type Parse struct {
	Comment core.Comment `json:"comment"`
	Source  bool         `json:"source"`
}

// Plugin contains all the options used to establish a plugin on File
type Plugin struct {
	Path string `json:"path"`
}

// Replace contains all the options used to establish a replace on File
type Regex struct {
	Find    string `json:"find"`
	Replace string `json:"replace"`
}

// Load attemps to open FILE and returns any errors from Validate()
func (c *Configuration) Load() error {
	jsonFile, err := os.Open(FILE)
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}
	if json.Unmarshal(byteValue, &c) != nil {
		return err
	}
	jsonFile.Close()
	return nil
}

// Validate returns all known validation errors at once, rather than one at a time
func (c *Configuration) Validate() []error {
	var errors []error
	scriptCount := 0
	taskCount := 0
	fileCount := 0
	if len(c.Task) == 0 {
		errors = append(errors, fmt.Errorf("`%s` must contain at least one task definition", FILE))
	}
	if len(c.File) == 0 {
		errors = append(errors, fmt.Errorf("`%s` must contain at least one file definition", FILE))
	} else {
		for _, file := range c.File {
			if len(file.Type) == 0 {
				fileCount++
			}
			if len(file.Parse.Comment.Line) == 0 {
				errors = append(errors, fmt.Errorf("file for `%s` type is missing a parse line comment definition", strings.Join(file.Type, ",")))
			}
			if len(file.Parse.Comment.Block.Start) == 0 {
				errors = append(errors, fmt.Errorf("file for `%s` type is missing a parse block comment start definition", strings.Join(file.Type, ",")))
			}
			if len(file.Parse.Comment.Block.End) == 0 {
				errors = append(errors, fmt.Errorf("file for `%s` type is missing a parse block comment end definition", strings.Join(file.Type, ",")))
			}
			for _, plugin := range file.Modify.Plugin {
				if len(plugin.Path) == 0 {
					errors = append(errors, fmt.Errorf("file for `%s` type is missing a parse modify plugin path definition", strings.Join(file.Type, ",")))
				}
			}
			for _, regex := range file.Modify.Regex {
				if len(regex.Find) == 0 {
					errors = append(errors, fmt.Errorf("file for `%s` type is missing a parse modify regex find definition", strings.Join(file.Type, ",")))
				}
			}
		}
		if fileCount > 0 {
			plural, plural_ := "s", "are"
			if fileCount == 1 {
				plural, plural_ = "", "is"
			}
			errors = append(errors, fmt.Errorf("%d file%s %s missing a type definition", fileCount, plural, plural_))
		}
	}
	for _, script := range c.Script {
		if len(script.Name) == 0 {
			scriptCount++
		}
		if len(script.Task) == 0 {
			scriptName := script.Name
			if len(script.Name) == 0 {
				scriptName = "`unknown`"
			}
			errors = append(errors, fmt.Errorf("`%s` script must contain at least one task definition", scriptName))
		} else {
			var seenTask []string
			for _, task := range script.Task {
				scriptName := script.Name
				if len(script.Name) == 0 {
					scriptName = "`unknown`"
				} else {
					taskSeen := false
					for _, seen := range seenTask {
						if seen == task {
							taskSeen = true
							break
						}
					}
					if taskSeen {
						errors = append(errors, fmt.Errorf("`%s` script referencing duplicate `%s` task definition", scriptName, task))
					} else {
						seenTask = append(seenTask, task)
					}
				}
				if c.FindTask(task) == nil {
					errors = append(errors, fmt.Errorf("`%s` script referencing unknown `%s` task definition", scriptName, task))
				}
			}
		}
	}
	if scriptCount > 0 {
		plural, plural_ := "s", "are"
		if scriptCount == 1 {
			plural, plural_ = "", "is"
		}
		errors = append(errors, fmt.Errorf("%d script%s %s missing a name definition", scriptCount, plural, plural_))
	}
	for _, task := range c.Task {
		if len(task.Name) == 0 {
			taskCount++
		}
		if len(task.Path.Include) == 0 {
			taskName := task.Name
			if len(task.Name) == 0 {
				taskName = "`unknown`"
			}
			errors = append(errors, fmt.Errorf("`%s` task must contain at least one path include definition", taskName))
		}
	}
	if taskCount > 0 {
		plural, plural_ := "s", "are"
		if taskCount == 1 {
			plural, plural_ = "", "is"
		}
		errors = append(errors, fmt.Errorf("%d task%s %s missing a name definition", taskCount, plural, plural_))
	}
	return errors
}

// FindTask returns the Task if found or nil if not found; used to validate Script Task references
func (c *Configuration) FindTask(name string) *Task {
	for _, task := range c.Task {
		if task.Name == name {
			return &task
		}
	}
	return nil
}

// FindScript returns the Script if found or nil if not found; used to validate Script references
func (c *Configuration) FindScript(name string) *Script {
	for _, script := range c.Script {
		if script.Name == name {
			return &script
		}
	}
	return nil
}
