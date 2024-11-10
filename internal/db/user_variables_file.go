package db

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
)

func updateVariablesFile(filename string, variablesToUpdate ...*UserVariable) error {
	byteValue, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading file: %s", filename)
	}
	var variables []map[string]interface{}
	json.Unmarshal(byteValue, &variables)
	userVariablesMap := make(map[string]*UserVariable)
	for _, variable := range variablesToUpdate {
		userVariablesMap[variable.Name] = variable
	}

	for _, variableInFile := range variables {
		variableName, exists := variableInFile["name"]
		if exists {
			variableName, ok := variableName.(string)
			if ok {
				updatedVariable, exists := userVariablesMap[variableName]
				if exists {
					variableInFile["value"] = updatedVariable.Value
				}
			}
		}
	}

	updatedData, err := json.MarshalIndent(variables, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %v", err)
	}

	slog.Info("updating file", "file", filename)
	if err := os.WriteFile(filename, updatedData, 0644); err != nil {
		return fmt.Errorf("error writing file: %v", err)
	}
	return nil
}

func deleteVariableFromFile(variableToDelete *UserVariable) error {
	filename := variableToDelete.Filename
	byteValue, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading file: %s", filename)
	}
	var variables []map[string]interface{}
	json.Unmarshal(byteValue, &variables)
	var notDeletedVariables []map[string]interface{}
	for _, variableInFile := range variables {
		variableName, exists := variableInFile["name"]
		if exists {
			variableName, ok := variableName.(string)
			if ok && variableName != variableToDelete.Name {
				notDeletedVariables = append(notDeletedVariables, variableInFile)
			}
		}
	}

	updatedData, err := json.MarshalIndent(notDeletedVariables, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %v", err)
	}

	slog.Info("updating file", "file", filename)
	if err := os.WriteFile(filename, updatedData, 0644); err != nil {
		return fmt.Errorf("error writing file: %v", err)
	}
	return nil

}
