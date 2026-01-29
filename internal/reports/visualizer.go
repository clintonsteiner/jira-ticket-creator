package reports

import (
	"fmt"
	"strings"

	"github.com/clintonsteiner/jira-ticket-creator/internal/storage"
)

// Visualizer creates visualizations of ticket dependencies
type Visualizer struct {
	repo storage.Repository
}

// NewVisualizer creates a new visualizer
func NewVisualizer(repo storage.Repository) *Visualizer {
	return &Visualizer{repo: repo}
}

// GenerateTree generates an ASCII tree visualization
func (v *Visualizer) GenerateTree() (string, error) {
	records, err := v.repo.GetAll()
	if err != nil {
		return "", err
	}

	if len(records) == 0 {
		return "No tickets found", nil
	}

	// Build dependency graph
	graph := make(map[string][]string) // key -> list of keys it blocks
	roots := make([]string, 0)

	for _, record := range records {
		graph[record.Key] = []string{}
	}

	// Find blocking relationships
	hasIncoming := make(map[string]bool)
	for _, record := range records {
		for _, blocker := range record.BlockedBy {
			graph[blocker] = append(graph[blocker], record.Key)
			hasIncoming[record.Key] = true
		}
	}

	// Find root tickets (not blocked by anything)
	for _, record := range records {
		if !hasIncoming[record.Key] {
			roots = append(roots, record.Key)
		}
	}

	var sb strings.Builder
	sb.WriteString("Ticket Dependency Tree\n")
	sb.WriteString("======================\n\n")

	// Render tree for each root
	visited := make(map[string]bool)
	for _, root := range roots {
		v.renderTreeNode(&sb, root, "", visited, graph)
	}

	// Render orphaned branches
	for _, record := range records {
		if !visited[record.Key] {
			sb.WriteString(fmt.Sprintf("\n%s (orphaned)\n", record.Key))
			visited[record.Key] = true
			for _, blocked := range graph[record.Key] {
				v.renderTreeNode(&sb, blocked, "  ", visited, graph)
			}
		}
	}

	return sb.String(), nil
}

// renderTreeNode recursively renders a tree node
func (v *Visualizer) renderTreeNode(sb *strings.Builder, key, indent string, visited map[string]bool, graph map[string][]string) {
	if visited[key] {
		sb.WriteString(fmt.Sprintf("%s├── %s (circular)\n", indent, key))
		return
	}

	visited[key] = true
	sb.WriteString(fmt.Sprintf("%s├── %s\n", indent, key))

	children := graph[key]
	for i, child := range children {
		if i == len(children)-1 {
			v.renderTreeNode(sb, child, indent+"    ", visited, graph)
		} else {
			v.renderTreeNode(sb, child, indent+"│   ", visited, graph)
		}
	}
}

// GenerateMermaid generates a Mermaid diagram
func (v *Visualizer) GenerateMermaid() (string, error) {
	records, err := v.repo.GetAll()
	if err != nil {
		return "", err
	}

	if len(records) == 0 {
		return "", fmt.Errorf("no tickets found")
	}

	var sb strings.Builder
	sb.WriteString("graph TD\n")

	// Add all tickets
	for _, record := range records {
		sb.WriteString(fmt.Sprintf("    %s[\"%s\"]\n", sanitizeID(record.Key), record.Key))
	}

	// Add relationships
	for _, record := range records {
		for _, blocker := range record.BlockedBy {
			sb.WriteString(fmt.Sprintf("    %s -->|blocks| %s\n", sanitizeID(blocker), sanitizeID(record.Key)))
		}
	}

	return sb.String(), nil
}

// GenerateDOT generates a Graphviz DOT file
func (v *Visualizer) GenerateDOT() (string, error) {
	records, err := v.repo.GetAll()
	if err != nil {
		return "", err
	}

	if len(records) == 0 {
		return "", fmt.Errorf("no tickets found")
	}

	var sb strings.Builder
	sb.WriteString("digraph TicketDependencies {\n")
	sb.WriteString("    rankdir=LR;\n")
	sb.WriteString("    node [shape=box, style=rounded];\n\n")

	// Add all tickets
	for _, record := range records {
		sb.WriteString(fmt.Sprintf("    \"%s\" [label=\"%s\"];\n", record.Key, record.Key))
	}

	sb.WriteString("\n")

	// Add relationships
	for _, record := range records {
		for _, blocker := range record.BlockedBy {
			sb.WriteString(fmt.Sprintf("    \"%s\" -> \"%s\" [label=\"blocks\"];\n", blocker, record.Key))
		}
	}

	sb.WriteString("}\n")

	return sb.String(), nil
}

// DetectCircularDependencies finds circular dependencies
func (v *Visualizer) DetectCircularDependencies() ([][]string, error) {
	records, err := v.repo.GetAll()
	if err != nil {
		return nil, err
	}

	// Build graph
	graph := make(map[string][]string)
	for _, record := range records {
		graph[record.Key] = record.BlockedBy
	}

	var cycles [][]string
	visited := make(map[string]bool)
	path := make(map[string]bool)

	var dfs func(string, []string)
	dfs = func(node string, currentPath []string) {
		visited[node] = true
		path[node] = true
		currentPath = append(currentPath, node)

		for _, neighbor := range graph[node] {
			if path[neighbor] {
				// Found cycle
				cycleStart := -1
				for i, n := range currentPath {
					if n == neighbor {
						cycleStart = i
						break
					}
				}
				if cycleStart >= 0 {
					cycle := append([]string{}, currentPath[cycleStart:]...)
					cycle = append(cycle, neighbor)
					cycles = append(cycles, cycle)
				}
			} else if !visited[neighbor] {
				dfs(neighbor, currentPath)
			}
		}

		path[node] = false
	}

	for _, record := range records {
		if !visited[record.Key] {
			dfs(record.Key, []string{})
		}
	}

	return cycles, nil
}

// sanitizeID sanitizes an ID for use in Mermaid diagrams
func sanitizeID(id string) string {
	return strings.ReplaceAll(id, "-", "")
}
