/*
Chaincode created for Oracle hackathon

*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

var projectsIndexStr = "GE::ABCConsulting"
var projectMilestonesStr = "::milestones"
var projectUsersStr = "::users::ABCConsulting"

//Data elements

// user entered time entry

type TimeEntry struct {
	ProjectName     string `json:"projectname"`
	TaskName        string `json:"taskname"`
	PersonName      string `json:"personname"`
	QuantityInHours string `json:"quantityhours"`
	ExpenseType	    string `json:"expensetype"`
	DerivedAmount   string `json:"totalamount"`
}

type ProjectMilestone struct {
	ProjectName   string `json:"projectname"`
	MilestoneName string `json:"milestonename"`
	PersonName    string `json:"personname"`
	Amount        string `json:"amount"`
}

type UserRate struct {
	User string `json:"user"`
	Rate string `json:"rate"`
}

// ============================================================================================================================
//  Main - main - Starts up the chaincode
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	err := stub.PutState("test", []byte(args[0]))
	if err != nil {
		return nil, err
	}

	_, err = t.initializeData(stub, args)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" { //initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	}else if function == "resourcetimeentry" {
		return t.EnterResourceTime(stub,args)
	}else if function == "completeprojectmilestone" {
		return t.CompleteProjectMilestone(stub,args)
	}

	fmt.Println("invoke did not find func: " + function) //error

	return nil, errors.New("Received unknown function invocation: " + function)
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "dummy_query" { //read a variable
		fmt.Println("hi there " + function) //error
		return nil, nil
	} else if function == "read" {
		return t.read(stub, args)
	}

	fmt.Println("query did not find func: " + function) //error

	return nil, errors.New("Received unknown function query: " + function)
}

// test method to return the keys and read values
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}

//Initilizing project Data

func (t *SimpleChaincode) initializeData(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	//Initilizing the sample projects (can be dynamically derived from DB in realtime)
	consultingProjects := []string{"Proj1", "Proj2", "Proj3"}

	jsonAsBytes, _ := json.Marshal(consultingProjects)
	err := stub.PutState(projectsIndexStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}

	//Initilizing the Project user rates (can be dynamically derived from DB in realtime) p1 -->[ {u1 100}, {u2 200 }]
	var projectUserRates []UserRate
	userrate := UserRate{}
	userrate.User = "Chandra"
	userrate.Rate = "110"
	projectUserRates = append(projectUserRates, userrate)

	userrate.User = "Sudheer"
	userrate.Rate = "100"
	projectUserRates = append(projectUserRates, userrate)

	userrate.User = "Sanjay"
	userrate.Rate = "80"
	projectUserRates = append(projectUserRates, userrate)

	jsonAsBytes, _ = json.Marshal(projectUserRates)
	//initialize user rates for Proj1
	err = stub.PutState(consultingProjects[0], jsonAsBytes)
	if err != nil {
		return nil, err
	}

	//Initilizing the Project user rates p1 -->[ {u1 100}, {u2 200 }]
	projectUserRates = []UserRate{}
	userrate = UserRate{}
	userrate.User = "Chandra"
	userrate.Rate = "105"
	projectUserRates = append(projectUserRates, userrate)

	userrate.User = "Sudheer"
	userrate.Rate = "110"
	projectUserRates = append(projectUserRates, userrate)

	userrate.User = "Sanjay"
	userrate.Rate = "75"
	projectUserRates = append(projectUserRates, userrate)

	jsonAsBytes, _ = json.Marshal(projectUserRates)
	//initialize user rates for Proj2
	err = stub.PutState(consultingProjects[1], jsonAsBytes)
	if err != nil {
		return nil, err
	}

//initialize user rates for Proj3
	err = stub.PutState(consultingProjects[2], jsonAsBytes)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (t *SimpleChaincode) EnterResourceTime(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	//       0              1         2           3                  4
	// "ProjectName", "TaskName", "PersonName", "QuantityInHours","ExpenseType"
	var rate int
  var hours int

	if len(args) != 5 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	//input sanitation
	fmt.Println("- start EnterResourceTime")
	if len(args[0]) <= 0 {
		return nil, errors.New("1st argument ProjectName must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return nil, errors.New("2nd argument TaskName must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return nil, errors.New("3rd argument PersonName must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return nil, errors.New("4th argument QuantityInHours must be a non-empty string")
	}
	if len(args[4]) <= 0 {
		return nil, errors.New("5th argument ExpenseType must be a non-empty string")
	}

 timeEntry := TimeEntry{}
 timeEntry.ProjectName = args[0]
 timeEntry.TaskName = args[1]
 timeEntry.PersonName = args[2]
 timeEntry.QuantityInHours = args[3]
 timeEntry.DerivedAmount = "0"
 timeEntry.ExpenseType = args[4]
// derive amount

	projectUserRatesAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return nil, errors.New("Failed to get project user rates")
	}
	projectUserRates := []UserRate{}
	json.Unmarshal(projectUserRatesAsBytes, &projectUserRates)

	for i := range projectUserRates{
		if strings.ToLower(projectUserRates[i].User) == strings.ToLower(args[2]) {
			hours,_ = strconv.Atoi(args[3])
			rate,_ = strconv.Atoi(projectUserRates[i].Rate)
       timeEntry.DerivedAmount = strconv.Itoa(hours * rate )
		}
	}

//get time entires for user and project
projectUserTimeEntryAsBytes, err := stub.GetState(args[0]+"::"+args[2])
if err != nil {
	return nil, errors.New("Failed to get project user time entry")
}

allProjectTimeEntries := []TimeEntry{}
json.Unmarshal(projectUserTimeEntryAsBytes, &allProjectTimeEntries)

//add current time entry to exisitng
allProjectTimeEntries = append(allProjectTimeEntries, timeEntry)

//put back all time entries
jsonAsBytes, _ := json.Marshal(allProjectTimeEntries)
err = stub.PutState(args[0]+"::"+args[2], jsonAsBytes)
if err != nil {
	return nil, err
}
//put project and unique users

projectUsersAsBytes, err := stub.GetState(args[0]+projectUsersStr)
if err != nil {
	return nil, errors.New("Failed to get project user time entry")
}

var projectActiveUsers []string
json.Unmarshal(projectUsersAsBytes, &projectActiveUsers)

//check if user already exits ..if exits then nothing else add to the list
var userExists bool = false

for i := range projectActiveUsers {
  if strings.ToLower(projectActiveUsers[i]) == strings.ToLower(args[2]) {
  	userExists = true
  }
}

if !userExists {
	projectActiveUsers = append(projectActiveUsers, args[2])

	jsonAsBytes, _ = json.Marshal(projectActiveUsers)
	err = stub.PutState(args[0]+projectUsersStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}
}

	return nil,nil
}

func (t *SimpleChaincode) CompleteProjectMilestone(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	//       0              1                  2        3
	// "ProjectName", "MilestoneName", "PersonName", "Amount"

	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	//input sanitation
	fmt.Println("- start CompleteProjectMilestone")
	if len(args[0]) <= 0 {
		return nil, errors.New("1st argument ProjectName must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return nil, errors.New("2nd argument MilestoneName must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return nil, errors.New("3rd argument PersonName must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return nil, errors.New("4th argument Amount must be a non-empty string")
	}

	projectMilestone := ProjectMilestone{}
  projectMilestone.ProjectName = args[0]
  projectMilestone.MilestoneName = args[1]
  projectMilestone.PersonName = args[2]
  projectMilestone.Amount = args[3]

	//get time entires for user and project
	projectMilestonesAsBytes, err := stub.GetState(args[0]+projectMilestonesStr)
	if err != nil {
		return nil, errors.New("Failed to get project Milestones")
	}

	allProjectMilestones := []ProjectMilestone{}
	json.Unmarshal(projectMilestonesAsBytes, &allProjectMilestones)

	//add current time entry to exisitng
	allProjectMilestones = append(allProjectMilestones, projectMilestone)

	//put back all time entries
	jsonAsBytes, _ := json.Marshal(allProjectMilestones)
	err = stub.PutState(args[0]+projectMilestonesStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}

	return nil,nil
}
