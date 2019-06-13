package utils

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strings"
	"time"
)

const (
	// MasterSession provides direct access to master database.
	MasterSession = "master"

	// MonotonicSession provides reads to slaves.
	MonotonicSession = "monotonic"
)

var (
	// Reference to the singleton.
	singleton mongoManager
)

type (
	// mongoConfiguration contains settings for initialization.
	mongoConfiguration struct {
		Hosts    string
		Database string
		UserName string
		Password string
	}

	// mongoManager contains dial and session information.
	mongoSession struct {
		mongoDBDialInfo *mgo.DialInfo
		mongoSession    *mgo.Session
	}

	// mongoManager manages a map of session.
	mongoManager struct {
		sessions map[string]mongoSession
	}

	// DBCall defines a type of function that can be used
	// to excecute code against MongoDB.
	DBCall func(*mgo.Collection) error
)

// Startup brings the manager to a running state.
func Startup(sessionID string) error {
	// If the system has already been started ignore the call.
	if singleton.sessions != nil {
		return nil
	}

	logs.Info(sessionID, "Startup")

	// Pull in the configuration.
	var config mongoConfiguration
	if err := envconfig.Process("mgo", &config); err != nil {
		logs.Error(err, sessionID, "Startup")
		return err
	}

	// Create the Mongo Manager.
	singleton = mongoManager{
		sessions: make(map[string]mongoSession),
	}
	if len(config.Hosts) <= 0 || config.Hosts == "" || strings.Count(config.Hosts, " ") == len(config.Hosts) {
		config.Hosts = beego.AppConfig.String("mgoHost")
		config.Database = beego.AppConfig.String("mgoDatabase")
	}

	// Log the mongodb connection straps.
	logs.Info(sessionID, "Startup", fmt.Sprintf("MongoDB : Hosts[%s]", config.Hosts))
	logs.Info(sessionID, "Startup", fmt.Sprintf("MongoDB : Database[%s]", config.Database))
	logs.Info(sessionID, "Startup", fmt.Sprintf("MongoDB : Username[%s]", config.UserName))

	hosts := strings.Split(config.Hosts, ",")

	// Create the strong session.
	if err := CreateSession(sessionID, "strong", MasterSession, hosts, config.Database, config.UserName, config.Password); err != nil {
		logs.Error(err, sessionID, "Startup")
		return err
	}

	// Create the monotonic session.
	if err := CreateSession(sessionID, "monotonic", MonotonicSession, hosts, config.Database, config.UserName, config.Password); err != nil {
		logs.Error(err, sessionID, "Startup")
		return err
	}
	initRedisConnect()
	logs.Info(sessionID, "Startup")
	return nil
}

func Shutdown(sessionID string) error {
	logs.Info(sessionID, "Shutdown")

	// Close the databases
	for _, session := range singleton.sessions {
		CloseSession(sessionID, session.mongoSession)
	}
	shutdownRedisConnect()
	logs.Info(sessionID, "Shutdown")
	return nil
}

// CreateSession creates a connection pool for use.
func CreateSession(sessionID string, mode string, sessionName string, hosts []string, databaseName string, username string, password string) error {
	logs.Debug(sessionID, "CreateSession", fmt.Sprintf("Mode[%s] SessionName[%s] Hosts[%s] DatabaseName[%s] Username[%s]", mode, sessionName, hosts, databaseName, username))
	var mSession mongoSession
	var err error
	if (len(username) == 0 || username == "") && len(hosts) > 0 {
		url := strings.Join(hosts, ",")
		mSession.mongoSession, err = mgo.Dial(url)
	} else {
		// Create the database object
		mSession = mongoSession{
			mongoDBDialInfo: &mgo.DialInfo{
				Addrs:    hosts,
				Timeout:  60 * time.Second,
				Database: databaseName,
				Username: username,
				Password: password,
			},
		}
		// Establish the master session.
		mSession.mongoSession, err = mgo.DialWithInfo(mSession.mongoDBDialInfo)
		if err != nil {
			logs.Error(err, sessionID, "CreateSession")
			return err
		}
	}

	switch mode {
	case "strong":
		// Reads and writes will always be made to the master server using a
		// unique connection so that reads and writes are fully consistent,
		// ordered, and observing the most up-to-date data.
		// http://godoc.org/github.com/finapps/mgo#Session.SetMode
		mSession.mongoSession.SetMode(mgo.Strong, true)
		break

	case "monotonic":
		// Reads may not be entirely up-to-date, but they will always see the
		// history of changes moving forward, the data read will be consistent
		// across sequential queries in the same session, and modifications made
		// within the session will be observed in following queries (read-your-writes).
		// http://godoc.org/github.com/finapps/mgo#Session.SetMode
		mSession.mongoSession.SetMode(mgo.Monotonic, true)
	}

	// Have the session check for errors.
	// http://godoc.org/github.com/finapps/mgo#Session.SetSafe
	mSession.mongoSession.SetSafe(&mgo.Safe{})

	// Add the database to the map.
	singleton.sessions[sessionName] = mSession

	logs.Debug(sessionID, "CreateSession")
	return nil
}

// CopyMasterSession makes a copy of the master session for client use.
func CopyMasterSession(sessionID string) (*mgo.Session, error) {
	return CopySession(sessionID, MasterSession)
}

// CopyMonotonicSession makes a copy of the monotonic session for client use.
func CopyMonotonicSession(sessionID string) (*mgo.Session, error) {
	return CopySession(sessionID, MonotonicSession)
}

// CopySession makes a copy of the specified session for client use.
func CopySession(sessionID string, useSession string) (*mgo.Session, error) {
	logs.Debug(sessionID, "CopySession", "UseSession[%s]", useSession)

	// Find the session object.
	session := singleton.sessions[useSession]

	if session.mongoSession == nil {
		err := fmt.Errorf("Unable To Locate Session %s", useSession)
		logs.Error(err, sessionID, "CopySession")
		return nil, err
	}

	// Copy the master session.
	mongoSession := session.mongoSession.Copy()

	logs.Debug(sessionID, "CopySession")
	return mongoSession, nil
}

// CloneMasterSession makes a clone of the master session for client use.
func CloneMasterSession(sessionID string) (*mgo.Session, error) {
	return CloneSession(sessionID, MasterSession)
}

// CloneMonotonicSession makes a clone of the monotinic session for client use.
func CloneMonotonicSession(sessionID string) (*mgo.Session, error) {
	return CloneSession(sessionID, MonotonicSession)
}

// CloneSession makes a clone of the specified session for client use.
func CloneSession(sessionID string, useSession string) (*mgo.Session, error) {
	logs.Debug(sessionID, "CloneSession", "UseSession[%s]", useSession)

	// Find the session object.
	session := singleton.sessions[useSession]

	if session.mongoSession == nil {
		err := fmt.Errorf("Unable To Locate Session %s", useSession)
		logs.Error(err, sessionID, "CloneSession")
		return nil, err
	}

	// Clone the master session.
	mongoSession := session.mongoSession.Clone()

	logs.Debug(sessionID, "CloneSession")
	return mongoSession, nil
}

// CloseSession puts the connection back into the pool.
func CloseSession(sessionID string, mongoSession *mgo.Session) {
	logs.Debug(sessionID, "CloseSession")
	mongoSession.Close()
	logs.Info(sessionID, "CloseSession")
}

// GetDatabase returns a reference to the specified database.
func GetDatabase(mongoSession *mgo.Session, useDatabase string) *mgo.Database {
	return mongoSession.DB(useDatabase)
}

// GetCollection returns a reference to a collection for the specified database and collection name.
func GetCollection(mongoSession *mgo.Session, useDatabase string, useCollection string) *mgo.Collection {
	return mongoSession.DB(useDatabase).C(useCollection)
}

// CollectionExists returns true if the collection name exists in the specified database.
func CollectionExists(sessionID string, mongoSession *mgo.Session, useDatabase string, useCollection string) bool {
	database := mongoSession.DB(useDatabase)
	collections, err := database.CollectionNames()

	if err != nil {
		return false
	}

	for _, collection := range collections {
		if collection == useCollection {
			return true
		}
	}

	return false
}

// ToString converts the quer map to a string.
func ToString(queryMap interface{}) string {
	json, err := json.Marshal(queryMap)
	if err != nil {
		return ""
	}

	return string(json)
}

// ToStringD converts bson.D to a string.
func ToStringD(queryMap bson.D) string {
	json, err := json.Marshal(queryMap)
	if err != nil {
		return ""
	}

	return string(json)
}

// Execute the MongoDB literal function.
func Execute(sessionID string, mongoSession *mgo.Session, databaseName string, collectionName string, dbCall DBCall) error {
	logs.Debug(sessionID, "Execute", "Database[%s] Collection[%s]", databaseName, collectionName)

	// Capture the specified collection.
	collection := GetCollection(mongoSession, databaseName, collectionName)
	if collection == nil {
		err := fmt.Errorf("Collection %s does not exist", collectionName)
		logs.Error(err, sessionID, "Execute")
		return err
	}

	// Execute the MongoDB call.
	err := dbCall(collection)
	if err != nil {
		logs.Error(err, sessionID, "Execute")
		return err
	}

	logs.Debug(sessionID, "Execute")
	return nil
}
