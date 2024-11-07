package db

import "strings"

type UserVariable struct {
	Name  string `json:"name"`
	Value any    `json:"value"`
}

type UserString string

func (s UserString) String() string {
	return string(s)
}

func GetUserVariable(name string) (*UserVariable, bool) {
	userVariable, exists := UserVariables.Load(name)
	if !exists {
		return nil, exists
	}

	userVariableValue, ok := userVariable.(*UserVariable)
	if ok {
		return userVariableValue, true
	}
	return nil, false
}

func AddUserVariable(name string, userVariable *UserVariable) {
	name = strings.ReplaceAll(name, " ", "_")
	UserVariables.Store(name, userVariable)
}

func DeleteUserVariable(name string) *UserVariable {
	uv, exists := GetUserVariable(name)

	if !exists {
		return nil
	}

	UserVariables.Delete(name)
	return uv
}
