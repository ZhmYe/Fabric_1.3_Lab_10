package smallbank

import (
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type SmallBankChainCode struct{}

func (t *SmallBankChainCode) accountInfoPrint(stub shim.ChaincodeStubInterface, account string) (int, int) {
	savings, err := stub.GetState(account + "_savings")
	if err != nil {
		fmt.Println("Failed to get savings state")
		return -1, -1
	}
	savingsValue, _ := strconv.Atoi(string(savings))
	checkings, err := stub.GetState(account + "_checkings")
	if err != nil {
		fmt.Println("Failed to get checkings state")
		return -1, -1
	}
	checkingsValue, _ := strconv.Atoi(string(checkings))
	fmt.Println("Account ", account)
	fmt.Println("	Savings: ", savingsValue)
	fmt.Println("	Checkings: ", checkingsValue)
	return savingsValue, checkingsValue
}

// Init 初始化Smallbank，这里初始化两个账户的savings和checkings
// 所有账户的savings用"账户_savings"表示，checkings用"账户_checkings"表示
// 初始化需要6个参数, 账户1 savings余额 checkings余额 账户2 savings余额 checkings余额
func (t *SmallBankChainCode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	var err error
	_, args := stub.GetFunctionAndParameters()
	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}

	// Initialize the chaincode
	accountA := args[0]
	accountASavings, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	accountACheckings, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	// 将账户A的相关信息存入
	// 存入savings
	err = stub.PutState(accountA+"_savings", []byte(strconv.Itoa(accountASavings)))
	// 存入checkings
	err = stub.PutState(accountA+"_checkings", []byte(strconv.Itoa(accountACheckings)))

	accountB := args[0]
	accountBSavings, err := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	accountBCheckings, err := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	// 存入savings
	err = stub.PutState(accountB+"_savings", []byte(strconv.Itoa(accountBSavings)))
	if err != nil {
		return shim.Error(err.Error())
	}
	// 存入checkings
	err = stub.PutState(accountB+"_checkings", []byte(strconv.Itoa(accountBCheckings)))
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("Init Account ", accountA, " ", accountB)
	t.accountInfoPrint(stub, accountA)
	t.accountInfoPrint(stub, accountB)
	return shim.Success(nil)
}
func (t *SmallBankChainCode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "create" {
		// 创建一个新的账户
		return t.create(stub, args)
	} else if function == "transactSavings" {
		// 向储蓄账户增加一定余额
		return t.transactSavings(stub, args)
	} else if function == "depositChecking" {
		// 向支票账户增加一定余额
		return t.depositChecking(stub, args)
	} else if function == "sendPayment" {
		// 在两个支票账户间转账
		return t.sendPayment(stub, args)
	} else if function == "WriteCheck" {
		// 减少一个支票账户
		return t.writeCheck(stub, args)
	} else if function == "Amalgamate" {
		// 将储蓄账户的资金全部转到支票账户
		return t.amalgamate(stub, args)
	} else if function == "query" {
		// 读取一个用户的支票账户以及储蓄账户
		return t.query(stub, args)
	}

	return shim.Error("Invalid invoke function name.")
}

// create 创建一个账户，根据传入的参数初始化saving、checking
// args: [账户, saving, checking]
func (t *SmallBankChainCode) create(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	account := args[0]

	// todo
	savings, _ := strconv.Atoi(args[1])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	checkings, _ := strconv.Atoi(args[2])
	if err != nil {
		return shim.Error("Expecting integer value for asset holding")
	}
	err = stub.PutState(account+"_savings", []byte(strconv.Itoa(savings)))
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(account+"_checkings", []byte(strconv.Itoa(checkings)))
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("Create Account ", account)
	t.accountInfoPrint(stub, account)
	return shim.Success(nil)
}

// 向储蓄账户增加一定余额
// 参数: [账户，金额]， 这里只需要账户, saving由"账户_savings"拼接得到
func (t *SmallBankChainCode) transactSavings(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	account := args[0]

	// todo
	amount, _ := strconv.Atoi(args[1])
	savings, err := stub.GetState(account + "_savings")
	if err != nil {
		fmt.Println("Failed to get savings state")
		return shim.Error("No such Account Saving!!!")
	}
	savingsValue, _ := strconv.Atoi(savings)
	savingsValue += amount
	err = stub.PutState(account+"_savings", []byte(strconv.Itoa(savingsValue)))
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("Update Account ", account)
	t.accountInfoPrint(stub, account)
	return shim.Success(nil)
}

// depositChecking 向支票账户增加一定余额
// 参数: [账户，金额]， 这里只需要账户, saving由"账户_savings"拼接得到
func (t *SmallBankChainCode) depositChecking(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	account := args[0]

	// todo
	amount, _ := strconv.Atoi(args[1])
	checkings, err := stub.GetState(account + "_savings")
	if err != nil {
		fmt.Println("Failed to get savings state")
		return shim.Error("No such Account Checking!!!")
	}
	checkingsValue, _ := strconv.Atoi(checkings)
	checkingsValue += amount
	err = stub.PutState(account+"_savings", []byte(strconv.Itoa(checkingsValue)))
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("Update Account ", account)
	t.accountInfoPrint(stub, account)
	return shim.Success(nil)
}

// sendPayment 在两个支票账户间转账
// 参数: [accountA, accountB, amount], A向B转账amount(在checkings里)
func (t *SmallBankChainCode) sendPayment(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	accountA, accountB := args[0], args[1]

	// todo
	amount, _ := strconv.Atoi(args[2])
	checkingsA, err := stub.GetState(accountA + "_savings")
	if err != nil {
		fmt.Println("Failed to get savings state")
		return shim.Error("No such Account Checking!!!")
	}
	checkingsValueA, _ := strconv.Atoi(checkingsA)
	if checkingsValueA < amount {
		return shim.Error("Account ", accountA, " don't have enough checkings!!!")
	}
	checkingsB, err := stub.GetState(accountA + "_savings")
	if err != nil {
		fmt.Println("Failed to get savings state")
		return shim.Error("No such Account Checking!!!")
	}
	checkingsValueB, _ := strconv.Atoi(checkingsB)
	checkingsValueA -= amount
	checkingsValueB += amount
	err = stub.PutState(accountA+"_checkings", []byte(strconv.Itoa(checkingsValueA)))
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(accountB+"_checkings", []byte(strconv.Itoa(checkingsValueB)))
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("Update Account ", accountA, " ", accountB)
	t.accountInfoPrint(stub, accountA)
	t.accountInfoPrint(stub, accountB)
	return shim.Success(nil)
}

// 减少一个支票账户
// 参数: [账户，金额]，扣除一笔钱
func (t *SmallBankChainCode) writeCheck(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}
	account := args[0]
	// todo
	amount, _ := strconv.Atoi(args[1])
	checkings, err := stub.GetState(account + "_savings")
	if err != nil {
		fmt.Println("Failed to get savings state")
		return shim.Error("No such Account Checking!!!")
	}
	checkingsValue, _ := strconv.Atoi(checkings)
	if checkingsValue < amount {
		return shim.Error("Account ", account, " don't have enough checkings!!!")
	}
	checkingsValue -= amount
	err = stub.PutState(account+"_checkings", []byte(strconv.Itoa(checkingsValue)))
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("Update Account ", account)
	t.accountInfoPrint(stub, account)
	return shim.Success(nil)
}

// amalgamate 将储蓄账户的资金全部转到支票账户
// 参数: [account]
func (t *SmallBankChainCode) amalgamate(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	account := args[0]
	// todo
	savings, err := stub.GetState(account + "_savings")
	if err != nil {
		return shim.Error("Failed to get savings state")
	}
	savingsValue, _ := strconv.Atoi(string(savings))
	checkings, err := stub.GetStage(account + "_checkings")
	if err != nil {
		return shim.Error("Failed to get checkings state")
	}
	checkingsValue, _ := strconv.Atoi(string(checkings))
	checkingsValue += savingsValue

	err = stub.PutState(account+"_checkings", []byte(strconv.Itoa(checkingsValue)))
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(account+"_savings", []byte("0"))
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("Update Account ", account)
	t.accountInfoPrint(stub, account)
	return shim.Success(nil)
}

// query 查询一个账户的余额信息
// 参数: [account]
// 这里就不用写了，简单调用下上面的accountInfoPrint即可完成
func (t *SmallBankChainCode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}
	account := args[0]
	s, c := t.accountInfoPrint(stub, account)
	return shim.Success("Savings Amount: " + strconv.Itoa(s) + ", Checkings Amount: " + strconv.Itoa(c))
}
