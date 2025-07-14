package graph

import (
	"context"
	"fmt"
	"log"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type GraphRepository struct {
	driver neo4j.DriverWithContext
}

func MustConnect(
	uri string,
	username string,
	password string,
	realm string,
) neo4j.DriverWithContext {
	driver, err := neo4j.NewDriverWithContext(
		uri,
		neo4j.BasicAuth(username, password, realm),
	)
	if err != nil {
		log.Fatal(err)
	}

	return driver
}

func NewGraphRepository(driver neo4j.DriverWithContext) *GraphRepository {
	return &GraphRepository{driver: driver}
}

func (r *GraphRepository) Close(ctx context.Context) {
	r.driver.Close(ctx)
}

func (r *GraphRepository) CreateUser(ctx context.Context, id, handle string) error {
	const op = "storage.graph.CreateUser"

	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `MERGE (:User {id: $id, handle: $handle})`
		_, err := tx.Run(ctx, query, map[string]any{"id": id, "handle": handle})
		return nil, err
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *GraphRepository) SendFriendRequest(ctx context.Context, fromID, toID string) error {
	const op = "storage.graph.SendFriendRequest"

	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (from:User {id: $fromID}), (to:User {id: $toID})
			MERGE (from)-[r:SENT_REQUEST]->(to)
			ON CREATE SET r.created_at = datetime()
		`
		_, err := tx.Run(ctx, query, map[string]any{"fromID": fromID, "toID": toID})
		return nil, err
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *GraphRepository) AcceptFriendRequest(ctx context.Context, userID, requesterID string) error {
	const op = "storage.graph.AcceptFriendRequest"

	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (a:User {id: $requesterID})-[r:SENT_REQUEST]->(b:User {id: $userID})
			DELETE r
			CREATE (a)-[:FRIEND {since: date()}]->(b)
			CREATE (b)-[:FRIEND {since: date()}]->(a)
		`
		_, err := tx.Run(ctx, query, map[string]any{
			"userID":      userID,
			"requesterID": requesterID,
		})
		return nil, err
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *GraphRepository) RemoveFriend(ctx context.Context, userID, friendID string) error {
	const op = "storage.graph.RemoveFriend"

	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (a:User {id: $userID})-[r1:FRIEND]->(b:User {id: $friendID})
			DELETE r1
			WITH a, b
			MATCH (b)-[r2:FRIEND]->(a)
			DELETE r2
		`
		_, err := tx.Run(ctx, query, map[string]any{"userID": userID, "friendID": friendID})
		return nil, err
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *GraphRepository) BlockUser(ctx context.Context, userID, targetID string) error {
	const op = "storage.graph.BlockUser"

	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (a:User {id: $userID})-[r]->(b:User {id: $targetID})
			DELETE r
			MERGE (a)-[:BLOCKED]->(b)
		`
		_, err := tx.Run(ctx, query, map[string]any{"userID": userID, "targetID": targetID})
		return nil, err
	})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *GraphRepository) ListFriends(ctx context.Context, userID string) ([]string, error) {
	const op = "storage.graph.ListFriends"

	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (:User {id: $userID})-[:FRIEND]->(f:User)
			RETURN f.id
		`
		rows, err := tx.Run(ctx, query, map[string]any{"userID": userID})
		if err != nil {
			return nil, err
		}

		var ids []string
		for rows.Next(ctx) {
			id, _ := rows.Record().Get("f.id")
			if str, ok := id.(string); ok {
				ids = append(ids, str)
			}
		}
		return ids, nil
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return result.([]string), nil
}

func (r *GraphRepository) ListMutualFriends(ctx context.Context, userA, userB string) ([]string, error) {
	const op = "storage.graph.ListMutualFriends"

	session := r.driver.NewSession(ctx, neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close(ctx)

	result, err := session.ExecuteRead(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		query := `
			MATCH (a:User {id: $userA})-[:FRIEND]->(f:User)<-[:FRIEND]-(b:User {id: $userB})
			RETURN f.id
		`
		rows, err := tx.Run(ctx, query, map[string]any{"userA": userA, "userB": userB})
		if err != nil {
			return nil, err
		}

		var ids []string
		for rows.Next(ctx) {
			id, _ := rows.Record().Get("f.id")
			if str, ok := id.(string); ok {
				ids = append(ids, str)
			}
		}
		return ids, nil
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return result.([]string), nil
}
