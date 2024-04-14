package dal

import (
	"context"
	"log"

	"github.com/puzzles/services/catalog/gen"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDAL struct {
	client *mongo.Client
	db     *mongo.Database
}

func NewMongoDAL(uri, dbName string) *MongoDAL {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	db := client.Database(dbName)
	return &MongoDAL{
		client: client,
		db:     db,
	}
}

func (dal *MongoDAL) GetAllPuzzles(ctx context.Context) ([]gen.Puzzle, *gen.Error) {
	var puzzles []gen.Puzzle
	collection := dal.db.Collection("puzzles")
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, &gen.Error{Code: 500, Message: err.Error()}
	}
	if err = cursor.All(ctx, &puzzles); err != nil {
		return nil, &gen.Error{Code: 500, Message: err.Error()}
	}
	return puzzles, nil
}

func (dal *MongoDAL) GetPuzzle(ctx context.Context, id string) (*gen.Puzzle, *gen.Error) {
	var puzzle gen.Puzzle
	collection := dal.db.Collection("puzzles")
	if err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&puzzle); err != nil {
		return nil, &gen.Error{Code: 404, Message: "Puzzle not found"}
	}
	return &puzzle, nil
}

func (dal *MongoDAL) AddPuzzle(ctx context.Context, puzzle gen.NewPuzzle) (*gen.Puzzle, *gen.Error) {
	collection := dal.db.Collection("puzzles")
	id := uuid.New().String()
	newPuzzle := gen.Puzzle{
		Id:          id,
		Name:        puzzle.Name,
		Description: puzzle.Description,
		Price:       puzzle.Price,
		Quantity:    puzzle.Quantity,
		Type:        puzzle.Type,
	}
	// Ignore result from insert operation, we don't need it
	_, err := collection.InsertOne(ctx, newPuzzle)
	if err != nil {
		return nil, &gen.Error{Code: 500, Message: err.Error()}
	}
	return &newPuzzle, nil
}

func (dal *MongoDAL) UpdatePuzzle(ctx context.Context, id string, updates gen.PuzzleUpdate) *gen.Error {
	collection := dal.db.Collection("puzzles")
	result, err := collection.UpdateOne(ctx, bson.M{"id": id}, bson.D{{Key: "$set", Value: updates}})
	if err != nil {
		return &gen.Error{Code: 500, Message: err.Error()}
	}
	if result.MatchedCount == 0 {
		return &gen.Error{Code: 404, Message: "Puzzle not found"}
	}
	return nil
}

func (dal *MongoDAL) DeletePuzzle(ctx context.Context, id string) *gen.Error {
	collection := dal.db.Collection("puzzles")
	// Ignore result from delete operation, we don't need it
	_, err := collection.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return &gen.Error{Code: 500, Message: err.Error()}
	}
	return nil
}
