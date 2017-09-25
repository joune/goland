package users

// load/query user profiles in Redis
// hide the redis implem behind UserRepository abstraction

import (
	"github.com/go-redis/redis"
	"github.com/golang/protobuf/proto"
	"github.com/joune/zenly/data"
	"log"
	"strconv"
)

type UserRepository struct {
	client *redis.Client
}

func Users() UserRepository {
	return UserRepository{
		client: redis.NewClient(&redis.Options{
			Addr:     "redis:6379",
			Password: "",
			DB:       0,
		}), //don't forget the trailing coma! :/
	}
}

func usrKey(usr data.User) string {
	return idKey(usr.GetId())
}

func idKey(id uint64) string {
	return strconv.Itoa(int(id))
}

func (repo UserRepository) Put(usr data.User) error {
	log.Printf("Register user %s\n", usr.String())
	bytes, err := proto.Marshal(&usr)
	if err != nil {
		return err
	}
	return repo.client.Set(usrKey(usr), bytes, 0).Err()
}

func (repo UserRepository) Get(id uint64) (data.User, error) {
	got := repo.client.Get(idKey(id))
	usr := data.User{}
	// error handling at every single line of code :(
	// there has to be a better way..
	if err := got.Err(); err != nil {
		return usr, err
	}
	bytes, err := got.Bytes()
	if err != nil {
		return usr, err
	}
	if err := proto.Unmarshal(bytes, &usr); err != nil {
		return usr, err
	}
	return usr, nil
}

func (repo UserRepository) ProvisionUsers(usrs []data.User) {
	for i := 0; i < len(usrs); i++ {
		if err := repo.Put(usrs[i]); err != nil {
			panic(err)
		}
	}

	// just to check I can read what I just wrote
	usr, err := repo.Get(usrs[0].GetId())
	if err != nil {
		panic(err)
	}
	log.Printf("Retrieved user %s\n", usr.String())
}
