package templates

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

// Template represents a ticket template
type Template struct {
	Name        string   `yaml:"name"`
	DisplayDesc string   `yaml:"description"` // Description of the template itself
	IssueType   string   `yaml:"issue_type"`
	Priority    string   `yaml:"priority"`
	Summary     string   `yaml:"summary"`
	Description string   `yaml:"template_description"` // Description for the issue
	Labels      []string `yaml:"labels,omitempty"`
	Components  []string `yaml:"components,omitempty"`
}

// Loader handles template loading
type Loader struct {
	templateDirs []string
}

// NewLoader creates a new template loader
func NewLoader() *Loader {
	dirs := []string{}

	// Add user template directory
	if homeDir, err := os.UserHomeDir(); err == nil {
		dirs = append(dirs, filepath.Join(homeDir, ".jira", "templates"))
	}

	return &Loader{
		templateDirs: dirs,
	}
}

// List lists all available templates
func (l *Loader) List() []Template {
	templates := make([]Template, 0)

	// Add built-in templates
	templates = append(templates, BuiltinTemplates...)

	// Add user templates
	for _, dir := range l.templateDirs {
		if entries, err := os.ReadDir(dir); err == nil {
			for _, entry := range entries {
				if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".yaml") {
					if t, err := l.Load(strings.TrimSuffix(entry.Name(), ".yaml")); err == nil {
						templates = append(templates, *t)
					}
				}
			}
		}
	}

	return templates
}

// Load loads a template by name
func (l *Loader) Load(name string) (*Template, error) {
	// Check built-in templates first
	if t := GetBuiltinTemplate(name); t != nil {
		return t, nil
	}

	// Check user templates
	for _, dir := range l.templateDirs {
		path := filepath.Join(dir, name+".yaml")
		if data, err := os.ReadFile(path); err == nil {
			var t Template
			if err := yaml.Unmarshal(data, &t); err == nil {
				return &t, nil
			}
		}
	}

	return nil, fmt.Errorf("template not found: %s", name)
}

// Render renders a template with the given variables
func (t *Template) Render(vars map[string]string) (*Template, error) {
	result := &Template{
		Name:        t.Name,
		Description: t.Description,
		IssueType:   t.IssueType,
		Priority:    t.Priority,
		Labels:      t.Labels,
		Components:  t.Components,
	}

	// Render summary
	summaryTpl, err := template.New("summary").Parse(t.Summary)
	if err != nil {
		return nil, fmt.Errorf("failed to parse summary template: %w", err)
	}

	var summaryBuf bytes.Buffer
	if err := summaryTpl.Execute(&summaryBuf, vars); err != nil {
		return nil, fmt.Errorf("failed to render summary: %w", err)
	}
	result.Summary = summaryBuf.String()

	// Render description
	descTpl, err := template.New("description").Parse(t.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to parse description template: %w", err)
	}

	var descBuf bytes.Buffer
	if err := descTpl.Execute(&descBuf, vars); err != nil {
		return nil, fmt.Errorf("failed to render description: %w", err)
	}
	result.Description = descBuf.String()

	return result, nil
}

// Save saves a template to disk
func (l *Loader) Save(t *Template) error {
	if len(l.templateDirs) == 0 {
		return fmt.Errorf("no template directories available")
	}

	// Create directory if it doesn't exist
	dir := l.templateDirs[0]
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create template directory: %w", err)
	}

	// Marshal template to YAML
	data, err := yaml.Marshal(t)
	if err != nil {
		return fmt.Errorf("failed to marshal template: %w", err)
	}

	// Write to file
	path := filepath.Join(dir, t.Name+".yaml")
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write template: %w", err)
	}

	return nil
}
