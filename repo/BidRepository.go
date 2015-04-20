package repo

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/luck02/dibbler/models"
)

type BidRepository interface {
	PlaceBid(models.Campaign) (models.Campaign, bool)
	GetCampaigns() []models.Campaign
}

type RedisBidRepository struct {
	pool *redis.Pool
}

func NewRedisBidRepository(server string) *RedisBidRepository {
	bidRepo := new(RedisBidRepository)
	bidRepo.pool = newPool(server)
	return bidRepo
}

func (r *RedisBidRepository) PlaceBid(campaign models.Campaign) (models.Campaign, bool) {
	// hincrbyfloat campaigns:1 remainingBudget -.00032
	return models.Campaign{}, true
}

func (r *RedisBidRepository) GetCampaigns() []models.Campaign {
	return nil
}

func (r *RedisBidRepository) saveCampaign(campaign models.Campaign) error {
	conn := r.pool.Get()
	// HMSET campaigns:1 id "100101" type "Placement" remainingBudget 25.50
	targetType := reflect.TypeOf(campaign.Targeting)
	targetJson, err := json.Marshal(campaign.Targeting)
	if err != nil {
		return err
	}
	fmt.Printf("%+v", campaign)
	_, err = conn.Do("HMSET", fmt.Sprintf("campaigns:%d", campaign.ID), campaign.ID, targetType, campaign.BidCpm, campaign.DailyBudget, campaign.RemainingBudget, targetJson)

	return err
}
func (r *RedisBidRepository) getCampaign(ID int32) (models.Campaign, error) {
	conn := r.pool.Get()
	campaign := models.Campaign{}
	var targetType string
	var targetJson string
	reply, err := redis.Values(conn.Do("HGETALL", fmt.Sprintf("campaigns:%d", ID)))
	if err != nil {
		return models.Campaign{}, err
	}

	fmt.Println(reply)
	_, err = redis.Scan(reply, &campaign.ID, &targetType, &campaign.BidCpm, &campaign.DailyBudget, &campaign.RemainingBudget, &targetJson)
	fmt.Println(targetType)
	if err != nil {
		return models.Campaign{}, err
	}
	campaign.Targeting, err = deserializeTarget(targetType, targetJson)
	if err != nil {
		return models.Campaign{}, err
	}
	return campaign, nil
}

func deserializeTarget(targetType, targetJson string) (interface{}, error) {
	switch targetType {
	case reflect.TypeOf(models.AdTarget{}).String():
		var adTarget models.AdTarget
		jsonBytes := []byte(targetJson)
		err := json.Unmarshal(jsonBytes, &adTarget)
		if err != nil {
			return nil, err
		}
		return adTarget, nil
	case reflect.TypeOf(models.PlacementTarget{}).String():
		var placementTarget models.PlacementTarget
		jsonBytes := []byte(targetJson)
		err := json.Unmarshal(jsonBytes, &placementTarget)
		if err != nil {
			return nil, err
		}
		return placementTarget, nil
	}

	return nil, errors.New(fmt.Sprintf("Unable to find type: %+v", targetType))
}

func newPool(server string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
