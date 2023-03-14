package consumer

import "github.com/kysion/base-library/utility/enum"

type SexEnum enum.IEnumCode[int]

type sex struct {
	In  SexEnum
	Out SexEnum
}

var SexType = sex{
	In:  enum.New[SexEnum](0, "女"),
	Out: enum.New[SexEnum](1, "男"),
}

func (e sex) New(code int, description string) SexEnum {
	if (code&SexType.In.Code()) == SexType.In.Code() ||
		(code&SexType.Out.Code()) == SexType.Out.Code() {
		return enum.New[SexEnum](code, description)
	}
	panic("consumerSex: error")
}
