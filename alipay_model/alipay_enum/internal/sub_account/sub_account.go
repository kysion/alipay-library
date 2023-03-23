package sub_account

import "github.com/kysion/base-library/utility/enum"

/*
   分账绑定请求响应参数 result_code:
       SUCCESS：分账关系绑定成功；
       FAIL：分账关系绑定失败。
*/

type SubAccountBindResEnum enum.IEnumCode[string]

type sub struct {
	Success SubAccountBindResEnum
	Fail    SubAccountBindResEnum
}

var SubAccountBindRes = sub{
	Success: enum.New[SubAccountBindResEnum]("SUCCESS", "分账关系绑定成功"),
	Fail:    enum.New[SubAccountBindResEnum]("FAIL", "分账关系绑定失败"),
}

func (e sub) New(code string) SubAccountBindResEnum {
	if code == SubAccountBindRes.Success.Code() {
		return e.Success
	}

	if code == SubAccountBindRes.Fail.Code() {
		return e.Fail
	}
	panic("consumerSubAccountBindRes: error")
}
