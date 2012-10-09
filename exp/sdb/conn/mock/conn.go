// This file was auto-generated using createmock. See the following page for
// more information:
//
//     https://github.com/jacobsa/oglemock
//

package mock_conn

import (
	conn "github.com/jacobsa/aws/exp/sdb/conn"
	fmt "fmt"
	oglemock "github.com/jacobsa/oglemock"
	runtime "runtime"
	unsafe "unsafe"
)

type MockConn interface {
	conn.Conn
	oglemock.MockObject
}

type mockConn struct {
	controller	oglemock.Controller
	description	string
}

func NewMockConn(
	c oglemock.Controller,
	desc string) MockConn {
	return &mockConn{
		controller:	c,
		description:	desc,
	}
}

func (m *mockConn) Oglemock_Id() uintptr {
	return uintptr(unsafe.Pointer(m))
}

func (m *mockConn) Oglemock_Description() string {
	return m.description
}

func (m *mockConn) SendRequest(p0 conn.Request) (o0 []uint8, o1 error) {
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
		panic(fmt.Sprintf("mockConn.SendRequest: invalid return values: %v", retVals))
	}

	// o0 []uint8
	if retVals[0] != nil {
		o0 = retVals[0].([]uint8)
	}

	// o1 error
	if retVals[1] != nil {
		o1 = retVals[1].(error)
	}

	return
}
