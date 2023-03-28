package sub_account

import "github.com/kysion/base-library/utility/enum"

/*
   分账绑定请求响应参数 result_code:
       SUCCESS：分账关系绑定成功；
       FAIL：分账关系绑定失败。
*/

type TradeSubAccountEnum enum.IEnumCode[string]

type tradeSubAccount struct {
	Success TradeSubAccountEnum
	//Fail    TradeSubAccountEnum
}

var TradeSubAccountRes = tradeSubAccount{
	Success: enum.New[TradeSubAccountEnum]("10000", "Success"),
	//Fail:    enum.New[TradeSubAccountEnum]("40004", "参数无效"),
}

func (e tradeSubAccount) New(code string) TradeSubAccountEnum {
	if code == TradeSubAccountRes.Success.Code() {
		return e.Success
	}

	//if code == TradeSubAccountRes.Fail.Code() {
	//    return e.Fail
	//}
	panic("consumerTradeSubAccountRes: error")
}
