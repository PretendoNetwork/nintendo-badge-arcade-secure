package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gocql/gocql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var cluster *gocql.ClusterConfig
var cassandraClusterSession *gocql.Session

var mongoClient *mongo.Client
var mongoContext context.Context

var accountDatabase *mongo.Database
var badgeArcadeDatabase *mongo.Database
var nexAccountsCollection *mongo.Collection
var sessionsCollection *mongo.Collection

func connectMongo() {
	mongoClient, _ = mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	mongoContext, _ = context.WithTimeout(context.Background(), 10*time.Second)
	_ = mongoClient.Connect(mongoContext)

	accountDatabase = mongoClient.Database("pretendo")
	nexAccountsCollection = accountDatabase.Collection("nexaccounts")

	badgeArcadeDatabase = mongoClient.Database("badge_arcade")
	sessionsCollection = badgeArcadeDatabase.Collection("sessions")

	sessionsCollection.DeleteMany(context.TODO(), bson.D{})
}

func connectCassandra() {
	// Connect to Cassandra

	var err error

	cluster = gocql.NewCluster(os.Getenv("CASSANDRA_ADDRESS"))
	cluster.Timeout = 30 * time.Second

	createKeyspace("pretendo_badge_arcade")

	cluster.Keyspace = "pretendo_badge_arcade"

	cassandraClusterSession, err = cluster.CreateSession()

	if err != nil {
		panic(err)
	}

	// Create tables if missing
	
	if err := cassandraClusterSession.Query(`CREATE TABLE IF NOT EXISTS pretendo_badge_arcade.free_play_data (
			data_id bigint PRIMARY KEY,
			owner_id int,
			meta_binary blob,
			created_time bigint,
			updated_time bigint,
			period smallint,
			flag int,
			referred_time bigint
		)`).Exec(); err != nil {
		fmt.Println("pretendo_badge_arcade.free_play_data")
		log.Fatal(err)
	}

	if err := cassandraClusterSession.Query(`CREATE TABLE IF NOT EXISTS pretendo_badge_arcade.persistence_info (
			data_id bigint PRIMARY KEY,
			pid int,
			slot smallint
		)`).Exec(); err != nil {
		fmt.Println("pretendo_badge_arcade.persistence_info")
		log.Fatal(err)
	}	

	if err := cassandraClusterSession.Query(`CREATE TABLE IF NOT EXISTS pretendo_badge_arcade.user_play_info (
		data_id bigint PRIMARY KEY,
		version int,
		size int
	)`).Exec(); err != nil {
		fmt.Println("pretendo_badge_arcade.user_play_info")
		log.Fatal(err)
	}

	fmt.Println("Connected to Cassandra")
}

// Adapted from gocql common_test.go
func createKeyspace(keyspace string) {
	flagRF := flag.Int("rf", 1, "replication factor for pretendo_badge_arcade keyspace")

	c := *cluster
	c.Keyspace = "system"
	c.Timeout = 30 * time.Second

	s, err := c.CreateSession()

	if err != nil {
		panic(err)
	}

	defer s.Close()

	if err := s.Query(fmt.Sprintf(`CREATE KEYSPACE IF NOT EXISTS %s
	WITH replication = {
		'class' : 'SimpleStrategy',
		'replication_factor' : %d
	}`, keyspace, *flagRF)).Exec(); err != nil {
		log.Fatal(err)
	}
}

////////////////////////////////
//                            //
// Cassandra database methods //
//                            //
////////////////////////////////

func getDataStorePersistenceInfo(ownerID uint32, persistenceSlotID uint16) uint64 {
	var dataID uint64
	_ = cassandraClusterSession.Query(`SELECT data_id FROM pretendo_badge_arcade.persistence_info WHERE pid=? AND slot=? ALLOW FILTERING`, ownerID, persistenceSlotID).Scan(&dataID)

	return dataID
}

func getVersionByDataID(dataID uint64) uint32 {
	var version uint32
	_ = cassandraClusterSession.Query(`SELECT version FROM pretendo_badge_arcade.user_play_info WHERE data_id=? ALLOW FILTERING`, dataID).Scan(&version)

	return version
}

func getSizeByDataID(dataID uint64) uint32 {
	var size uint32
	_ = cassandraClusterSession.Query(`SELECT size FROM pretendo_badge_arcade.user_play_info WHERE data_id=? ALLOW FILTERING`, dataID).Scan(&size)

	return size
}

func getFreePlayDataMetaInfoByOwnerID(ownerID uint32) (uint64, []byte, uint64, uint64, uint16, uint32, uint64)  {
	var dataID uint64
	var metaBinary []byte
	var createdTime uint64
	var updatedTime uint64
	var period uint16
	var flag uint32
	var referredTime uint64
	_ = cassandraClusterSession.Query(`SELECT data_id, meta_binary, created_time, updated_time, period, flag, referred_time FROM pretendo_badge_arcade.free_play_data WHERE owner_id=? ALLOW FILTERING`, ownerID).Scan(&dataID, &metaBinary, &createdTime, &updatedTime, &period, &flag, &referredTime)

	return dataID, metaBinary, createdTime, updatedTime, period, flag, referredTime
}

func postFreePlayDataMetaInfo(dataID uint64, ownerID uint32, metaBinary []byte, createdTime uint64, period uint16, flag uint32) {
	_ = cassandraClusterSession.Query(`INSERT INTO pretendo_badge_arcade.free_play_data(
		data_id,
		owner_id,
		meta_binary,
		created_time,
		updated_time,
		period,
		flag,
		referred_time
	)
	VALUES (
		?,
		?,
		?,
		?,
		?,
		?,
		?,
		?
	) IF NOT EXISTS`, dataID, ownerID, metaBinary, createdTime, createdTime, period, flag, createdTime).Exec()
}

func postPersistenceInfo(dataID uint64, ownerID uint32, slot uint16) {
	_ = cassandraClusterSession.Query(`INSERT INTO pretendo_badge_arcade.persistence_info(
		data_id,
		pid,
		slot
	)
	VALUES (
		?,
		?,
		?
	) IF NOT EXISTS`, dataID, ownerID, slot).Exec()
}

func postUserPlayInfo(dataID uint64, version uint32, size uint32) {
	_ = cassandraClusterSession.Query(`INSERT INTO pretendo_badge_arcade.user_play_info(
		data_id,
		version,
		size
	)
	VALUES (
		?,
		?,
		?
	) IF NOT EXISTS`, dataID, version, size).Exec()
}

func updateFreePlayDataMetaBinary(dataID uint64, metaBinary []byte, updatedTime uint64) {
	_ = cassandraClusterSession.Query(`UPDATE pretendo_badge_arcade.free_play_data SET meta_binary=?, updated_time=? WHERE data_id=? ALLOW FILTERING`, metaBinary, updatedTime, dataID).Exec()
}

func updateUserPlayInfoSize(dataID uint64, size uint32) {
	_ = cassandraClusterSession.Query(`UPDATE pretendo_badge_arcade.user_play_info SET size=? WHERE data_id=? ALLOW FILTERING`, size, dataID).Exec()
}

func updateUserPlayInfoVersion(dataID uint64, version uint32) {
	_ = cassandraClusterSession.Query(`UPDATE pretendo_badge_arcade.user_play_info SET version=? WHERE data_id=? ALLOW FILTERING`, version, dataID).Exec()
}

//////////////////////////////
//                          //
// MongoDB database methods //
//                          //
//////////////////////////////

func addPlayerSession(pid uint32, urls []string, ip string, port string) {
	_, err := sessionsCollection.InsertOne(context.TODO(), bson.D{{"pid", pid}, {"urls", urls}, {"ip", ip}, {"port", port}})
	if err != nil {
		panic(err)
	}
}

func getAllSessionPIDs() []uint32 {
	var result []bson.M
	output := []uint32{}

	c, _ := sessionsCollection.Find(context.TODO(), bson.D{})
	c.All(context.TODO(), &result)
	for _, i := range result {
		output = append(output, uint32(i["pid"].(int64)))
	}
	return output
}

func doesSessionExist(pid uint32) bool {
	var result bson.M

	err := sessionsCollection.FindOne(context.TODO(), bson.D{{"pid", pid}}, options.FindOne()).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		} else {
			panic(err)
		}
	} else {
		return true
	}
}

func updatePlayerSessionAll(pid uint32, urls []string, ip string, port string) {
	_, err := sessionsCollection.UpdateOne(context.TODO(), bson.D{{"pid", pid}}, bson.D{{"$set", bson.D{{"pid", pid}, {"urls", urls}, {"ip", ip}, {"port", port}}}})
	if err != nil {
		panic(err)
	}
}

func updatePlayerSessionPort(pid uint32, port string) {
	_, err := sessionsCollection.UpdateOne(context.TODO(), bson.D{{"pid", pid}}, bson.D{{"$set", bson.D{{"port", port}}}})
	if err != nil {
		panic(err)
	}
}

func deletePlayerSession(pid uint32) {
	_, err := sessionsCollection.DeleteOne(context.TODO(), bson.D{{"pid", pid}})
	if err != nil {
		panic(err)
	}
}
