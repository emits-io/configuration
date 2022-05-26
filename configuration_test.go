package configuration_test

import (
	"github.com/emits-io/configuration"
	"github.com/emits-io/core"
	"testing"
)

func TestConfiguration_Write(t *testing.T) {
	c := &configuration.Configuration{
		Name:        "Name",
		Description: "Description",
		Author:      "Daniel Esquivias",
		License:     "MIT",
		Version:     "1.0.0",
		Task: []*configuration.Task{
			{
				Name: "lorem",
				Path: &configuration.Path{
					Include: []string{"*"},
				},
			},
			{
				Name: "ipsum",
				Path: &configuration.Path{
					Include: []string{"*"},
				},
			},
		},
		Script: []*configuration.Script{
			{
				Name: "foo",
				Task: []string{"lorem", "ipsum"},
			},
		},
		File: []*configuration.File{
			{
				Type: []string{"go"},
				Parse: &configuration.Parse{
					Comment: &core.Comment{
						Line: "//",
						Block: &core.CommentBlock{
							Start: "/*",
							End:   "*/",
						},
					},
					Source: true,
				},
				Modify: &configuration.Modify{
					Plugin: []*configuration.Plugin{
						{
							"./foo.js",
						},
						{
							"./bar.js",
						},
					},
					Regex: []*core.RegularExpression{
						{
							Find:    "foo",
							Replace: "bar",
						},
					},
				},
			},
		},
	}
	err := c.Write()
	if err != nil {
		t.Errorf("Expecting nil, got %v", err)
	}
}

func TestConfiguration_Load(t *testing.T) {
	TestConfiguration_Write(t)
	c := &configuration.Configuration{}
	err := c.Load()
	if err != nil {
		t.Errorf("Expecting nil, got %v", err)
	}
}

func TestConfiguration_Validate(t *testing.T) {
	c := &configuration.Configuration{}
	err := c.Validate()
	if err == nil {
		t.Errorf("Expecting error, nil")
	}
	c.Task = []*configuration.Task{
		{
			Name: "",
		},
	}
	c.Script = []*configuration.Script{
		{
			Name: "",
		},
	}
	c.File = []*configuration.File{
		{
			Type: []string{""},
		},
	}
	err = c.Validate()
	if err == nil {
		t.Errorf("Expecting error, nil")
	}
}

func TestConfiguration_ValidateTaskDefinitionExists(t *testing.T) {
	c := &configuration.Configuration{}
	err := c.ValidateTaskDefinitionExists()
	if err == nil {
		t.Errorf("Expecting error, got nil")
	}
	c.Task = []*configuration.Task{
		{
			Name: "test",
		},
	}
	err = c.ValidateTaskDefinitionExists()
	if err != nil {
		t.Errorf("Expecting nil, got %v", err)
	}
}

func TestConfiguration_ValidateFileDefinitionExists(t *testing.T) {
	c := configuration.Configuration{}
	err := c.ValidateFileDefinitionExists()
	if err == nil {
		t.Errorf("Expecting error, got nil")
	}
	c.File = []*configuration.File{
		{
			Type: []string{"test"},
		},
	}
	err = c.ValidateFileDefinitionExists()
	if err != nil {
		t.Errorf("Expecting nil, got %v", err)
	}
}

func TestFile_Validate(t *testing.T) {
	f := &configuration.File{}
	err := f.Validate()
	if err == nil {
		t.Errorf("Expecting error, got nil")
	}
	f.Modify = &configuration.Modify{
		Plugin: []*configuration.Plugin{
			{
				Path: "",
			},
		},
		Regex: []*core.RegularExpression{
			{
				Find: "",
			},
		},
	}
	err = f.Validate()
	if err == nil {
		t.Errorf("Expecting error, got nil")
	}
}

func TestParse_Validate(t *testing.T) {
	f := &configuration.File{
		Type: []string{"test"},
	}
	p := &configuration.Parse{}
	err := p.Validate(f)
	if err == nil {
		t.Errorf("Expecting error, got nil")
	}
	p.Comment = &core.Comment{
		Line: "//",
	}
	err = p.Validate(f)
	if err != nil {
		t.Errorf("Expecting nil, got %v", err)
	}
	p.Comment = &core.Comment{
		Block: &core.CommentBlock{
			Start: "/*",
			End:   "*/",
		},
	}
	err = p.Validate(f)
	if err != nil {
		t.Errorf("Expecting nil, got %v", err)
	}
	p.Comment = &core.Comment{
		Block: &core.CommentBlock{
			Start: "/*",
		},
	}
	err = p.Validate(f)
	if err == nil {
		t.Errorf("Expecting error, got nil")
	}
	p.Comment = &core.Comment{
		Block: &core.CommentBlock{
			End: "*/",
		},
	}
	err = p.Validate(f)
	if err == nil {
		t.Errorf("Expecting error, got nil")
	}
}

func TestTask_Validate(t *testing.T) {
	task := &configuration.Task{}
	err := task.Validate()
	if err == nil {
		t.Errorf("Expecting error, got nil")
	}
	task.Path = &configuration.Path{
		Include: []string{""},
	}
	err = task.Validate()
	if err == nil {
		t.Errorf("Expecting error, got nil")
	}
	task.Path = &configuration.Path{
		Include: []string{"test"},
	}
	err = task.Validate()
	if err != nil {
		t.Errorf("Expecting nil, got %v", err)
	}
	task.Path = &configuration.Path{
		Exclude: []string{""},
	}
	err = task.Validate()
	if err == nil {
		t.Errorf("Expecting error, got nil")
	}
}

func TestScript_Validate(t *testing.T) {
	c := &configuration.Configuration{
		Task: []*configuration.Task{
			{
				Name: "test",
			},
		},
	}
	s := &configuration.Script{}
	err := s.Validate(c)
	if err == nil {
		t.Errorf("Expecting error, got nil")
	}
	s.Task = []string{"no", "no"}
	err = s.Validate(c)
	if err == nil {
		t.Errorf("Expecting error, got nil")
	}
}

func TestConfiguration_FindTask(t *testing.T) {
	c := &configuration.Configuration{
		Task: []*configuration.Task{
			{
				Name: "test",
			},
		},
	}
	task := c.FindTask("test")
	if task == nil {
		t.Errorf("Expecting task, got nil")
	}
	task = c.FindTask("foo")
	if task != nil {
		t.Errorf("Expecting nil, got task %v", task)
	}
}

func TestConfiguration_FindScript(t *testing.T) {
	c := &configuration.Configuration{
		Script: []*configuration.Script{
			{
				Name: "test",
			},
		},
	}
	script := c.FindScript("test")
	if script == nil {
		t.Errorf("Expecting script, got nil")
	}
	script = c.FindScript("foo")
	if script != nil {
		t.Errorf("Expecting nil, got script %v", script)
	}
}
