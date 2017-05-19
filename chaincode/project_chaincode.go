/*
Chaincode created for Oracle hackathon

*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

//can be derived from DB in real application
var organizationStr = "Blockbuster Studios"
var consultingOrgStr = "Primetime Editing Services"
var projectMilestonesStr = "::milestones::" + consultingOrgStr
var projectUsersStr = "::users::" + consultingOrgStr
var timeFormat = "02-Jan-2006"
var initialAmount = "200000"

//Data elements

// user entered time entry

type TimeEntry struct {
	ProjectName     string `json:"projectname"`
	TaskName        string `json:"taskname"`
	PersonName      string `json:"personname"`
	QuantityInHours string `json:"quantityhours"`
	ExpenseType     string `json:"expensetype"`
	DerivedAmount   string `json:"derivedamount"`
	EntryDate       string `json:entrydate`
}

type ProjectMilestone struct {
	ProjectName   string `json:"projectname"`
	MilestoneName string `json:"milestonename"`
	PersonName    string `json:"personname"`
	Amount        string `json:"amount"`
	DateActual    string `json:dateactual`
}

type UserRate struct {
	User string `json:"user"`
	Rate string `json:"rate"`
}

type OrgResult struct {
	CompletedWorkAmount   string           `json:"completedworkamount"`
	PendingContractAmount string           `json:"pendingcontractamount"`
	AmountPaid            string           `json:"amountpaid"`
	BalanceTobePaid       string           `json:"balancetobepaid"`
	ProjectResults        []ProjectResult  `json:"projectresults"`
	AmountPaidLog         []AmountTransfer `json:"amounttransactions"`
}

type ProjectResult struct {
	Name          string `json:"name"`
	Date          string `json:"date"`
	DerivedAmount string `json:"derivedamount"`
	ProjectName   string `json:"projectname"`
	TaskName      string `json:"taskname"`
}

type AmountTransfer struct {
	AmountPaid string `json:"amountpaid"`
	Date       string `json:"date"`
	ProjectName   string `json:"projectname"`
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
	} else if function == "resource_time_entry" {
		return t.EnterResourceTime(stub, args)
	} else if function == "complete_project_milestone" {
		return t.CompleteProjectMilestone(stub, args)
	} else if function == "write" {
		return t.Write(stub, args)
	} else if function == "pay_amount" {
		return t.PayAmount(stub, args)
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
	} else if function == "get_pending_amount" {
		return t.GetOrgOverview(stub, args)
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
	consultingProjects := []string{"Wonders of Galactica Project", "Making of Big Labowski Project", "Mission to Pluto"}

	jsonAsBytes, _ := json.Marshal(consultingProjects)
	err := stub.PutState(organizationStr+"::"+consultingOrgStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}

	//Initilizing the Project user rates (can be dynamically derived from DB in realtime) p1 -->[ {u1 100}, {u2 200 }]
	var projectUserRates []UserRate
	userrate := UserRate{}
	userrate.User = "Connor Horton"
	userrate.Rate = "200"
	projectUserRates = append(projectUserRates, userrate)

	userrate.User = "Lisa James"
	userrate.Rate = "200"
	projectUserRates = append(projectUserRates, userrate)

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
	userrate.User = "Connor Horton"
	userrate.Rate = "200"
	projectUserRates = append(projectUserRates, userrate)

	userrate.User = "Lisa James"
	userrate.Rate = "200"
	projectUserRates = append(projectUserRates, userrate)

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

	//initialize initial contract amount

	err = stub.PutState(consultingOrgStr, []byte(initialAmount))
	if err != nil {
		return nil, err
	}

	err = stub.PutState(consultingOrgStr+"::amount_paid", []byte("0"))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (t *SimpleChaincode) EnterResourceTime(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	//       0              1         2           3                  4            5
	// "ProjectName", "TaskName", "PersonName", "QuantityInHours","ExpenseType","EntryDate"
	var rate int
	var hours int

	if len(args) != 6 {
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
	if len(args[5]) <= 0 {
		return nil, errors.New("5th argument EntryDate must be a non-empty string")
	}

	timeEntry := TimeEntry{}
	timeEntry.ProjectName = args[0]
	timeEntry.TaskName = args[1]
	timeEntry.PersonName = args[2]
	timeEntry.QuantityInHours = args[3]
	timeEntry.DerivedAmount = "0"
	timeEntry.ExpenseType = args[4]
	//timeEntry.EntryDate = time.Now().Format(timeFormat)
	timeEntry.EntryDate = args[5]
	// derive amount

	projectUserRatesAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return nil, errors.New("Failed to get project user rates")
	}
	projectUserRates := []UserRate{}
	json.Unmarshal(projectUserRatesAsBytes, &projectUserRates)

	for i := range projectUserRates {
		if strings.ToLower(projectUserRates[i].User) == strings.ToLower(args[2]) {
			hours, _ = strconv.Atoi(args[3])
			rate, _ = strconv.Atoi(projectUserRates[i].Rate)
			timeEntry.DerivedAmount = strconv.Itoa(hours * rate)
		}
	}

	//get time entires for user and project
	projectUserTimeEntryAsBytes, err := stub.GetState(args[0] + "::" + args[2])
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

	projectUsersAsBytes, err := stub.GetState(args[0] + projectUsersStr)
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

	return nil, nil
}

func (t *SimpleChaincode) CompleteProjectMilestone(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	//       0              1                  2        3         4
	// "ProjectName", "MilestoneName", "PersonName", "Amount" , "Date"

	if len(args) != 5 {
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
	if len(args[4]) <= 0 {
		return nil, errors.New("4th argument Date must be a non-empty string")
	}

	projectMilestone := ProjectMilestone{}
	projectMilestone.ProjectName = args[0]
	projectMilestone.MilestoneName = args[1]
	projectMilestone.PersonName = args[2]
	projectMilestone.Amount = args[3]
	//projectMilestone.DateActual = time.Now().Format(timeFormat)
	projectMilestone.DateActual = args[4]

	//get time entires for user and project
	projectMilestonesAsBytes, err := stub.GetState(args[0] + projectMilestonesStr)
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

	return nil, nil
}

func (t *SimpleChaincode) GetOrgOverview(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	//    0        1
	// "GE", "ABCConsulting"

	var completedworkamount int64
	var pendingcontractamount int64
	var amountpaid int64
	var balancetobepaid int64

	orgResult := OrgResult{}
	projectResults := []ProjectResult{}
	projectResult := ProjectResult{}

fmt.Println("- start GetOrgOverview")

	//get projects
	projectsAsBytes, err := stub.GetState(organizationStr + "::" + consultingOrgStr)
	if err != nil {
		return nil, err
	}

	var projects []string
	json.Unmarshal(projectsAsBytes, &projects)

	for i := range projects {
		//for each project get time entires and milestones

		//get users for each project
		projectUsersAsBytes, err := stub.GetState(projects[i] + projectUsersStr)
		if err != nil {
			return nil, err
		}

		var projectUsers []string
		json.Unmarshal(projectUsersAsBytes, &projectUsers)

		for j := range projectUsers {
			//for each user get time entries

			timeEntriesAsBytes, err := stub.GetState(projects[i] + "::" + projectUsers[j])
			if err != nil {
				return nil, err
			}

			timeEntries := []TimeEntry{}
			json.Unmarshal(timeEntriesAsBytes, &timeEntries)

			for x := range timeEntries {
				//for each time entry add amount and add to list
				Aval, _ := strconv.ParseInt(timeEntries[x].DerivedAmount, 10, 32)
				completedworkamount += Aval

				projectResult.Name = timeEntries[x].PersonName + " reported " + timeEntries[x].QuantityInHours + " Hours "
				projectResult.Date = timeEntries[x].EntryDate
				projectResult.DerivedAmount = timeEntries[x].DerivedAmount
				projectResult.ProjectName = projects[i]
				projectResult.TaskName = timeEntries[x].TaskName
				projectResults = append(projectResults, projectResult)
			} //time entires

		} //users

		//get milestones for each project
		projectMilestonesAsBytes, err := stub.GetState(projects[i] + projectMilestonesStr)
		if err != nil {
			return nil, err
		}

		projectMilestones := []ProjectMilestone{}
		json.Unmarshal(projectMilestonesAsBytes, &projectMilestones)

		for y := range projectMilestones {
			//for each milestone add amount and add to list
			Aval, _ := strconv.ParseInt(projectMilestones[y].Amount, 10, 32)
			completedworkamount += Aval

			projectResult.Name = "Milestone " + projectMilestones[y].MilestoneName + " Completed "
			projectResult.Date = projectMilestones[y].DateActual
			projectResult.DerivedAmount = projectMilestones[y].Amount
			projectResult.ProjectName = projects[i]
			projectResult.TaskName = ""
			projectResults = append(projectResults, projectResult)
		} //time entires

	} //projects

	//do calculations comvert int to string ,return result
	orgResult.ProjectResults = projectResults
	orgResult.CompletedWorkAmount = strconv.FormatInt(completedworkamount, 10)

	initialAmountAsBytes, err := stub.GetState(consultingOrgStr)
	if err != nil {
		return nil, err
	}
	initialAmount, _ := strconv.ParseInt(string(initialAmountAsBytes), 10, 32)

	amountPaidAsBytes, err := stub.GetState(consultingOrgStr + "::amount_paid")
	if err != nil {
		return nil, err
	}
	amountpaid, _ = strconv.ParseInt(string(amountPaidAsBytes), 10, 32)

	pendingcontractamount = initialAmount - completedworkamount

	orgResult.PendingContractAmount = strconv.FormatInt(pendingcontractamount, 10)
	orgResult.AmountPaid = string(amountPaidAsBytes)

	balancetobepaid = completedworkamount - amountpaid
	orgResult.BalanceTobePaid = strconv.FormatInt(balancetobepaid, 10)

	//append amount paid transactions

	fmt.Println("- adding amount paid log")

	transfersAsBytes, err := stub.GetState(consultingOrgStr + "::amount_paid_log")
	if err != nil {
		return nil, errors.New("Failed to get the first account")
	}

	amountPaidLog := []AmountTransfer{}
	json.Unmarshal(transfersAsBytes, &amountPaidLog)

	orgResult.AmountPaidLog = amountPaidLog

	return json.Marshal(orgResult)
}

func (t *SimpleChaincode) Write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	if len(args[0]) <= 0 {
		return nil, errors.New("1st argument key must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return nil, errors.New("2nd argument value must be a non-empty string")
	}

	err = stub.PutState(args[0], []byte(args[1])) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (t *SimpleChaincode) PayAmount(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	//     0           1         2       3         4
	// "GE", "ABCConsulting", "1000" , "Date","Project Name"

	var err error
	fmt.Println("running PayAmount()")

	if len(args) < 5 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}

	if len(args[0]) <= 0 {
		return nil, errors.New("1st argument Organization must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return nil, errors.New("2nd argument Consulting Organization must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return nil, errors.New("3rd argument amount must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return nil, errors.New("4th argument Date must be a non-empty string")
	}
	if len(args[4]) <= 0 {
		return nil, errors.New("4th argument Project Name must be a non-empty string")
	}

	accountAsBytes, err := stub.GetState(args[1] + "::amount_paid")
	if err != nil {
		return nil, errors.New("Failed to get the first account")
	}

	Aval, err := strconv.Atoi(string(accountAsBytes))
	newAval, err := strconv.Atoi(args[2])

	Aval += newAval

	err = stub.PutState(args[1]+"::amount_paid", []byte(strconv.Itoa(Aval))) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}

	//log each entry
	transfersAsBytes, err := stub.GetState(args[1] + "::amount_paid_log")
	if err != nil {
		return nil, errors.New("Failed to get the first account")
	}

	amountPaidLog := []AmountTransfer{}
	json.Unmarshal(transfersAsBytes, &amountPaidLog)

	amountpaid := AmountTransfer{}
	amountpaid.AmountPaid = args[2]
	amountpaid.Date = args[3]
	amountpaid.ProjectName = args[4]
	amountPaidLog = append(amountPaidLog, amountpaid)

	jsonAsBytes, _ := json.Marshal(amountPaidLog)

	err = stub.PutState(args[1]+"::amount_paid_log", jsonAsBytes) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}

	return nil, nil
}
