package responses

type ErrorCode struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// InternalErrorCodes groups error codes related to internal errors
var InternalErrorCodes = struct {
	ResponseError ErrorCode
}{
	ResponseError: ErrorCode{Code: 1000, Message: "Error creating response"},
}

// GroupErrorCodes groups error codes related to groups
var GroupErrorCodes = struct {
	NotFound      ErrorCode
	InvalidSchema ErrorCode
	UpdateFailed  ErrorCode
	DeleteFailed  ErrorCode
}{
	NotFound:      ErrorCode{Code: 2001, Message: "Group not found"},
	InvalidSchema: ErrorCode{Code: 2002, Message: "Group invalid schema"},
	UpdateFailed:  ErrorCode{Code: 2003, Message: "Group update failed"},
	DeleteFailed:  ErrorCode{Code: 2004, Message: "Group delete failed"},
}

// NodeErrorCodes groups error codes related to nodes
var NodeErrorCodes = struct {
	NotFound      ErrorCode
	InvalidSchema ErrorCode
	UpdateFailed  ErrorCode
	DeleteFailed  ErrorCode
}{
	NotFound:      ErrorCode{Code: 3001, Message: "Node not found"},
	InvalidSchema: ErrorCode{Code: 3002, Message: "Node invalid schema"},
	UpdateFailed:  ErrorCode{Code: 3003, Message: "Node update failed"},
	DeleteFailed:  ErrorCode{Code: 3004, Message: "Node delete failed"},
}

// SensorErrorCodes groups error codes related to sensors
var SensorErrorCodes = struct {
	NotFound      ErrorCode
	InvalidSchema ErrorCode
	UpdateFailed  ErrorCode
	DeleteFailed  ErrorCode
}{
	NotFound:      ErrorCode{Code: 4001, Message: "Sensor not found"},
	InvalidSchema: ErrorCode{Code: 4002, Message: "Sensor invalid schema"},
	UpdateFailed:  ErrorCode{Code: 4003, Message: "Sensor update failed"},
	DeleteFailed:  ErrorCode{Code: 4004, Message: "Sensor delete failed"},
}

// SensorErrorCodes groups error codes related to sensors
var ProcessorErrorCodes = struct {
	NotFound      ErrorCode
	InvalidSchema ErrorCode
	UpdateFailed  ErrorCode
	DeleteFailed  ErrorCode
	AlreadyExists ErrorCode
}{
	NotFound:      ErrorCode{Code: 4001, Message: "Processor not found"},
	InvalidSchema: ErrorCode{Code: 4002, Message: "Processor invalid schema"},
	UpdateFailed:  ErrorCode{Code: 4003, Message: "Processor update failed"},
	DeleteFailed:  ErrorCode{Code: 4004, Message: "Processor delete failed"},
	AlreadyExists: ErrorCode{Code: 4005, Message: "Processor already exists"},
}
