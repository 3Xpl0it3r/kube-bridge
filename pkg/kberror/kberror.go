package kberror

type KubeBridgeError interface {
	AddFiled(string,string)KubeBridgeError
	AddError(err error)KubeBridgeError
	Error()error
	AddMessage(string)KubeBridgeError
	AddStackTrack([]byte)KubeBridgeError
	StackTrack()[]byte
	Fields()map[string]string
}

