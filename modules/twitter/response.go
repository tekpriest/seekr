package twitter

type TwitterRecentSearchResponseData struct {
	TweetAPI
	Attachments       []Attachment `json:"attachments"`
	PossiblySensitive bool         `json:"possibly_sensitive"`
	InReplyToUserID   string       `json:"in_reply_to_user_id"`
	Geo               Geo          `json:"geo"`
}

type ResponseMeta struct {
	NewestID    string
	OldestID    string
	ResultCount int
	NextToken   string
}

type TwitterRecentSearchResponse struct {
	Data     TwitterRecentSearchResponseData `json:"data"`
	Includes []Include                       `json:"includes"`
	Meta     ResponseMeta                    `json:"meta"`
}
