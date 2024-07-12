package globle

import (
	"IM/db/query"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

const (
	Project = "JL-IM"
)

var (
	Db             *query.Query
	Logger         *logrus.Logger
	RocketProducer rocketmq.Producer
	Rdb            *redis.Client
)
