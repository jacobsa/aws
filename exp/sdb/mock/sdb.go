// This file was auto-generated using createmock. See the following page for
// more information:
//
//     https://github.com/jacobsa/oglemock
//

package mock_sdb

import (
	fmt "fmt"
	oglemock "github.com/jacobsa/oglemock"
	runtime "runtime"
	sdb "github.com/jacobsa/aws/exp/sdb"
	unsafe "unsafe"
)

type MockSimpleDB interface {
	sdb.SimpleDB
	oglemock.MockObject
}

type mockSimpleDB struct {
	controller	oglemock.Controller
	description	string
}

func NewMockSimpleDB(
	c oglemock.Controller,
	desc string) MockSimpleDB {
	return &mockSimpleDB{
		controller:	c,
		description:	desc,
	}
}

func (m *mockSimpleDB) Oglemock_Id() uintptr {
	return uintptr(unsafe.Pointer(m))
}

func (m *mockSimpleDB) Oglemock_Description() string {
	return m.description
}

func (m *mockSimpleDB) DeleteDomain(p0 sdb.Domain) (o0 error) {
	// Get a file name and line number for the caller.
	_, file, line, _ := runtime.Caller(1)

	// Hand the call off to the controller, which does most of the work.
	retVals := m.controller.HandleMethodCall(
		m,
		"DeleteDomain",
		file,
		line,
		[]interface{}{p0})

	if len(retVals) != 1 {
		panic(fmt.Sprintf("mockSimpleDB.DeleteDomain: invalid return values: %v", retVals))
	}

	// o0 error
	if retVals[0] != nil {
		o0 = retVals[0].(error)
	}

	return
}

func (m *mockSimpleDB) OpenDomain(p0 string) (o0 sdb.Domain, o1 error) {
	// Get a file name and line number for the caller.
	_, file, line, _ := runtime.Caller(1)

	// Hand the call off to the controller, which does most of the work.
	retVals := m.controller.HandleMethodCall(
		m,
		"OpenDomain",
		file,
		line,
		[]interface{}{p0})

	if len(retVals) != 2 {
		panic(fmt.Sprintf("mockSimpleDB.OpenDomain: invalid return values: %v", retVals))
	}

	// o0 sdb.Domain
	if retVals[0] != nil {
		o0 = retVals[0].(sdb.Domain)
	}

	// o1 error
	if retVals[1] != nil {
		o1 = retVals[1].(error)
	}

	return
}

func (m *mockSimpleDB) Select(p0 string, p1 bool, p2 []uint8) (o0 map[sdb.ItemName][]sdb.Attribute, o1 []uint8, o2 error) {
	// Get a file name and line number for the caller.
	_, file, line, _ := runtime.Caller(1)

	// Hand the call off to the controller, which does most of the work.
	retVals := m.controller.HandleMethodCall(
		m,
		"Select",
		file,
		line,
		[]interface{}{p0, p1, p2})

	if len(retVals) != 3 {
		panic(fmt.Sprintf("mockSimpleDB.Select: invalid return values: %v", retVals))
	}

	// o0 map[sdb.ItemName][]sdb.Attribute
	if retVals[0] != nil {
		o0 = retVals[0].(map[sdb.ItemName][]sdb.Attribute)
	}

	// o1 []uint8
	if retVals[1] != nil {
		o1 = retVals[1].([]uint8)
	}

	// o2 error
	if retVals[2] != nil {
		o2 = retVals[2].(error)
	}

	return
}
