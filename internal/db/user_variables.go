package db

import (
	"log/slog"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserVariable struct {
	Id       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name" bson:"name"`
	Value    any                `json:"value" bson:"value"`
	Filename string             `json:"filename" bson:"filename"`
}

func (v *UserVariable) GetId() string {
	return v.Id.Hex()
}

func (v *UserVariable) setUnmutableFields(existingFields *UserVariable) {
	v.Filename = existingFields.Filename
	v.Name = existingFields.Name
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

func GetUserVariableById(id string) (*UserVariable, bool) {
	var variableToReturn *UserVariable
	exists := false
	UserVariables.Range(func(_ any, value any) bool {
		variable, ok := value.(*UserVariable)
		if ok {
			if variable.Id.Hex() == id {
				variableToReturn = variable
				exists = true
				return false
			}
		}
		return true
	})
	return variableToReturn, exists

}

func StoreUserVariableInMemory(name string, userVariable *UserVariable) {
	name = strings.ReplaceAll(name, " ", "_")
	UserVariables.Store(name, userVariable)
}

func AddUserVariable(name string, userVariable *UserVariable, insertToDb bool) (*UserVariable, error) {
	StoreUserVariableInMemory(name, userVariable)
	if insertToDb {
		userVariable, err := Database.insertVariable(userVariable)
		return userVariable, err
	}
	return userVariable, nil
}

func UpdateUserVariable(name string, userVariable *UserVariable, insertToDb bool) (*UserVariable, error) {
	name = strings.ReplaceAll(name, " ", "_")
	existingUserVariable, exists := UserVariables.Load(name)
	if exists {
		if existingUserVariable, ok := existingUserVariable.(*UserVariable); ok {
			userVariable.setUnmutableFields(existingUserVariable)
			UserVariables.Store(name, userVariable)
		}
	}

	if insertToDb {
		userVariable, err := Database.updateVariable(userVariable)
		if err != nil {
			return nil, err
		}
		err = updateVariablesFile(userVariable.Filename, userVariable)
		return userVariable, err
	}
	return userVariable, nil
}

func AddUserVariables(variables ...*UserVariable) {
	filter := buildNameFilterForVariables(variables)
	results, err := Database.findUserVariables(filter)
	if err != nil {
		slog.Error("fail to find variables", "variables", variables)
		return
	}

	existingVariablesMap := make(map[string]*UserVariable)
	for _, result := range results {
		existingVariablesMap[result.Name] = &result
	}

	var variablesToInsert []interface{}
	for _, variable := range variables {
		existingVariable, exists := existingVariablesMap[variable.Name]
		if !exists {
			variablesToInsert = append(variablesToInsert, variable)
		} else {
			variable.Id = existingVariable.Id
		}
		StoreUserVariableInMemory(variable.Name, variable)
	}

	if len(variablesToInsert) > 0 {
		insertAndUpdateVariablesIds(variablesToInsert)
	}
}

func insertAndUpdateVariablesIds(variablesToInsert []interface{}) error {
	result, err := Database.insertVariables(variablesToInsert)
	if err != nil {
		slog.Error("error inserting variables to database", "error", err)
		return err
	}
	for i, variable := range variablesToInsert {
		variable, ok := variable.(*UserVariable)
		if ok {
			id, ok := result.InsertedIDs[i].(primitive.ObjectID)
			if ok {
				variable.Id = id
			}
		}
	}
	return nil
}

func DeleteUserVariable(id string) (*UserVariable, error) {
	uv, exists := GetUserVariableById(id)
	if !exists {
		return nil, nil
	}

	UserVariables.Delete(uv.Name)
	deletedVariable, err := Database.deleteVariable(id)
	if err != nil {
		return nil, err
	}

	err = deleteVariableFromFile(uv)
	if err != nil {
		return nil, err
	}
	return deletedVariable, nil
}
