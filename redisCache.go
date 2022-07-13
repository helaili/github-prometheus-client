package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/go-github/v43/github"
)

type RedisCache struct {
	AbstractCache
	client *redis.Client
}

func NewRedisCache(redisAddress string, readisPassword string, app_id int64, private_key []byte) *RedisCache {
	log.Printf("Using the Redis cache at %s with password %s\n", redisAddress, readisPassword)
	return &RedisCache{
		*NewAbstractCache(app_id, private_key),
		redis.NewClient(&redis.Options{
			Addr:     redisAddress,
			Password: readisPassword,
			DB:       0, // use default DB
		}),
	}
}

func (m RedisCache) set(event *github.WorkflowRunEvent) {
	key := fmt.Sprintf("%d-%d", event.GetInstallation().GetID(), event.GetWorkflowRun().GetID())
	err := m.client.Set(context.Background(), key, event.GetWorkflow().GetName(), time.Hour*24*35)
	if err != nil {
		log.Printf("error setting key %s in Redis: err=%s\n", key, err)
	}
}

func (m RedisCache) get(event *github.WorkflowJobEvent) string {
	key := fmt.Sprintf("%d-%d", event.GetInstallation().GetID(), event.GetWorkflowJob().GetRunID())
	worfklowName := m.client.Get(context.Background(), key)

	// Cache miss, we need to retrieve the workflow name from github
	if worfklowName == nil {
		workflowName := m.getWorkflowNameFromGitHub(event)
		err := m.client.Set(context.Background(), key, workflowName, time.Hour*24*35)
		if err != nil {
			log.Printf("error setting key %s in Redis after a missed cache: err=%s\n", key, err)
		}
		return workflowName
	} else {
		return worfklowName.String()
	}
}
