package serv

import (
	"context"
	"fmt"
	"strings"

	"github.com/dosco/graphjin/serv/v3/internal/mcpcompat/mcp"
)

type guidanceToolResult struct {
	Title         string `json:"title"`
	GuideMarkdown string `json:"guide_markdown"`
}

// registerGuidanceTools is retained as an inert compatibility hook.
func (ms *mcpServer) registerGuidanceTools() {
}

func promptRequestFromToolArgs(args map[string]any) (mcp.GetPromptRequest, error) {
	out := make(map[string]string, len(args))
	for key, val := range args {
		if val == nil {
			continue
		}
		s, ok := val.(string)
		if !ok {
			return mcp.GetPromptRequest{}, fmt.Errorf("%s must be a string", key)
		}
		if strings.TrimSpace(s) == "" {
			continue
		}
		out[key] = s
	}

	return mcp.GetPromptRequest{
		Params: mcp.GetPromptParams{
			Arguments: out,
		},
	}, nil
}

func promptResultText(res *mcp.GetPromptResult) string {
	if res == nil {
		return ""
	}

	var sb strings.Builder
	for _, msg := range res.Messages {
		text, ok := msg.Content.(mcp.TextContent)
		if !ok {
			continue
		}
		if sb.Len() > 0 {
			sb.WriteString("\n\n")
		}
		sb.WriteString(text.Text)
	}
	return sb.String()
}

func (ms *mcpServer) guidanceToolResult(tool string, args map[string]any, res *mcp.GetPromptResult) (*mcp.CallToolResult, error) {
	return ms.toolResultJSON(tool, args, guidanceToolResult{
		Title:         res.Description,
		GuideMarkdown: promptResultText(res),
	})
}

func (ms *mcpServer) handleWriteQueryTool(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	promptReq, err := promptRequestFromToolArgs(req.GetArguments())
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	res, err := ms.handleWriteQuery(ctx, promptReq)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return ms.guidanceToolResult("write_query", req.GetArguments(), res)
}

func (ms *mcpServer) handleWriteMutationTool(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	promptReq, err := promptRequestFromToolArgs(req.GetArguments())
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	res, err := ms.handleWriteMutation(ctx, promptReq)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	return ms.guidanceToolResult("write_mutation", req.GetArguments(), res)
}

func (ms *mcpServer) handleFixQueryErrorTool(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	promptReq, err := promptRequestFromToolArgs(req.GetArguments())
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	query := promptReq.Params.Arguments["query"]
	errorMsg := promptReq.Params.Arguments["error"]
	if query == "" {
		return mcp.NewToolResultError("query argument is required"), nil
	}
	if errorMsg == "" {
		return mcp.NewToolResultError("error argument is required"), nil
	}

	repair := buildFixQueryErrorRepair(query, errorMsg, ms.analyticsModeOn())
	return ms.toolResultJSON("fix_query_error", req.GetArguments(), repair)
}
