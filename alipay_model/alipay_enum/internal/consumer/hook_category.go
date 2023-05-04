package consumer

import "github.com/kysion/base-library/utility/enum"

// CategoryEnum 类别：消费者 + 平台与用户
type CategoryEnum enum.IEnumCode[int]

type category struct {
	Consumer     CategoryEnum
	PlatFormUser CategoryEnum
}

var CategoryType = category{
	Consumer:     enum.New[CategoryEnum](1, "消费者"),
	PlatFormUser: enum.New[CategoryEnum](2, "平台与用户"),
}

func (e category) New(code int, description string) CategoryEnum {
	if (code & CategoryType.Consumer.Code()) == CategoryType.Consumer.Code() {
		return e.Consumer
	}
	if (code & CategoryType.PlatFormUser.Code()) == CategoryType.PlatFormUser.Code() {
		return e.PlatFormUser
	}
	return enum.New[CategoryEnum](code, description)
}
