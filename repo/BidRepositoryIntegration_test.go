package repo

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/luck02/dibbler/fixtures"
	"github.com/luck02/dibbler/models"
)

func newPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3, // Get from config
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			/*if _, err := c.Do("AUTH", password); err != nil {
				c.Close()
				return nil, err
			}*/
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func TestRedisBasicFunctions(t *testing.T) {
	pool := newPool("localhost:6379", "")

	conn := pool.Get()
	defer conn.Close()

	n, err := conn.Do("SET", "derp", "DERPS")

	fmt.Println(n)
	fmt.Println(err)

	b, err := json.Marshal(fixtures.CampaignTests[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))

	sampleCampaign := models.Campaign{}
	err = json.Unmarshal(b, &sampleCampaign)
	fmt.Println(err)
	fmt.Printf("%+v\n", sampleCampaign)

}

func TestPlaceBidSuccess(t *testing.T) {
	//redisBidRepo := RedisBidRepository{}

}
