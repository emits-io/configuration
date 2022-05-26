package configuration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/emits-io/core"
)

const (
	// ConfigFile constant for configuration file name
	ConfigFile = "emits.json"
)

// Configuration contains all options used to establish processing of ConfigFile
type Configuration struct {
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	Author      string    `json:"author,omitempty"`
	License     string    `json:"license,omitempty"`
	Version     string    `json:"version,omitempty"`
	Task        []*Task   `json:"task,omitempty"`
	Script      []*Script `json:"script,omitempty"`
	File        []*File   `json:"file,omitempty"`
}

// Script contains all the options used to establish a script on Configuration
type Script struct {
	Name string   `json:"name,omitempty"`
	Task []string `json:"task,omitempty"`
}

// Task contains all the options used to establish a task on Configuration
type Task struct {
	Name string `json:"name,omitempty"`
	Path *Path  `json:"path,omitempty"`
}

// Path contains all the options used to establish a path on Task
type Path struct {
	Include []string `json:"include,omitempty"`
	Exclude []string `json:"exclude,omitempty"`
}

// File contains all the options used to establish a file on Configuration
type File struct {
	Type   []string `json:"type,omitempty"`
	Parse  *Parse   `json:"parse,omitempty"`
	Modify *Modify  `json:"modify,omitempty"`
}

// Modify contains all the options used to establish a modify on File
type Modify struct {
	Plugin []*Plugin                 `json:"plugin,omitempty"`
	Regex  []*core.RegularExpression `json:"regex,omitempty"`
}

// Parse contains all the options used to establish a parse on File
type Parse struct {
	Comment *core.Comment `json:"comment,omitempty"`
	Source  bool          `json:"source,omitempty"`
}

// Plugin contains all the options used to establish a plugin on File
type Plugin struct {
	Path string `json:"path,omitempty"`
}

func (c *Configuration) Write() error {
	data, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		return err
	}
	err = os.WriteFile(ConfigFile, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

// Load attempts to open ConfigFile and returns any errors from Validate()
func (c *Configuration) Load() error {
	jsonFile, err := os.Open(ConfigFile)
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
	err := c.ValidateTaskDefinitionExists()
	if err != nil {
		errors = append(errors, err)
	}
	err = c.ValidateFileDefinitionExists()
	if err != nil {
		errors = append(errors, err)
	}
	for _, task := range c.Task {
		errTaskDefinition := task.Validate()
		if errTaskDefinition != nil {
			errors = append(errors, errTaskDefinition...)
		}
	}
	for _, file := range c.File {
		errFileDefinition := file.Validate()
		if errFileDefinition != nil {
			errors = append(errors, errFileDefinition...)
		}
	}
	for _, script := range c.Script {
		errScriptDefinition := script.Validate(c)
		if errScriptDefinition != nil {
			errors = append(errors, errScriptDefinition...)
		}
	}
	return errors
}

func (c *Configuration) ValidateTaskDefinitionExists() error {
	if len(c.Task) == 0 {
		return fmt.Errorf("`%s` must contain at least one task definition", ConfigFile)
	}
	return nil
}

func (c *Configuration) ValidateFileDefinitionExists() error {
	if len(c.File) == 0 {
		return fmt.Errorf("`%s` must contain at least one file definition", ConfigFile)
	}
	return nil
}

func (f *File) Validate() []error {
	var errors []error
	if len(f.Type) == 0 {
		f.Type = []string{fmt.Sprintf("%v", &f)}
		errors = append(errors, fmt.Errorf("`%s` file missing type definition", strings.Join(f.Type, ",")))
	}
	errParseDefinition := f.Parse.Validate(f)
	if errParseDefinition != nil {
		errors = append(errors, errParseDefinition...)
	}
	if f.Modify != nil {
		if f.Modify.Plugin != nil {
			for i, plugin := range f.Modify.Plugin {
				if len(plugin.Path) == 0 {
					errors = append(errors, fmt.Errorf("`%s` file modify plugin path definition at index `%v` is empty", strings.Join(f.Type, ","), i))
				}
			}
		}
		if f.Modify.Regex != nil {
			for i, regex := range f.Modify.Regex {
				if len(regex.Find) == 0 {
					errors = append(errors, fmt.Errorf("`%s` file modify find definition at index `%v` is empty", strings.Join(f.Type, ","), i))
				}
			}
		}
	}
	return errors
}

func (p *Parse) Validate(f *File) []error {
	var errors []error
	if p == nil {
		errors = append(errors, fmt.Errorf("file `%s` type missing parse definition", strings.Join(f.Type, ",")))
	} else {
		if p.Comment == nil || p.Comment != nil && len(p.Comment.Line) == 0 && p.Comment.Block == nil {
			errors = append(errors, fmt.Errorf("file `%s` type missing parse comment definition", strings.Join(f.Type, ",")))
		} else if p.Comment.Block != nil {
			if len(p.Comment.Block.Start) == 0 {
				errors = append(errors, fmt.Errorf("file `%s` type missing parse block comment start definition", strings.Join(f.Type, ",")))
			}
			if len(p.Comment.Block.End) == 0 {
				errors = append(errors, fmt.Errorf("file `%s` type missing parse block comment end definition", strings.Join(f.Type, ",")))
			}
		}
	}
	return errors
}

func (t *Task) Validate() []error {
	var errors []error
	if len(t.Name) == 0 {
		t.Name = fmt.Sprintf("%v", &t)
		errors = append(errors, fmt.Errorf("`%s` task missing name definition", t.Name))
	}
	if t.Path != nil {
		if t.Path.Include == nil {
			errors = append(errors, fmt.Errorf("`%s` task missing path include definition", t.Name))
		}
		for i, include := range t.Path.Include {
			if len(strings.TrimSpace(include)) == 0 {
				errors = append(errors, fmt.Errorf("`%s` task path include definition at index `%v` is empty", t.Name, i))
			}
		}
		for i, exclude := range t.Path.Exclude {
			if len(strings.TrimSpace(exclude)) == 0 {
				errors = append(errors, fmt.Errorf("`%s` task path exclude definition at index `%v` is empty", t.Name, i))
			}
		}
	} else {
		errors = append(errors, fmt.Errorf("`%s` task missing path definition", t.Name))
	}
	return errors
}

func (s *Script) Validate(c *Configuration) []error {
	var errors []error
	if len(s.Name) == 0 {
		s.Name = fmt.Sprintf("%v", &s)
		errors = append(errors, fmt.Errorf("`%s` script missing name definition", s.Name))
	}
	if len(s.Task) == 0 {
		errors = append(errors, fmt.Errorf("`%s` script must contain at least one task definition", s.Name))
	} else {
		var seenTask []string
		for _, task := range s.Task {
			taskSeen := false
			for _, seen := range seenTask {
				if seen == task {
					taskSeen = true
					break
				}
			}
			if taskSeen {
				errors = append(errors, fmt.Errorf("`%s` script referencing duplicate `%s` task definition", s.Name, task))
			} else {
				seenTask = append(seenTask, task)
			}
			if c.FindTask(task) == nil {
				errors = append(errors, fmt.Errorf("`%s` script referencing unknown `%s` task definition", s.Name, task))
			}
		}
	}
	return errors
}

// FindTask returns the Task if found or nil if not found; used to validate Script Task references
func (c *Configuration) FindTask(name string) *Task {
	for _, t := range c.Task {
		if t.Name == name {
			return t
		}
	}
	return nil
}

// FindScript returns the Script if found or nil if not found; used to validate Script references
func (c *Configuration) FindScript(name string) *Script {
	for _, s := range c.Script {
		if s.Name == name {
			return s
		}
	}
	return nil
}
