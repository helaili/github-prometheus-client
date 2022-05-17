package main

import (
	"fmt"
	"log"

	"github.com/gomodule/redigo/redis"
	"github.com/google/go-github/v43/github"
)

type RedisCache struct {
	AbstractCache
	pool *redis.Pool
}

func NewRedisCache(redisAddress string, app_id int64, private_key []byte) *RedisCache {
	log.Printf("Using the Redis cache at %s\n", redisAddress)
	return &RedisCache{
		*NewAbstractCache(app_id, private_key),
		&redis.Pool{
			MaxIdle:   80,
			MaxActive: 12000,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", redisAddress)
				if err != nil {
					panic(err.Error())
				} else {
					log.Println("Connected to Redis")
				}
				return c, err
			},
		},
	}
}

func (m RedisCache) set(event *github.WorkflowRunEvent) {
	client := m.pool.Get()
	defer client.Close()

	key := fmt.Sprintf("%d-%d", event.GetInstallation().GetID(), event.GetWorkflowRun().GetID())
	_, err := client.Do("SET", key, event.GetWorkflow().GetName())
	if err != nil {
		log.Printf("error setting key %s in Redis: err=%s\n", key, err)
	}
}

func (m RedisCache) get(event *github.WorkflowJobEvent) string {
	client := m.pool.Get()
	defer client.Close()

	key := fmt.Sprintf("%d-%d", event.GetInstallation().GetID(), event.GetWorkflowJob().GetRunID())
	worfklowName, err := client.Do("GET", key)

	if err != nil {
		workflowName := m.getWorkflowNameFromGitHub(event)
		_, err := client.Do("SET", key, workflowName)
		if err != nil {
			log.Printf("error setting key %s in Redis after a missed cache: err=%s\n", key, err)
		}
		return workflowName
	} else {
		return fmt.Sprintf("%s", worfklowName)
	}
}
