Use atlassian mcp tool to read jira ticket {ticket}

# 參數
- ticket: {ticket id} or {ticket url}

# response format
- JSON string , no markdown ```
```
{
    "id": "id",
    "key":"key",
    "title":"title",
    "description":"description"
}
```

# exec flow
- Read {ticket} and give me a JSON content.
- Break down the steps for this ticket.