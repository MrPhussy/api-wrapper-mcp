# MCP Server Factory Deployment

- [x] Scaffolding & Config Strategy
    - [x] Clone generic wrapper
    - [x] Implement config loading
    - [x] Dockerfile & Railway setup
- [ ] Execution & Ingestion
    - [x] Create tool configs (Notion, Stripe, Jina)
    - [x] Fix healthcheck (PORT)
    - [x] Refactor for SSE support (mark3labs/mcp-go)
    - [x] Fix SSE implementation (Custom Bridge) <!-- id: 0 -->
    - [x] Phase 3: Custom Server Implementation
        - [x] Rewrite main.go (Manual MCP Protocol)
        - [x] Implement initialize & tools handlers
        - [x] Verify deployment
