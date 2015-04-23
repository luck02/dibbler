package repo

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
	"github.com/luck02/dibbler/models"
)

// BidRepository allows us to get campaigns and place bids
type BidRepository interface {
	PlaceBid(models.Campaign) (models.Campaign, bool, error)
	GetCampaigns() ([]models.Campaign, error)
}

// RedisBidRepository is the concrete implementation of BidRepository
type RedisBidRepository struct {
	pool *redis.Pool
}

// NewRedisBidRepository is the constructor for RedisBidRepository
func NewRedisBidRepository(server string) *RedisBidRepository {
	bidRepo := new(RedisBidRepository)
	bidRepo.pool = newPool(server)
	return bidRepo
}

// PlaceBid will decrement the current value of the campaign to
func (r *RedisBidRepository) PlaceBid(campaign models.Campaign) (models.Campaign, bool, error) {
	campaignIdKey := fmt.Sprintf("campaigns:%d", campaign.ID)

	conn := r.pool.Get()
	defer conn.Close()

	conn.Send("WATCH", campaignIdKey)
	conn.Send("HGET", campaignIdKey, "RemainingBudget")
	conn.Flush()
	conn.Receive()

	remainingBudget, err := redis.Float64(conn.Receive())
	if err != nil {
		return models.Campaign{}, false, err
	}

	if float32(remainingBudget) < campaign.BidCpm/1000 {
		campaign.RemainingBudget = 0
		err = r.SaveCampaign(campaign)
		if err != nil {
			return campaign, false, err
		}

		campaign, err := r.getCampaign(campaign.ID)
		return campaign, false, err
	}

	conn.Send("MULTI")
	conn.Send("HINCRBYFLOAT", campaignIdKey, "RemainingBudget", -campaign.BidCpm/1000)
	_, err = conn.Do("EXEC")
	if err != nil {
		return models.Campaign{}, false, err
	}

	if err != nil {
		return models.Campaign{}, false, err
	}

	campaign, err = r.getCampaign(campaign.ID)
	if err != nil {
		return models.Campaign{}, false, err
	}
	logrus.WithFields(logrus.Fields{
		"Event":      "PlaceBid",
		"Success":    true,
		"campaignId": campaign.ID}).Info("Bid placed successfully")

	return campaign, true, nil
}

func (r *RedisBidRepository) GetCampaigns() ([]models.Campaign, error) {
	conn := r.pool.Get()
	defer conn.Close()

	campaignsKeyCount, err := redis.Int(conn.Do("SCARD", "campaignIdList"))
	if err != nil {
		return nil, err
	}

	reply, err := redis.Values(conn.Do("SMEMBERS", "campaignIdList"))
	if err != nil {
		return nil, err
	}

	campaigns := make([]models.Campaign, campaignsKeyCount, campaignsKeyCount)

	var campaignId int64
	for i := 0; i < campaignsKeyCount; i++ {
		reply, err = redis.Scan(reply, &campaignId)
		campaign, err := r.getCampaign(int32(campaignId))
		if err != nil {
			return nil, err
		}
		campaigns[i] = campaign
	}

	return campaigns, nil
}

func (r *RedisBidRepository) SaveCampaign(campaign models.Campaign) error {
	conn := r.pool.Get()
	defer conn.Close()

	campaignJSON, err := json.Marshal(campaign)
	if err != nil {
		return err
	}

	targetType := reflect.TypeOf(campaign.Targeting)
	targetJson, err := json.Marshal(campaign.Targeting)
	if err != nil {
		return err
	}

	_, err = conn.Do("SADD", "campaignIdList", campaign.ID)
	if err != nil {
		return err
	}
	_, err = conn.Do("HMSET", fmt.Sprintf("campaigns:%d", campaign.ID), "CampaignJson", campaignJSON, targetType, targetJson, "RemainingBudget", campaign.RemainingBudget, "BidCPM", campaign.BidCpm)

	return err
}

func (r *RedisBidRepository) getCampaign(ID int32) (models.Campaign, error) {
	conn := r.pool.Get()
	defer conn.Close()

	campaign := models.Campaign{}
	var campaignJSONKey string
	var campaignJSON string
	var targetType string
	var targetJson string
	var remainingBudget float32
	var remainingBudgetKey string // these key strings seem to be placeholders and are required.  There's probably a smarter way to make this work.

	reply, err := redis.Values(conn.Do("HGETALL", fmt.Sprintf("campaigns:%d", ID)))
	if err != nil {
		return models.Campaign{}, err
	}

	_, err = redis.Scan(reply, &campaignJSONKey, &campaignJSON, &targetType, &targetJson, &remainingBudgetKey, &remainingBudget)
	if err != nil {
		return models.Campaign{}, err
	}
	err = json.Unmarshal([]byte(campaignJSON), &campaign)
	if err != nil {
		return models.Campaign{}, err
	}
	campaign.RemainingBudget = remainingBudget

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
	case reflect.TypeOf(models.CountryTarget{}).String():
		var countryTarget models.CountryTarget
		jsonBytes := []byte(targetJson)
		err := json.Unmarshal(jsonBytes, &countryTarget)
		if err != nil {
			return nil, err
		}
		return countryTarget, nil
	case reflect.TypeOf(models.OSTarget{}).String():
		var osTarget models.OSTarget
		jsonBytes := []byte(targetJson)
		err := json.Unmarshal(jsonBytes, &osTarget)
		if err != nil {
			return nil, err
		}
		return osTarget, nil
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
