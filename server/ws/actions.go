package ws

// Subscribed Messages
const (
	JoinUserAction        = "joinUser"
	JoinGuildAction       = "joinGuild"
	JoinChannelAction     = "joinChannel"
	JoinVoiceAction       = "joinVoice"
	LeaveGuildAction      = "leaveGuild"
	LeaveRoomAction       = "leaveRoom"
	LeaveVoiceAction      = "leaveVoice"
	StartTypingAction     = "startTyping"
	StopTypingAction      = "stopTyping"
	ToggleOnlineAction    = "toggleOnline"
	ToggleOfflineAction   = "toggleOffline"
	GetRequestCountAction = "getRequestCount"
)

// Emitted Messages
const (
	NewMessageAction        = "new_message"
	EditMessageAction       = "edit_message"
	DeleteMessageAction     = "delete_message"
	AddChannelAction        = "add_channel"
	AddPrivateChannelAction = "add_private_channel"
	EditChannelAction       = "edit_channel"
	DeleteChannelAction     = "delete_channel"
	EditGuildAction         = "edit_guild"
	DeleteGuildAction       = "delete_guild"
	RemoveFromGuildAction   = "remove_from_guild"
	AddMemberAction         = "add_member"
	RemoveMemberAction      = "remove_member"
	NewDMNotificationAction = "new_dm_notification"
	NewNotificationAction   = "new_notification"
	ToggleOnlineEmission    = "toggle_online"
	ToggleOfflineEmission   = "toggle_offline"
	AddToTypingAction       = "addToTyping"
	RemoveFromTypingAction  = "removeFromTyping"
	SendRequestAction       = "send_request"
	AddRequestAction        = "add_request"
	AddFriendAction         = "add_friend"
	RemoveFriendAction      = "remove_friend"
	PushToTopAction         = "push_to_top"
	RequestCountEmission    = "requestCount"
	VoiceSignal             = "voice-signal"
	ToggleMute              = "toggle-mute"
	ToggleDeafen            = "toggle-deafen"
)
