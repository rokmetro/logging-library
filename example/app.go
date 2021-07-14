package main

import (
	"encoding/json"

	log "github.com/rokmetro/logging-library/loglib"
)

func main() {
	//Instantiate a logger for each service
	var logger = log.NewLogger("health-service")

	var random = 1234
	logger.Infof("%d", random)
	logger.InfoWithFields("ENV_VAR", log.Fields{"name": "test", "val": 123})

	//Instantiate a new log object for every request
	var logObj = logger.NewLog("12345", "8392")
	logObj.MissingArg("clientID")

	var userData log.Fields
	response := []byte(`{"uid":"123456789", "name":"John Doe"}`)

	if err := json.Unmarshal(response, &userData); err != nil {
		logger.Error("Failed to unmarshal")
	}
	//Add unstructured context like userData or tokenID to log
	logObj.AddContext("user_data", userData)
	logObj.AddContext("token_id", "aw901Q2jnk123")

	logObj.Info("Log object is working")

	// logObj.InvalidArg("tokenID", 4567)
	logObj.PrintContext()

}
