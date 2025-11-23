# AnythingLLM MCP Server Setup

This guide explains how to use the AnythingLLM API wrapper as an MCP server to integrate your self-hosted AnythingLLM knowledge base with Claude and other MCP clients.

## Prerequisites

- Self-hosted AnythingLLM instance (e.g., on Railway)
- AnythingLLM API key
- Base URL of your AnythingLLM instance

## Configuration

The `anythingllm.yaml` configuration file provides the following tools:

### Knowledge Base Tools
1. **anythingllm-query-workspace** - Query workspace using RAG (Retrieval-Augmented Generation)
2. **anythingllm-chat-workspace** - Chat using general LLM knowledge
3. **anythingllm-get-chat-history** - Retrieve conversation history
4. **anythingllm-list-workspaces** - List all available workspaces
5. **anythingllm-get-workspace** - Get workspace details

### Thread Management Tools
6. **anythingllm-create-thread** - Create new conversation thread
7. **anythingllm-query-thread** - Query specific thread using RAG
8. **anythingllm-get-thread-history** - Get thread conversation history

### Document Management
9. **anythingllm-list-documents** - List workspace documents

## Environment Variables

You need to set the following environment variables:

```bash
export ANYTHINGLLM_API_KEY="your-api-key-here"
export ANYTHINGLLM_BASE_URL="https://your-instance.railway.app"
```

**Note**: The `ANYTHINGLLM_BASE_URL` should NOT include a trailing slash and should NOT include `/api/v1` (this is added automatically in the endpoint configurations).

## Getting Your API Key

1. Log in to your AnythingLLM instance
2. Navigate to Settings → API Keys
3. Generate a new API key
4. Copy the key and save it securely

## Usage Examples

### Running the Server

```bash
# Build the server
go build -o api_wrapper

# Run with AnythingLLM config
ANYTHINGLLM_API_KEY="your-key" \
ANYTHINGLLM_BASE_URL="https://your-instance.railway.app" \
./api_wrapper config/anythingllm.yaml
```

### Claude Desktop Integration

Add to your `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "anythingllm": {
      "command": "/path/to/api_wrapper",
      "args": ["/path/to/config/anythingllm.yaml"],
      "env": {
        "ANYTHINGLLM_API_KEY": "your-api-key",
        "ANYTHINGLLM_BASE_URL": "https://your-instance.railway.app"
      }
    }
  }
}
```

### Docker Deployment

```bash
docker run -p 8080:8080 \
  -e ANYTHINGLLM_API_KEY="your-key" \
  -e ANYTHINGLLM_BASE_URL="https://your-instance.railway.app" \
  -v $(pwd)/config/anythingllm.yaml:/config/anythingllm.yaml \
  your-image-name /config/anythingllm.yaml
```

## Finding Your Workspace Slug

The workspace slug is part of your AnythingLLM workspace URL:
```
https://your-instance.railway.app/workspace/my-workspace-slug
                                              ^^^^^^^^^^^^^^^^^^
```

Or use the `anythingllm-list-workspaces` tool to get all workspace slugs.

## Query vs Chat Mode

- **Query Mode** (RAG): Searches your embedded documents for accurate, context-based answers. Best for knowledge base queries.
- **Chat Mode**: Uses general LLM knowledge without strict document retrieval. Best for general conversation.

For knowledge base use cases, **always use the query mode** tools.

## Thread Management

Threads allow you to maintain separate conversation contexts:

1. Create a thread using `anythingllm-create-thread`
2. Use the thread slug for subsequent queries with `anythingllm-query-thread`
3. Retrieve thread history with `anythingllm-get-thread-history`

This is useful for:
- Multi-user environments
- Session isolation
- Organizing different conversation topics

## API Documentation

For complete API documentation, visit your AnythingLLM instance at:
```
https://your-instance.railway.app/api/docs
```

## Troubleshooting

### Authentication Errors
- Verify your API key is correct
- Check that the API key has not expired
- Ensure the `Authorization: Bearer` header is being set correctly

### Connection Errors
- Verify your `ANYTHINGLLM_BASE_URL` is correct
- Check that your AnythingLLM instance is accessible
- Ensure there's no trailing slash in the base URL

### Workspace Not Found
- Verify the workspace slug exists using `anythingllm-list-workspaces`
- Check that the workspace hasn't been deleted or renamed

## Resources

- [AnythingLLM Documentation](https://docs.anythingllm.com)
- [AnythingLLM API Documentation](https://docs.useanything.com/features/api)
- [AnythingLLM GitHub](https://github.com/Mintplex-Labs/anything-llm)
