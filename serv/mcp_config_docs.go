package serv

import (
	_ "embed"
	"strings"
)

// devConfigTemplate / prodConfigTemplate / agenticConfigTemplate are the
// heavily-commented counterparts to the configs that `graphjin serve new`
// writes for new apps. They are the canonical reference for what a GraphJin
// config.yml looks like: every option is documented inline with a comment
// explaining what it does and what the reasonable defaults are. We expose
// them via the `get_config_docs` MCP tool so an LLM (or a human) can read them
// without having to find the template files inside the repo.
//
//go:embed config_docs/dev.yml
var devConfigTemplate string

//go:embed config_docs/prod.yml
var prodConfigTemplate string

//go:embed config_docs/agentic.yml
var agenticConfigTemplate string

// ConfigDocsTemplate returns the annotated example config template for the
// given variant ("dev", "prod", or "agentic"); any other value returns dev.
// These are the annotated counterparts to the templates `graphjin serve new`
// scaffolds, exposed for `graphjin config docs` (offline, no server).
func ConfigDocsTemplate(variant string) string {
	switch strings.ToLower(strings.TrimSpace(variant)) {
	case "prod", "production":
		return prodConfigTemplate
	case "agentic", "agent":
		return agenticConfigTemplate
	default:
		return devConfigTemplate
	}
}

// registerConfigDocsTool registers the `get_config_docs` MCP tool. Unlike
// `get_current_config`, this is *static documentation* — it's safe to expose
// in production because it never reveals anything about the running server's
// configuration (no secrets, no DB names, no schemas — just example YAML
// with explanatory comments).
// registerConfigDocsTool is retained as an inert compatibility hook.
func (ms *mcpServer) registerConfigDocsTool() {
}
