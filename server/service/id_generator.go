package service

import (
	"github.com/bwmarrin/snowflake"
	"github.com/sentrionic/valkyrie/model/apperrors"
	"log"
)

// GenerateId generates a snowflake id
func GenerateId() (string, error) {
	node, err := snowflake.NewNode(1)
	if err != nil {
		log.Printf("Failed to genenerate an snowflake id: %v\n", err.Error())
		return "", apperrors.NewInternal()
	}

	// Generate a snowflake ID.
	id := node.Generate()

	return id.String(), nil
}
