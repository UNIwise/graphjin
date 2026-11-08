package serv

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/dosco/graphjin/core/v3"
	"github.com/dosco/graphjin/serv/v3/internal/mcpcompat/mcp"
)

// registerQueryDiscoveryTools is retained as an inert compatibility hook.
func (ms *mcpServer) registerQueryDiscoveryTools() {
}

// handleListSavedQueries returns all saved queries
func (ms *mcpServer) handleListSavedQueries(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ctx = ms.effectiveContext(ctx)
	if err := ms.requireDB(); err != nil {
		return err, nil
	}

	args := req.GetArguments()
	namespace, _ := args["namespace"].(string)

	queries, err := ms.service.listSavedQueriesForContext(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to list queries: %v", err)), nil
	}

	// Filter by namespace if provided
	if namespace != "" {
		filtered := make([]core.SavedQueryInfo, 0)
		for _, q := range queries {
			if q.Namespace == namespace {
				filtered = append(filtered, q)
			}
		}
		queries = filtered
	}

	// core.ListSavedQueries returns results sorted by (namespace, name).

	result := struct {
		Queries []core.SavedQueryInfo `json:"queries"`
		Count   int                   `json:"count"`
	}{
		Queries: queries,
		Count:   len(queries),
	}
	return ms.toolResultJSON("list_saved_queries", args, result)
}

// handleSearchSavedQueries searches queries by name
func (ms *mcpServer) handleSearchSavedQueries(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
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

	queries, err := ms.service.listSavedQueriesForContext(ctx)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to list queries: %v", err)), nil
	}

	// Simple fuzzy search - match if search term is contained in name
	searchTerm := strings.ToLower(searchQuery)
	type scoredQuery struct {
		Query core.SavedQueryInfo
		Score int
	}

	scored := make([]scoredQuery, 0)
	for _, q := range queries {
		name := strings.ToLower(q.Name)
		score := fuzzyScore(searchTerm, name)
		if score > 0 {
			scored = append(scored, scoredQuery{Query: q, Score: score})
		}
	}

	// Sort by score (higher is better). Input is already in (namespace, name)
	// order from core.ListSavedQueries, so SliceStable preserves that as the
	// tiebreaker.
	sort.SliceStable(scored, func(i, j int) bool {
		return scored[i].Score > scored[j].Score
	})

	// Limit results
	if len(scored) > limit {
		scored = scored[:limit]
	}

	// Extract just the query info
	results := make([]core.SavedQueryInfo, len(scored))
	for i, s := range scored {
		results[i] = s.Query
	}

	result := struct {
		Queries []core.SavedQueryInfo `json:"queries"`
		Count   int                   `json:"count"`
	}{
		Queries: results,
		Count:   len(results),
	}
	return ms.toolResultJSON("search_saved_queries", args, result)
}

// handleGetSavedQuery returns details of a specific saved query
func (ms *mcpServer) handleGetSavedQuery(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	ctx = ms.effectiveContext(ctx)
	if err := ms.requireDB(); err != nil {
		return err, nil
	}

	args := req.GetArguments()
	name, _ := args["name"].(string)
	namespace, _ := args["namespace"].(string)

	if name == "" {
		return mcp.NewToolResultError("query name is required"), nil
	}

	details, _, err := ms.service.getSavedQueryForContext(ctx, qualifyAllowListName(namespace, name))
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to get query: %v", err)), nil
	}
	return ms.toolResultJSON("get_saved_query", args, details)
}

// fuzzyScore returns a score for how well the search term matches the target
// Higher score = better match
func fuzzyScore(search, target string) int {
	// Exact match
	if search == target {
		return 100
	}

	// Starts with
	if strings.HasPrefix(target, search) {
		return 90
	}

	// Contains
	if strings.Contains(target, search) {
		return 70
	}

	// Word boundary match
	words := strings.FieldsFunc(target, func(r rune) bool {
		return r == '_' || r == '-' || r == '.'
	})
	for _, word := range words {
		if strings.HasPrefix(word, search) {
			return 60
		}
	}

	// Character-by-character fuzzy match
	searchIdx := 0
	matches := 0
	for i := 0; i < len(target) && searchIdx < len(search); i++ {
		if target[i] == search[searchIdx] {
			matches++
			searchIdx++
		}
	}

	if searchIdx == len(search) {
		// All characters found in order
		return 50 * matches / len(target)
	}

	return 0
}

func qualifyAllowListName(namespace, name string) string {
	namespace = strings.TrimSpace(namespace)
	name = strings.TrimSpace(name)

	if name == "" || namespace == "" || strings.Contains(name, ".") {
		return name
	}

	return namespace + "." + name
}
