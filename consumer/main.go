package main

import (
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/joune/zenly/data"
	"github.com/joune/zenly/db"
	"github.com/joune/zenly/geo"
	"github.com/joune/zenly/users"
	"github.com/joune/zenly/util"
	"log"
	"os"
	"os/signal"
	"time"
)

type handler func(msg data.Session)

type unmarshaler func(bytes []byte) (data.Session, error)

var (
	userRepo = users.Users()
	cql      = db.InitDB()
)

func main() {
	defer cql.Close()

	consume("sessions", mkSession, sessionHandler)
}

func sessionHandler(session data.Session) {
	log.Printf("Got session %s\n", session.String())

	user1, err := userRepo.Get(session.GetUser1())
	if err != nil {
		log.Printf("ERROR! %v\n", err)
		return
	}

	user2, err := userRepo.Get(session.GetUser2())
	if err != nil {
		log.Printf("ERROR! %v\n", err)
		return
	}

	loc1 := geo.IdentifyLocation(user1, user2, session.GetLatitude(), session.GetLongitude())
	log.Printf("User1: %v @ %s\n", user1.GetId(), loc1.String())
	loc2 := geo.IdentifyLocation(user2, user1, session.GetLatitude(), session.GetLongitude())
	log.Printf("User2: %v @ %s\n", user2.GetId(), loc2.String())

	sStart, err := ptypes.Timestamp(session.GetStartingDate())
	if err != nil {
		log.Printf("ERROR! %v\n", err)
		return
	}
	sEnd, err := ptypes.Timestamp(session.GetEndDate())
	if err != nil {
		log.Printf("ERROR! %v\n", err)
		return
	}

	duration := sEnd.Sub(sStart)
	isNight := util.IsNight(sStart, sEnd)

	log.Printf("Duration: %s, isNight: %v\n", duration.String(), isNight)

	db.InsertSession(cql,
		user1.GetId(), user2.GetId(),
		session.GetLatitude(), session.GetLongitude(),
		loc1, loc2,
		sStart, sEnd,
		duration, isNight)

	// insert a duplicate row with switched users (to allow simpler/faster queries for clients)
	db.InsertSession(cql,
		user2.GetId(), user1.GetId(),
		session.GetLatitude(), session.GetLongitude(),
		loc2, loc1,
		sStart, sEnd,
		duration, isNight)

	// also insert/compute aggregated time spent per couple
	db.UpdateDuration(cql, user1.GetId(), user2.GetId(), duration)
	db.UpdateDuration(cql, user2.GetId(), user1.GetId(), duration)
}

func mkSession(bytes []byte) (data.Session, error) {
	session := data.Session{}
	if err := proto.Unmarshal(bytes, &session); err != nil {
		return session, err
	}
	return session, nil
}

func consume(topic string, mkSession unmarshaler, handle handler) {
	log.Printf("Consumer started\n")

	//FIXME? for some reason this 'lazy connection' is NOT working and I have to restart the consumer
	// I'm not sure if it's docker's fault or mine
	var consumer sarama.Consumer = nil
	var err error = nil
	for consumer, err = sarama.NewConsumer([]string{"kafka:9092"}, nil); err != nil; {
		log.Printf("Failed to connect to Kafka. Will retry in 5 seconds. %s", err)
		time.Sleep(5 * time.Second) // to let containers start
	}
	log.Printf("Starting kafka listener\n")

	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	//FIXME: kafka partitions should be 'tuned' for scalability
	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	// Trap SIGINT to trigger a shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	consumed := 0
ConsumerLoop:
	for {
		select {
		case msg := <-partitionConsumer.Messages():
			log.Printf("Consumed message offset %d\n", msg.Offset)

			session, err := mkSession(msg.Value)
			if err != nil {
				log.Printf("ERROR! %s\n", err)
			} else {
				handle(session)
			}
			consumed++
		case <-signals:
			break ConsumerLoop
		}
	}

	log.Printf("Consumed: %d\n", consumed)
}
