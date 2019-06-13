package services

import (
	"github.com/astaxie/beego/logs"
	"gopkg.in/mgo.v2"
	mongo "mul/utils"
)

type Service struct {
	MongoSession *mgo.Session
	UserID       string
}

type ServiceAction interface {
	DBAction(databaseName string, collectionName string, dbCall mongo.DBCall) (err error)
}

func (service *Service) Prepare() (err error) {
	logs.Debug("service.UserID", service.UserID)
	service.MongoSession, err = mongo.CopyMonotonicSession(service.UserID)
	if err != nil {
		logs.Error(err, service.UserID, "Service.Prepare")
		return err
	}

	return err
}

func (service *Service) Finish() (err error) {
	// 捕获所有异常
	defer mongo.CatchPanic(&err, service.UserID, "Service.Finish")

	if service.MongoSession != nil {
		mongo.CloseSession(service.UserID, service.MongoSession)
		service.MongoSession = nil
	}

	return err
}

// DBAction executes the MongoDB literal function
func (service *Service) DBAction(databaseName string, collectionName string, dbCall mongo.DBCall) (err error) {
	return mongo.Execute(service.UserID, service.MongoSession, databaseName, collectionName, dbCall)
}
