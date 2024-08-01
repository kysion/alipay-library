package consumer

import "github.com/kysion/base-library/utility/enum"

//用户性别: 0未知、1男、2女

type SexEnum enum.IEnumCode[int]

type sex struct {
	Unknown SexEnum
	Woman   SexEnum
	Man     SexEnum
}

var Sex = sex{
	Unknown: enum.New[SexEnum](0, "未知"),
	Man:     enum.New[SexEnum](1, "男"),
	Woman:   enum.New[SexEnum](2, "女"),
}

func (e sex) New(code int) SexEnum {
	if code == Sex.Unknown.Code() {
		return Sex.Unknown
	}

	if code == Sex.Man.Code() {
		return Sex.Man
	}

	if code == Sex.Woman.Code() {
		return Sex.Woman
	}

	panic("Sex.New: error")
}
