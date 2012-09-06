// This file was auto-generated using createmock. See the following page for
// more information:
//
//     https://github.com/jacobsa/oglemock
//

package mock_auth

import (
	auth "github.com/jacobsa/aws/s3/auth"
	fmt "fmt"
	http "github.com/jacobsa/aws/s3/http"
	oglemock "github.com/jacobsa/oglemock"
	runtime "runtime"
	unsafe "unsafe"
)

type MockSigner interface {
	auth.Signer
	oglemock.MockObject
}

type mockSigner struct {
	controller	oglemock.Controller
	description	string
}

func NewMockSigner(
	c oglemock.Controller,
	desc string) MockSigner {
	return &mockSigner{
		controller:	c,
		description:	desc,
	}
}

func (m *mockSigner) Oglemock_Id() uintptr {
	return uintptr(unsafe.Pointer(m))
}

func (m *mockSigner) Oglemock_Description() string {
	return m.description
}

func (m *mockSigner) Sign(p0 *http.Request) (o0 error) {
	// Get a file name and line number for the caller.
	_, file, line, _ := runtime.Caller(1)

	// Hand the call off to the controller, which does most of the work.
	retVals := m.controller.HandleMethodCall(
		m,
		"Sign",
		file,
		line,
		[]interface{}{p0})

	if len(retVals) != 1 {
		panic(fmt.Sprintf("mockSigner.Sign: invalid return values: %v", retVals))
	}

	// o0 error
	if retVals[0] != nil {
		o0 = retVals[0].(error)
	}

	return
}
