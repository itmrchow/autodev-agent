package dto

var _ PmToolReqInterface = &JiraWebhookDto{}

// JiraWebhookDto represents the main webhook payload from Jira
type JiraWebhookDto struct {
	Issue        JiraIssue     `json:"issue"`
	Changelog    JiraChangelog `json:"changelog"`
	WebhookEvent string        `json:"webhookEvent"`
}

// JiraIssue represents the issue information in the webhook
type JiraIssue struct {
	ID     string          `json:"id"`
	Self   string          `json:"self"`
	Key    string          `json:"key"`
	Fields JiraIssueFields `json:"fields"`
}

// JiraIssueFields represents the fields within an issue
type JiraIssueFields struct {
	AggregateProgress             JiraProgress     `json:"aggregateprogress"`
	AggregateTimeEstimate         *int64           `json:"aggregatetimeestimate"`
	AggregateTimeOriginalEstimate *int64           `json:"aggregatetimeoriginalestimate"`
	AggregateTimeSpent            *int64           `json:"aggregatetimespent"`
	Assignee                      *JiraUser        `json:"assignee"`
	Attachment                    []JiraAttachment `json:"attachment"`
	Components                    []JiraComponent  `json:"components"`
	Created                       string           `json:"created"`
	Creator                       JiraUser         `json:"creator"`
	Summary                       string           `json:"summary"`
	Description                   string           `json:"description"`
	Labels                        []string         `json:"labels"`
	Priority                      *JiraPriority    `json:"priority"`
}

// JiraPriority represents priority information
type JiraPriority struct {
	Self        string `json:"self"`
	IconUrl     string `json:"iconUrl"`
	Name        string `json:"name"`
	ID          string `json:"id"`
	StatusColor string `json:"statusColor"`
}

// JiraProgress represents progress information
type JiraProgress struct {
	Progress int `json:"progress"`
	Total    int `json:"total"`
}

// JiraAttachment represents attachment information
type JiraAttachment struct {
	// Define attachment fields as needed
}

// JiraComponent represents component information
type JiraComponent struct {
	// Define component fields as needed
}

// JiraUser represents user information
type JiraUser struct {
	AccountId   string         `json:"accountId"`
	AccountType string         `json:"accountType"`
	Active      bool           `json:"active"`
	AvatarUrls  JiraAvatarUrls `json:"avatarUrls"`
	DisplayName string         `json:"displayName"`
	Self        string         `json:"self"`
	TimeZone    string         `json:"timeZone"`
	// Legacy fields for backward compatibility
	Name         string `json:"name,omitempty"`
	Key          string `json:"key,omitempty"`
	EmailAddress string `json:"emailAddress,omitempty"`
}

// JiraAvatarUrls represents avatar URL information
type JiraAvatarUrls struct {
	Size16x16 string `json:"16x16"`
	Size24x24 string `json:"24x24"`
	Size32x32 string `json:"32x32"`
	Size48x48 string `json:"48x48"`
}

// JiraChangelog represents changelog information
type JiraChangelog struct {
	Items []JiraChangelogItem `json:"items"`
	ID    string              `json:"id"`
}

// JiraChangelogItem represents individual changelog items
type JiraChangelogItem struct {
	Field            string `json:"field"`
	FieldId          string `json:"fieldId"`
	FieldType        string `json:"fieldtype"`
	From             string `json:"from"`
	FromString       string `json:"fromString"`
	TmpFromAccountId string `json:"tmpFromAccountId"`
	TmpToAccountId   string `json:"tmpToAccountId"`
	To               string `json:"to"`
	ToString         string `json:"toString"`
}

func (dto *JiraWebhookDto) IsAssignUpdate() bool {
	// - webhookEvent = "jira:issue_updated"
	if dto.WebhookEvent != "jira:issue_updated" {
		return false
	}

	// - changelog.items.field = "assignee"
	for _, item := range dto.Changelog.Items {
		if item.Field == "assignee" {
			return true
		}
	}

	return false
}
