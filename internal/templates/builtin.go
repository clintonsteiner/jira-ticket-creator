package templates

// Built-in templates for common ticket types

var (
	// BugTemplate is a template for bug reports
	BugTemplate = Template{
		Name:        "bug",
		DisplayDesc: "Report a bug",
		IssueType:   "Bug",
		Priority:    "Medium",
		Summary:     "Bug: {{.title}}",
		Description: `## Description
{{.description}}

## Steps to Reproduce
1.

## Expected Behavior
{{.expected}}

## Actual Behavior
{{.actual}}

## Environment
- Version: {{.version}}
- OS: {{.os}}`,
	}

	// StoryTemplate is a template for user stories
	StoryTemplate = Template{
		Name:        "story",
		DisplayDesc: "Create a user story",
		IssueType:   "Story",
		Priority:    "Medium",
		Summary:     "Story: {{.title}}",
		Description: `## As a
{{.persona}}

## I want to
{{.action}}

## So that
{{.benefit}}

## Acceptance Criteria
- {{.criteria1}}
- {{.criteria2}}`,
	}

	// TaskTemplate is a template for tasks
	TaskTemplate = Template{
		Name:        "task",
		DisplayDesc: "Create a task",
		IssueType:   "Task",
		Priority:    "Medium",
		Summary:     "Task: {{.title}}",
		Description: `## Objective
{{.objective}}

## Deliverables
- {{.deliverable1}}
- {{.deliverable2}}

## Notes
{{.notes}}`,
	}

	// BuiltinTemplates contains all built-in templates
	BuiltinTemplates = []Template{
		BugTemplate,
		StoryTemplate,
		TaskTemplate,
	}
)

// GetBuiltinTemplate returns a built-in template by name
func GetBuiltinTemplate(name string) *Template {
	for _, t := range BuiltinTemplates {
		if t.Name == name {
			return &t
		}
	}
	return nil
}

// GetBuiltinTemplateNames returns all built-in template names
func GetBuiltinTemplateNames() []string {
	names := make([]string, len(BuiltinTemplates))
	for i, t := range BuiltinTemplates {
		names[i] = t.Name
	}
	return names
}
