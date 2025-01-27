// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package share_dao

import (
    "github.com/kysion/alipay-library/alipay_model/alipay_dao/internal"
	"github.com/kysion/base-library/utility/daoctl/dao_interface"
)

type AlipayConsumerConfigDao = dao_interface.TIDao[internal.AlipayConsumerConfigColumns]

func NewAlipayConsumerConfig(dao ...dao_interface.IDao) AlipayConsumerConfigDao {
	return (AlipayConsumerConfigDao)(internal.NewAlipayConsumerConfigDao(dao...))
}

var (
	AlipayConsumerConfig = NewAlipayConsumerConfig()
)
