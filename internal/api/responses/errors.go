package responses

type ErrorCode struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// InternalErrorCodes groups error codes related to internal errors
var InternalErrorCodes = struct {
	ResponseError   ErrorCode
	UnexpectedError ErrorCode
}{
	ResponseError:   ErrorCode{Code: 1001, Message: "Error creating response"},
	UnexpectedError: ErrorCode{Code: 1001, Message: "Unexpected error"},
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
	NotFound:      ErrorCode{Code: 5001, Message: "Processor not found"},
	InvalidSchema: ErrorCode{Code: 5002, Message: "Processor invalid schema"},
	UpdateFailed:  ErrorCode{Code: 5003, Message: "Processor update failed"},
	DeleteFailed:  ErrorCode{Code: 5004, Message: "Processor delete failed"},
	AlreadyExists: ErrorCode{Code: 5005, Message: "Processor already exists"},
}

// SensorErrorCodes groups error codes related to sensors
var TriggerErrorCodes = struct {
	NotFound                ErrorCode
	InvalidSchema           ErrorCode
	CreationFailed          ErrorCode
	UpdateFailed            ErrorCode
	DeleteFailed            ErrorCode
	AlreadyExists           ErrorCode
	ConditionDoesNotCompile ErrorCode
}{
	NotFound:                ErrorCode{Code: 6001, Message: "Trigger not found"},
	InvalidSchema:           ErrorCode{Code: 6002, Message: "Trigger invalid schema"},
	CreationFailed:          ErrorCode{Code: 6003, Message: "Trigger creation failed"},
	UpdateFailed:            ErrorCode{Code: 6004, Message: "Trigger update failed"},
	DeleteFailed:            ErrorCode{Code: 6005, Message: "Trigger delete failed"},
	AlreadyExists:           ErrorCode{Code: 6006, Message: "Trigger already exists"},
	ConditionDoesNotCompile: ErrorCode{Code: 6007, Message: "Condition does not compile"},
}
