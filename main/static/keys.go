package static

type Key string

const (
	//general
	UUIDKey Key = "uuid"
	//Discord
	CallerIDKey Key = "calledID"
	//Player
	PlayerIDKey Key = "playerID"
	UsernameKey Key = "username"
	//Etc
	ErrorKey Key = "error"
)
