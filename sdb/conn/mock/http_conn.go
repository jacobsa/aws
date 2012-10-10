// This file was auto-generated using createmock. See the following page for
// more information:
//
//     https://github.com/jacobsa/oglemock
//

package mock_conn

import (
	fmt "fmt"
	conn "github.com/jacobsa/aws/exp/sdb/conn"
	oglemock "github.com/jacobsa/oglemock"
	runtime "runtime"
	unsafe "unsafe"
)

type MockHttpConn interface {
	conn.HttpConn
	oglemock.MockObject
}

type mockHttpConn struct {
	controller  oglemock.Controller
	description string
}

func NewMockHttpConn(
	c oglemock.Controller,
	desc string) MockHttpConn {
	return &mockHttpConn{
		controller:  c,
		description: desc,
	}
}

func (m *mockHttpConn) Oglemock_Id() uintptr {
	return uintptr(unsafe.Pointer(m))
}

func (m *mockHttpConn) Oglemock_Description() string {
	return m.description
}

func (m *mockHttpConn) SendRequest(p0 conn.Request) (o0 *conn.HttpResponse, o1 error) {
	// Get a file name and line number for the caller.
	_, file, line, _ := runtime.Caller(1)

	// Hand the call off to the controller, which does most of the work.
	retVals := m.controller.HandleMethodCall(
		m,
		"SendRequest",
		file,
		line,
		[]interface{}{p0})

	if len(retVals) != 2 {
		panic(fmt.Sprintf("mockHttpConn.SendRequest: invalid return values: %v", retVals))
	}

	// o0 *conn.HttpResponse
	if retVals[0] != nil {
		o0 = retVals[0].(*conn.HttpResponse)
	}

	// o1 error
	if retVals[1] != nil {
		o1 = retVals[1].(error)
	}

	return
}
