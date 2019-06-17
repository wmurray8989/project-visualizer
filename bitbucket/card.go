package bitbucket

// Card is a struct containing information about a bitbucket card
type Card struct {
	ID     string `json:"id"`
	Key    string `json:"key"`
	Self   string `json:"self"`
	Fields struct {
		Summary string `json:"summary"`
		Epic    string `json:"customfield_10009"`
		Status  struct {
			Name string `json:"name"`
		} `json:"status"`
		Creator struct {
			Name        string `json:"name"`
			DisplayName string `json:"displayName"`
		} `json:"creator"`
		Reporter struct {
			Name        string `json:"name"`
			DisplayName string `json:"displayName"`
		} `json:"reporter"`
		Assignee struct {
			Name        string `json:"name"`
			DisplayName string `json:"displayName"`
		} `json:"assignee"`
		IssueLinks []struct {
			ID string `json:"id"`
		} `json:"issuelinks"`
	} `json:"fields"`
}
