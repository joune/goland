package main

// load users in Redis
// produce sample sessions to publish to kafka

import (
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/joune/zenly/data"
	"github.com/joune/zenly/geo"
	"github.com/joune/zenly/users"
	"log"
	"strconv"
	"time"
)

func main() {
	//let servers start!
	time.Sleep(20 * time.Second)

	users.Users().ProvisionUsers(Users)

	producer := producer()
	defer func() {
		if err := producer.Close(); err != nil {
			panic(err)
		}
	}()

	for i := 0; i < len(Sessions); i++ {
		pushSession(producer, Sessions[i])
	}
}

//sample user ids
const (
	joune uint64 = iota + 1
	mary
	soy
	sylvain
	oliv
	aicha
)

var ( //why not const?!

	//sample places
	Lods       = geo.Place{48.875364, 2.377096}
	Soys       = geo.Place{48.878345, 2.311961}
	Olivs      = geo.Place{48.837557, 2.232180}
	Aichas     = geo.Place{48.849017, 2.355579}
	Kyriba     = geo.Place{48.841953, 2.220330}
	La_defense = geo.Place{48.892570, 2.236523}
	Cafe       = geo.Place{48.874115, 2.374966}
	Les_halles = geo.Place{48.862688, 2.344135}
	Luxembourg = geo.Place{48.848159, 2.336731}
	Opera      = geo.Place{48.872430, 2.331604}
	Resto      = geo.Place{48.874797, 2.325844}

	//sample users
	Users = []data.User{
		user(joune, Lods, Kyriba),       // 1
		user(mary, Lods, La_defense),    // 2 - joune's wife
		user(soy, Soys, Soys),           // 3 - works from home
		user(sylvain, Lods, Lods),       // 4 - soy's best friend
		user(oliv, Olivs, Kyriba),       // 5 - joune's colleague
		user(aicha, Aichas, Les_halles), // 6 - soy's crush
	}

	//sample sessions
	Sessions = []data.Session{
		//day1
		session(joune, mary, Lods, -7, 0*time.Hour, 8*time.Hour),
		session(joune, oliv, Kyriba, -7, 10*time.Hour, 8*time.Hour),
		session(joune, sylvain, Cafe, -7, 22*time.Hour, 1*time.Hour),
		session(soy, sylvain, Cafe, -7, 20*time.Hour, 5*time.Hour),
		session(joune, mary, Lods, -7, 23*time.Hour, 9*time.Hour),
		session(soy, aicha, Soys, -7, 23*time.Hour, 7*time.Hour),
		//day2
		session(soy, sylvain, Les_halles, -6, 12*time.Hour, 1*time.Hour),
		session(joune, oliv, Kyriba, -6, 10*time.Hour, 8*time.Hour),
		session(mary, soy, Opera, -6, 19*time.Hour, 1*time.Hour),
		session(mary, soy, Opera, -6, 19*time.Hour, 1*time.Hour),
		session(soy, sylvain, Cafe, -6, 20*time.Hour, 5*time.Hour),
		session(joune, mary, Lods, -6, 23*time.Hour, 9*time.Hour),
		session(soy, aicha, Soys, -6, 23*time.Hour, 7*time.Hour),
		//day3
		session(soy, sylvain, Resto, -5, 15*time.Hour, 2*time.Hour),
		session(joune, oliv, Kyriba, -5, 10*time.Hour, 8*time.Hour),
		session(mary, sylvain, Resto, -5, 13*time.Hour, 1*time.Hour),
		session(mary, soy, Lods, -5, 20*time.Hour, 5*time.Hour),
		session(joune, soy, Lods, -5, 20*time.Hour, 5*time.Hour),
		session(soy, sylvain, Cafe, -5, 19*time.Hour, 1*time.Hour),
		session(joune, mary, Lods, -5, 23*time.Hour, 9*time.Hour),
		session(soy, aicha, Aichas, -5, 23*time.Hour, 7*time.Hour),
	}
)

func user(id uint64, home, work geo.Place) data.User {
	return data.User{
		Id:       id,
		HomeCell: home.CellId(), //why the renaming convention? :(
		WorkCell: work.CellId(),
	}
}

func session(u1, u2 uint64, where geo.Place, daysAgo int, startTime, dur time.Duration) data.Session {
	now := time.Now()
	// set our users timezone on the other side of the world
	// to make sure we compute is night correctly
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.FixedZone("Asia/Vladivostok", 10))
	start := today.AddDate(0, 0, daysAgo).Add(startTime)
	startts, err := ptypes.TimestampProto(start)
	if err != nil {
		panic(err)
	}
	endts, err := ptypes.TimestampProto(start.Add(dur))
	if err != nil {
		panic(err)
	}
	return data.Session{
		u1, u2,
		startts, endts,
		where.Lat, where.Lng,
	}
}

func producer() sarama.AsyncProducer {
	config := sarama.NewConfig()
	config.Producer.Retry.Max = 5
	config.Producer.RequiredAcks = sarama.WaitForAll
	brokers := []string{"kafka:9092"}
	producer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		panic(err)
	}
	return producer
}

func pushSession(producer sarama.AsyncProducer, s data.Session) {
	bytes, err := proto.Marshal(&s)
	if err != nil {
		panic(err)
	}
	strTime := strconv.Itoa(int(time.Now().Unix()))

	callback := &sarama.ProducerMessage{
		Topic: "sessions",
		Key:   sarama.StringEncoder(strTime),
		Value: sarama.ByteEncoder(bytes),
	}
	select {
	case producer.Input() <- callback:
		log.Printf("Session sent %s\n", strTime)
	case err := <-producer.Errors():
		log.Printf("Failed to send Session at %s: %s\n", strTime, err)
	}
}
