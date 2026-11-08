package serv

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/dosco/graphjin/core/v3"
	"github.com/dosco/graphjin/serv/v3/internal/mcpcompat/mcp"
)

// Workflow metadata is stored as a JSON header comment in the JS file:
//
//	// @graphjin-workflow {"description":"...","tags":["a","b"],"variables":[{"name":"customer_id","type":"number","required":true}]}
//
// followed by the JS code.
const workflowMetaPrefix = "// @graphjin-workflow "

// WorkflowMeta holds discoverable metadata for a saved workflow.
type WorkflowMeta struct {
	Description string             `json:"description"`
	Tags        []string           `json:"tags,omitempty"`
	Variables   []WorkflowVariable `json:"variables,omitempty"`
	CreatedAt   string             `json:"created_at,omitempty"`
	UpdatedAt   string             `json:"updated_at,omitempty"`
}

// WorkflowVariable describes one input variable expected by a saved workflow.
type WorkflowVariable struct {
	Name        string `json:"name"`
	Type        string `json:"type,omitempty"`
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required,omitempty"`
}

// WorkflowInfo is returned by list_workflows.
type WorkflowInfo struct {
	Name             string             `json:"name"`
	Description      string             `json:"description"`
	Tags             []string           `json:"tags,omitempty"`
	Variables        []WorkflowVariable `json:"variables,omitempty"`
	CreatedAt        string             `json:"created_at,omitempty"`
	UpdatedAt        string             `json:"updated_at,omitempty"`
	WorkflowRevision string             `json:"workflow_revision,omitempty"`
}

var workflowNameRe = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9_-]{0,63}$`)
var workflowVariableNameRe = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

// registerWorkflowMgmtTools is retained as an inert compatibility hook.
func (ms *mcpServer) registerWorkflowMgmtTools() {
}

// handleListWorkflows returns all workflows with their metadata.
func (ms *mcpServer) handleListWorkflows(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ctx = ms.effectiveContext(ctx)
	args := req.GetArguments()
	snap, err := ms.service.workflowSnapshotForContext(ctx, ms.service.workflowTimeoutSeconds())
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	workflows := make([]WorkflowInfo, 0, len(snap.workflows))
	for _, wf := range snap.workflows {
		info := WorkflowInfo{
			Name:             wf.Name,
			Description:      wf.Description,
			Tags:             append([]string(nil), wf.Tags...),
			Variables:        workflowVariablesFromCatalog(wf.Variables),
			CreatedAt:        wf.CreatedAt,
			UpdatedAt:        wf.UpdatedAt,
			WorkflowRevision: snap.revision,
		}
		workflows = append(workflows, info)
	}

	result := map[string]any{
		"workflows": workflows,
		"count":     len(workflows),
	}
	return ms.toolResultJSON("list_workflows", args, result)
}

// handleSaveWorkflow saves LLM-authored JS code as a workflow file.
func (ms *mcpServer) handleSaveWorkflow(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ctx = ms.effectiveContext(ctx)
	if !ms.service.conf.MCP.AllowWorkflowUpdates {
		return mcp.NewToolResultError("workflow updates are not allowed. Enable allow_workflow_updates in MCP config."), nil
	}

	args := req.GetArguments()
	row, err := newControlPlaneGraphQL(ms.service).mutateWorkflow(ctx, core.ManagedMutationRoot{
		Operation: "insert",
		Input:     args,
	})
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	result := map[string]any{
		"saved":       true,
		"name":        row["name"],
		"path":        row["path"],
		"description": row["description"],
		"tags":        row["tags"],
		"variables":   row["variables"],
		"source_hash": row["source_hash"],
		"created_at":  row["created_at"],
		"updated_at":  row["updated_at"],
		"hint":        "Now run gj_workflow_execution(insert), or call execute_workflow when the legacy MCP tool is enabled, with name: " + fmt.Sprint(row["name"]),
	}
	return ms.toolResultJSON("save_workflow", args, result)
}

// parseWorkflowMeta extracts metadata from the first line of a workflow file.
func parseWorkflowMeta(src string) (WorkflowMeta, bool) {
	firstLine := src
	if idx := strings.IndexByte(src, '\n'); idx >= 0 {
		firstLine = src[:idx]
	}

	if !strings.HasPrefix(firstLine, workflowMetaPrefix) {
		return WorkflowMeta{}, false
	}

	jsonStr := strings.TrimPrefix(firstLine, workflowMetaPrefix)
	var meta WorkflowMeta
	if err := json.Unmarshal([]byte(jsonStr), &meta); err != nil {
		return WorkflowMeta{}, false
	}
	return meta, true
}

func parseWorkflowVariables(raw any) ([]WorkflowVariable, error) {
	if raw == nil {
		return nil, nil
	}

	items, ok := raw.([]any)
	if !ok {
		return nil, fmt.Errorf("variables must be an array of objects")
	}

	vars := make([]WorkflowVariable, 0, len(items))
	seen := make(map[string]struct{}, len(items))

	for i, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("variables[%d] must be an object", i)
		}

		v, err := parseWorkflowVariable(m)
		if err != nil {
			return nil, fmt.Errorf("variables[%d]: %w", i, err)
		}
		if _, exists := seen[v.Name]; exists {
			return nil, fmt.Errorf("variables[%d]: duplicate variable name %q", i, v.Name)
		}
		seen[v.Name] = struct{}{}
		vars = append(vars, v)
	}

	return vars, nil
}

func parseWorkflowVariable(m map[string]any) (WorkflowVariable, error) {
	var v WorkflowVariable

	name, _ := m["name"].(string)
	name = strings.TrimSpace(name)
	if name == "" {
		return v, fmt.Errorf("name is required")
	}
	if !workflowVariableNameRe.MatchString(name) {
		return v, fmt.Errorf("invalid variable name %q", name)
	}
	v.Name = name

	if kind, ok := m["type"].(string); ok {
		v.Type = strings.TrimSpace(kind)
	}
	if description, ok := m["description"].(string); ok {
		v.Description = strings.TrimSpace(description)
	}
	if required, ok := m["required"].(bool); ok {
		v.Required = required
	}

	return v, nil
}

func workflowVariablesFromCatalog(vars []core.CatalogWorkflowVariable) []WorkflowVariable {
	out := make([]WorkflowVariable, len(vars))
	for i, v := range vars {
		out[i] = WorkflowVariable{
			Name:        v.Name,
			Type:        v.Type,
			Description: v.Description,
			Required:    v.Required,
		}
	}
	return out
}
