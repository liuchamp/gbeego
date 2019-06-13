package services

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	mongo "mul/utils"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"mul/models"
)

type userConfiguration struct {
	Database string
}

var Config userConfiguration

func init() {
	// Pull in the configuration.
	if err := envconfig.Process("user", &Config); err != nil {
		logs.Error(err, mongo.MainGoRoutine, "Init")
	}
	if isBlank := mongo.CheckStringIsBlank(Config.Database); isBlank {
		Config.Database = beego.AppConfig.String("userDatabase")
	}
}

func FindUserById(service *Service, userId string) (*models.User, error) {

	var user models.User
	if vx := mongo.CheckStringIsBlank(userId); vx {
		return nil, mgo.ErrNotFound
	}
	f := func(collection *mgo.Collection) error {
		queryMap := bson.M{"_id": userId}

		return collection.Find(queryMap).One(&user)
	}
	if err := service.DBAction(Config.Database, models.UserDb, f); err != nil {
		if err != nil {
			logs.Error(err, service.UserID, "FindUserById")
			return nil, err
		}
	}
	logs.Info(service.UserID, "findUserById", "users %+v", &user)
	return &user, nil
}

func MoveMoney(service *Service, td *models.TransferDTO) error {

	fromU, es := FindUserById(service, td.FromUserId)
	if es != nil {
		return es
	}
	fromT, es := FindUserById(service, td.ToUserId)
	if es != nil {
		return es
	}

	fromU.UserMoney -= td.Money
	fromT.UserMoney += td.Money
	updateUserU := func(collection *mgo.Collection) error {
		return collection.UpdateId(fromU.UserId, fromU)
	}
	updateUserT := func(collection *mgo.Collection) error {
		return collection.UpdateId(fromT.UserId, fromT)
	}
	if err := service.DBAction(Config.Database, models.UserDb, updateUserT); err != nil {
		if err != nil {
			logs.Error(err, service.UserID, "FindUserById")

			return err
		}
	}
	if err := service.DBAction(Config.Database, models.UserDb, updateUserU); err != nil {
		if err != nil {
			logs.Error(err, service.UserID, "FindUserById")
			// 操作失败后，回退
			return err
		}
	}

	return nil
}

// 创建用户
func CreateUser(service *Service, user *models.User) (*models.User, error) {
	// 验证数据合法
	// 生成id
	// 保存数据

	if se := genNewUser(user); se != nil {
		logs.Error("Check new User meeting error->", se)
		return nil, se
	}
	f := func(collection *mgo.Collection) error {

		return collection.Insert(user)
	}
	if err := service.DBAction(Config.Database, models.UserDb, f); err != nil {
		if err != mgo.ErrNotFound {
			logs.Error(err, service.UserID, "CreateUser")
			return nil, err
		}
	}
	var uq models.User
	q := func(collection *mgo.Collection) error {
		queryMap := bson.M{"_id": user.UserId}
		return collection.Find(queryMap).One(&uq)
	}
	if err := service.DBAction(Config.Database, models.UserDb, q); err != nil {
		if err != mgo.ErrNotFound {
			logs.Error(err, service.UserID, "CreateUser")
			return nil, err
		}
	}
	return &uq, nil
}

func genNewUser(user *models.User) error {
	if isb := mongo.CheckStringIsBlank(user.UserId); isb {
		user.UserId = ""
	}
	if user.UserMoney < 0 {
		user.UserMoney = 0
	}
	if isb := mongo.CheckStringIsBlank(user.UserId); isb {
		user.UserPhone = "+8559696996"
	}

	return nil
}
