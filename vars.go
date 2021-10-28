package xmusemail

import "errors"

//Xmus errors
var (
	ErrClientIsNotValid   = errors.New("Client is not valid")
	ErrUnauthorized       = errors.New("auth is not valid")
	ErrYourEmailIsInvalid = errors.New("Your email is invalid")
	ErrRCPTConnection     = errors.New("rcpt connection error")
	ErrOnDataCommand      = errors.New("on data command error")
	ErrOnWriteData        = errors.New("on write data error")
	ErrOnCloseDataPipe    = errors.New("on close data pipe error")
)
