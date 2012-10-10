// This file was auto-generated using createmock. See the following page for
// more information:
//
//     https://github.com/jacobsa/oglemock
//

package mock_sdb

import (
	fmt "fmt"
	sdb "github.com/jacobsa/aws/exp/sdb"
	oglemock "github.com/jacobsa/oglemock"
	runtime "runtime"
	unsafe "unsafe"
)

type MockDomain interface {
	sdb.Domain
	oglemock.MockObject
}

type mockDomain struct {
	controller  oglemock.Controller
	description string
}

func NewMockDomain(
	c oglemock.Controller,
	desc string) MockDomain {
	return &mockDomain{
		controller:  c,
		description: desc,
	}
}

func (m *mockDomain) Oglemock_Id() uintptr {
	return uintptr(unsafe.Pointer(m))
}

func (m *mockDomain) Oglemock_Description() string {
	return m.description
}

func (m *mockDomain) BatchDeleteAttributes(p0 sdb.BatchDeleteMap) (o0 error) {
	// Get a file name and line number for the caller.
	_, file, line, _ := runtime.Caller(1)

	// Hand the call off to the controller, which does most of the work.
	retVals := m.controller.HandleMethodCall(
		m,
		"BatchDeleteAttributes",
		file,
		line,
		[]interface{}{p0})

	if len(retVals) != 1 {
		panic(fmt.Sprintf("mockDomain.BatchDeleteAttributes: invalid return values: %v", retVals))
	}

	// o0 error
	if retVals[0] != nil {
		o0 = retVals[0].(error)
	}

	return
}

func (m *mockDomain) BatchPutAttributes(p0 sdb.BatchPutMap) (o0 error) {
	// Get a file name and line number for the caller.
	_, file, line, _ := runtime.Caller(1)

	// Hand the call off to the controller, which does most of the work.
	retVals := m.controller.HandleMethodCall(
		m,
		"BatchPutAttributes",
		file,
		line,
		[]interface{}{p0})

	if len(retVals) != 1 {
		panic(fmt.Sprintf("mockDomain.BatchPutAttributes: invalid return values: %v", retVals))
	}

	// o0 error
	if retVals[0] != nil {
		o0 = retVals[0].(error)
	}

	return
}

func (m *mockDomain) DeleteAttributes(p0 sdb.ItemName, p1 []sdb.DeleteUpdate, p2 []sdb.Precondition) (o0 error) {
	// Get a file name and line number for the caller.
	_, file, line, _ := runtime.Caller(1)

	// Hand the call off to the controller, which does most of the work.
	retVals := m.controller.HandleMethodCall(
		m,
		"DeleteAttributes",
		file,
		line,
		[]interface{}{p0, p1, p2})

	if len(retVals) != 1 {
		panic(fmt.Sprintf("mockDomain.DeleteAttributes: invalid return values: %v", retVals))
	}

	// o0 error
	if retVals[0] != nil {
		o0 = retVals[0].(error)
	}

	return
}

func (m *mockDomain) GetAttributes(p0 sdb.ItemName, p1 bool, p2 []string) (o0 []sdb.Attribute, o1 error) {
	// Get a file name and line number for the caller.
	_, file, line, _ := runtime.Caller(1)

	// Hand the call off to the controller, which does most of the work.
	retVals := m.controller.HandleMethodCall(
		m,
		"GetAttributes",
		file,
		line,
		[]interface{}{p0, p1, p2})

	if len(retVals) != 2 {
		panic(fmt.Sprintf("mockDomain.GetAttributes: invalid return values: %v", retVals))
	}

	// o0 []sdb.Attribute
	if retVals[0] != nil {
		o0 = retVals[0].([]sdb.Attribute)
	}

	// o1 error
	if retVals[1] != nil {
		o1 = retVals[1].(error)
	}

	return
}

func (m *mockDomain) Name() (o0 string) {
	// Get a file name and line number for the caller.
	_, file, line, _ := runtime.Caller(1)

	// Hand the call off to the controller, which does most of the work.
	retVals := m.controller.HandleMethodCall(
		m,
		"Name",
		file,
		line,
		[]interface{}{})

	if len(retVals) != 1 {
		panic(fmt.Sprintf("mockDomain.Name: invalid return values: %v", retVals))
	}

	// o0 string
	if retVals[0] != nil {
		o0 = retVals[0].(string)
	}

	return
}

func (m *mockDomain) PutAttributes(p0 sdb.ItemName, p1 []sdb.PutUpdate, p2 []sdb.Precondition) (o0 error) {
	// Get a file name and line number for the caller.
	_, file, line, _ := runtime.Caller(1)

	// Hand the call off to the controller, which does most of the work.
	retVals := m.controller.HandleMethodCall(
		m,
		"PutAttributes",
		file,
		line,
		[]interface{}{p0, p1, p2})

	if len(retVals) != 1 {
		panic(fmt.Sprintf("mockDomain.PutAttributes: invalid return values: %v", retVals))
	}

	// o0 error
	if retVals[0] != nil {
		o0 = retVals[0].(error)
	}

	return
}
