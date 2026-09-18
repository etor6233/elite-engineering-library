# XERJ MCP wiring

Top-level `.mcp.json` is the official `xerj init` MCP wiring for agent tooling (see https://xerj.org/llms.txt). Admitted in `$releaseRootMetadata`.

Release form uses `command: xerj` on PATH and `XERJ_URL=http://localhost:9200` so VERIFY does not see a machine-local home path. A local absolute `xerj.exe` under the user profile may be rewritten by `xerj init` for the host; do not commit Users\\ paths.
