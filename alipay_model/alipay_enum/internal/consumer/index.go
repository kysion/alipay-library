package consumer

type consumer struct {
	SexEnum    sex
	ActionEnum action
	Category   category
}

var Consumer = consumer{
	SexEnum:    SexType,
	ActionEnum: ActionType,
	Category:   CategoryType,
}
