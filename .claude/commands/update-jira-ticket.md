Use atlassian mcp tool to update {ticket} jira ticket {status} & {message}

# Args
- {ticket}: {ticket id} or {ticket url} , Required for new conversations.
- {status}: update to status , Required
- {message}: add or update description , Required for new conversations.

# exec flow
- if not exist conversation , Read {ticket} and give me a JSON content.
- update ticket status to {status}
- determine message is update or add in description
