package main

import (
	"database/sql"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var postgres *sql.DB

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

func connectPostgres() {
	// Connect to Postgres

	var err error

	postgres, err = sql.Open("postgres", os.Getenv("DATABASE_URI"))
	if err != nil {
		panic(err)
	}

	_, err = postgres.Exec(`CREATE SCHEMA IF NOT EXISTS pretendo_badge_arcade`)
	if err != nil {
		fmt.Println("pretendo_badge_arcade")
		log.Fatal(err)
	}

	// Create tables if missing
	
	_, err = postgres.Exec(`CREATE TABLE IF NOT EXISTS pretendo_badge_arcade.free_play_data (
			data_id bigint PRIMARY KEY,
			owner_id int,
			meta_binary bytea,
			created_time bigint,
			updated_time bigint,
			period smallint,
			flag int,
			referred_time bigint
		)`)
	if err != nil {
		fmt.Println("pretendo_badge_arcade.free_play_data")
		log.Fatal(err)
	}

	_, err = postgres.Exec(`CREATE TABLE IF NOT EXISTS pretendo_badge_arcade.persistence_info (
			data_id bigint PRIMARY KEY,
			pid int,
			slot smallint
		)`)
	if err != nil {
		fmt.Println("pretendo_badge_arcade.persistence_info")
		log.Fatal(err)
	}	

	_, err = postgres.Exec(`CREATE TABLE IF NOT EXISTS pretendo_badge_arcade.user_play_info (
		data_id bigint PRIMARY KEY,
		version int,
		size int
	)`)
	if err != nil {
		fmt.Println("pretendo_badge_arcade.user_play_info")
		log.Fatal(err)
	}

	fmt.Println("Connected to Postgres")
}

////////////////////////////////
//                            //
// Postgres database methods  //
//                            //
////////////////////////////////

func getDataStorePersistenceInfo(ownerID uint32, persistenceSlotID uint16) uint64 {
	var dataID uint64
	err := postgres.QueryRow(`SELECT data_id FROM pretendo_badge_arcade.persistence_info WHERE pid=$1 AND slot=$2`, ownerID, persistenceSlotID).Scan(&dataID)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)
	}

	return dataID
}

func getVersionByDataID(dataID uint64) uint32 {
	var version uint32
	err := postgres.QueryRow(`SELECT version FROM pretendo_badge_arcade.user_play_info WHERE data_id=$1`, dataID).Scan(&version)
	if err != nil {
		log.Fatal(err)
	}

	return version
}

func getSizeByDataID(dataID uint64) uint32 {
	var size uint32
	err := postgres.QueryRow(`SELECT size FROM pretendo_badge_arcade.user_play_info WHERE data_id=$1`, dataID).Scan(&size)
	if err != nil {
		log.Fatal(err)
	}

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
	err := postgres.QueryRow(`SELECT data_id, meta_binary, created_time, updated_time, period, flag, referred_time FROM pretendo_badge_arcade.free_play_data WHERE owner_id=$1`, ownerID).Scan(&dataID, &metaBinary, &createdTime, &updatedTime, &period, &flag, &referredTime)
	if err != nil && err != sql.ErrNoRows {
		log.Fatal(err)
	}

	return dataID, metaBinary, createdTime, updatedTime, period, flag, referredTime
}

func postFreePlayDataMetaInfo(dataID uint64, ownerID uint32, metaBinary []byte, createdTime uint64, period uint16, flag uint32) {
	var err error
	_, err = postgres.Exec(`INSERT INTO pretendo_badge_arcade.free_play_data(
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
		$1,
		$2,
		$3,
		$4,
		$4,
		$5,
		$6,
		$4
	) ON CONFLICT DO NOTHING`, dataID, ownerID, metaBinary, createdTime, period, flag)
	if err != nil {
		log.Fatal(err)
	}
}

func postPersistenceInfo(dataID uint64, ownerID uint32, slot uint16) {
	var err error
	_, err = postgres.Exec(`INSERT INTO pretendo_badge_arcade.persistence_info(
		data_id,
		pid,
		slot
	)
	VALUES (
		$1,
		$2,
		$3
	) ON CONFLICT DO NOTHING`, dataID, ownerID, slot)
	if err != nil {
		log.Fatal(err)
	}
}

func postUserPlayInfo(dataID uint64, version uint32, size uint32) {
	var err error
	_, err = postgres.Exec(`INSERT INTO pretendo_badge_arcade.user_play_info(
		data_id,
		version,
		size
	)
	VALUES (
		$1,
		$2,
		$3
	) ON CONFLICT DO NOTHING`, dataID, version, size)
	if err != nil {
		log.Fatal(err)
	}
}

func updateFreePlayDataMetaBinary(dataID uint64, metaBinary []byte, updatedTime uint64) {
	var err error
	_, err = postgres.Exec(`UPDATE pretendo_badge_arcade.free_play_data SET meta_binary=$1, updated_time=$2 WHERE data_id=$3`, metaBinary, updatedTime, dataID)
	if err != nil {
		log.Fatal(err)
	}
}

func updateUserPlayInfoSize(dataID uint64, size uint32) {
	var err error
	_, err = postgres.Exec(`UPDATE pretendo_badge_arcade.user_play_info SET size=$1 WHERE data_id=$2`, size, dataID)
	if err != nil {
		log.Fatal(err)
	}
}

func updateUserPlayInfoVersion(dataID uint64, version uint32) {
	var err error
	_, err = postgres.Exec(`UPDATE pretendo_badge_arcade.user_play_info SET version=$1 WHERE data_id=$2`, version, dataID)
	if err != nil {
		log.Fatal(err)
	}
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
