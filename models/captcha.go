package models

//** 验证码models
type VerCode struct {
	Time int64  `json:"time"`
	Code string `json:"code"`
}
type SMSVerCode struct {
	VerCode
	TimeOut int64 `json:"timeout"`
}

func (vc *SMSVerCode) toVerCode() VerCode {
	return vc.VerCode
}
