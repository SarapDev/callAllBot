package internal

type Updates struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	Id      int64   `json:"update_id"`
	Message Message `json:"message"`
}

type Message struct {
	Id       int64    `json:"id"`
	Chat     Chat     `json:"chat"`
	Date     int64    `json:"date"`
	Entities []Entity `json:"entities"`
	From     User     `json:"from"`
	Text     string   `json:"text"`
}

type Chat struct {
	Id        int64  `json:"id"`
	Title     string `json:"title"`
	GroupType string `json:"type"`
}

type Entity struct {
	Length     int    `json:"length"`
	Offset     int    `json:"offset"`
	EntityType string `json:"type"`
}

type User struct {
	Id        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	IsBot     bool   `json:"is_bot"`
}

type AllAdminResponse struct {
	Ok     bool    `json:"ok"`
	Result []Admin `json:"result"`
}

type Admin struct {
	CanBeEdited         bool   `json:"can_be_edited"`
	CanChangeInfo       bool   `json:"can_change_info"`
	CanDeleteMessages   bool   `json:"can_delete_messages"`
	CanInviteUsers      bool   `json:"can_invite_users"`
	CanManageChat       bool   `json:"can_manage_chat"`
	CanManageVideoChats bool   `json:"can_manage_video_chats"`
	CanManageVoiceChats bool   `json:"can_manage_voice_chats"`
	CanPinMessages      bool   `json:"can_pin_messages"`
	CanPromoteMembers   bool   `json:"can_promote_members"`
	CanRestrictMembers  bool   `json:"can_restrict_members"`
	CustomTitle         string `json:"custom_title"`
	IsAnonymous         bool   `json:"is_anonymous"`
	Status              string `json:"status"`
	User                User   `json:"user"`
}
