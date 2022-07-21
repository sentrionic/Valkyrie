package service

import (
	"github.com/bwmarrin/snowflake"
	"log"
)

var node *snowflake.Node

func init() {
	const nodeID int64 = 1

	var err error
	node, err = snowflake.NewNode(nodeID)
	if err != nil {
		log.Fatalf("failed to init snowflake node: %v", err.Error())
	}
}

// GenerateId generates a snowflake id
func GenerateId() string {
	// Generate a snowflake ID.
	id := node.Generate()

	return id.String()
}
