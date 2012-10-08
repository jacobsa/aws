// This file was auto-generated using createmock. See the following page for
// more information:
//
//     https://github.com/jacobsa/oglemock
//

package mock_s3

import (
	fmt "fmt"
	oglemock "github.com/jacobsa/oglemock"
	runtime "runtime"
	s3 "github.com/jacobsa/aws/s3"
	unsafe "unsafe"
)

type MockBucket interface {
	s3.Bucket
	oglemock.MockObject
}

type mockBucket struct {
	controller	oglemock.Controller
	description	string
}

func NewMockBucket(
	c oglemock.Controller,
	desc string) MockBucket {
	return &mockBucket{
		controller:	c,
		description:	desc,
	}
}

func (m *mockBucket) Oglemock_Id() uintptr {
	return uintptr(unsafe.Pointer(m))
}

func (m *mockBucket) Oglemock_Description() string {
	return m.description
}

func (m *mockBucket) DeleteObject(p0 string) (o0 error) {
	// Get a file name and line number for the caller.
	_, file, line, _ := runtime.Caller(1)

	// Hand the call off to the controller, which does most of the work.
	retVals := m.controller.HandleMethodCall(
		m,
		"DeleteObject",
		file,
		line,
		[]interface{}{p0})

	if len(retVals) != 1 {
		panic(fmt.Sprintf("mockBucket.DeleteObject: invalid return values: %v", retVals))
	}

	// o0 error
	if retVals[0] != nil {
		o0 = retVals[0].(error)
	}

	return
}

func (m *mockBucket) GetObject(p0 string) (o0 []uint8, o1 error) {
	// Get a file name and line number for the caller.
	_, file, line, _ := runtime.Caller(1)

	// Hand the call off to the controller, which does most of the work.
	retVals := m.controller.HandleMethodCall(
		m,
		"GetObject",
		file,
		line,
		[]interface{}{p0})

	if len(retVals) != 2 {
		panic(fmt.Sprintf("mockBucket.GetObject: invalid return values: %v", retVals))
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

func (m *mockBucket) ListKeys(p0 string) (o0 []string, o1 error) {
	// Get a file name and line number for the caller.
	_, file, line, _ := runtime.Caller(1)

	// Hand the call off to the controller, which does most of the work.
	retVals := m.controller.HandleMethodCall(
		m,
		"ListKeys",
		file,
		line,
		[]interface{}{p0})

	if len(retVals) != 2 {
		panic(fmt.Sprintf("mockBucket.ListKeys: invalid return values: %v", retVals))
	}

	// o0 []string
	if retVals[0] != nil {
		o0 = retVals[0].([]string)
	}

	// o1 error
	if retVals[1] != nil {
		o1 = retVals[1].(error)
	}

	return
}

func (m *mockBucket) StoreObject(p0 string, p1 []uint8) (o0 error) {
	// Get a file name and line number for the caller.
	_, file, line, _ := runtime.Caller(1)

	// Hand the call off to the controller, which does most of the work.
	retVals := m.controller.HandleMethodCall(
		m,
		"StoreObject",
		file,
		line,
		[]interface{}{p0, p1})

	if len(retVals) != 1 {
		panic(fmt.Sprintf("mockBucket.StoreObject: invalid return values: %v", retVals))
	}

	// o0 error
	if retVals[0] != nil {
		o0 = retVals[0].(error)
	}

	return
}
