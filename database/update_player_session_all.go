package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"github.com/PretendoNetwork/nintendo-badge-arcade-secure/globals"
)

func UpdatePlayerSessionAll(pid uint32, urls []string, ip string, port string) {
	_, err := sessionsCollection.UpdateOne(context.TODO(), bson.D{{"pid", pid}}, bson.D{{"$set", bson.D{{"pid", pid}, {"urls", urls}, {"ip", ip}, {"port", port}}}})
	if err != nil {
		globals.Logger.Error(err.Error())
	}
}
