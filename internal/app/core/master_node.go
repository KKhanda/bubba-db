package core

import (
	"log"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

var master Entity

func getFilePath(fileName string) string {
	path, err := filepath.Abs("./" + fileName)
	if err != nil {
		log.Panic("Error at master_node getFilePath")
	}
	pathElems := strings.Split(path, "/")
	res := 0
	for i := range pathElems {
		if pathElems[i] == "databases" {
			res = i
		}
	}
	return "/" + strings.Join(pathElems[res:],"/")
}

func SendReadData(entity *Entity, file *os.File, offset int, id int, bs *[]byte) {
	var reply Reply
	var attempts = 0
	fileAbsPath := getFilePath(file.Name())
	requestedData := RequestedData{
		File: fileAbsPath,
		Offset: offset,
		Bs: *bs,
		Id: id,
	}
	for attempts < 5 {
		err = nil
		request := RPCRequest {requestedData }
		err = entity.Connector.Call("Entity.Read", &request, &reply)
		if err != nil {
			log.Panic("Error in master_node SendReadData ", err)
			err = errors.New("problems in requestSlaveStatus")
			attempts++
			continue
		}
		if reply.Message == "success" {
			attempts = 5
			*bs = reply.Data
		}
	}
}

func SendWriteData(entity *Entity, file *os.File, offset int, id int, bs []byte) error {
	var reply Reply
	var attempts = 0

	fileAbsPath := getFilePath(file.Name())
	requestedData := RequestedData{
		File: fileAbsPath,
		Offset: offset,
		Id: id,
		Bs: bs,
	}
	for attempts < 5 {
		err = nil
		request := RPCRequest{ requestedData }
		err = entity.Connector.Call("Entity.Write", &request, &reply)
		if err != nil {
			log.Panic("Error in master_node SendWriteData ", err)
			err = errors.New("problems in requestSlaveStatus")
			attempts++
			continue
		}
		if reply.Message == "success" {
			attempts = 5
		}
	}
	return nil
}

func SendSwitchDatabaseStructure(entity *Entity, newStructure *string) error {
	var reply Reply
	var attempts = 0
	for attempts < 5 {
		err = nil
		request := RPCRequest{ RequestedData{ Payload: *newStructure } }
		err = entity.Connector.Call("Entity.SwitchDatabaseStructure", &request, &reply)
		if err != nil {
			log.Panic("Error in master_node SendSwitchDatabaseStructure ", err)
			err = errors.New("problems in requestSlaveStatus")
			attempts++
			continue
		}
		if reply.Message == "success" {
			println("Slave " + entity.Ip + ":" + entity.Port + " switched status on " + *newStructure)
			attempts = 5
		}
	}
	return nil
}

func RequestSlaveStatus(entity *Entity) error {
	var reply Reply
	var attempts = 0
	for attempts < 5 {
		err = nil
		var request RPCRequest
		err = entity.Connector.Call("Entity.SendStatus", &request, &reply)
		if err != nil {
			log.Panic("Error in master_node RequestSlaveStatus ", err)
			err = errors.New("problems in requestSlaveStatus")
			attempts++
			continue
		}
		if reply.Message == "success" {
			println("Slave " + entity.Ip + ":" + entity.Port + " is available")
			attempts = 5
		}
	}
	return nil
}

func SendDeploy(entity *Entity) error {
	var reply Reply
	var attempts = 0
	for attempts < 5 {
		err = nil
		request := RPCRequest{*new(RequestedData) }
		err = entity.Connector.Call("Entity.Deploy", &request, &reply)
		if err != nil {
			log.Panic("Error in master_node SendDeploy ", err)
			err = errors.New("problems in requestSlaveStatus")
			attempts++
			continue
		}
		if reply.Message == "success" {
			attempts = 5
		}
	}
	return err
}

func SendInitDatabaseStructure(entity *Entity, dbName *string) error {
	var reply Reply
	var attempts = 0
	for attempts < 5 {
		log.Printf("Try to SendInitDatabaseStructure (attempts %d) to %s:%s\n", attempts, entity.Ip, entity.Port)
		err = nil
		var requestedData = RequestedData{ Payload: *dbName}
		request := RPCRequest{ requestedData }
		err = entity.Connector.Call("Entity.InitDatabaseStructure", &request, &reply)
		if err != nil {
			log.Panic("Error in master_node SendInitDatabaseStructure ", err)
			err = errors.New("problems in requestSlaveStatus")
			attempts++
			continue
		}
		if reply.Message == "success" {
			attempts = 5
		}
	}
	return nil
}

func SendDropDatabase(entity *Entity, dbName *string) error {
	var reply Reply
	var attempts = 0
	for attempts < 5 {
		err = nil
		request := RPCRequest{ RequestedData{ Payload: *dbName } }
		err = entity.Connector.Call("Entity.DropDatabase", &request, &reply)
		if err != nil {
			log.Panic("Error in master_node SendDropDatabase ", err)
			err = errors.New("problems in requestSlaveStatus")
			attempts++
			continue
		}
		if reply.Message == "success" {
			attempts = 5
		}
	}
	return err
}