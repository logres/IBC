package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv" 
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)


type SmartContract struct {
	contractapi.Contract
}


type StateMemory struct {
    Is_available bool `json:"Is_available"`
	Invoice bool `json:"Invoice"`
	Need_external_provider bool `json:"Need_external_provider"`
}

type InitParameters struct {
    Participant_1080bkg Participant `json:"Participant_1080bkg"`
	Participant_0sktaei Participant `json:"Participant_0sktaei"`
	Participant_1gcdqza Participant `json:"Participant_1gcdqza"`
}

type ContractInstance struct {
	// Incremental ID
	InstanceID string `json:"InstanceID"`
	// global Memory
	InstanceStateMemory StateMemory `json:"stateMemory"`
	// map type from string to Message、Gateway、ActionEvent
	InstanceElements map[string]interface{} `json:"InstanceElements"`
	// state of the instance
	InstanceState InstanceState `json:"InstanceState"`
}

type ElementState int

const (
	DISABLED = iota
	ENABLED
	WAITINGFORCONFIRMATION
	COMPLETED
)

type InstanceState int

const (
	TOBEREGISTERED = iota
	READY
)

type Participant struct {
	MSP          string            `json:"msp"`
	Attributes   map[string]string `json:"attributes"`
	IsMulti      bool              `json:"isMulti"`
	MultiMaximum int               `json:"multiMaximum"`
	MultiMinimum int               `json:"multiMinimum"`

	X509 string `json:"x509"`
}

type Message struct {
	MessageID            string       `json:"messageID"`
	SendParticipantID    string       `json:"sendMspID"`
	ReceiveParticipantID string       `json:"receiveMspID"`
	FireflyTranID        string       `json:"fireflyTranID"`
	MsgState             ElementState `json:"msgState"`
	Format               string       `json:"format"`
}

type Gateway struct {
	GatewayID    string       `json:"gatewayID"`
	GatewayState ElementState `json:"gatewayState"`
}

type ActionEvent struct {
	EventID    string       `json:"eventID"`
	EventState ElementState `json:"eventState"`
}

func (cc *SmartContract) CreateParticipant(ctx contractapi.TransactionContextInterface, instanceID string, participantID string, msp string, attributes map[string]string, IsMulti bool, MultiMaximum int, MultiMinimum int) (*Participant, error) {
	stub := ctx.GetStub()

	existingData, err := stub.GetState(instanceID)
	if err != nil {
		return nil, fmt.Errorf("获取状态数据时出错: %v", err)
	}
	if existingData == nil {
		return nil, fmt.Errorf("实例 %s 不存在", instanceID)
	}

	// 从现有实例中读取
	var instance ContractInstance
	err = json.Unmarshal(existingData, &instance)
	if err != nil {
		return nil, fmt.Errorf("反序列化实例数据时出错: %v", err)
	}

	// 创建参与者对象
	instance.InstanceElements[participantID] = &Participant{
		MSP:          msp,
		Attributes:   attributes,
		IsMulti:      IsMulti,
		MultiMaximum: MultiMaximum,
		MultiMinimum: MultiMinimum,
	}

	instanceJson, err := json.Marshal(instance)
	if err != nil {
		return nil, fmt.Errorf("序列化实例数据时出错: %v", err)
	}
	// 将参与者对象序列化为JSON字符串并保存在状态数据库中
	err = stub.PutState(instanceID, instanceJson)
	if err != nil {
		return nil, fmt.Errorf("保存实例数据时出错: %v", err)
	}

	returnParticipant, ok := instance.InstanceElements[participantID].(*Participant)
	if !ok {
		return nil, fmt.Errorf("无法将实例元素转换为Participant")
	}

	return returnParticipant, nil

}

func (cc *SmartContract) CreateMessage(ctx contractapi.TransactionContextInterface, instanceID string, messageID string, sendParticipantID string, receiveParticipantID string, fireflyTranID string, msgState ElementState, format string) (*Message, error) {
	stub := ctx.GetStub()

	existingData, err := stub.GetState(instanceID)
	if err != nil {
		return nil, fmt.Errorf("获取状态数据时出错: %v", err)
	}
	if existingData == nil {
		return nil, fmt.Errorf("实例 %s 不存在", instanceID)
	}

	// read from the existing instance
	var instance ContractInstance
	err = json.Unmarshal(existingData, &instance)
	if err != nil {
		return nil, fmt.Errorf("反序列化实例数据时出错: %v", err)
	}

	// 创建消息对象
	instance.InstanceElements[messageID] = &Message{
		MessageID:            messageID,
		SendParticipantID:    sendParticipantID,
		ReceiveParticipantID: receiveParticipantID,
		FireflyTranID:        fireflyTranID,
		MsgState:             msgState,
		Format:               format,
	}
	instanceJson, err := json.Marshal(instance)
	if err != nil {
		return nil, fmt.Errorf("序列化实例数据时出错: %v", err)
	}
	// 将消息对象序列化为JSON字符串并保存在状态数据库中
	err = stub.PutState(instanceID, instanceJson)
	if err != nil {
		return nil, fmt.Errorf("保存实例数据时出错: %v", err)
	}

	returnMessage, ok := instance.InstanceElements[messageID].(*Message)
	if !ok {
		return nil, fmt.Errorf("无法将实例元素转换为Message")
	}

	return returnMessage, nil
}

func (cc *SmartContract) CreateGateway(ctx contractapi.TransactionContextInterface, instanceID string, gatewayID string, gatewayState ElementState) (*Gateway, error) {
	stub := ctx.GetStub()

	// 检查是否存在具有相同ID的记录
	existingData, err := stub.GetState(instanceID)
	if err != nil {
		return nil, fmt.Errorf("获取状态数据时出错: %v", err)
	}
	if existingData == nil {
		return nil, fmt.Errorf("实例 %s 不存在", instanceID)
	}

	// 从现有实例中读取
	var instance ContractInstance
	err = json.Unmarshal(existingData, &instance)
	if err != nil {
		return nil, fmt.Errorf("反序列化实例数据时出错: %v", err)
	}

	// 创建网关对象
	instance.InstanceElements[gatewayID] = &Gateway{
		GatewayID:    gatewayID,
		GatewayState: gatewayState,
	}

	instanceJson, err := json.Marshal(instance)
	if err != nil {
		return nil, fmt.Errorf("序列化实例数据时出错: %v", err)
	}
	// 将网关对象序列化为JSON字符串并保存在状态数据库中
	err = stub.PutState(instanceID, instanceJson)
	if err != nil {
		return nil, fmt.Errorf("保存实例数据时出错: %v", err)
	}

	returnGateway, ok := instance.InstanceElements[gatewayID].(*Gateway)
	if !ok {
		return nil, fmt.Errorf("无法将实例元素转换为Gateway")
	}

	return returnGateway, nil
}

func (cc *SmartContract) CreateActionEvent(ctx contractapi.TransactionContextInterface, instanceID string, eventID string, eventState ElementState) (*ActionEvent, error) {
	stub := ctx.GetStub()

	// 检查是否存在具有相同ID的记录
	existingData, err := stub.GetState(instanceID)
	if err != nil {
		return nil, fmt.Errorf("获取状态数据时出错: %v", err)
	}
	if existingData == nil {
		return nil, fmt.Errorf("实例 %s 不存在", instanceID)
	}

	// 从现有实例中读取
	var instance ContractInstance
	err = json.Unmarshal(existingData, &instance)
	if err != nil {
		return nil, fmt.Errorf("反序列化实例数据时出错: %v", err)
	}

	// 创建事件对象
	instance.InstanceElements[eventID] = &ActionEvent{
		EventID:    eventID,
		EventState: eventState,
	}

	returnEvent, ok := instance.InstanceElements[eventID].(*ActionEvent)
	if !ok {
		return nil, fmt.Errorf("无法将实例元素转换为ActionEvent")
	}

	return returnEvent, nil

}

// Read function
func (c *SmartContract) ReadMsg(ctx contractapi.TransactionContextInterface, instanceID string, messageID string) (*Message, error) {
	instanceJson, err := ctx.GetStub().GetState(instanceID)
	if err != nil {
		return nil, err
	}
	if instanceJson == nil {
		errorMessage := fmt.Sprintf("Instance %s does not exist", instanceID)
		fmt.Println(errorMessage)
		return nil, errors.New(errorMessage)
	}

	var instance ContractInstance
	err = json.Unmarshal(instanceJson, &instance)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	msg, ok := instance.InstanceElements[messageID].(*Message)
	if !ok {
		errorMessage := fmt.Sprintf("Message %s does not exist", messageID)
		fmt.Println(errorMessage)
		return nil, errors.New(errorMessage)
	}

	return msg, nil
}

func (c *SmartContract) ReadGtw(ctx contractapi.TransactionContextInterface, instanceID string, gatewayID string) (*Gateway, error) {

	instanceJson, err := ctx.GetStub().GetState(instanceID)
	if err != nil {
		return nil, err
	}
	if instanceJson == nil {
		errorMessage := fmt.Sprintf("Instance %s does not exist", instanceID)
		fmt.Println(errorMessage)
		return nil, errors.New(errorMessage)
	}

	var instance ContractInstance
	err = json.Unmarshal(instanceJson, &instance)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	gtw, ok := instance.InstanceElements[gatewayID].(*Gateway)
	if !ok {
		errorMessage := fmt.Sprintf("Gateway %s does not exist", gatewayID)
		fmt.Println(errorMessage)
		return nil, errors.New(errorMessage)
	}

	return gtw, nil

}

func (c *SmartContract) ReadEvent(ctx contractapi.TransactionContextInterface, instanceID string, eventID string) (*ActionEvent, error) {

	instanceJson, err := ctx.GetStub().GetState(instanceID)
	if err != nil {
		return nil, err
	}
	if instanceJson == nil {
		errorMessage := fmt.Sprintf("Instance %s does not exist", instanceID)
		fmt.Println(errorMessage)
		return nil, errors.New(errorMessage)
	}

	var instance ContractInstance
	err = json.Unmarshal(instanceJson, &instance)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	actionEvent, ok := instance.InstanceElements[eventID].(*ActionEvent)
	if !ok {
		errorMessage := fmt.Sprintf("Event %s does not exist", eventID)
		fmt.Println(errorMessage)
		return nil, errors.New(errorMessage)
	}

	return actionEvent, nil

}

// Change State  function
func (c *SmartContract) ChangeMsgState(ctx contractapi.TransactionContextInterface, instanceID string, messageID string, msgState ElementState) error {

	stub := ctx.GetStub()

	instanceJson, err := stub.GetState(instanceID)
	if err != nil {
		return err
	}
	if instanceJson == nil {
		errorMessage := fmt.Sprintf("Instance %s does not exist", instanceID)
		fmt.Println(errorMessage)
		return errors.New(errorMessage)
	}

	var instance ContractInstance
	err = json.Unmarshal(instanceJson, &instance)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	msg, ok := instance.InstanceElements[messageID].(*Message)
	if !ok {
		errorMessage := fmt.Sprintf("Message %s does not exist", messageID)
		fmt.Println(errorMessage)
		return errors.New(errorMessage)
	}

	msg.MsgState = msgState

	instanceJson, err = json.Marshal(instance)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	err = stub.PutState(instanceID, instanceJson)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil

}

func (c *SmartContract) ChangeMsgFireflyTranID(ctx contractapi.TransactionContextInterface, instanceID string, messageID string, fireflyTranID string) error {

	stub := ctx.GetStub()

	instanceJson, err := stub.GetState(instanceID)
	if err != nil {
		return err
	}
	if instanceJson == nil {
		errorMessage := fmt.Sprintf("Instance %s does not exist", instanceID)
		fmt.Println(errorMessage)
		return errors.New(errorMessage)
	}

	var instance ContractInstance
	err = json.Unmarshal(instanceJson, &instance)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	msg, ok := instance.InstanceElements[messageID].(*Message)
	if !ok {
		errorMessage := fmt.Sprintf("Message %s does not exist", messageID)
		fmt.Println(errorMessage)
		return errors.New(errorMessage)
	}

	msg.FireflyTranID = fireflyTranID

	instanceJson, err = json.Marshal(instance)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	err = stub.PutState(instanceID, instanceJson)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil

}

func (c *SmartContract) ChangeGtwState(ctx contractapi.TransactionContextInterface, instanceID string, gatewayID string, gtwState ElementState) error {

	stub := ctx.GetStub()

	instanceJson, err := stub.GetState(instanceID)
	if err != nil {
		return err
	}
	if instanceJson == nil {
		errorMessage := fmt.Sprintf("Instance %s does not exist", instanceID)
		fmt.Println(errorMessage)
		return errors.New(errorMessage)
	}

	var instance ContractInstance
	err = json.Unmarshal(instanceJson, &instance)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	gtw, ok := instance.InstanceElements[gatewayID].(*Gateway)
	if !ok {
		errorMessage := fmt.Sprintf("Gateway %s does not exist", gatewayID)
		fmt.Println(errorMessage)
		return errors.New(errorMessage)
	}

	gtw.GatewayState = gtwState

	instanceJson, err = json.Marshal(instance)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	err = stub.PutState(instanceID, instanceJson)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil

}

func (c *SmartContract) ChangeEventState(ctx contractapi.TransactionContextInterface, instanceID string, eventID string, eventState ElementState) error {

	stub := ctx.GetStub()

	instanceJson, err := stub.GetState(instanceID)
	if err != nil {
		return err
	}
	if instanceJson == nil {
		errorMessage := fmt.Sprintf("Instance %s does not exist", instanceID)
		fmt.Println(errorMessage)
		return errors.New(errorMessage)
	}

	var instance ContractInstance
	err = json.Unmarshal(instanceJson, &instance)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	actionEvent, ok := instance.InstanceElements[eventID].(*ActionEvent)
	if !ok {
		errorMessage := fmt.Sprintf("Event %s does not exist", eventID)
		fmt.Println(errorMessage)
		return errors.New(errorMessage)
	}

	actionEvent.EventState = eventState

	instanceJson, err = json.Marshal(instance)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	err = stub.PutState(instanceID, instanceJson)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil

}

//get all message

func (cc *SmartContract) GetAllMessages(ctx contractapi.TransactionContextInterface, instanceID string) ([]*Message, error) {
	instanceJson, err := ctx.GetStub().GetState(instanceID)
	if err != nil {
		return nil, err
	}
	if instanceJson == nil {
		errorMessage := fmt.Sprintf("Instance %s does not exist", instanceID)
		fmt.Println(errorMessage)
		return nil, errors.New(errorMessage)
	}

	var instance ContractInstance
	err = json.Unmarshal(instanceJson, &instance)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	var messages []*Message
	for _, element := range instance.InstanceElements {
		msg, ok := element.(*Message)
		if ok {
			messages = append(messages, msg)
		}
	}

	return messages, nil
}

func (cc *SmartContract) GetAllGateways(ctx contractapi.TransactionContextInterface, instanceID string) ([]*Gateway, error) {

	instanceJson, err := ctx.GetStub().GetState(instanceID)
	if err != nil {
		return nil, err
	}
	if instanceJson == nil {
		errorMessage := fmt.Sprintf("Instance %s does not exist", instanceID)
		fmt.Println(errorMessage)
		return nil, errors.New(errorMessage)
	}

	var instance ContractInstance
	err = json.Unmarshal(instanceJson, &instance)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	var gateways []*Gateway
	for _, element := range instance.InstanceElements {
		gtw, ok := element.(*Gateway)
		if ok {
			gateways = append(gateways, gtw)
		}
	}

	return gateways, nil
}

func (cc *SmartContract) GetAllActionEvents(ctx contractapi.TransactionContextInterface, instanceID string) ([]*ActionEvent, error) {

	instanceJson, err := ctx.GetStub().GetState(instanceID)
	if err != nil {
		return nil, err
	}
	if instanceJson == nil {
		errorMessage := fmt.Sprintf("Instance %s does not exist", instanceID)
		fmt.Println(errorMessage)
		return nil, errors.New(errorMessage)
	}

	var instance ContractInstance
	err = json.Unmarshal(instanceJson, &instance)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	var actionEvents []*ActionEvent
	for _, element := range instance.InstanceElements {
		event, ok := element.(*ActionEvent)
		if ok {
			actionEvents = append(actionEvents, event)
		}
	}

	return actionEvents, nil

}

func (cc *SmartContract) ReadGlobalVariable(ctx contractapi.TransactionContextInterface, instanceID string) (*StateMemory, error) {

	instanceJson, err := ctx.GetStub().GetState(instanceID)
	if err != nil {
		return nil, err
	}
	if instanceJson == nil {
		errorMessage := fmt.Sprintf("Instance %s does not exist", instanceID)
		fmt.Println(errorMessage)
		return nil, errors.New(errorMessage)
	}

	var instance ContractInstance
	err = json.Unmarshal(instanceJson, &instance)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &instance.InstanceStateMemory, nil

}

func (cc *SmartContract) SetGlobalVariable(ctx contractapi.TransactionContextInterface, instanceID string, globalVariable *StateMemory) error {

	stub := ctx.GetStub()

	instanceJson, err := stub.GetState(instanceID)
	if err != nil {
		return err
	}
	if instanceJson == nil {
		errorMessage := fmt.Sprintf("Instance %s does not exist", instanceID)
		fmt.Println(errorMessage)
		return errors.New(errorMessage)
	}

	var instance ContractInstance
	err = json.Unmarshal(instanceJson, &instance)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	instance.InstanceStateMemory = *globalVariable

	instanceJson, err = json.Marshal(instance)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	err = stub.PutState(instanceID, instanceJson)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil

}

func (cc *SmartContract) ReadParticipant(ctx contractapi.TransactionContextInterface, instanceID string, participantID string) (*Participant, error) {

	instanceJson, err := ctx.GetStub().GetState(instanceID)
	if err != nil {
		return nil, err
	}
	if instanceJson == nil {
		errorMessage := fmt.Sprintf("Instance %s does not exist", instanceID)
		fmt.Println(errorMessage)
		return nil, errors.New(errorMessage)
	}

	var instance ContractInstance
	err = json.Unmarshal(instanceJson, &instance)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	participant, ok := instance.InstanceElements[participantID].(*Participant)
	if !ok {
		errorMessage := fmt.Sprintf("Participant %s does not exist", participantID)
		fmt.Println(errorMessage)
		return nil, errors.New(errorMessage)
	}

	return participant, nil

}

func (cc *SmartContract) WriteParticipant(ctx contractapi.TransactionContextInterface, instanceID string, participantID string, participant *Participant) error {
	stub := ctx.GetStub()

	instanceJson, err := stub.GetState(instanceID)
	if err != nil {
		return err
	}
	if instanceJson == nil {
		errorMessage := fmt.Sprintf("Instance %s does not exist", instanceID)
		fmt.Println(errorMessage)
		return errors.New(errorMessage)
	}

	var instance ContractInstance
	err = json.Unmarshal(instanceJson, &instance)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	instance.InstanceElements[participantID] = participant

	instanceJson, err = json.Marshal(instance)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	err = stub.PutState(instanceID, instanceJson)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil

}

func (cc *SmartContract) check_msp(ctx contractapi.TransactionContextInterface, instanceID string, target_participant string) bool {
	// Read the target participant's msp
	targetParticipant, err := cc.ReadParticipant(ctx, instanceID, target_participant)
	if err != nil {
		return false
	}
	mspID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return false
	}
	return mspID == targetParticipant.MSP
}

func (cc *SmartContract) check_attribute(ctx contractapi.TransactionContextInterface, instanceID string, target_participant string, attributeName string) bool {
	targetParticipant, err := cc.ReadParticipant(ctx, instanceID, target_participant)
	if err != nil {
		return false
	}
	if ctx.GetClientIdentity().AssertAttributeValue(attributeName, targetParticipant.Attributes[attributeName]) != nil {
		return false
	}

	return true
}

func (cc *SmartContract) check_participant(ctx contractapi.TransactionContextInterface, instanceID string, target_participant string) bool {
	// Read the target participant's msp
	targetParticipant, err := cc.ReadParticipant(ctx, instanceID, target_participant)
	if err != nil {
		return false
	}

	if !targetParticipant.IsMulti {
		// check X509 = MSPID + @ + ID
		mspID, _ := ctx.GetClientIdentity().GetMSPID()
		pid, _ := ctx.GetClientIdentity().GetID()
		if targetParticipant.X509 == pid+"@"+mspID {
			return true
		} else {
			return false
		}
	}

	// check MSP if msp!=''
	if targetParticipant.MSP != "" && cc.check_msp(ctx, instanceID, target_participant) == false {
		return false
	}

	// check all attributes
	for key, _ := range targetParticipant.Attributes {
		if cc.check_attribute(ctx, instanceID, target_participant, key) == false {
			return false
		}
	}

	return true
}

func (cc *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	stub := ctx.GetStub()

	// isInited in state
	isInitedBytes, err := stub.GetState("isInited")
	if err != nil {
		return fmt.Errorf("Failed to get isInited: %v", err)
	}
	if isInitedBytes != nil {
		errorMessage := "Chaincode has already been initialized"
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	stub.PutState("currentInstanceId", []byte("0"))

	stub.PutState("isInited", []byte("true"))

	stub.SetEvent("initContractEvent", []byte("Contract has been initialized successfully"))
	return nil
}


func (cc *SmartContract) CheckRegister(ctx contractapi.TransactionContextInterface, instanceID string) (bool, error) {
	stub := ctx.GetStub()

	// Check if the instance has been registered
	instanceBytes, err := stub.GetState(instanceID)
	if err != nil {
		return false, fmt.Errorf("Failed to read from world state. %s", err.Error())
	}
	var instance ContractInstance
	err = json.Unmarshal(instanceBytes, &instance)
	if err != nil {
		return false, fmt.Errorf("Failed to unmarshal. %s", err.Error())
	}

	if instance.InstanceState == READY {
		return true, nil
	}

	// set State depend on Participant with IsMulti=true

	for element, value := range instance.InstanceElements {
		participant, ok := value.(*Participant)
		if ok {
			if !participant.IsMulti && participant.X509 == "" {
				return false, fmt.Errorf("The participant %s is not registered.", element)
			}
		}
	}

	// set State depend on Participant with IsMulti=false
	instance.InstanceState = READY
	instanceBytes, err = json.Marshal(instance)
	if err != nil {
		return false, fmt.Errorf("Failed to marshal. %s", err.Error())
	}

	err = stub.PutState(instanceID, instanceBytes)
	if err != nil {
		return false, fmt.Errorf("Failed to put state. %s", err.Error())
	}

	return true, nil
}

func (cc *SmartContract) RegisterParticipant(ctx contractapi.TransactionContextInterface, instanceID string, targetParticipantID string) error {
	{
		// check if the participant is single
		var targetParticipant Participant
		participant, _ := cc.ReadParticipant(ctx, instanceID, targetParticipantID)
		if participant.IsMulti {
			{
				return fmt.Errorf("The participant is not multi")
			}
		}

		// check ACL

		if !cc.check_participant(ctx, instanceID, targetParticipantID) {
			return fmt.Errorf("The participant is not allowed to be registered")
		}

		// Read the identity of invoker ,and binding it's identity to the participant

		// Get the identity of the invoker
		invokerIdentity, err := ctx.GetClientIdentity().GetID()
		mspIndentity, err := ctx.GetClientIdentity().GetMSPID()

		X509 := invokerIdentity + "@" + mspIndentity

		// save the identity to the participant
		targetParticipant.X509 = X509

		// save the participant
		err = cc.WriteParticipant(ctx, instanceID, targetParticipantID, &targetParticipant)
		if err != nil {
			{
				return err
			}
		}

		return nil
	}
}

func (cc *SmartContract) CreateInstance(ctx contractapi.TransactionContextInterface, initParametersBytes string) (string, error) {
	stub := ctx.GetStub()

	isInitedBytes, err := stub.GetState("isInited")
	if err != nil {
		return "", fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	if isInitedBytes != nil {
		return "", fmt.Errorf("The instance has been initialized.")
	}

	var isInited bool
	err = json.Unmarshal(isInitedBytes, &isInited)
	if err != nil {
		return "", fmt.Errorf("Failed to unmarshal. %s", err.Error())
	}

	if !isInited {
		return "", fmt.Errorf("The instance has not been initialized.")
	}

	// get the instanceID

	var instanceID string
	instanceIDString, err := stub.GetState("currentInstanceID")
	if err != nil {
		return "", fmt.Errorf("Failed to read from world state. %s", err.Error())
	}

	err = json.Unmarshal(instanceIDString, &instanceID)
	if err != nil {
		return "", fmt.Errorf("Failed to unmarshal. %s", err.Error())
	}

	// Create the instance with the data from the InitParameters
	var initParameters InitParameters
	err = json.Unmarshal([]byte(initParametersBytes), &initParameters)
	if err != nil {
		return "", fmt.Errorf("Failed to unmarshal. %s", err.Error())
	}

	instance := ContractInstance{
		InstanceID:          instanceID,
		InstanceStateMemory: StateMemory{},
		InstanceElements:    make(map[string]interface{}),
	}

	// Save the instance
	instanceBytes, err := json.Marshal(instance)
	if err != nil {
		return "", fmt.Errorf("Failed to marshal. %s", err.Error())
	}

	err = stub.PutState(instanceID, instanceBytes)
	if err != nil {
		return "", fmt.Errorf("Failed to put state. %s", err.Error())
	}

	// Update the currentInstanceID

		cc.CreateParticipant(ctx, instanceID, "Participant_1080bkg", initParameters.Participant_1080bkg.MSP, initParameters.Participant_1080bkg.Attributes, false, 0, 0)
	cc.CreateParticipant(ctx, instanceID, "Participant_0sktaei", initParameters.Participant_0sktaei.MSP, initParameters.Participant_0sktaei.Attributes, false, 0, 0)
	cc.CreateParticipant(ctx, instanceID, "Participant_1gcdqza", initParameters.Participant_1gcdqza.MSP, initParameters.Participant_1gcdqza.Attributes, false, 0, 0)
	cc.CreateActionEvent(ctx, instanceID, "Event_1jtgn3j", ENABLED)

	cc.CreateActionEvent(ctx, instanceID, "Event_0366pfz", DISABLED)

	cc.CreateActionEvent(ctx, instanceID, "Event_08edp7f", DISABLED)

	cc.CreateActionEvent(ctx, instanceID, "Event_146eii4", DISABLED)

	cc.CreateMessage(ctx, instanceID, "Message_1qbk325", "Participant_1gcdqza", "Participant_0sktaei", "", DISABLED, `{"properties":{"product Id":{"type":"string","description":"Delivered product id"}},"required":["product Id"],"files":{},"file required":[]}`)
	cc.CreateMessage(ctx, instanceID, "Message_1q05nnw", "Participant_0sktaei", "Participant_1gcdqza", "", DISABLED, `{"properties":{"payment amount":{"type":"number","description":"payment amount"}},"required":["payment amount"],"files":{},"file required":[]}`)
	cc.CreateMessage(ctx, instanceID, "Message_1i8rlqn", "Participant_0sktaei", "Participant_1gcdqza", "", DISABLED, `{"properties":{"external service Id":{"type":"string","description":"The requested external service information"}},"required":["external service Id"],"files":{},"file required":[]}`)
	cc.CreateMessage(ctx, instanceID, "Message_0m9p3da", "Participant_1080bkg", "Participant_0sktaei", "", DISABLED, `{"properties":{"invoice":{"type":"boolean","description":"Do you need an invoice?"}},"required":["invoice"],"files":{},"file required":[]}`)
	cc.CreateMessage(ctx, instanceID, "Message_1etcmvl", "Participant_0sktaei", "Participant_1080bkg", "", DISABLED, `{"properties":{"invoice_id":{"type":"string","description":"Invoice Id"},"invoice_data":{"type":"number","description":"Date of invoice issuance"}},"required":["invoice_id"],"files":{"invoice":{"type":"file","description":"Invoice documents"}},"file required":["invoice"]}`)
	cc.CreateMessage(ctx, instanceID, "Message_1joj7ca", "Participant_1080bkg", "Participant_0sktaei", "", DISABLED, `{"properties":{"invoice information":{"type":"string","description":"Invoice related information"}},"required":["invoice information"],"files":{},"file required":[]}`)
	cc.CreateMessage(ctx, instanceID, "Message_1ljlm4g", "Participant_0sktaei", "Participant_1080bkg", "", DISABLED, `{"properties":{"delivered_product_id":{"type":"string","description":"delivered_product_id"}},"required":["delivered_product_id"],"files":{},"file required":[]}`)
	cc.CreateMessage(ctx, instanceID, "Message_1xm9dxy", "Participant_1080bkg", "Participant_0sktaei", "", DISABLED, `{"properties":{"motivation":{"type":"string","description":"Motivation for Canceling orders"}},"required":["motivation"],"files":{},"file required":[]}`)
	cc.CreateMessage(ctx, instanceID, "Message_0o8eyir", "Participant_1080bkg", "Participant_0sktaei", "", DISABLED, `{"properties":{"payment amount":{"type":"number","description":"payment amount"},"orderID":{"type":"number","description":"The order id of payment"}},"required":["payment amount","orderID"],"files":{},"file required":[]}`)
	cc.CreateMessage(ctx, instanceID, "Message_1nlagx2", "Participant_1080bkg", "Participant_0sktaei", "", DISABLED, `{"properties":{"confirmation":{"type":"boolean","description":"Whether to accept the service plan"}},"required":["confirmation"],"files":{},"file required":[]}`)
	cc.CreateMessage(ctx, instanceID, "Message_1em0ee4", "Participant_0sktaei", "Participant_1080bkg", "", DISABLED, `{"properties":{"service plan":{"type":"string","description":"service plan"},"price_quotation":{"type":"number","description":"Price quotation"},"need_external_provider":{"type":"boolean","description":"Whether external service providers are required"}},"required":["service plan","price_quotation","need_external_provider"],"files":{},"file required":[]}`)
	cc.CreateMessage(ctx, instanceID, "Message_0r9lypd", "Participant_0sktaei", "Participant_1080bkg", "", DISABLED, `{"properties":{"is_available":{"type":"boolean","description":"Is the service available?"}},"required":["is_available"],"files":{},"file required":[]}`)
	cc.CreateMessage(ctx, instanceID, "Message_045i10y", "Participant_1080bkg", "Participant_0sktaei", "", DISABLED, `{"properties":{"serviceId":{"type":"string","description":"The required service id"}},"required":["serviceId"],"files":{},"file required":[]}`)
	cc.CreateGateway(ctx, instanceID, "ExclusiveGateway_106je4z", DISABLED)

	cc.CreateGateway(ctx, instanceID, "ExclusiveGateway_0hs3ztq", DISABLED)

	cc.CreateGateway(ctx, instanceID, "ExclusiveGateway_0nzwv7v", DISABLED)

	cc.CreateGateway(ctx, instanceID, "Gateway_1bhtapl", DISABLED)

	cc.CreateGateway(ctx, instanceID, "Gateway_04h9e6e", DISABLED)

	cc.CreateGateway(ctx, instanceID, "EventBasedGateway_1fxpmyn", DISABLED)


	instanceIDInt, err := strconv.Atoi(instanceID)
	if err != nil {
		return "", fmt.Errorf("Failed to convert instanceID to int. %s", err.Error())
	}

	instanceIDInt++
	instanceID = strconv.Itoa(instanceIDInt)

	instanceIDBytes, err := json.Marshal(instanceID)
	if err != nil {
		return "", fmt.Errorf("Failed to marshal instanceID. %s", err.Error())
	}

	err = stub.PutState("currentInstanceID", instanceIDBytes)
	if err != nil {
		return "", fmt.Errorf("Failed to put state. %s", err.Error())
	}

	return instanceID, nil

}

func (cc *SmartContract) Event_1jtgn3j(ctx contractapi.TransactionContextInterface, instanceID string) error {
	stub := ctx.GetStub()
	isRegistered, err := cc.CheckRegister(ctx, instanceID)
	if !isRegistered {
		errorMessage := fmt.Sprintf("Instance %s is not registered fully", instanceID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}


	actionEvent, err := cc.ReadEvent(ctx, instanceID, "Event_1jtgn3j")
	if err != nil {
		return err
	}

	if actionEvent.EventState != ENABLED {
		errorMessage := fmt.Sprintf("Event state %s is not allowed", actionEvent.EventID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	cc.ChangeEventState(ctx, instanceID, "Event_1jtgn3j", COMPLETED)
	stub.SetEvent("Event_1jtgn3j", []byte("Contract has been started successfully"))
	
	    cc.ChangeGtwState(ctx, instanceID, "ExclusiveGateway_0hs3ztq", ENABLED)
	
	return nil
}

func (cc *SmartContract) Message_045i10y_Send(ctx contractapi.TransactionContextInterface, instanceID string, fireflyTranID string ) error {
	stub := ctx.GetStub()
	msg, err := cc.ReadMsg(ctx, instanceID, "Message_045i10y")
	if err != nil {
		return err
	}

	//
	if cc.check_participant(ctx, instanceID, msg.SendParticipantID) == false{
		errorMessage := fmt.Sprintf("Participant %s is not allowed to send the message", msg.SendParticipantID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	if msg.MsgState != ENABLED {
		errorMessage := fmt.Sprintf("Message state %s is not allowed", msg.MessageID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	cc.ChangeMsgFireflyTranID(ctx, instanceID, fireflyTranID, msg.MessageID)
	cc.ChangeMsgState(ctx, instanceID, msg.MessageID, WAITINGFORCONFIRMATION)
	
	stub.SetEvent("Message_045i10y", []byte("Message is waiting for confirmation"))

	
	return nil
}

func (cc *SmartContract) Message_045i10y_Complete(ctx contractapi.TransactionContextInterface, instanceID string) error {
	stub := ctx.GetStub()
	msg, err := cc.ReadMsg(ctx, instanceID, "Message_045i10y")
	if err != nil {
		return err
	}

	if cc.check_participant(ctx, instanceID, msg.ReceiveParticipantID) == false{
		errorMessage := fmt.Sprintf("Participant %s is not allowed to send the message", msg.SendParticipantID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	if msg.MsgState != WAITINGFORCONFIRMATION {
		errorMessage := fmt.Sprintf("Event state %s is not allowed", msg.MessageID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	cc.ChangeMsgState(ctx, instanceID, msg.MessageID, COMPLETED)
	stub.SetEvent("Message_045i10y", []byte("Message has been done"))

	
	    cc.ChangeMsgState(ctx, instanceID, "Message_0r9lypd", ENABLED)

	
	return nil
}

func (cc *SmartContract) Message_0r9lypd_Send(ctx contractapi.TransactionContextInterface, instanceID string, fireflyTranID string , Is_available bool) error {
	stub := ctx.GetStub()
	msg, err := cc.ReadMsg(ctx, instanceID, "Message_0r9lypd")
	if err != nil {
		return err
	}

	//
	if cc.check_participant(ctx, instanceID, msg.SendParticipantID) == false{
		errorMessage := fmt.Sprintf("Participant %s is not allowed to send the message", msg.SendParticipantID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	if msg.MsgState != ENABLED {
		errorMessage := fmt.Sprintf("Message state %s is not allowed", msg.MessageID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	cc.ChangeMsgFireflyTranID(ctx, instanceID, fireflyTranID, msg.MessageID)
	cc.ChangeMsgState(ctx, instanceID, msg.MessageID, WAITINGFORCONFIRMATION)
		globalMemory,readGloabolError := cc.ReadGlobalVariable(ctx, instanceID)
	if readGloabolError != nil {
		fmt.Println(readGloabolError.Error())
		return readGloabolError
	}
	globalMemory.Is_available = Is_available
	setGloabolErrror :=cc.SetGlobalVariable(ctx, instanceID, globalMemory)
	if setGloabolErrror != nil {
		fmt.Println(setGloabolErrror.Error())
		return setGloabolErrror
	}
	stub.SetEvent("Message_0r9lypd", []byte("Message is waiting for confirmation"))

	
	return nil
}

func (cc *SmartContract) Message_0r9lypd_Complete(ctx contractapi.TransactionContextInterface, instanceID string) error {
	stub := ctx.GetStub()
	msg, err := cc.ReadMsg(ctx, instanceID, "Message_0r9lypd")
	if err != nil {
		return err
	}

	if cc.check_participant(ctx, instanceID, msg.ReceiveParticipantID) == false{
		errorMessage := fmt.Sprintf("Participant %s is not allowed to send the message", msg.SendParticipantID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	if msg.MsgState != WAITINGFORCONFIRMATION {
		errorMessage := fmt.Sprintf("Event state %s is not allowed", msg.MessageID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	cc.ChangeMsgState(ctx, instanceID, msg.MessageID, COMPLETED)
	stub.SetEvent("Message_0r9lypd", []byte("Message has been done"))

	
	    cc.ChangeGtwState(ctx, instanceID, "ExclusiveGateway_106je4z", ENABLED)

	
	return nil
}

func (cc *SmartContract) ExclusiveGateway_106je4z(ctx contractapi.TransactionContextInterface, instanceID string) error {
	stub := ctx.GetStub()
	gtw, err := cc.ReadGtw(ctx, instanceID, "ExclusiveGateway_106je4z")
	if err != nil {
		return err
	}

	if gtw.GatewayState != ENABLED {
		errorMessage := fmt.Sprintf("Gateway state %s is not allowed", gtw.GatewayID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	cc.ChangeGtwState(ctx, instanceID, gtw.GatewayID, COMPLETED)
	stub.SetEvent("ExclusiveGateway_106je4z", []byte("ExclusiveGateway has been done"))

    
    	currentMemory, err := cc.ReadGlobalVariable(ctx, instanceID)
	if err != nil {
		return err
	}

    Is_available:=currentMemory.Is_available

if Is_available==true {
	    cc.ChangeMsgState(ctx, instanceID, "Message_1em0ee4", ENABLED)
}
if Is_available==false {
	    cc.ChangeGtwState(ctx, instanceID, "ExclusiveGateway_0hs3ztq", ENABLED)
}
    

	return nil
}

func (cc *SmartContract) Message_1em0ee4_Send(ctx contractapi.TransactionContextInterface, instanceID string, fireflyTranID string , Need_external_provider bool) error {
	stub := ctx.GetStub()
	msg, err := cc.ReadMsg(ctx, instanceID, "Message_1em0ee4")
	if err != nil {
		return err
	}

	//
	if cc.check_participant(ctx, instanceID, msg.SendParticipantID) == false{
		errorMessage := fmt.Sprintf("Participant %s is not allowed to send the message", msg.SendParticipantID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	if msg.MsgState != ENABLED {
		errorMessage := fmt.Sprintf("Message state %s is not allowed", msg.MessageID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	cc.ChangeMsgFireflyTranID(ctx, instanceID, fireflyTranID, msg.MessageID)
	cc.ChangeMsgState(ctx, instanceID, msg.MessageID, WAITINGFORCONFIRMATION)
		globalMemory,readGloabolError := cc.ReadGlobalVariable(ctx, instanceID)
	if readGloabolError != nil {
		fmt.Println(readGloabolError.Error())
		return readGloabolError
	}
	globalMemory.Need_external_provider = Need_external_provider
	setGloabolErrror :=cc.SetGlobalVariable(ctx, instanceID, globalMemory)
	if setGloabolErrror != nil {
		fmt.Println(setGloabolErrror.Error())
		return setGloabolErrror
	}
	stub.SetEvent("Message_1em0ee4", []byte("Message is waiting for confirmation"))

	
	return nil
}

func (cc *SmartContract) Message_1em0ee4_Complete(ctx contractapi.TransactionContextInterface, instanceID string) error {
	stub := ctx.GetStub()
	msg, err := cc.ReadMsg(ctx, instanceID, "Message_1em0ee4")
	if err != nil {
		return err
	}

	if cc.check_participant(ctx, instanceID, msg.ReceiveParticipantID) == false{
		errorMessage := fmt.Sprintf("Participant %s is not allowed to send the message", msg.SendParticipantID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	if msg.MsgState != WAITINGFORCONFIRMATION {
		errorMessage := fmt.Sprintf("Event state %s is not allowed", msg.MessageID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	cc.ChangeMsgState(ctx, instanceID, msg.MessageID, COMPLETED)
	stub.SetEvent("Message_1em0ee4", []byte("Message has been done"))

	
	    cc.ChangeMsgState(ctx, instanceID, "Message_1nlagx2", ENABLED)

	
	return nil
}

func (cc *SmartContract) Message_1nlagx2_Send(ctx contractapi.TransactionContextInterface, instanceID string, fireflyTranID string ) error {
	stub := ctx.GetStub()
	msg, err := cc.ReadMsg(ctx, instanceID, "Message_1nlagx2")
	if err != nil {
		return err
	}

	//
	if cc.check_participant(ctx, instanceID, msg.SendParticipantID) == false{
		errorMessage := fmt.Sprintf("Participant %s is not allowed to send the message", msg.SendParticipantID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	if msg.MsgState != ENABLED {
		errorMessage := fmt.Sprintf("Message state %s is not allowed", msg.MessageID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	cc.ChangeMsgFireflyTranID(ctx, instanceID, fireflyTranID, msg.MessageID)
	cc.ChangeMsgState(ctx, instanceID, msg.MessageID, WAITINGFORCONFIRMATION)
	
	stub.SetEvent("Message_1nlagx2", []byte("Message is waiting for confirmation"))

	
	return nil
}

func (cc *SmartContract) Message_1nlagx2_Complete(ctx contractapi.TransactionContextInterface, instanceID string) error {
	stub := ctx.GetStub()
	msg, err := cc.ReadMsg(ctx, instanceID, "Message_1nlagx2")
	if err != nil {
		return err
	}

	if cc.check_participant(ctx, instanceID, msg.ReceiveParticipantID) == false{
		errorMessage := fmt.Sprintf("Participant %s is not allowed to send the message", msg.SendParticipantID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	if msg.MsgState != WAITINGFORCONFIRMATION {
		errorMessage := fmt.Sprintf("Event state %s is not allowed", msg.MessageID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	cc.ChangeMsgState(ctx, instanceID, msg.MessageID, COMPLETED)
	stub.SetEvent("Message_1nlagx2", []byte("Message has been done"))

	
	    cc.ChangeGtwState(ctx, instanceID, "EventBasedGateway_1fxpmyn", ENABLED)

	
	return nil
}

func (cc *SmartContract) ExclusiveGateway_0hs3ztq(ctx contractapi.TransactionContextInterface, instanceID string) error {
	stub := ctx.GetStub()
	gtw, err := cc.ReadGtw(ctx, instanceID, "ExclusiveGateway_0hs3ztq")
	if err != nil {
		return err
	}

	if gtw.GatewayState != ENABLED {
		errorMessage := fmt.Sprintf("Gateway state %s is not allowed", gtw.GatewayID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	cc.ChangeGtwState(ctx, instanceID, gtw.GatewayID, COMPLETED)
	stub.SetEvent("ExclusiveGateway_0hs3ztq", []byte("ExclusiveGateway has been done"))

    
        cc.ChangeMsgState(ctx, instanceID, "Message_045i10y", ENABLED)
    

	return nil
}

func (cc *SmartContract) EventBasedGateway_1fxpmyn(ctx contractapi.TransactionContextInterface, instanceID string) error { 
	stub := ctx.GetStub()
	gtw, err := cc.ReadGtw(ctx, instanceID, "EventBasedGateway_1fxpmyn")
	if err != nil {
		return err
	}

	if gtw.GatewayState != ENABLED {
		errorMessage := fmt.Sprintf("Gateway state %s is not allowed", gtw.GatewayID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	cc.ChangeGtwState(ctx, instanceID, gtw.GatewayID, COMPLETED)
	stub.SetEvent("EventBasedGateway_1fxpmyn", []byte("EventbasedGateway has been done"))

 	
        cc.ChangeMsgState(ctx, instanceID, "Message_0o8eyir", ENABLED)
    cc.ChangeMsgState(ctx, instanceID, "Message_1xm9dxy", ENABLED)
    

    return nil
}

func (cc *SmartContract) Message_0o8eyir_Send(ctx contractapi.TransactionContextInterface, instanceID string, fireflyTranID string ) error {
	stub := ctx.GetStub()
	msg, err := cc.ReadMsg(ctx, instanceID, "Message_0o8eyir")
	if err != nil {
		return err
	}

	//
	if cc.check_participant(ctx, instanceID, msg.SendParticipantID) == false{
		errorMessage := fmt.Sprintf("Participant %s is not allowed to send the message", msg.SendParticipantID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	if msg.MsgState != ENABLED {
		errorMessage := fmt.Sprintf("Message state %s is not allowed", msg.MessageID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	cc.ChangeMsgFireflyTranID(ctx, instanceID, fireflyTranID, msg.MessageID)
	cc.ChangeMsgState(ctx, instanceID, msg.MessageID, WAITINGFORCONFIRMATION)
	
	stub.SetEvent("Message_0o8eyir", []byte("Message is waiting for confirmation"))

	    cc.ChangeMsgState(ctx, instanceID, "Message_1xm9dxy", DISABLED)
	return nil
}

func (cc *SmartContract) Message_0o8eyir_Complete(ctx contractapi.TransactionContextInterface, instanceID string) error {
	stub := ctx.GetStub()
	msg, err := cc.ReadMsg(ctx, instanceID, "Message_0o8eyir")
	if err != nil {
		return err
	}

	if cc.check_participant(ctx, instanceID, msg.ReceiveParticipantID) == false{
		errorMessage := fmt.Sprintf("Participant %s is not allowed to send the message", msg.SendParticipantID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	if msg.MsgState != WAITINGFORCONFIRMATION {
		errorMessage := fmt.Sprintf("Event state %s is not allowed", msg.MessageID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	cc.ChangeMsgState(ctx, instanceID, msg.MessageID, COMPLETED)
	stub.SetEvent("Message_0o8eyir", []byte("Message has been done"))

	
	    cc.ChangeGtwState(ctx, instanceID, "Gateway_1bhtapl", ENABLED)

	
	return nil
}

func (cc *SmartContract) Message_1xm9dxy_Send(ctx contractapi.TransactionContextInterface, instanceID string, fireflyTranID string ) error {
	stub := ctx.GetStub()
	msg, err := cc.ReadMsg(ctx, instanceID, "Message_1xm9dxy")
	if err != nil {
		return err
	}

	//
	if cc.check_participant(ctx, instanceID, msg.SendParticipantID) == false{
		errorMessage := fmt.Sprintf("Participant %s is not allowed to send the message", msg.SendParticipantID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	if msg.MsgState != ENABLED {
		errorMessage := fmt.Sprintf("Message state %s is not allowed", msg.MessageID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	cc.ChangeMsgFireflyTranID(ctx, instanceID, fireflyTranID, msg.MessageID)
	cc.ChangeMsgState(ctx, instanceID, msg.MessageID, WAITINGFORCONFIRMATION)
	
	stub.SetEvent("Message_1xm9dxy", []byte("Message is waiting for confirmation"))

	    cc.ChangeMsgState(ctx, instanceID, "Message_0o8eyir", DISABLED)
	return nil
}

func (cc *SmartContract) Message_1xm9dxy_Complete(ctx contractapi.TransactionContextInterface, instanceID string) error {
	stub := ctx.GetStub()
	msg, err := cc.ReadMsg(ctx, instanceID, "Message_1xm9dxy")
	if err != nil {
		return err
	}

	if cc.check_participant(ctx, instanceID, msg.ReceiveParticipantID) == false{
		errorMessage := fmt.Sprintf("Participant %s is not allowed to send the message", msg.SendParticipantID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	if msg.MsgState != WAITINGFORCONFIRMATION {
		errorMessage := fmt.Sprintf("Event state %s is not allowed", msg.MessageID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	cc.ChangeMsgState(ctx, instanceID, msg.MessageID, COMPLETED)
	stub.SetEvent("Message_1xm9dxy", []byte("Message has been done"))

	
	    cc.ChangeEventState(ctx, instanceID, "Event_0366pfz", ENABLED)

	
	return nil
}

func (cc *SmartContract) Event_0366pfz(ctx contractapi.TransactionContextInterface, instanceID string) error {
	stub := ctx.GetStub()
	event, err := cc.ReadEvent(ctx, instanceID, "Event_0366pfz")
	if err != nil {
		return err
	}

	if event.EventState != ENABLED {
		errorMessage := fmt.Sprintf("Event state %s is not allowed", event.EventID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	cc.ChangeEventState(ctx, instanceID, event.EventID, COMPLETED) 
	stub.SetEvent("Event_0366pfz", []byte("EndEvent has been done"))
	
	return nil
}

func (cc *SmartContract) Message_1ljlm4g_Send(ctx contractapi.TransactionContextInterface, instanceID string, fireflyTranID string ) error {
	stub := ctx.GetStub()
	msg, err := cc.ReadMsg(ctx, instanceID, "Message_1ljlm4g")
	if err != nil {
		return err
	}

	//
	if cc.check_participant(ctx, instanceID, msg.SendParticipantID) == false{
		errorMessage := fmt.Sprintf("Participant %s is not allowed to send the message", msg.SendParticipantID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	if msg.MsgState != ENABLED {
		errorMessage := fmt.Sprintf("Message state %s is not allowed", msg.MessageID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	cc.ChangeMsgFireflyTranID(ctx, instanceID, fireflyTranID, msg.MessageID)
	cc.ChangeMsgState(ctx, instanceID, msg.MessageID, WAITINGFORCONFIRMATION)
	
	stub.SetEvent("Message_1ljlm4g", []byte("Message is waiting for confirmation"))

	
	return nil
}

func (cc *SmartContract) Message_1ljlm4g_Complete(ctx contractapi.TransactionContextInterface, instanceID string) error {
	stub := ctx.GetStub()
	msg, err := cc.ReadMsg(ctx, instanceID, "Message_1ljlm4g")
	if err != nil {
		return err
	}

	if cc.check_participant(ctx, instanceID, msg.ReceiveParticipantID) == false{
		errorMessage := fmt.Sprintf("Participant %s is not allowed to send the message", msg.SendParticipantID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	if msg.MsgState != WAITINGFORCONFIRMATION {
		errorMessage := fmt.Sprintf("Event state %s is not allowed", msg.MessageID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	cc.ChangeMsgState(ctx, instanceID, msg.MessageID, COMPLETED)
	stub.SetEvent("Message_1ljlm4g", []byte("Message has been done"))

	
	    cc.ChangeMsgState(ctx, instanceID, "Message_0m9p3da", ENABLED)

	
	return nil
}

func (cc *SmartContract) Message_0m9p3da_Send(ctx contractapi.TransactionContextInterface, instanceID string, fireflyTranID string , Invoice bool) error {
	stub := ctx.GetStub()
	msg, err := cc.ReadMsg(ctx, instanceID, "Message_0m9p3da")
	if err != nil {
		return err
	}

	//
	if cc.check_participant(ctx, instanceID, msg.SendParticipantID) == false{
		errorMessage := fmt.Sprintf("Participant %s is not allowed to send the message", msg.SendParticipantID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	if msg.MsgState != ENABLED {
		errorMessage := fmt.Sprintf("Message state %s is not allowed", msg.MessageID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	cc.ChangeMsgFireflyTranID(ctx, instanceID, fireflyTranID, msg.MessageID)
	cc.ChangeMsgState(ctx, instanceID, msg.MessageID, WAITINGFORCONFIRMATION)
		globalMemory,readGloabolError := cc.ReadGlobalVariable(ctx, instanceID)
	if readGloabolError != nil {
		fmt.Println(readGloabolError.Error())
		return readGloabolError
	}
	globalMemory.Invoice = Invoice
	setGloabolErrror :=cc.SetGlobalVariable(ctx, instanceID, globalMemory)
	if setGloabolErrror != nil {
		fmt.Println(setGloabolErrror.Error())
		return setGloabolErrror
	}
	stub.SetEvent("Message_0m9p3da", []byte("Message is waiting for confirmation"))

	
	return nil
}

func (cc *SmartContract) Message_0m9p3da_Complete(ctx contractapi.TransactionContextInterface, instanceID string) error {
	stub := ctx.GetStub()
	msg, err := cc.ReadMsg(ctx, instanceID, "Message_0m9p3da")
	if err != nil {
		return err
	}

	if cc.check_participant(ctx, instanceID, msg.ReceiveParticipantID) == false{
		errorMessage := fmt.Sprintf("Participant %s is not allowed to send the message", msg.SendParticipantID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	if msg.MsgState != WAITINGFORCONFIRMATION {
		errorMessage := fmt.Sprintf("Event state %s is not allowed", msg.MessageID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	cc.ChangeMsgState(ctx, instanceID, msg.MessageID, COMPLETED)
	stub.SetEvent("Message_0m9p3da", []byte("Message has been done"))

	
	    cc.ChangeGtwState(ctx, instanceID, "ExclusiveGateway_0nzwv7v", ENABLED)

	
	return nil
}

func (cc *SmartContract) ExclusiveGateway_0nzwv7v(ctx contractapi.TransactionContextInterface, instanceID string) error {
	stub := ctx.GetStub()
	gtw, err := cc.ReadGtw(ctx, instanceID, "ExclusiveGateway_0nzwv7v")
	if err != nil {
		return err
	}

	if gtw.GatewayState != ENABLED {
		errorMessage := fmt.Sprintf("Gateway state %s is not allowed", gtw.GatewayID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	cc.ChangeGtwState(ctx, instanceID, gtw.GatewayID, COMPLETED)
	stub.SetEvent("ExclusiveGateway_0nzwv7v", []byte("ExclusiveGateway has been done"))

    
    	currentMemory, err := cc.ReadGlobalVariable(ctx, instanceID)
	if err != nil {
		return err
	}

    Invoice:=currentMemory.Invoice

if Invoice==false {
	    cc.ChangeEventState(ctx, instanceID, "Event_08edp7f", ENABLED)
}
if Invoice==true {
	    cc.ChangeMsgState(ctx, instanceID, "Message_1joj7ca", ENABLED)
}
    

	return nil
}

func (cc *SmartContract) Event_08edp7f(ctx contractapi.TransactionContextInterface, instanceID string) error {
	stub := ctx.GetStub()
	event, err := cc.ReadEvent(ctx, instanceID, "Event_08edp7f")
	if err != nil {
		return err
	}

	if event.EventState != ENABLED {
		errorMessage := fmt.Sprintf("Event state %s is not allowed", event.EventID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	cc.ChangeEventState(ctx, instanceID, event.EventID, COMPLETED) 
	stub.SetEvent("Event_08edp7f", []byte("EndEvent has been done"))
	
	return nil
}

func (cc *SmartContract) Message_1joj7ca_Send(ctx contractapi.TransactionContextInterface, instanceID string, fireflyTranID string ) error {
	stub := ctx.GetStub()
	msg, err := cc.ReadMsg(ctx, instanceID, "Message_1joj7ca")
	if err != nil {
		return err
	}

	//
	if cc.check_participant(ctx, instanceID, msg.SendParticipantID) == false{
		errorMessage := fmt.Sprintf("Participant %s is not allowed to send the message", msg.SendParticipantID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	if msg.MsgState != ENABLED {
		errorMessage := fmt.Sprintf("Message state %s is not allowed", msg.MessageID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	cc.ChangeMsgFireflyTranID(ctx, instanceID, fireflyTranID, msg.MessageID)
	cc.ChangeMsgState(ctx, instanceID, msg.MessageID, WAITINGFORCONFIRMATION)
	
	stub.SetEvent("Message_1joj7ca", []byte("Message is waiting for confirmation"))

	
	return nil
}

func (cc *SmartContract) Message_1joj7ca_Complete(ctx contractapi.TransactionContextInterface, instanceID string) error {
	stub := ctx.GetStub()
	msg, err := cc.ReadMsg(ctx, instanceID, "Message_1joj7ca")
	if err != nil {
		return err
	}

	if cc.check_participant(ctx, instanceID, msg.ReceiveParticipantID) == false{
		errorMessage := fmt.Sprintf("Participant %s is not allowed to send the message", msg.SendParticipantID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	if msg.MsgState != WAITINGFORCONFIRMATION {
		errorMessage := fmt.Sprintf("Event state %s is not allowed", msg.MessageID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	cc.ChangeMsgState(ctx, instanceID, msg.MessageID, COMPLETED)
	stub.SetEvent("Message_1joj7ca", []byte("Message has been done"))

	
	    cc.ChangeMsgState(ctx, instanceID, "Message_1etcmvl", ENABLED)

	
	return nil
}

func (cc *SmartContract) Message_1etcmvl_Send(ctx contractapi.TransactionContextInterface, instanceID string, fireflyTranID string ) error {
	stub := ctx.GetStub()
	msg, err := cc.ReadMsg(ctx, instanceID, "Message_1etcmvl")
	if err != nil {
		return err
	}

	//
	if cc.check_participant(ctx, instanceID, msg.SendParticipantID) == false{
		errorMessage := fmt.Sprintf("Participant %s is not allowed to send the message", msg.SendParticipantID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	if msg.MsgState != ENABLED {
		errorMessage := fmt.Sprintf("Message state %s is not allowed", msg.MessageID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	cc.ChangeMsgFireflyTranID(ctx, instanceID, fireflyTranID, msg.MessageID)
	cc.ChangeMsgState(ctx, instanceID, msg.MessageID, WAITINGFORCONFIRMATION)
	
	stub.SetEvent("Message_1etcmvl", []byte("Message is waiting for confirmation"))

	
	return nil
}

func (cc *SmartContract) Message_1etcmvl_Complete(ctx contractapi.TransactionContextInterface, instanceID string) error {
	stub := ctx.GetStub()
	msg, err := cc.ReadMsg(ctx, instanceID, "Message_1etcmvl")
	if err != nil {
		return err
	}

	if cc.check_participant(ctx, instanceID, msg.ReceiveParticipantID) == false{
		errorMessage := fmt.Sprintf("Participant %s is not allowed to send the message", msg.SendParticipantID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	if msg.MsgState != WAITINGFORCONFIRMATION {
		errorMessage := fmt.Sprintf("Event state %s is not allowed", msg.MessageID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	cc.ChangeMsgState(ctx, instanceID, msg.MessageID, COMPLETED)
	stub.SetEvent("Message_1etcmvl", []byte("Message has been done"))

	
	    cc.ChangeEventState(ctx, instanceID, "Event_146eii4", ENABLED)

	
	return nil
}

func (cc *SmartContract) Event_146eii4(ctx contractapi.TransactionContextInterface, instanceID string) error {
	stub := ctx.GetStub()
	event, err := cc.ReadEvent(ctx, instanceID, "Event_146eii4")
	if err != nil {
		return err
	}

	if event.EventState != ENABLED {
		errorMessage := fmt.Sprintf("Event state %s is not allowed", event.EventID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	cc.ChangeEventState(ctx, instanceID, event.EventID, COMPLETED) 
	stub.SetEvent("Event_146eii4", []byte("EndEvent has been done"))
	
	return nil
}

func (cc *SmartContract) Message_1i8rlqn_Send(ctx contractapi.TransactionContextInterface, instanceID string, fireflyTranID string ) error {
	stub := ctx.GetStub()
	msg, err := cc.ReadMsg(ctx, instanceID, "Message_1i8rlqn")
	if err != nil {
		return err
	}

	//
	if cc.check_participant(ctx, instanceID, msg.SendParticipantID) == false{
		errorMessage := fmt.Sprintf("Participant %s is not allowed to send the message", msg.SendParticipantID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	if msg.MsgState != ENABLED {
		errorMessage := fmt.Sprintf("Message state %s is not allowed", msg.MessageID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	cc.ChangeMsgFireflyTranID(ctx, instanceID, fireflyTranID, msg.MessageID)
	cc.ChangeMsgState(ctx, instanceID, msg.MessageID, WAITINGFORCONFIRMATION)
	
	stub.SetEvent("Message_1i8rlqn", []byte("Message is waiting for confirmation"))

	
	return nil
}

func (cc *SmartContract) Message_1i8rlqn_Complete(ctx contractapi.TransactionContextInterface, instanceID string) error {
	stub := ctx.GetStub()
	msg, err := cc.ReadMsg(ctx, instanceID, "Message_1i8rlqn")
	if err != nil {
		return err
	}

	if cc.check_participant(ctx, instanceID, msg.ReceiveParticipantID) == false{
		errorMessage := fmt.Sprintf("Participant %s is not allowed to send the message", msg.SendParticipantID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	if msg.MsgState != WAITINGFORCONFIRMATION {
		errorMessage := fmt.Sprintf("Event state %s is not allowed", msg.MessageID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	cc.ChangeMsgState(ctx, instanceID, msg.MessageID, COMPLETED)
	stub.SetEvent("Message_1i8rlqn", []byte("Message has been done"))

	
	    cc.ChangeMsgState(ctx, instanceID, "Message_1q05nnw", ENABLED)

	
	return nil
}

func (cc *SmartContract) Gateway_1bhtapl(ctx contractapi.TransactionContextInterface, instanceID string) error {
	stub := ctx.GetStub()
	gtw, err := cc.ReadGtw(ctx, instanceID, "Gateway_1bhtapl")
	if err != nil {
		return err
	}

	if gtw.GatewayState != ENABLED {
		errorMessage := fmt.Sprintf("Gateway state %s is not allowed", gtw.GatewayID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	cc.ChangeGtwState(ctx, instanceID, gtw.GatewayID, COMPLETED)
	stub.SetEvent("Gateway_1bhtapl", []byte("ExclusiveGateway has been done"))

    
    	currentMemory, err := cc.ReadGlobalVariable(ctx, instanceID)
	if err != nil {
		return err
	}

    Need_external_provider:=currentMemory.Need_external_provider

if Need_external_provider==true {
	    cc.ChangeMsgState(ctx, instanceID, "Message_1i8rlqn", ENABLED)
}
if Need_external_provider==false {
	    cc.ChangeGtwState(ctx, instanceID, "Gateway_04h9e6e", ENABLED)
}
    

	return nil
}

func (cc *SmartContract) Gateway_04h9e6e(ctx contractapi.TransactionContextInterface, instanceID string) error {
	stub := ctx.GetStub()
	gtw, err := cc.ReadGtw(ctx, instanceID, "Gateway_04h9e6e")
	if err != nil {
		return err
	}

	if gtw.GatewayState != ENABLED {
		errorMessage := fmt.Sprintf("Gateway state %s is not allowed", gtw.GatewayID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	cc.ChangeGtwState(ctx, instanceID, gtw.GatewayID, COMPLETED)
	stub.SetEvent("Gateway_04h9e6e", []byte("ExclusiveGateway has been done"))

    
        cc.ChangeMsgState(ctx, instanceID, "Message_1ljlm4g", ENABLED)
    

	return nil
}

func (cc *SmartContract) Message_1q05nnw_Send(ctx contractapi.TransactionContextInterface, instanceID string, fireflyTranID string ) error {
	stub := ctx.GetStub()
	msg, err := cc.ReadMsg(ctx, instanceID, "Message_1q05nnw")
	if err != nil {
		return err
	}

	//
	if cc.check_participant(ctx, instanceID, msg.SendParticipantID) == false{
		errorMessage := fmt.Sprintf("Participant %s is not allowed to send the message", msg.SendParticipantID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	if msg.MsgState != ENABLED {
		errorMessage := fmt.Sprintf("Message state %s is not allowed", msg.MessageID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	cc.ChangeMsgFireflyTranID(ctx, instanceID, fireflyTranID, msg.MessageID)
	cc.ChangeMsgState(ctx, instanceID, msg.MessageID, WAITINGFORCONFIRMATION)
	
	stub.SetEvent("Message_1q05nnw", []byte("Message is waiting for confirmation"))

	
	return nil
}

func (cc *SmartContract) Message_1q05nnw_Complete(ctx contractapi.TransactionContextInterface, instanceID string) error {
	stub := ctx.GetStub()
	msg, err := cc.ReadMsg(ctx, instanceID, "Message_1q05nnw")
	if err != nil {
		return err
	}

	if cc.check_participant(ctx, instanceID, msg.ReceiveParticipantID) == false{
		errorMessage := fmt.Sprintf("Participant %s is not allowed to send the message", msg.SendParticipantID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	if msg.MsgState != WAITINGFORCONFIRMATION {
		errorMessage := fmt.Sprintf("Event state %s is not allowed", msg.MessageID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	cc.ChangeMsgState(ctx, instanceID, msg.MessageID, COMPLETED)
	stub.SetEvent("Message_1q05nnw", []byte("Message has been done"))

	
	    cc.ChangeMsgState(ctx, instanceID, "Message_1qbk325", ENABLED)

	
	return nil
}

func (cc *SmartContract) Message_1qbk325_Send(ctx contractapi.TransactionContextInterface, instanceID string, fireflyTranID string ) error {
	stub := ctx.GetStub()
	msg, err := cc.ReadMsg(ctx, instanceID, "Message_1qbk325")
	if err != nil {
		return err
	}

	//
	if cc.check_participant(ctx, instanceID, msg.SendParticipantID) == false{
		errorMessage := fmt.Sprintf("Participant %s is not allowed to send the message", msg.SendParticipantID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	if msg.MsgState != ENABLED {
		errorMessage := fmt.Sprintf("Message state %s is not allowed", msg.MessageID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	cc.ChangeMsgFireflyTranID(ctx, instanceID, fireflyTranID, msg.MessageID)
	cc.ChangeMsgState(ctx, instanceID, msg.MessageID, WAITINGFORCONFIRMATION)
	
	stub.SetEvent("Message_1qbk325", []byte("Message is waiting for confirmation"))

	
	return nil
}

func (cc *SmartContract) Message_1qbk325_Complete(ctx contractapi.TransactionContextInterface, instanceID string) error {
	stub := ctx.GetStub()
	msg, err := cc.ReadMsg(ctx, instanceID, "Message_1qbk325")
	if err != nil {
		return err
	}

	if cc.check_participant(ctx, instanceID, msg.ReceiveParticipantID) == false{
		errorMessage := fmt.Sprintf("Participant %s is not allowed to send the message", msg.SendParticipantID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	if msg.MsgState != WAITINGFORCONFIRMATION {
		errorMessage := fmt.Sprintf("Event state %s is not allowed", msg.MessageID)
		fmt.Println(errorMessage)
		return fmt.Errorf(errorMessage)
	}

	cc.ChangeMsgState(ctx, instanceID, msg.MessageID, COMPLETED)
	stub.SetEvent("Message_1qbk325", []byte("Message has been done"))

	
	    cc.ChangeGtwState(ctx, instanceID, "Gateway_04h9e6e", ENABLED)

	
	return nil
}