/*
	Copyright IBM Corp. 2016 All Rights Reserved.

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

/*
  @Author:
    XIN LIU
  @Description:
    This is the chaincode on the housing project, including 3 data structures:
    Tenancy Contract Table:
    Insurance Contract Table;
    Request Table;
  @History:
  |-------------------------------------------------------------------------|
  |    Date    |      Owner    |            Comments                        |
  |-------------------------------------------------------------------------|
  | 05/15/2017 |      SEAN     |            Creation                        |
  |-------------------------------------------------------------------------|

*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/core/util"
)

//============================================================================
// Insurance Package Types:
//   Two different packages defined with different premium and items
//============================================================================
const InsurancePkg_Standard = 0
const InsurancePkg_Standard_Plus = 1

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

//============================================================================
// TenancyContract - A table structure to save tenancy contract
//                   between Landlord and Tenancy
//============================================================================
type TenancyContract struct {
	TenancyContractID string `json:"TenancyContractID"`
	LandlordID        string `json:"LandlordID"`
	TenantID          string `json:"TenantID"`
	Address           string `json:"Address"`
}

//============================================================================
// InsuranceContract - A table structure to save insurance contract
//                     between Landlord and insurance company
//============================================================================
type InsuranceContract struct {
	InsuranceContractID string `json:"InsuranceContractID"`
	InsuranceCompany    string `json:"InsuranceCompany"`
	LandlordID          string `json:"LandlordID"`
	Address             string `json:"Address"`
	InsurancePkg        int32  `json:"InsurancePkg"`
	StartDate           string `json:"StartDate"`
	EndDate             string `json:"EndDate"`
}

//============================================================================
// Request - A table structure to save repair request
//           among Landlord, Tenant and Service Provider
//============================================================================
type Request struct {
	RequestID         string `json:"RequestID"`
	LandlordID        string `json:"LandlordID"`
	TenantID          string `json:"TenantID"`
	ServiceProviderID string `json:"ServiceProviderID"`
	Address           string `json:"Address"`
	RequestStatus     string `json:"RequestStatus"`
	GoodsType         string `json:"GoodsType"`
	GoodsBrand        string `json:"GoodsBrand"`
	GoodsModel        string `json:"GoodsModel"`
	GoodsDescription  string `json:"GoodsDescription"`
	Price             string `json:"Price"`
	Signature         string `json:"Signature"`
	Receipt           string `json:"Receipt"`
}

//==============================================================================================================================
// Init Function of a Chaincode
//==============================================================================================================================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if len(args) != 0 {
		return nil, errors.New("incorrect number of arguments. Expecting 0")
	}

	return nil, nil
}

//==============================================================================================================================
// Invoke Function of a Chaincode
//                  - CreateTableTenancyContract
//					- CreateTableInsuranceContract
//					- CreateTableRequest
//                  - InsertRowRequest
//                  - DeleteRowRequest
//					- UpdateRowRequest
//					- ...
//==============================================================================================================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if function == "createTableTenancyContract" {

		table1Args := []string{"Tenancy_Contract"}
		return t.createTable(stub, table1Args)

	} else if function == "createTableInsuranceContract" {

		table2Args := []string{"Insurance_Contract"}
		return t.createTable(stub, table2Args)

	} else if function == "createTableRequest" {

		table3Args := []string{"Request"}
		return t.createTable(stub, table3Args)

	} else if function == "deleteTable" {
		return t.deleteTable(stub, args)

	} else if function == "insertRowTenancyContract" {
		return t.insertRowTenancyContract(stub, args)

	} else if function == "insertRowInsuranceContract" {
		return t.insertRowInsuranceContract(stub, args)

	} else if function == "insertRowRequest" {
		return t.insertRowRequest(stub, args)

	} else if function == "deleteRowTenancyContract" {
		return t.deleteRowTenancyContract(stub, args)

	} else if function == "deleteRowInsuranceContract" {
		return t.deleteRowInsuranceContract(stub, args)

	} else if function == "deleteRowRequest" {
		return t.deleteRowRequest(stub, args)

	} else if function == "updateRowRequest" {
		return t.updateRowRequest(stub, args)
	}

	return nil, errors.New("Unsupported Invoke Functions [" + function + "]")
}

//==============================================================================================================================
// Query Function of a Chaincode
//==============================================================================================================================
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	if function == "getTenancyContract" {
		return t.getTenancyContract(stub, args)

	} else if function == "getTenancyContractsByID" {
		return t.getTenancyContractsByID(stub, args)

	} else if function == "getAllTenancyContracts" {
		return t.getAllTenancyContracts(stub, args)

	} else if function == "getInsuranceContract" {
		return t.getInsuranceContract(stub, args)

	} else if function == "getInsuranceContractsByID" {
		return t.getInsuranceContractsByID(stub, args)

	} else if function == "getAllInsuranceContracts" {
		return t.getAllTenancyContracts(stub, args)

	} else if function == "getRequest" {
		return t.getRequest(stub, args)

	} else if function == "getRequestsByID" {
		return t.getRequestsByID(stub, args)

	} else if function == "getAllRequests" {
		return t.getAllRequests(stub, args)
	}

	return nil, errors.New("Received unsupported query parameter [" + function + "]")
}

//==============================================================================================================================
// createTable - Create a table
//    input - a table name (Tenancy_Contract, Insurance_Contract, Request)
//==============================================================================================================================
func (t *SimpleChaincode) createTable(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 1 {
		return nil, errors.New("incorrect number of auguments, Expencting 1, the name of the table")
	}

	tableName := args[0]

	fmt.Println("Start to create table " + tableName + " ...")

	var columnDefs []*shim.ColumnDefinition

	if tableName == "Tenancy_Contract" {

		columnDef1 := shim.ColumnDefinition{Name: "TenancyContractID", Type: shim.ColumnDefinition_STRING, Key: true}
		columnDef2 := shim.ColumnDefinition{Name: "LandlordID", Type: shim.ColumnDefinition_STRING, Key: false}
		columnDef3 := shim.ColumnDefinition{Name: "TenantID", Type: shim.ColumnDefinition_STRING, Key: false}
		columnDef4 := shim.ColumnDefinition{Name: "Address", Type: shim.ColumnDefinition_STRING, Key: false}

		columnDefs = append(columnDefs, &columnDef1)
		columnDefs = append(columnDefs, &columnDef2)
		columnDefs = append(columnDefs, &columnDef3)
		columnDefs = append(columnDefs, &columnDef4)

	} else if tableName == "Insurance_Contract" {

		columnDef1 := shim.ColumnDefinition{Name: "InsuranceContractID", Type: shim.ColumnDefinition_STRING, Key: true}
		columnDef2 := shim.ColumnDefinition{Name: "InsuranceCompany", Type: shim.ColumnDefinition_STRING, Key: false}
		columnDef3 := shim.ColumnDefinition{Name: "LandlordID", Type: shim.ColumnDefinition_STRING, Key: false}
		columnDef4 := shim.ColumnDefinition{Name: "Address", Type: shim.ColumnDefinition_STRING, Key: false}
		columnDef5 := shim.ColumnDefinition{Name: "InsurancePkg", Type: shim.ColumnDefinition_INT32, Key: false}
		columnDef6 := shim.ColumnDefinition{Name: "StartDate", Type: shim.ColumnDefinition_STRING, Key: false}
		columnDef7 := shim.ColumnDefinition{Name: "EndDate", Type: shim.ColumnDefinition_STRING, Key: false}

		columnDefs = append(columnDefs, &columnDef1)
		columnDefs = append(columnDefs, &columnDef2)
		columnDefs = append(columnDefs, &columnDef3)
		columnDefs = append(columnDefs, &columnDef4)
		columnDefs = append(columnDefs, &columnDef5)
		columnDefs = append(columnDefs, &columnDef6)
		columnDefs = append(columnDefs, &columnDef7)

	} else if tableName == "Request" {

		columnDef1 := shim.ColumnDefinition{Name: "RequestID", Type: shim.ColumnDefinition_STRING, Key: true}
		columnDef2 := shim.ColumnDefinition{Name: "LandlordID", Type: shim.ColumnDefinition_STRING, Key: false}
		columnDef3 := shim.ColumnDefinition{Name: "TenantID", Type: shim.ColumnDefinition_STRING, Key: false}
		columnDef4 := shim.ColumnDefinition{Name: "ServiceProviderID", Type: shim.ColumnDefinition_STRING, Key: false}
		columnDef5 := shim.ColumnDefinition{Name: "Address", Type: shim.ColumnDefinition_STRING, Key: false}
		columnDef6 := shim.ColumnDefinition{Name: "RequestStatus", Type: shim.ColumnDefinition_STRING, Key: false}
		columnDef7 := shim.ColumnDefinition{Name: "GoodsType", Type: shim.ColumnDefinition_STRING, Key: false}
		columnDef8 := shim.ColumnDefinition{Name: "GoodsBrand", Type: shim.ColumnDefinition_STRING, Key: false}
		columnDef9 := shim.ColumnDefinition{Name: "GoodsModel", Type: shim.ColumnDefinition_STRING, Key: false}
		columnDef10 := shim.ColumnDefinition{Name: "GoodsDescription", Type: shim.ColumnDefinition_STRING, Key: false}
		columnDef11 := shim.ColumnDefinition{Name: "Price", Type: shim.ColumnDefinition_STRING, Key: false}
		columnDef12 := shim.ColumnDefinition{Name: "Signature", Type: shim.ColumnDefinition_STRING, Key: false}
		columnDef13 := shim.ColumnDefinition{Name: "Receipt", Type: shim.ColumnDefinition_STRING, Key: false}

		columnDefs = append(columnDefs, &columnDef1)
		columnDefs = append(columnDefs, &columnDef2)
		columnDefs = append(columnDefs, &columnDef3)
		columnDefs = append(columnDefs, &columnDef4)
		columnDefs = append(columnDefs, &columnDef5)
		columnDefs = append(columnDefs, &columnDef6)
		columnDefs = append(columnDefs, &columnDef7)
		columnDefs = append(columnDefs, &columnDef8)
		columnDefs = append(columnDefs, &columnDef9)
		columnDefs = append(columnDefs, &columnDef10)
		columnDefs = append(columnDefs, &columnDef11)
		columnDefs = append(columnDefs, &columnDef12)
		columnDefs = append(columnDefs, &columnDef13)

	} else {
		return nil, errors.New("Unsupported table name " + tableName + "!!!")
	}

	err := stub.CreateTable(tableName, columnDefs)

	if err != nil {
		return nil, fmt.Errorf("Create a table operation failed. %s", err)
	}

	fmt.Println("End to create table " + tableName + " ...")
	return nil, nil
}

//==============================================================================================================================
// deleteTable - Delete a table
//    input - a table name (Tenancy_Contract, Insurance_Contract, Request)
//==============================================================================================================================
func (t *SimpleChaincode) deleteTable(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of auguments, Expencting 1, the name of the table")
	}

	tableName := args[0]

	err := stub.DeleteTable(tableName)

	if err != nil {
		return nil, fmt.Errorf("Delete a table operation failed. %s", err)
	}

	return nil, nil
}

//==============================================================================================================================
// insertRowTenancyContract - Insert a row into the table Tenancy_Contract
//    input - array of key:Value
//            [
//            "LandlordID" : "xxx";
//            "TenantID"   : "yyy";
//            "Address"    : "zzz zzz zzz"
//            ]
//==============================================================================================================================
func (t *SimpleChaincode) insertRowTenancyContract(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 6 {
		return nil, errors.New("Incorrect number of arguments. Expecting 6")
	}

	var landlordID string
	var tenantID string
	var address string

	for i := 0; i < len(args); i = i + 2 {
		switch args[i] {

		case "LandlordID":
			landlordID = args[i+1]
		case "TenantID":
			tenantID = args[i+1]
		case "Address":
			address = args[i+1]

		default:
			return nil, errors.New("Unsupported Parameter " + args[i] + "!!!")
		}
	}

	tenancyContractID := util.GenerateUUID()

	var columns []*shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: tenancyContractID}}
	col2 := shim.Column{Value: &shim.Column_String_{String_: landlordID}}
	col3 := shim.Column{Value: &shim.Column_String_{String_: tenantID}}
	col4 := shim.Column{Value: &shim.Column_String_{String_: address}}
	columns = append(columns, &col1)
	columns = append(columns, &col2)
	columns = append(columns, &col3)
	columns = append(columns, &col4)

	row := shim.Row{Columns: columns}
	ok, err := stub.InsertRow("Tenancy_Contract", row)

	if err != nil {
		return nil, fmt.Errorf("Insert a row in Tenancy_Contract table operation failed. %s", err)
	}

	if !ok {
		return nil, errors.New("Insert a row in Tenancy_Contract table operation failed. Row with given key already exists")
	}

	fmt.Println("TenancyContractID inserted is [" + tenancyContractID + "] ...")
	return []byte(tenancyContractID), err
}

//==============================================================================================================================
// insertRowInsuranceContract - Insert a row into the table Insurance_Contract
//    input - array of key:Value
//            [
//            "InsuranceCompany" : "China Ping'An";
//            "LandlordID"   : "L000001";
//            "Address"      : "Keyuan Road  #399, Shanghai";
//            "InsurancePkg" : "Standard_Plus";
//            "StartDate"    : "01/01/2017";
//            "EndDate"      : "12/31/2017";
//            ]
//==============================================================================================================================
func (t *SimpleChaincode) insertRowInsuranceContract(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 12 {
		return nil, errors.New("Incorrect number of arguments. Expecting 12")
	}

	var insuranceCompany string
	var landlordID string
	var address string
	var insurancePkg int32
	var startDate string
	var endDate string

	for i := 0; i < len(args); i = i + 2 {
		switch args[i] {

		case "InsuranceCompany":
			insuranceCompany = args[i+1]
		case "LandlordID":
			landlordID = args[i+1]
		case "Address":
			address = args[i+1]
		case "InsurancePkg":
			var pkgValue int
			if args[i+1] == "Standard" {
				pkgValue = InsurancePkg_Standard
			} else if args[i+1] == "Standard_Plus" {
				pkgValue = InsurancePkg_Standard_Plus
			} else {
				return nil, errors.New("Unsupported Insurance Package " + args[i+1] + "!!!")
			}
			insurancePkg = int32(pkgValue)
		case "StartDate":
			startDate = args[i+1]
		case "EndDate":
			endDate = args[i+1]
		default:
			return nil, errors.New("Unsupported Parameter " + args[i] + "!!!")
		}
	}

	insuranceContractID := util.GenerateUUID()

	var columns []*shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: insuranceContractID}}
	col2 := shim.Column{Value: &shim.Column_String_{String_: insuranceCompany}}
	col3 := shim.Column{Value: &shim.Column_String_{String_: landlordID}}
	col4 := shim.Column{Value: &shim.Column_String_{String_: address}}
	col5 := shim.Column{Value: &shim.Column_Int32{Int32: insurancePkg}}
	col6 := shim.Column{Value: &shim.Column_String_{String_: startDate}}
	col7 := shim.Column{Value: &shim.Column_String_{String_: endDate}}

	columns = append(columns, &col1)
	columns = append(columns, &col2)
	columns = append(columns, &col3)
	columns = append(columns, &col4)
	columns = append(columns, &col5)
	columns = append(columns, &col6)
	columns = append(columns, &col7)

	row := shim.Row{Columns: columns}
	ok, err := stub.InsertRow("Insurance_Contract", row)

	if err != nil {
		return nil, fmt.Errorf("Insert a row in Insurance_Contract table operation failed. %s", err)
	}

	if !ok {
		return nil, errors.New("Insert a row in Insurance_Contract table operation failed. Row with given key already exists")
	}

	fmt.Println("InsuranceContractID inserted is [" + insuranceContractID + "] ...")
	return []byte(insuranceContractID), err
}

//==============================================================================================================================
// insertRowRequest- Insert a row into the table Request
//    input - array of key:Value
//            [
//            "LandlordID"   : "L000001";
//            "TenantID  "   : "T000100";
//            "Address"      : "Keyuan Road  #399, Shanghai";
//            "GoodsType"    : "TV";
//            "GoodsBrand"   : "Changhong";
//            "GoodsModel"   : "LCD50";
//            "GoodsDescription"   :  "Display picture blurry.";
//            ]
//==============================================================================================================================
func (t *SimpleChaincode) insertRowRequest(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) < 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting >= 4, at least including LandlordID and TenantID")
	}

	var landlordID string
	var tenantID string
	var address string
	var requestStatus string
	var goodsType string
	var goodsBrand string
	var goodsModel string
	var goodsDescription string

	requestStatus = "New"

	serviceProviderID := ""
	price := ""
	signature := ""
	receipt := ""

	// when a new request is created, there will be no price, signature and receipt info.
	// var price            int32
	// var signature        string
	// var receipt          string
	//
	// switch case statement below
	// case "Price":
	//     priceInt, err := strconv.ParseInt(args[i + 1], 10, 32)
	//     if err != nil {
	//         return nil, errors.New("Unsupported Parameter " + args[i] " : " + args[i + 1] + ", must be convertable to int32")
	//     }
	//     price = int32(priceInt)
	// case "Signature":
	//     signature = args[i + 1]
	// case "Receipt":
	//     receipt = args[i + 1]

	for i := 0; i < len(args); i = i + 2 {
		switch args[i] {

		case "LandlordID":
			landlordID = args[i+1]
		case "TenantID":
			tenantID = args[i+1]
		case "ServiceProviderID":
			serviceProviderID = args[i+1]
		case "Address":
			address = args[i+1]
		case "RequestStatus":
			requestStatus = args[i+1]
		case "GoodsType":
			goodsType = args[i+1]
		case "GoodsBrand":
			goodsBrand = args[i+1]
		case "GoodsModel":
			goodsModel = args[i+1]
		case "GoodsDescription":
			goodsDescription = args[i+1]
		default:
			return nil, errors.New("Unsupported parameter " + args[i])
		}
	}

	requestID := util.GenerateUUID()

	var columns []*shim.Column
	col1 := shim.Column{Value: &shim.Column_String_{String_: requestID}}
	col2 := shim.Column{Value: &shim.Column_String_{String_: landlordID}}
	col3 := shim.Column{Value: &shim.Column_String_{String_: tenantID}}
	col4 := shim.Column{Value: &shim.Column_String_{String_: serviceProviderID}}
	col5 := shim.Column{Value: &shim.Column_String_{String_: address}}
	col6 := shim.Column{Value: &shim.Column_String_{String_: requestStatus}}
	col7 := shim.Column{Value: &shim.Column_String_{String_: goodsType}}
	col8 := shim.Column{Value: &shim.Column_String_{String_: goodsBrand}}
	col9 := shim.Column{Value: &shim.Column_String_{String_: goodsModel}}
	col10 := shim.Column{Value: &shim.Column_String_{String_: goodsDescription}}
	col11 := shim.Column{Value: &shim.Column_String_{String_: price}}
	col12 := shim.Column{Value: &shim.Column_String_{String_: signature}}
	col13 := shim.Column{Value: &shim.Column_String_{String_: receipt}}

	columns = append(columns, &col1)
	columns = append(columns, &col2)
	columns = append(columns, &col3)
	columns = append(columns, &col4)
	columns = append(columns, &col5)
	columns = append(columns, &col6)
	columns = append(columns, &col7)
	columns = append(columns, &col8)
	columns = append(columns, &col9)
	columns = append(columns, &col10)
	columns = append(columns, &col11)
	columns = append(columns, &col12)
	columns = append(columns, &col13)

	row := shim.Row{Columns: columns}
	ok, err := stub.InsertRow("Request", row)

	if err != nil {
		return nil, fmt.Errorf("Insert a row in Request table operation failed. %s", err)
	}

	if !ok {
		return nil, errors.New("Insert a row in Request table operation failed. Row with given key already exists")
	}

	fmt.Println("RequestID inserted is [" + requestID + "] ...")
	return []byte(requestID), err
}

//==============================================================================================================================
// deleteRowTenancyContract - Delete a row from the table Tenancy_Contract
//    input - key:value
//            "TenancyContractID" : "xxxxxx" <uuid>
//==============================================================================================================================
func (t *SimpleChaincode) deleteRowTenancyContract(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments, expecting 2")
	}

	if args[0] != "TenancyContractID" {
		return nil, errors.New("Incorrect argument name, expecting \"TenancyContractID\"")
	}

	tenancyContractID := args[1]

	err := stub.DeleteRow(
		"Tenancy_Contract",
		[]shim.Column{shim.Column{Value: &shim.Column_String_{String_: tenancyContractID}}},
	)

	if err != nil {
		return nil, fmt.Errorf("Deleting a row in Tenancy_Contract table operation failed. %s", err)
	}

	return nil, nil
}

//==============================================================================================================================
// deleteRowInsuranceContract - Delete a row from the table Insurance_Contract
//    input - key:value
//            "InsuranceContractID" : "xxxxxx" <uuid>
//==============================================================================================================================
func (t *SimpleChaincode) deleteRowInsuranceContract(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments, expecting 2 ")
	}

	if args[0] != "InsuranceContractID" {
		return nil, errors.New("Incorrect argument name, expecting \"InsuranceContractID\" ")
	}

	insuranceContractID := args[1]

	err := stub.DeleteRow(
		"Insurance_Contract",
		[]shim.Column{shim.Column{Value: &shim.Column_String_{String_: insuranceContractID}}},
	)

	if err != nil {
		return nil, fmt.Errorf("Deleting a row in Insurance_Contract table operation failed. %s", err)
	}

	return nil, nil
}

//==============================================================================================================================
// deleteRowRequest- Delete a row from the table Request
//    input - key:value
//            "requestID" : "xxxxxx" <uuid>
//==============================================================================================================================
func (t *SimpleChaincode) deleteRowRequest(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments, expecting 2 ")
	}

	if args[0] != "RequestID" {
		return nil, errors.New("Incorrect augement name, expecting \"RequestID\" : \"xxxxxx\" ")
	}

	requestID := args[1]

	var columns []shim.Column
	keyCol1 := shim.Column{Value: &shim.Column_String_{String_: requestID}}

	columns = append(columns, keyCol1)

	err := stub.DeleteRow("Request", columns)

	if err != nil {
		return nil, fmt.Errorf("Deleting a row in Request table operation failed. %s", err)
	}

	return nil, nil
}

//==============================================================================================================================
// getTenancyContract - Query a row in the table Tenancy_Contract
//    input  - key:value
//						 "TenancyContractID"  : "T000100"
//		output - row
//==============================================================================================================================
func (t *SimpleChaincode) getTenancyContract(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2, \"TenancyContractID\" : \"xxxxxx\"")
	}

	if args[0] != "TenancyContractID" {
		return nil, errors.New("Unsupoprted query arguments [" + args[0] + "]")
	}

	fmt.Println("Start to query a TenancyContract ...")

	tenancyContractID := args[1]

	var columns []shim.Column
	keyCol1 := shim.Column{Value: &shim.Column_String_{String_: tenancyContractID}}

	columns = append(columns, keyCol1)

	row, err := stub.GetRow("Tenancy_Contract", columns)

	if err != nil {
		return nil, fmt.Errorf("Query a Tenancy Contract (TenancyContractID = %s) failed ", tenancyContractID)
	}

	if len(row.Columns) == 0 {
		return nil, errors.New("Tenancy Contract was NOT found ")
	}

	fmt.Println("Query a Tenancy Contract (TenancyContractID = " + tenancyContractID + " successfully ...")

	// Convert to the structure TenancyContract, the returns would be key1:value1,key2:value2,key3:value3, ...
	tContract := &TenancyContract{
		row.Columns[0].GetString_(),
		row.Columns[1].GetString_(),
		row.Columns[2].GetString_(),
		row.Columns[3].GetString_(),
	}

	jsonRow, err := json.Marshal(tContract)

	if err != nil {
		return nil, errors.New("getTenancyContract() json marshal error")
	}

	fmt.Println("End to query a TenancyContract ...")
	return jsonRow, nil
}

//==============================================================================================================================
// getTenancyContractsByID - Query rows in the table Tenancy_Contract owned by a specific ID
//    input  - key:value
//						 "LandlordID"  : "L000001"
//						 or
//						 "TenantID"    : "T000100"
//		output - rows
//==============================================================================================================================
func (t *SimpleChaincode) getTenancyContractsByID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2, \"LandlordID\" or \"TenantID\"")
	}

	fmt.Println("Start to query TenancyContractsByID ...")

	var colName string
	var colValue string

	colName = args[0]
	colValue = args[1]

	if colName != "LandlordID" && colName != "TenantID" {
		return nil, errors.New("Unsupoprted query arguments + [" + colName + "]")
	}

	var emptyArgs []string

	jsonTContracts, err := t.getAllTenancyContracts(stub, emptyArgs)

	if err != nil {
		return nil, fmt.Errorf("Query TenancyContractsByID failed. %s", err)
	}

	var allTContracts []TenancyContract
	unMarError := json.Unmarshal(jsonTContracts, &allTContracts)

	if unMarError != nil {
		return nil, fmt.Errorf("Error unmarshalling TenancyContracts: %s", err)
	}

	var returnTContracts []TenancyContract
	for i := 0; i < len(allTContracts); i = i + 1 {
		tContract := allTContracts[i]

		if colName == "LandlordID" {
			if tContract.LandlordID == colValue {
				returnTContracts = append(returnTContracts, tContract)
			}
		} else if colName == "TenantID" {
			if tContract.TenantID == colValue {
				returnTContracts = append(returnTContracts, tContract)
			}
		}
	}

	jsonReturnTContracts, jsonErr := json.Marshal(returnTContracts)
	if jsonErr != nil {
		return nil, fmt.Errorf("Query TenancyContractsIDs failed. Error marshaling JSON: %s", jsonErr)
	}

	fmt.Println("End to query TenancyContractsByID ...")

	return jsonReturnTContracts, nil
}

//==============================================================================================================================
// getAllTenancyContracts - Query all rows in the table Insurance_Contract
//		input  - nil
//
//		output - rows in JSON
//==============================================================================================================================
func (t *SimpleChaincode) getAllTenancyContracts(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) > 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting 0, no input needed ")
	}

	fmt.Println("Start to query All TenancyContracts ...")

	var columns []shim.Column

	rowChannel, err := stub.GetRows("Tenancy_Contract", columns)

	if err != nil {
		return nil, fmt.Errorf("Query all Tenancy Contracts failed. %s", err)
	}

	var rows []shim.Row
	for {
		select {
		case row, ok := <-rowChannel:
			if !ok {
				rowChannel = nil
			} else {
				rows = append(rows, row)
			}
		}
		if rowChannel == nil {
			break
		}
	}

	var tContracts []TenancyContract

	for i := 0; i < len(rows); i = i + 1 {
		oneRow := rows[i]
		tContract := &TenancyContract{
			oneRow.Columns[0].GetString_(),
			oneRow.Columns[1].GetString_(),
			oneRow.Columns[2].GetString_(),
			oneRow.Columns[3].GetString_(),
		}
		tContracts = append(tContracts, *tContract)
	}

	jsonRows, err := json.Marshal(tContracts)
	if err != nil {
		return nil, fmt.Errorf("Query all Tenancy Contracts failed. Error marshaling JSON: %s", err)
	}

	fmt.Println("End to query All TenancyContracts ...")

	return jsonRows, nil
}

//==============================================================================================================================
// getInsuranceContract - Query a row in the table Insurance_Contract
//    input  - key:value
//						 "InsuranceContractID"  : "I000100"
//		output - row
//==============================================================================================================================
func (t *SimpleChaincode) getInsuranceContract(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2, \"InsuranceContractID\" : \"xxxxxx\"")
	}

	if args[0] != "InsuranceContractID" {
		return nil, errors.New("Unsupoprted query arguments [" + args[0] + "]")
	}

	fmt.Println("Start to query a InsuranceContract ...")

	insuranceContractID := args[1]

	var columns []shim.Column
	keyCol1 := shim.Column{Value: &shim.Column_String_{String_: insuranceContractID}}

	columns = append(columns, keyCol1)

	row, err := stub.GetRow("Insurance_Contract", columns)

	if err != nil {
		return nil, fmt.Errorf("Query an Insurance Contract (InsuranceContractID = %s) failed", insuranceContractID)
	}

	if len(row.Columns) == 0 {
		return nil, errors.New("Insurance Contract was NOT found")
	}

	// Convert to the structure InsuranceContract, the returns would be key1:value1,key2:value2,key3:value3, ...
	iContract := &InsuranceContract{
		row.Columns[0].GetString_(),
		row.Columns[1].GetString_(),
		row.Columns[2].GetString_(),
		row.Columns[3].GetString_(),
		row.Columns[4].GetInt32(),
		row.Columns[5].GetString_(),
		row.Columns[6].GetString_(),
	}

	returnIContract, err := json.Marshal(iContract)

	if err != nil {
		return nil, errors.New("GetInsuranceContract() json marshal error")
	}

	fmt.Println("End to query a InsuranceContract ...")

	return returnIContract, nil
}

//==============================================================================================================================
// getInsuranceContractID - Query rows in the table Tenancy_Contract owned by a specific LandlordID
//    input  - key:value
//						 "LandlordID"  : "L000100"
//		output - row
//==============================================================================================================================
func (t *SimpleChaincode) getInsuranceContractsByID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2, \"LandlordID\" : \"Lxxxxxx\"")
	}

	fmt.Println("Start to query InsuranceContractsByID ...")

	var colName string
	var colValue string

	colName = args[0]
	colValue = args[1]

	if colName != "LandlordID" {
		return nil, errors.New("unsupoprted query arguments [" + colName + "]")
	}

	var emptyArgs []string

	jsonAllRows, err := t.getAllInsuranceContracts(stub, emptyArgs)

	if err != nil {
		return nil, fmt.Errorf("Query InsuranceContractsByID failed. %s", err)
	}

	var allIContracts []InsuranceContract
	unMarError := json.Unmarshal(jsonAllRows, &allIContracts)

	if unMarError != nil {
		return nil, fmt.Errorf("Error unmarshalling rows: %s", err)
	}

	var returnIContracts []InsuranceContract

	for i := 0; i < len(allIContracts); i = i + 1 {
		iContract := allIContracts[i]
		if iContract.LandlordID == colValue {
			returnIContracts = append(returnIContracts, iContract)
		}
	}

	jsonReturnIContracts, jsonErr := json.Marshal(returnIContracts)
	if jsonErr != nil {
		return nil, fmt.Errorf("Query InsuranceContractsByID failed. Error marshaling JSON: %s", jsonErr)
	}

	fmt.Println("End to query InsuranceContractsByID ...")

	return jsonReturnIContracts, nil

	// Return rows but not structures, code below ...
	// ------------------------------------------------------------
	// var allRows []shim.Row
	// unMarError := json.Unmarshal(jsonRows, &allRows)

	// if unMarError != nil {
	// 	return nil, fmt.Errorf("Error unmarshalling rows: %s", err)
	// }

	// var returnRows []shim.Row
	// for i := 0; i < len(allRows); i = i + 1 {
	// 	row := allRows[i]
	// 	if row.Columns[colIndex].GetString_() == colValue {
	// 		returnRows = append(returnRows, row)
	// 	}
	// }

	// jsonReturnRows, jsonErr := json.Marshal(returnRows)
	// if jsonErr != nil {
	// 	return nil, fmt.Errorf("Query InsuranceContractsByID failed. Error marshaling JSON: %s", jsonErr)
	// }

	// fmt.Printf("Query InsuranceContractsByID successfully ")

	// return jsonReturnRows, nil
}

//==============================================================================================================================
// getAllInsuranceContracts - Query all rows in the table Insurance_Contract
//		input  - nil
//
//		output - rows in JSON
//==============================================================================================================================
func (t *SimpleChaincode) getAllInsuranceContracts(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) > 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting 0, no input needed ")
	}

	fmt.Println("Start to query all InsuranceContracts ...")
	var columns []shim.Column

	rowChannel, err := stub.GetRows("Insurance_Contract", columns)

	if err != nil {
		return nil, fmt.Errorf("Query all Insurance Contracts failed. %s", err)
	}

	var rows []shim.Row
	for {
		select {
		case row, ok := <-rowChannel:
			if !ok {
				rowChannel = nil
			} else {
				rows = append(rows, row)
			}
		}
		if rowChannel == nil {
			break
		}
	}

	var iContracts []InsuranceContract

	for i := 0; i < len(rows); i = i + 1 {
		oneRow := rows[i]
		iContract := &InsuranceContract{
			oneRow.Columns[0].GetString_(),
			oneRow.Columns[1].GetString_(),
			oneRow.Columns[2].GetString_(),
			oneRow.Columns[3].GetString_(),
			oneRow.Columns[4].GetInt32(),
			oneRow.Columns[5].GetString_(),
			oneRow.Columns[6].GetString_(),
		}
		iContracts = append(iContracts, *iContract)
	}

	jsonAllIContracts, err := json.Marshal(iContracts)
	if err != nil {
		return nil, fmt.Errorf("Query all Insurance Contracts failed. Error marshaling JSON: %s", err)
	}

	fmt.Println("End to query all InsuranceContracts ...")

	return jsonAllIContracts, nil
}

//==============================================================================================================================
// getRequest - Query a row in the table Request
//    input  - key:value
//						 "RequestID"  : "R000001"
//		output - rows
//==============================================================================================================================
func (t *SimpleChaincode) getRequest(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2, \"RequestID\" : \"xxxxxx\"")
	}

	if args[0] != "RequestID" {
		return nil, errors.New("Unsupoprted query arguments [" + args[0] + "]")
	}

	fmt.Println("Start to query a Request ...")

	requestID := args[1]

	var columns []shim.Column
	keyCol1 := shim.Column{Value: &shim.Column_String_{String_: requestID}}

	columns = append(columns, keyCol1)

	row, err := stub.GetRow("Request", columns)

	if err != nil {
		return nil, fmt.Errorf("Query one request (RequestID = %s) in the table Request failed", requestID)
	}

	if len(row.Columns) == 0 {
		return nil, errors.New("Request was NOT found")
	}

	fmt.Println("Query one request (RequestID = " + requestID + ") in the table Request successfully...")

	// Convert to the structure Request, the returns would be key1:value1,key2:value2,key3:value3, ...
	request := &Request{
		row.Columns[0].GetString_(),
		row.Columns[1].GetString_(),
		row.Columns[2].GetString_(),
		row.Columns[3].GetString_(),
		row.Columns[4].GetString_(),
		row.Columns[5].GetString_(),
		row.Columns[6].GetString_(),
		row.Columns[7].GetString_(),
		row.Columns[8].GetString_(),
		row.Columns[9].GetString_(),
		row.Columns[10].GetString_(),
		row.Columns[11].GetString_(),
		row.Columns[12].GetString_(),
	}

	returnRequest, err := json.Marshal(request)

	if err != nil {
		return nil, errors.New("getRequest() json marshal error")
	}

	fmt.Println("End to query a Request ...")

	return returnRequest, nil
}

//==============================================================================================================================
// getRequestsByID - Query rows in the table Request owned by a specific ID
//    input  - key:value
//						 "LandlordID"  : "L000001"
//						 or
//						 "TenantID"    : "T000100"
//						 or
//						 "ServiceProviderID"    : "SP000010"
//		output - rows
//==============================================================================================================================
func (t *SimpleChaincode) getRequestsByID(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2, \"LandlordID\" : \"Lxxxxxx\" or \"TenantID\" : \"Txxxxxx\" or \"ServiceProviderID\" : \"SPxxxxxx\" ")
	}

	fmt.Println("Start to query RequestsByID ...")

	var colName string
	var colValue string

	colName = args[0]
	colValue = args[1]

	if colName != "LandlordID" && colName != "TenantID" && colName != "ServiceProviderID" {
		return nil, errors.New("Unsupoprted query arguments [" + colName + "]")
	}

	var emptyArgs []string

	jsonRows, err := t.getAllRequests(stub, emptyArgs)

	if err != nil {
		return nil, fmt.Errorf("Query RequestsByID failed. %s", err)
	}

	var allRequests []Request
	unMarError := json.Unmarshal(jsonRows, &allRequests)

	if unMarError != nil {
		return nil, fmt.Errorf("Error unmarshalling rows: %s", err)
	}

	var returnRequests []Request

	for i := 0; i < len(allRequests); i = i + 1 {
		request := allRequests[i]
		if colName == "LandlordID" {
			if request.LandlordID == colValue {
				returnRequests = append(returnRequests, request)
			}
		} else if colName == "TenantID" {
			if request.TenantID == colValue {
				returnRequests = append(returnRequests, request)
			}
		} else if colName == "ServiceProviderID" {
			if request.ServiceProviderID == colValue {
				returnRequests = append(returnRequests, request)
			}
		}
	}

	jsonReturnRequests, jsonErr := json.Marshal(returnRequests)
	if jsonErr != nil {
		return nil, fmt.Errorf("Query RequestsByID failed. Error marshaling JSON: %s", jsonErr)
	}

	fmt.Println("End to query RequestsByID ...")

	return jsonReturnRequests, nil

	// var returnRows []shim.Row
	// for i := 0; i < len(allRows); i = i + 1 {
	// 	row := allRows[i]
	// 	if row.Columns[colIndex].GetString_() == colValue {
	// 		returnRows = append(returnRows, row)
	// 	}
	// }

	// jsonReturnRows, jsonErr := json.Marshal(returnRows)
	// if jsonErr != nil {
	// 	return nil, fmt.Errorf("Query RequestsByID failed. Error marshaling JSON: %s", jsonErr)
	// }

	// fmt.Printf("Query RequestsByID successfully ")

	// return jsonReturnRows, nil
}

//==============================================================================================================================
// getAllRequests - Query all rows in the table Request
//		input  - nil
//
//		output - rows in JSON
//==============================================================================================================================
func (t *SimpleChaincode) getAllRequests(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args) > 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting 0, no input needed")
	}

	fmt.Println("Start to query all requests.")

	var columns []shim.Column

	rowChannel, err := stub.GetRows("Request", columns)

	if err != nil {
		return nil, fmt.Errorf("Query all requests failed. %s", err)
	}

	var rows []shim.Row
	for {
		select {
		case row, ok := <-rowChannel:
			if !ok {
				rowChannel = nil
			} else {
				rows = append(rows, row)
			}
		}
		if rowChannel == nil {
			break
		}
	}

	var allRequests []Request

	for i := 0; i < len(rows); i = i + 1 {
		oneRow := rows[i]
		request := &Request{
			oneRow.Columns[0].GetString_(),
			oneRow.Columns[1].GetString_(),
			oneRow.Columns[2].GetString_(),
			oneRow.Columns[3].GetString_(),
			oneRow.Columns[4].GetString_(),
			oneRow.Columns[5].GetString_(),
			oneRow.Columns[6].GetString_(),
			oneRow.Columns[7].GetString_(),
			oneRow.Columns[8].GetString_(),
			oneRow.Columns[9].GetString_(),
			oneRow.Columns[10].GetString_(),
			oneRow.Columns[11].GetString_(),
			oneRow.Columns[12].GetString_(),
		}
		allRequests = append(allRequests, *request)
	}

	jsonAllRequests, err := json.Marshal(allRequests)
	if err != nil {
		return nil, fmt.Errorf("Query all requests operation failed. Error marshaling JSON: %s", err)
	}

	fmt.Println("End to query all requests.")

	return jsonAllRequests, nil
}

//==============================================================================================================================
// updateRowRequest - Update a row in the table Request
//		input  - key:value
//				 "RequestID" : "R000001"
//		output - nil
//==============================================================================================================================
func (t *SimpleChaincode) updateRowRequest(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	if len(args)%2 != 0 || len(args) == 2 {
		return nil, errors.New("Incorrect number of arguments")
	}

	if args[0] != "RequestID" {
		return nil, errors.New("The first argument should be RequestID")
	}

	fmt.Println("Start to update a request ...")
	requestID := args[1]

	row, err := t.extractRowRequest(stub, requestID)

	if err != nil {
		return nil, fmt.Errorf("UpdateRequest failed. %s", err)
	}

	for i := 2; i < len(args); i = i + 2 {
		colName := args[i]

		switch colName {

		case "LandlordID":
			cellValue := args[i+1]
			row.Columns[1] = &shim.Column{Value: &shim.Column_String_{String_: cellValue}}
		case "TenantID":
			cellValue := args[i+1]
			row.Columns[2] = &shim.Column{Value: &shim.Column_String_{String_: cellValue}}
		case "ServiceProviderID":
			cellValue := args[i+1]
			row.Columns[3] = &shim.Column{Value: &shim.Column_String_{String_: cellValue}}
		case "Address":
			cellValue := args[i+1]
			row.Columns[4] = &shim.Column{Value: &shim.Column_String_{String_: cellValue}}
		case "RequestStatus":
			cellValue := args[i+1]
			row.Columns[5] = &shim.Column{Value: &shim.Column_String_{String_: cellValue}}
		case "GoodsType":
			cellValue := args[i+1]
			row.Columns[6] = &shim.Column{Value: &shim.Column_String_{String_: cellValue}}
		case "GoodsBrand":
			cellValue := args[i+1]
			row.Columns[7] = &shim.Column{Value: &shim.Column_String_{String_: cellValue}}
		case "GoodsModel":
			cellValue := args[i+1]
			row.Columns[8] = &shim.Column{Value: &shim.Column_String_{String_: cellValue}}
		case "GoodsDescription":
			cellValue := args[i+1]
			row.Columns[9] = &shim.Column{Value: &shim.Column_String_{String_: cellValue}}
		case "Price":
			cellValue := args[i+1]
			row.Columns[10] = &shim.Column{Value: &shim.Column_String_{String_: cellValue}}
			// cellValueInt, strConErr := strconv.Atoi(cellValue)
			// if strConErr != nil {
			// 	return nil, fmt.Errorf("String conversion to Int failed. %s", strConErr)
			// }
			// row.Columns[10] = &shim.Column{Value: &shim.Column_Int32{Int32: int32(cellValueInt)}}
		case "Signature":
			cellValue := args[i+1]
			row.Columns[11] = &shim.Column{Value: &shim.Column_String_{String_: cellValue}}
		case "Receipt":
			cellValue := args[i+1]
			row.Columns[12] = &shim.Column{Value: &shim.Column_String_{String_: cellValue}}

		default:
			return nil, errors.New("Unsupported Parameter " + args[i])
		}
	}

	ok, err := stub.ReplaceRow("Request", row)

	if !ok && err == nil {
		return nil, errors.New("Failed to replace row with RequestID " + requestID)
	}

	fmt.Println("End to update a request ...")
	return nil, nil
}

//==============================================================================================================================
// getRowinRequest - Internal function, used in updateRowRequest a row in the table Request, return a Row but NOT []byte
//		input  - key:value
//				 "RequestID" : "R000001"
//		output - row
//==============================================================================================================================
func (t *SimpleChaincode) extractRowRequest(stub shim.ChaincodeStubInterface, requestID string) (shim.Row, error) {

	var columns []shim.Column

	keyCol1 := shim.Column{Value: &shim.Column_String_{String_: requestID}}
	columns = append(columns, keyCol1)

	return stub.GetRow("Request", columns)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
