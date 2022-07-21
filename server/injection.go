package main

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
	"github.com/sentrionic/valkyrie/config"
	"github.com/sentrionic/valkyrie/handler"
	"github.com/sentrionic/valkyrie/handler/middleware"
	"github.com/sentrionic/valkyrie/model"
	"github.com/sentrionic/valkyrie/repository"
	"github.com/sentrionic/valkyrie/service"
	"github.com/sentrionic/valkyrie/ws"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"
	"log"
	"net/http"
	"time"
)

func inject(d *dataSources, cfg config.Config) (*gin.Engine, error) {
	log.Println("Injecting data sources")

	// Repository layer
	userRepository := repository.NewUserRepository(d.DB)
	friendRepository := repository.NewFriendRepository(d.DB)
	guildRepository := repository.NewGuildRepository(d.DB)
	channelRepository := repository.NewChannelRepository(d.DB)
	messageRepository := repository.NewMessageRepository(d.DB)

	fileRepository := repository.NewFileRepository(d.S3Session, cfg.BucketName)
	redisRepository := repository.NewRedisRepository(d.RedisClient)

	mailRepository := repository.NewMailRepository(cfg.GmailUser, cfg.GmailPassword, cfg.CorsOrigin)

	// Service Layer
	userService := service.NewUserService(&service.USConfig{
		UserRepository:  userRepository,
		FileRepository:  fileRepository,
		RedisRepository: redisRepository,
		MailRepository:  mailRepository,
	})

	friendService := service.NewFriendService(&service.FSConfig{
		UserRepository:   userRepository,
		FriendRepository: friendRepository,
	})

	guildService := service.NewGuildService(&service.GSConfig{
		UserRepository:    userRepository,
		FileRepository:    fileRepository,
		RedisRepository:   redisRepository,
		GuildRepository:   guildRepository,
		ChannelRepository: channelRepository,
	})

	channelService := service.NewChannelService(&service.CSConfig{
		ChannelRepository: channelRepository,
		GuildRepository:   guildRepository,
	})

	messageService := service.NewMessageService(&service.MSConfig{
		MessageRepository: messageRepository,
		FileRepository:    fileRepository,
	})

	// initialize gin.Engine
	router := gin.Default()

	// set cors settings
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{cfg.CorsOrigin},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
	})
	router.Use(c)

	redisURL := d.RedisClient.Options().Addr
	password := d.RedisClient.Options().Password

	// initialize session store
	store, err := redis.NewStore(10, "tcp", redisURL, password, []byte(cfg.SessionSecret))

	if err != nil {
		return nil, fmt.Errorf("could not initialize redis session store: %w", err)
	}

	store.Options(sessions.Options{
		Domain:   cfg.Domain,
		MaxAge:   60 * 60 * 24 * 7, // 7 days
		Secure:   gin.Mode() == gin.ReleaseMode,
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})

	router.Use(sessions.Sessions(model.CookieName, store))

	// add rate limit
	rate := limiter.Rate{
		Period: 1 * time.Hour,
		Limit:  1500,
	}

	limitStore, _ := sredis.NewStore(d.RedisClient)

	rateLimiter := mgin.NewMiddleware(limiter.New(limitStore, rate))
	router.Use(rateLimiter)

	// Websockets Setup
	hub := ws.NewWebsocketHub(&ws.Config{
		UserService:    userService,
		GuildService:   guildService,
		ChannelService: channelService,
		Redis:          d.RedisClient,
	})
	go hub.Run()

	router.GET("/ws", middleware.AuthUser(), func(c *gin.Context) {
		ws.ServeWs(hub, c)
	})

	socketService := service.NewSocketService(&service.SSConfig{
		Hub:               *hub,
		GuildRepository:   guildRepository,
		ChannelRepository: channelRepository,
	})

	handler.NewHandler(&handler.Config{
		R:               router,
		UserService:     userService,
		FriendService:   friendService,
		GuildService:    guildService,
		ChannelService:  channelService,
		MessageService:  messageService,
		SocketService:   socketService,
		TimeoutDuration: time.Duration(cfg.HandlerTimeOut) * time.Second,
		MaxBodyBytes:    cfg.MaxBodyBytes,
	})

	return router, nil
}
