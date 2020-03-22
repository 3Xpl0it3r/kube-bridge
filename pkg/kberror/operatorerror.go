package kberror



type kubeOperatorError struct {
	message string
	fields map[string]string
	stackTrack []byte
	err error
}


func NewKubeOpeatorError()*kubeOperatorError{
	return &kubeOperatorError{
		message:    "",
		fields:      nil,
		stackTrack: nil,
		err:      nil,
	}
}

func(e *kubeOperatorError)AddFiled(k,v string)KubeBridgeError{
	e.fields[k] = v
	return e
}

func(e *kubeOperatorError)AddError(err error)KubeBridgeError{
	e.err = err
	return e
}

func(e *kubeOperatorError)Error()error{
	return e.err
}

func(e *kubeOperatorError)AddMessage(mesage string)KubeBridgeError{
	e.message = mesage
	return e
}

func(e *kubeOperatorError)AddStackTrack([]byte)KubeBridgeError{
	return e
}

func(e *kubeOperatorError)StackTrack()[]byte{
	return e.StackTrack()
}
func(e *kubeOperatorError)Fields()map[string]string{
	return e.fields
}