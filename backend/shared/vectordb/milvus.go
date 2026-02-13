// Package vectordb는 Milvus 벡터 데이터베이스 클라이언트를 제공합니다.
package vectordb

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

// MilvusClient wraps the Milvus SDK client.
type MilvusClient struct {
	client     client.Client
	collection string
	dimension  int
}

// NewMilvusClient creates a connection to Milvus and ensures the collection exists.
func NewMilvusClient(addr, collection string, dimension int) (*MilvusClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c, err := client.NewClient(ctx, client.Config{
		Address: addr,
	})
	if err != nil {
		return nil, fmt.Errorf("milvus connect failed: %w", err)
	}

	mc := &MilvusClient{
		client:     c,
		collection: collection,
		dimension:  dimension,
	}

	// Ensure collection exists
	if err := mc.ensureCollection(ctx); err != nil {
		c.Close()
		return nil, err
	}

	return mc, nil
}

func (m *MilvusClient) ensureCollection(ctx context.Context) error {
	exists, err := m.client.HasCollection(ctx, m.collection)
	if err != nil {
		return fmt.Errorf("check collection: %w", err)
	}
	if exists {
		return nil
	}

	schema := &entity.Schema{
		CollectionName: m.collection,
		Fields: []*entity.Field{
			{
				Name:       "id",
				DataType:   entity.FieldTypeVarChar,
				PrimaryKey: true,
				AutoID:     false,
				TypeParams: map[string]string{"max_length": "128"},
			},
			{
				Name:       "session_id",
				DataType:   entity.FieldTypeVarChar,
				TypeParams: map[string]string{"max_length": "128"},
			},
			{
				Name:       "vector",
				DataType:   entity.FieldTypeFloatVector,
				TypeParams: map[string]string{"dim": fmt.Sprintf("%d", m.dimension)},
			},
		},
	}

	if err := m.client.CreateCollection(ctx, schema, 2); err != nil {
		return fmt.Errorf("create collection: %w", err)
	}

	// Create IVF_FLAT index with cosine similarity
	idx, _ := entity.NewIndexIvfFlat(entity.COSINE, 128)
	if err := m.client.CreateIndex(ctx, m.collection, "vector", idx, false); err != nil {
		log.Printf("[milvus] index creation warning: %v", err)
	}

	// Load collection into memory for search
	if err := m.client.LoadCollection(ctx, m.collection, false); err != nil {
		log.Printf("[milvus] load collection warning: %v", err)
	}

	return nil
}

// Insert stores a vector with its ID and session ID.
func (m *MilvusClient) Insert(ctx context.Context, id, sessionID string, vector []float32) error {
	ids := []string{id}
	sessionIDs := []string{sessionID}
	vectors := [][]float32{vector}

	idCol := entity.NewColumnVarChar("id", ids)
	sessionCol := entity.NewColumnVarChar("session_id", sessionIDs)
	vectorCol := entity.NewColumnFloatVector("vector", m.dimension, vectors)

	_, err := m.client.Insert(ctx, m.collection, "", idCol, sessionCol, vectorCol)
	return err
}

// Search finds the topK most similar vectors using cosine similarity.
func (m *MilvusClient) Search(ctx context.Context, vector []float32, topK int) ([]SearchResult, error) {
	sp, _ := entity.NewIndexIvfFlatSearchParam(16)

	results, err := m.client.Search(
		ctx, m.collection, nil, "",
		[]string{"id", "session_id"},
		[]entity.Vector{entity.FloatVector(vector)},
		"vector", entity.COSINE, topK, sp,
	)
	if err != nil {
		return nil, fmt.Errorf("milvus search: %w", err)
	}

	var searchResults []SearchResult
	for _, result := range results {
		for i := 0; i < result.ResultCount; i++ {
			idCol, _ := result.Fields.GetColumn("id").GetAsString(i)
			sessionCol, _ := result.Fields.GetColumn("session_id").GetAsString(i)
			searchResults = append(searchResults, SearchResult{
				ID:        idCol,
				SessionID: sessionCol,
				Score:     result.Scores[i],
			})
		}
	}
	return searchResults, nil
}

// SearchResult holds a single search result from Milvus.
type SearchResult struct {
	ID        string
	SessionID string
	Score     float32
}

// Close closes the Milvus connection.
func (m *MilvusClient) Close() error {
	return m.client.Close()
}

// Health checks if Milvus is reachable.
func (m *MilvusClient) Health(ctx context.Context) error {
	_, err := m.client.HasCollection(ctx, m.collection)
	return err
}
