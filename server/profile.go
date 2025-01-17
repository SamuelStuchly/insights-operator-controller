/*
Copyright © 2019 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package server

import (
	"github.com/redhatinsighs/insights-operator-controller/logging"
	"github.com/redhatinsighs/insights-operator-controller/storage"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Read list of configuration profiles.
func listConfigurationProfiles(writer http.ResponseWriter, request *http.Request, storage storage.Storage) {
	profiles, err := storage.ListConfigurationProfiles()
	if err == nil {
		addJSONHeader(writer)
		addJSON(writer, profiles)
	} else {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, err.Error())
	}
}

// Read profile specified by its ID
func getConfigurationProfile(writer http.ResponseWriter, request *http.Request, storage storage.Storage) {
	id, err := retrieveIDRequestParameter(request)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, "Error reading profile ID from request\n")
		return
	}

	profile, err := storage.GetConfigurationProfile(int(id))
	if err == nil {
		addJSONHeader(writer)
		addJSON(writer, profile)
	} else {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, err.Error())
	}
}

// Create new configuration profile
func newConfigurationProfile(writer http.ResponseWriter, request *http.Request, storage storage.Storage, splunk logging.Client) {
	username, foundUsername := request.URL.Query()["username"]
	description, foundDescription := request.URL.Query()["description"]

	if !foundUsername {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, "User name needs to be specified\n")
		return
	}

	if !foundDescription {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, "Description needs to be specified\n")
		return
	}

	configuration, err := ioutil.ReadAll(request.Body)
	if err != nil || len(configuration) == 0 {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, "Configuration needs to be provided in the request body")
		return
	}

	splunk.LogAction("NewConfigurationProfile", username[0], string(configuration))
	profiles, err := storage.StoreConfigurationProfile(username[0], description[0], string(configuration))
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		io.WriteString(writer, err.Error())
	} else {
		addJSONHeader(writer)
		writer.WriteHeader(http.StatusCreated)
		addJSON(writer, profiles)
	}
}

// Delete configuration profile
func deleteConfigurationProfile(writer http.ResponseWriter, request *http.Request, storage storage.Storage, splunk logging.Client) {
	id, err := retrieveIDRequestParameter(request)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, "Error reading profile ID from request\n")
		return
	}

	splunk.LogAction("DeleteConfigurationProfile", "tester", strconv.Itoa(int(id)))
	profiles, err := storage.DeleteConfigurationProfile(int(id))
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, err.Error())
	} else {
		addJSONHeader(writer)
		writer.WriteHeader(http.StatusOK)
		addJSON(writer, profiles)
	}
}

// Change configuration profile
func changeConfigurationProfile(writer http.ResponseWriter, request *http.Request, storage storage.Storage, splunk logging.Client) {
	id, err := retrieveIDRequestParameter(request)
	username, foundUsername := request.URL.Query()["username"]
	description, foundDescription := request.URL.Query()["description"]

	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, "Error reading profile ID from request\n")
		return
	}

	if !foundUsername {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, "User name needs to be specified\n")
		return
	}

	if !foundDescription {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, "Description needs to be specified\n")
		return
	}

	configuration, err := ioutil.ReadAll(request.Body)
	if err != nil || len(configuration) == 0 {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, "Configuration needs to be provided in the request body")
		return
	}

	splunk.LogAction("ChangeConfigurationProfile", username[0], string(configuration))
	profiles, err := storage.ChangeConfigurationProfile(int(id), username[0], description[0], string(configuration))
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		io.WriteString(writer, err.Error())
	} else {
		addJSONHeader(writer)
		writer.WriteHeader(http.StatusAccepted)
		addJSON(writer, profiles)
	}
}
