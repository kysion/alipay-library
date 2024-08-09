package consumer

type consumer struct {
	SexEnum    sex
	ActionEnum action
	Category   category
	AuthState  authState
}

var Consumer = consumer{
	SexEnum:    Sex,
	ActionEnum: ActionType,
	Category:   CategoryType,
	AuthState:  AuthState,
}
