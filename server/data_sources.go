package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-redis/redis/v8"
	"github.com/sentrionic/valkyrie/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

type dataSources struct {
	DB          *gorm.DB
	RedisClient *redis.Client
	S3Session   *session.Session
}

// InitDS establishes connections to fields in dataSources
func initDS() (*dataSources, error) {
	log.Printf("Initializing data sources\n")
	dbUrl := os.Getenv("DATABASE_URL")

	log.Printf("Connecting to Postgresql\n")
	db, err := gorm.Open(postgres.Open(dbUrl))

	if err != nil {
		return nil, fmt.Errorf("error opening db: %w", err)
	}

	// Migrate models and setup join tables
	if err := db.AutoMigrate(
		&model.User{},
		&model.Guild{},
		&model.Member{},
		&model.Channel{},
		&model.DMMember{},
		&model.Message{},
		&model.Attachment{},
		&model.VCMember{},
	); err != nil {
		return nil, fmt.Errorf("error migrating models: %w", err)
	}

	if err := db.SetupJoinTable(&model.Guild{}, "Members", &model.Member{}); err != nil {
		return nil, fmt.Errorf("error creating join table: %w", err)
	}

	if err := db.SetupJoinTable(&model.Guild{}, "VCMembers", &model.VCMember{}); err != nil {
		return nil, fmt.Errorf("error creating join table: %w", err)
	}

	// Initialize redis connection
	redisURL := os.Getenv("REDIS_URL")
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		panic(err)
	}

	log.Printf("Connecting to Redis\n")
	rdb := redis.NewClient(opt)

	// verify redis connection
	_, err = rdb.Ping(context.Background()).Result()

	if err != nil {
		return nil, fmt.Errorf("error connecting to redis: %w", err)
	}

	// Initialize S3 Session
	accessKey := os.Getenv("AWS_ACCESS_KEY")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_S3_REGION")

	sess, err := session.NewSession(
		&aws.Config{
			Credentials: credentials.NewStaticCredentials(
				accessKey,
				secretKey,
				"",
			),
			Region: aws.String(region),
		},
	)

	if err != nil {
		return nil, fmt.Errorf("error creating s3 session: %w", err)
	}

	return &dataSources{
		DB:          db,
		RedisClient: rdb,
		S3Session:   sess,
	}, nil
}

// close to be used in graceful server shutdown
func (d *dataSources) close() error {
	if err := d.RedisClient.Close(); err != nil {
		return fmt.Errorf("error closing Redis Client: %w", err)
	}

	return nil
}
