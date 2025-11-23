# Implementation Plan - Antigravity MCP Configuration

## Goal
Configure Antigravity to use the newly deployed remote MCP servers on Railway:
1.  **Notion MCP**
2.  **Web Scraper MCP (Jina)**

## Proposed Changes
### Configuration
#### [NEW] [mcp_config.json](file:///c:/Users/Qntm/.gemini/antigravity/mcp_config.json)
- Create a new JSON configuration file.
- Define `mcpServers` object.
- Add `remote-notion` pointing to the Notion service.
- Add `remote-scraper` pointing to the Scraper service.
- Use placeholders for `serverUrl` (e.g., `https://<your-notion-service-url>/sse`) as the exact Railway URLs are not known.

## Verification Plan
### Manual Verification
1.  User needs to replace the placeholder URLs with their actual Railway service URLs.
2.  User restarts Antigravity/VS Code.
3.  User verifies that Notion and Scraper tools are available in the agent interface.
