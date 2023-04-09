package globals

import (
	"github.com/PretendoNetwork/nex-go"
	"github.com/PretendoNetwork/plogger-go"
	"github.com/aws/aws-sdk-go/service/s3"
)

var Logger = plogger.NewLogger()
var NEXServer *nex.Server
var S3Client *s3.S3
