package slackevents


type SlackEvent struct {
	Token      string `json:"token"`
	TeamId     string `json:"team_id"`
	Type       string `json:"type"`
	ApiAppId       string `json:"api_app_id"`
	eventId       string `json:"event_id"`
	eventTime       string `json:"event_time"`
	Event struct {
		Type       string  `json:"type"`
		User         string   `json:"user"`
		Text         string   `json:"text"`
		ClientMsgId         string   `json:"client_msg_id"`
		Ts         string   `json:"ts"`
		Channel         string   `json:"channel"`
		EventTs         string   `json:"event_ts"`
		ChannelType         string   `json:"channel_type"`

	 } `json:"event"`
	AuthedUsers []string `json:"authed_users"`
}



type Message  struct {
	User         string   `json:"user"`
	From         string   `json:"from"`
	Text         string   `json:"text"`
	Source         string   `json:"source"`
	ReplyTo struct {		
		Channel         string   `json:"address"`		
	} `json:"reply_to"`
}