package consumer

import "github.com/kysion/base-library/utility/enum"

// AuthStateEnum 授权状态：1授权、2取消授权
type AuthStateEnum enum.IEnumCode[int]

type authState struct {
	Auth   AuthStateEnum
	UnAuth AuthStateEnum
}

var AuthState = authState{
	Auth:   enum.New[AuthStateEnum](1, "关注"),
	UnAuth: enum.New[AuthStateEnum](2, "取消关注"),
}

func (e authState) New(code int, description string) AuthStateEnum {
	if (code & AuthState.Auth.Code()) == AuthState.Auth.Code() {
		return e.Auth
	}
	if (code & AuthState.UnAuth.Code()) == AuthState.UnAuth.Code() {
		return e.UnAuth
	}
	return enum.New[AuthStateEnum](code, description)
}
