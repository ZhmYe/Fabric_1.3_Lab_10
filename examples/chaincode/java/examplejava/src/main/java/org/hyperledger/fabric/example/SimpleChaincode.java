package org.hyperledger.fabric.example;

import java.util.List;

import com.google.protobuf.ByteString;
import io.netty.handler.ssl.OpenSsl;
import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.hyperledger.fabric.shim.ChaincodeBase;
import org.hyperledger.fabric.shim.ChaincodeStub;

import static java.nio.charset.StandardCharsets.UTF_8;

public class SimpleChaincode extends ChaincodeBase {

    private static Log _logger = LogFactory.getLog(SimpleChaincode.class);
    // 这里初始化合约状态，初始化了两个账户及其对应的余额
    @Override
    public Response init(ChaincodeStub stub) {
        try {
            _logger.info("Init java simple chaincode");
            String func = stub.getFunction();
            if (!func.equals("init")) {
                return newErrorResponse("function other than init is not supported");
            }
            List<String> args = stub.getParameters();
            if (args.size() != 4) {
                newErrorResponse("Incorrect number of arguments. Expecting 4");
            }
            // Initialize the chaincode
            String account1Key = args.get(0);
            int account1Value = Integer.parseInt(args.get(1));
            String account2Key = args.get(2);
            int account2Value = Integer.parseInt(args.get(3));

            _logger.info(String.format("account %s, value = %s; account %s, value %s", account1Key, account1Value, account2Key, account2Value));
            stub.putStringState(account1Key, args.get(1));
            stub.putStringState(account2Key, args.get(3));

            return newSuccessResponse();
        } catch (Throwable e) {
            return newErrorResponse(e);
        }
    }
    // public的invoke接口用于解析功能，根据命令行传入的要调用的函数调用相关的函数
    @Override
    public Response invoke(ChaincodeStub stub) {
        try {
            _logger.info("Invoke java simple chaincode");
            String func = stub.getFunction();
            List<String> params = stub.getParameters();
            if (func.equals("invoke")) {
                return invoke(stub, params);
            }
            if (func.equals("delete")) {
                return delete(stub, params);
            }
            if (func.equals("query")) {
                return query(stub, params);
            }
            return newErrorResponse("Invalid invoke function name. Expecting one of: [\"invoke\", \"delete\", \"query\"]");
        } catch (Throwable e) {
            return newErrorResponse(e);
        }
    }
    // invoke，这里的invoke可以理解是转账功能
    private Response invoke(ChaincodeStub stub, List<String> args) {
        // 该函数接收三个参数，["转账的人的key", "收账的人的key", "金额"]
        // 读取一个状态可以用stub.getStringState(string),这里的参数是key
        // 写入一个状态可以用stub.putStringState(string, string),这里的参数是key和amount, amount注意要转化为string
        if (args.size() != 3) {
            return newErrorResponse("Incorrect number of arguments. Expecting 3");
        }
        //todo
        // 这里记转账者、收账者为accountFromKey,accountToKey, 转账后二者的余额分别为accountFromValue,accountToValue
        return newSuccessResponse("invoke finished successfully", ByteString.copyFrom(accountFromKey + ": " + accountFromValue + " " + accountToKey + ": " + accountToValue, UTF_8).toByteArray());
    }

    // Deletes an entity from state, 删除一个状态
    private Response delete(ChaincodeStub stub, List<String> args) {
        // 该函数接收一个参数: 删除的key
        // 删除状态可以调用stub.delState(string)
        if (args.size() != 1) {
            return newErrorResponse("Incorrect number of arguments. Expecting 1");
        }
        //todo
        // 这里的返回值直接调用newSuccessResponse()即可
        return newSuccessResponse()
    }

    // query callback representing the query of a chaincode
    // 查询一个状态
    private Response query(ChaincodeStub stub, List<String> args) {
        if (args.size() != 1) {
            return newErrorResponse("Incorrect number of arguments. Expecting name of the person to query");
        }
        //todo
        // 这里记查询的状态余额为val
        return newSuccessResponse(val, ByteString.copyFrom(val, UTF_8).toByteArray());
    }

    public static void main(String[] args) {
        System.out.println("OpenSSL avaliable: " + OpenSsl.isAvailable());
        new SimpleChaincode().start(args);
    }

}

