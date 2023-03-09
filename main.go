package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"fabric-go-sdk-base/sdkInit"
	"time"
)

// 定义链码名称与版本
const (
	cc_name    = "simplecc"
	cc_version = "1.0.0"
)

var App sdkInit.Application

func main() {
	err := os.Setenv("DISCOVERY_AS_LOCALHOST", "true")
	if err != nil {
		fmt.Printf("Error setting DISCOVERY_AS_LOCALHOST environemnt variable: %v", err)
		return
	}
	// 初始化组织信息
	orgs := []*sdkInit.OrgInfo{
		{
			OrgAdminUser:  "Admin",
			OrgName:       "Org1",
			OrgMspId:      "Org1MSP",
			OrgUser:       "User1",
			OrgPeerNum:    2,
			OrgAnchorFile: "/home/usr/GolandProjects/fabric-go-sdk-base/fixtures/channel-artifacts/Org1MSPanchors.tx",
		},
		// {
		// 	OrgAdminUser:  "Admin",
		// 	OrgName:       "Org2",
		// 	OrgMspId:      "Org2MSP",
		// 	OrgUser:       "User1",
		// 	OrgPeerNum:    1,
		// 	OrgAnchorFile: "/home/usr/GolandProjects/fabric-go-sdk-base/fixtures/channel-artifacts/Org2MSPanchors.tx",
		// },
	}

	// 初始化sdk相关信息
	info := sdkInit.SdkEnvInfo{
		ChannelID:        "mychannel",
		ChannelConfig:    "/home/usr/GolandProjects/fabric-go-sdk-base/fixtures/channel-artifacts/channel.tx",
		Orgs:             orgs,
		OrdererAdminUser: "Admin",
		OrdererOrgName:   "OrdererOrg",
		OrdererEndpoint:  "orderer.example.com",
		ChaincodeID:      cc_name,
		ChaincodePath:    "/home/usr/GolandProjects/fabric-go-sdk-base/chaincode/",
		ChaincodeVersion: cc_version,
	}

	// 调用setup方法将sdk初始化
	sdk, err := sdkInit.Setup("config.yaml", &info)
	if err != nil {
		fmt.Println(">> SDK setup error:", err)
		os.Exit(-1)
	}

	// 调用CreateAndJoinChannel方法，创建并加入通道
	if err := sdkInit.CreateAndJoinChannel(&info); err != nil {
		fmt.Println(">> Create channel and join error:", err)
		os.Exit(-1)
	}

	// 调用CreateCCLifecycle方法实现链码生命周期
	if err := sdkInit.CreateCCLifecycle(&info, 1, false, sdk); err != nil {
		fmt.Printf(">> create chaincode lifecycle error: %v\n", err)
		os.Exit(-1)
	}

	// invoke chaincode set status
	fmt.Println(">> 通过链码外部服务设置链码状态......")

	if err := info.InitService(info.ChaincodeID, info.ChannelID, info.Orgs[0], sdk); err != nil {

		fmt.Println("InitService successful")
		os.Exit(-1)
	}

	App = sdkInit.Application{
		SdkEnvInfo: &info,
	}
	fmt.Println(">> 设置链码状态完成")

	defer info.EvClient.Unregister(sdkInit.BlockListener(info.EvClient))
	defer info.EvClient.Unregister(sdkInit.ChainCodeEventListener(info.EvClient, info.ChaincodeID))

	fmt.Println("==========command format==========")
	fmt.Println("input data: set [key] [value]")
	fmt.Println("query data: get [key]")
	fmt.Println("exit: exit")
	fmt.Println("Please input command: ")
	for {
		inputReader := bufio.NewReader(os.Stdin)
		input, _, err := inputReader.ReadLine()
		if err != nil {
			fmt.Println("input command error:", err)
		}
		cmd := strings.Split(string(input), " ")
		switch cmd[0] {
		case "set":
			ret, err := App.Set(cmd)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("<--- 添加信息　--->：", ret)
		case "get":
			response, err := App.Get(cmd)
			fmt.Println(cmd)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("<--- 查询信息　--->：", response)
		case "exit":
			break
		}
	}
	fmt.Println("==========program end==========")
	time.Sleep(time.Second * 10)
}
