package serv

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/dosco/graphjin/core/v3"
	"github.com/dosco/graphjin/serv/v3/internal/mcpcompat/mcp"
)

// registerFragmentTools is retained as an inert compatibility hook.
func (ms *mcpServer) registerFragmentTools() {
}

// handleListFragments returns all available fragments
func (ms *mcpServer) handleListFragments(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ctx = ms.effectiveContext(ctx)
	if err := ms.requireDB(); err != nil {
		return err, nil
	}

	args := req.GetArguments()
	namespace, _ := args["namespace"].(string)

	fragments, err := ms.service.listFragmentsForContext(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to list fragments: %v", err)), nil
	}

	// Filter by namespace if provided
	if namespace != "" {
		filtered := make([]core.FragmentInfo, 0)
		for _, f := range fragments {
			if f.Namespace == namespace {
				filtered = append(filtered, f)
			}
		}
		fragments = filtered
	}

	// core.ListFragments returns results sorted by (namespace, name).

	result := struct {
		Fragments []core.FragmentInfo `json:"fragments"`
		Count     int                 `json:"count"`
		Usage     string              `json:"usage"`
	}{
		Fragments: fragments,
		Count:     len(fragments),
		Usage:     `To use a fragment, add: #import "./fragments/<name>" at the top of your query, then use ...FragmentName in your selection set`,
	}
	return ms.toolResultJSON("list_fragments", args, result)
}

// handleGetFragment returns details of a specific fragment
func (ms *mcpServer) handleGetFragment(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ctx = ms.effectiveContext(ctx)
	if err := ms.requireDB(); err != nil {
		return err, nil
	}

	args := req.GetArguments()
	name, _ := args["name"].(string)
	namespace, _ := args["namespace"].(string)

	if name == "" {
		return mcp.NewToolResultError("fragment name is required"), nil
	}

	qualifiedName := qualifyAllowListName(namespace, name)
	details, _, err := ms.service.getFragmentForContext(ctx, qualifiedName)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get fragment: %v", err)), nil
	}

	importName := qualifiedName
	if importName == "" {
		importName = details.Name
		if details.Namespace != "" {
			importName = details.Namespace + "." + details.Name
		}
	}

	// Add usage example
	result := struct {
		*core.FragmentDetails
		ImportDirective string `json:"import_directive"`
		UsageExample    string `json:"usage_example"`
	}{
		FragmentDetails: details,
		ImportDirective: fmt.Sprintf(`#import "./fragments/%s"`, importName),
		UsageExample:    fmt.Sprintf("query { %s { ...%s } }", details.On, details.Name),
	}
	return ms.toolResultJSON("get_fragment", args, result)
}

// handleSearchFragments searches fragments by name
func (ms *mcpServer) handleSearchFragments(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ctx = ms.effectiveContext(ctx)
	if err := ms.requireDB(); err != nil {
		return err, nil
	}

	args := req.GetArguments()
	searchQuery, _ := args["query"].(string)
	limitFloat, _ := args["limit"].(float64)

	if searchQuery == "" {
		return mcp.NewToolResultError("search query is required"), nil
	}

	limit := int(limitFloat)
	if limit <= 0 {
		limit = 10
	}

	fragments, err := ms.service.listFragmentsForContext(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to list fragments: %v", err)), nil
	}

	// Simple fuzzy search - reuse fuzzyScore from mcp_search.go
	searchTerm := strings.ToLower(searchQuery)
	type scoredFragment struct {
		Fragment core.FragmentInfo
		Score    int
	}

	scored := make([]scoredFragment, 0)
	for _, f := range fragments {
		name := strings.ToLower(f.Name)
		score := fuzzyScore(searchTerm, name)
		if score > 0 {
			scored = append(scored, scoredFragment{Fragment: f, Score: score})
		}
	}

	// Sort by score (higher is better). Input is already in (namespace, name)
	// order from core.ListFragments, so SliceStable preserves that as the
	// tiebreaker.
	sort.SliceStable(scored, func(i, j int) bool {
		return scored[i].Score > scored[j].Score
	})

	// Limit results
	if len(scored) > limit {
		scored = scored[:limit]
	}

	// Extract just the fragment info
	results := make([]core.FragmentInfo, len(scored))
	for i, s := range scored {
		results[i] = s.Fragment
	}

	result := struct {
		Fragments []core.FragmentInfo `json:"fragments"`
		Count     int                 `json:"count"`
	}{
		Fragments: results,
		Count:     len(results),
	}
	return ms.toolResultJSON("search_fragments", args, result)
}
