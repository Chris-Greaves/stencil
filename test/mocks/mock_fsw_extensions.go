package mocks

import "os"

// Call the original os package functions as if it wasn't mocked.
// Due to the arguments being passed directly to the os package, mock.Anything is not supported.
func (_c *MockFsWrapper_GetHomeDirectory_Call) Passthrough() *MockFsWrapper_GetHomeDirectory_Call {
	_c.Call.Return(os.UserHomeDir)
	return _c
}

// Call the original os package functions as if it wasn't mocked.
// Due to the arguments being passed directly to the os package, mock.Anything is not supported.
func (_c *MockFsWrapper_MkdirAll_Call) Passthrough() *MockFsWrapper_MkdirAll_Call {
	_c.Call.Return(os.MkdirAll)
	return _c
}

// Call the original os package functions as if it wasn't mocked.
// Due to the arguments being passed directly to the os package, mock.Anything is not supported.
func (_c *MockFsWrapper_ReadFile_Call) Passthrough() *MockFsWrapper_ReadFile_Call {
	_c.Call.Return(os.ReadFile)
	return _c
}

// Call the original os package functions as if it wasn't mocked.
// Due to the arguments being passed directly to the os package, mock.Anything is not supported.
func (_c *MockFsWrapper_WriteFile_Call) Passthrough() *MockFsWrapper_WriteFile_Call {
	_c.Call.Return(os.WriteFile)
	return _c
}

// Call the original os package functions as if it wasn't mocked.
// Due to the arguments being passed directly to the os package, mock.Anything is not supported.
func (_c *MockFsWrapper_Lstat_Call) Passthrough() *MockFsWrapper_Lstat_Call {
	_c.Call.Return(os.Lstat)
	return _c
}
