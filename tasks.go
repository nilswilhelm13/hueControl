package main

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
	"log"
	"time"
)

const tasksBucket = "tasksBucket"

type Task struct {
	// ID
	ID string `json:"id"`
	// StartTime for fading
	StartTime time.Time `json:"start_time"`
	// TimerDuration as duration for fading
	TimerDuration int `json:"timer_duration"`
	// LightParameters
	LightParameters LightParameters `json:"light_parameters"`
	// Days of week for task
	Days []time.Weekday `json:"days"`
	// Hours
	Hours int `json:"hours"`
	// Minutes
	Minutes int `json:"minutes"`
	//Enabled
	Enabled bool `json:"enabled"`
}

// SaveOrUpdate saves or updates a task in the db
func SaveOrUpdate(db *bolt.DB, task Task) error {
	// Store the user model in the user bucket using the username as the key.
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(tasksBucket))
		if err != nil {
			return err
		}
		if task.ID == "" {
			task.ID = uuid.New().String()
		}
		encoded, err := json.Marshal(task)
		if err != nil {
			return err
		}
		return b.Put([]byte(task.ID), encoded)
	})
	return err
}

// GetTasksFromDB returns all tasks stored in the db
func GetTasksFromDB(db *bolt.DB) ([]Task, error) {
	var results []Task
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(tasksBucket))
		if err := b.ForEach(func(k, v []byte) error {
			var result Task
			if err := json.Unmarshal(v, &result); err != nil {
				return err
			}
			results = append(results, result)
			return nil
		}); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return []Task{}, err
	}
	return results, nil
}

// UpdateSchedule removes all cronjobs and reloads all tasks from db
func UpdateSchedule(db *bolt.DB, scheduler *cron.Cron) {
	tasks, err := GetTasksFromDB(db)
	if err != nil {
		log.Fatal("could not load tasks")
	}

	//RemoveCronEntries(scheduler)
	scheduler.Stop()
	scheduler = cron.New()

	for _, task := range tasks {
		if !task.Enabled {
			continue
		}
		_, err := scheduler.AddFunc(cronExpression(task), buildCronJob(task))
		if err != nil {
			log.Fatal(err)
		}
	}

	scheduler.Start()
}

// DeleteTask deletes a task from the db
func DeleteTask(db *bolt.DB, task Task) error {
	if err := db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket([]byte(tasksBucket)).Delete([]byte(task.ID))
	}); err != nil {
		return err
	}
	return nil
}

// RemoveCronEntries removes every cron entry
func RemoveCronEntries(scheduler *cron.Cron) {
	for _, entry := range scheduler.Entries() {
		scheduler.Remove(entry.ID)
	}
}

// buildCronJob returns a function to use with a cron job
func buildCronJob(task Task) func() {
	return func() {
		fmt.Println("Hello")
		fadeIn(task)
	}
}

func openDb() *bolt.DB {
	db, err := bolt.Open("tasks.db", 0600, nil)
	if err != nil {
		panic("could not open db")
	}
	if err := db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(tasksBucket)); err != nil {
			return err
		}
		return nil
	}); err != nil {
		panic(err)
	}
	return db
}
