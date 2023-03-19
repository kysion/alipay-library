package consumer

type consumer struct {
	SexEnum    sex
	ActionEnum action
}

var Consumer = consumer{
	SexEnum:    SexType,
	ActionEnum: ActionType,
}
