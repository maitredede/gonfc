package gonfc

import (
	"errors"
	"fmt"
)

var (
	NFC_EIO          error = &nfcIOError{}
	NFC_EINVARG      error = errors.New("NFC_EINVARG")
	NFC_EDEVNOTSUPP  error = errors.New("NFC_EDEVNOTSUPP")
	NFC_ENOTSUCHDEV  error = errors.New("NFC_ENOTSUCHDEV")
	NFC_EOVFLOW      error = errors.New("NFC_EOVFLOW")
	NFC_ETIMEOUT     error = &nfcTimeoutError{}
	NFC_EOPABORTED   error = errors.New("NFC_EOPABORTED")
	NFC_ENOTIMPL     error = errors.New("NFC_ENOTIMPL")
	NFC_ETGRELEASED  error = errors.New("NFC_ETGRELEASED")
	NFC_ERFTRANS     error = errors.New("NFC_ERFTRANS")
	NFC_EMFCAUTHFAIL error = errors.New("NFC_EMFCAUTHFAIL")
	NFC_ESOFT        error = errors.New("NFC_ESOFT")
	NFC_ECHIP        error = errors.New("NFC_ECHIP")
)

type nfcTimeoutError struct{}

func (nfcTimeoutError) Error() string {
	return "NFC_ETIMEOUT"
}

func (nfcTimeoutError) Timeout() bool {
	// compatibility with os.IsTimeout(err)
	return true
}

type nfcIOError struct {
	w error
	m string
}

func (e *nfcIOError) Error() string {
	return e.m
}

func (e *nfcIOError) Is(err error) bool {
	_, ok := err.(*nfcIOError)
	return err != nil && ok
}

func (e *nfcIOError) Unwrap() error {
	return e.w
}

func BuildNFC_EIO(err error) error {
	var m string
	if err == nil {
		m = "NFC_EIO"
	} else {
		m = fmt.Sprintf("NFC_EIO: %s", err.Error())
	}
	return &nfcIOError{
		w: err,
		m: m,
	}
}
