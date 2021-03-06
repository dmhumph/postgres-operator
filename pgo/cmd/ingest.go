package cmd

/*
 Copyright 2018 Crunchy Data Solutions, Inc.
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

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	msgs "github.com/crunchydata/postgres-operator/apiservermsgs"
	"github.com/crunchydata/postgres-operator/pgo/api"
	"io/ioutil"
	"os"
)

type IngestConfigFile struct {
	WatchDir        string `json:"WatchDir"`
	DBHost          string `json:"DBHost"`
	DBPort          string `json:"DBPort"`
	DBName          string `json:"DBName"`
	DBSecret        string `json:"DBSecret"`
	DBTable         string `json:"DBTable"`
	DBColumn        string `json:"DBColumn"`
	MaxJobs         int    `json:"MaxJobs"`
	PVCName         string `json:"PVCName"`
	SecurityContext string `json:"SecurityContext"`
}

func createIngest(args []string) {

	if len(args) == 0 {
		fmt.Println("Error: An ingest name argument is required.")
		return
	}

	r, err := parseRequest(IngestConfig, args[0])
	if err != nil {
		fmt.Println("Error: Problem parsing ingest configuration file.")
		fmt.Println("Error: ", err)
		return
	}

	response, err := api.CreateIngest(httpclient, &SessionCredentials, &r)
	if err != nil {
		fmt.Println("Error: " + err.Error())
		os.Exit(2)
	}

	if response.Status.Code == msgs.Ok {
		fmt.Println("Created ingest.")
	} else {
		fmt.Println("Error: ", response.Status.Msg)
		os.Exit(2)
	}

}

func deleteIngest(args []string) {
	log.Debugf("deleteIngest called %v", args)

	for _, v := range args {

		response, err := api.DeleteIngest(httpclient, v, &SessionCredentials)

		if err != nil {
			fmt.Println("Error: ", response.Status.Msg)
			os.Exit(2)
		}

		if response.Status.Code == msgs.Ok {
			if len(response.Results) == 0 {
				fmt.Println("No ingests found.")
				return
			}
			for k := range response.Results {
				fmt.Println("Deleted ingest " + response.Results[k])
			}
		} else {
			fmt.Println("Error: ", response.Status.Msg)
			os.Exit(2)
		}

	}

}

func showIngest(args []string) {
	log.Debugf("showIngest called %v", args)

	for _, v := range args {

		response, err := api.ShowIngest(httpclient, v, &SessionCredentials)

		if err != nil {
			fmt.Println("Error: " + err.Error())
			os.Exit(2)
		}

		if response.Status.Code == msgs.Error {
			fmt.Println("Error: " + response.Status.Msg)
			os.Exit(2)
		}

		if len(response.Details) == 0 {
			fmt.Println("no ingests found")
			return
		}

		log.Debugf("response = %v", response)
		for _, d := range response.Details {
			showIngestItem(&d)
		}

	}

}

func showIngestItem(detail *msgs.ShowIngestResponseDetail) {
	fmt.Printf("%s%s\n", "", "")
	fmt.Printf("%s%s\n", "", "pgingest : "+detail.Ingest.Spec.Name)
	fmt.Printf("%s%s\n", TreeBranch, "name : "+detail.Ingest.Spec.Name)
	fmt.Printf("%s%s\n", TreeBranch, "watchdir : "+detail.Ingest.Spec.WatchDir)
	fmt.Printf("%s%s\n", TreeBranch, "dbhost : "+detail.Ingest.Spec.DBHost)
	fmt.Printf("%s%s\n", TreeBranch, "dbport : "+detail.Ingest.Spec.DBPort)
	fmt.Printf("%s%s\n", TreeBranch, "dbname : "+detail.Ingest.Spec.DBName)
	fmt.Printf("%s%s\n", TreeBranch, "dbsecret : "+detail.Ingest.Spec.DBSecret)
	fmt.Printf("%s%s\n", TreeBranch, "dbtable : "+detail.Ingest.Spec.DBTable)
	fmt.Printf("%s%s\n", TreeBranch, "dbcolumn : "+detail.Ingest.Spec.DBColumn)
	fmt.Printf("%s%s%d\n", TreeBranch, "maxjobs : ", detail.Ingest.Spec.MaxJobs)
	fmt.Printf("%s%s%d\n", TreeBranch, "Running Jobs : ", detail.JobCountRunning)
	fmt.Printf("%s%s%d\n", TreeBranch, "Completed Jobs : ", detail.JobCountCompleted)

	fmt.Println("")

}

func parseRequest(configFilePath, name string) (msgs.CreateIngestRequest, error) {
	var err error

	r := msgs.CreateIngestRequest{}
	r.Name = name

	raw, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		fmt.Println("Error: ", err)
		return r, err
	}

	c := IngestConfigFile{}

	json.Unmarshal(raw, &c)

	r.WatchDir = c.WatchDir
	r.DBHost = c.DBHost
	r.DBPort = c.DBPort
	r.DBName = c.DBName
	r.DBSecret = c.DBSecret
	r.DBTable = c.DBTable
	r.DBColumn = c.DBColumn
	r.MaxJobs = c.MaxJobs
	r.PVCName = c.PVCName
	r.SecurityContext = c.SecurityContext
	return r, err
}
