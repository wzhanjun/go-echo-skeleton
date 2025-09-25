package enum

type ErrCode int32

//go:generate stringer -type ErrCode -linecomment -output errcode_string.go
const (
	Success             ErrCode = 200 // success
	ErrCodeUnauthorized ErrCode = 401 // unauthorized
	ErrCodeNotFound     ErrCode = 404 // not found

	Error       ErrCode = 1000 // error
	ParamsError ErrCode = 1001 // params error

)

func (s ErrCode) Error() string {
	return s.String()
}

type ApiError struct {
	Code ErrCode
	Msg  string
}

func (s ApiError) Error() string {
	if s.Msg == "" {
		return s.Code.String()
	}
	return s.Msg
}
