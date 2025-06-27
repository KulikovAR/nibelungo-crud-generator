package mongodb

import (
	"context"
	"time"
	"github.com/KulikovAR/github.com/KulikovAR/some-project/internal/domain"
	"github.com/KulikovAR/github.com/KulikovAR/some-project/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type productRepository struct {
	collection *mongo.Collection
}

func NewproductRepository(collection *mongo.Collection) repository.productRepository {
	return &productRepository{collection: collection}
}

func (r *productRepository) Create(ctx context.Context, entity *domain.product) error {
	entity.CreatedAt = time.Now()
	entity.UpdatedAt = time.Now()
	
	_, err := r.collection.InsertOne(ctx, entity)
	return err
}

func (r *productRepository) Get(ctx context.Context, id string) (*domain.product, error) {
	var entity domain.product
	err := r.collection.FindOne(ctx, bson.M{"id": id}).Decode(&entity)
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *productRepository) Update(ctx context.Context, entity *domain.product) error {
	entity.UpdatedAt = time.Now()
	
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"id": entity.ID},
		bson.M{"$set": entity},
	)
	return err
}

func (r *productRepository) Delete(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"id": id})
	return err
}

func (r *productRepository) List(ctx context.Context) ([]*domain.product, error) {
	opts := options.Find().SetSort(bson.D{{"created_at", -1}})
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	
	var entities []*domain.product
	if err = cursor.All(ctx, &entities); err != nil {
		return nil, err
	}
	return entities, nil
}
