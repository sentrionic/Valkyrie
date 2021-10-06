package apperrors

// Guild Errors
const (
	NotAMember             = "Not a member of the guild"
	AlreadyMember          = "Already a member of the guild"
	GuildLimitReached      = "The guild limit is 100"
	MustBeOwner            = "Must be the owner for that"
	InvalidImageType       = "imageFile must be 'image/jpeg' or 'image/png'"
	MustBeMemberInvite     = "Must be a member to fetch an invite"
	IsPermanentError       = "isPermanent is not a boolean"
	InvalidateInvitesError = "Only the owner can invalidate invites"
	InvalidInviteError     = "Invalid Link or the server got deleted"
	BannedFromServer       = "You are banned from this server"
	DeleteGuildError       = "Only the owner can delete their server"
	OwnerCantLeave         = "The owner cannot leave their server"
	BanYourselfError       = "You cannot ban yourself"
	KickYourselfError      = "You cannot kick yourself"
	UnbanYourselfError     = "You cannot unban yourself"
	OneChannelRequired     = "A server needs at least one channel"
	ChannelLimitError      = "The channel limit is 50"
	DMYourselfError        = "You cannot dm yourself"
)

// Account Errors
const (
	InvalidOldPassword  = "Invalid old password"
	InvalidCredentials  = "Invalid email and password combination"
	DuplicateEmail      = "An account with that email already exists"
	PasswordsDoNotMatch = "Passwords do not match"
	InvalidResetToken   = "Invalid reset token"
)

// Friend Errors
const (
	AddYourselfError    = "You cannot add yourself"
	RemoveYourselfError = "You cannot remove yourself"
	AcceptYourselfError = "You cannot accept yourself"
	CancelYourselfError = "You cannot cancel yourself"
	UnableAddError      = "Unable to add user as friend. Try again later"
	UnableRemoveError   = "Unable to remove the user. Try again later"
	UnableAcceptError   = "Unable to accept the request. Try again later"
)

// Generic Errors
const (
	InvalidSession = "Provided session is invalid"
	ServerError    = "Something went wrong. Try again later"
	Unauthorized   = "Not Authorized"
)

// Message Errors
const (
	MessageOrFileRequired = "Either a message or a file is required"
	EditMessageError      = "Only the author can edit the message"
	DeleteMessageError    = "Only the author or owner can delete the message"
	DeleteDMMessageError  = "Only the author can delete the message"
)
