package consumer

type consumer struct {
	SexEnum    sex
	ActionEnum action
	Category   category
	AuthState  authState
}

var Consumer = consumer{
	SexEnum:    SexType,
	ActionEnum: ActionType,
	Category:   CategoryType,
	AuthState:  AuthState,
}
