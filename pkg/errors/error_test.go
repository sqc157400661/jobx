package errors

import (
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"testing"
)

func TestCode(t *testing.T) {
	//var base *Error
	err := Newf(http.StatusBadRequest, "reason", "message")
	err1 := errors.New("sdfsdafsd")
	err2 := errors.Wrap(err, err1.Error())
	fmt.Println(err2, Code(err2))
}
