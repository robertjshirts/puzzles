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

func (dal *MongoDAL) GetAllPuzzles(ctx context.Context) ([]gen.Puzzle, error) {
	var puzzles []gen.Puzzle
	collection := dal.db.Collection("puzzles")
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &puzzles); err != nil {
		return nil, err
	}
	return puzzles, nil
}

func (dal *MongoDAL) GetPuzzle(ctx context.Context, id string) (*gen.Puzzle, error) {
	var puzzle gen.Puzzle
	collection := dal.db.Collection("puzzles")
	if err := collection.FindOne(ctx, bson.M{"id": id}).Decode(&puzzle); err != nil {
		return nil, err
	}
	return &puzzle, nil
}

func (dal *MongoDAL) AddPuzzle(ctx context.Context, puzzle gen.NewPuzzle) (*gen.Puzzle, error) {
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
	_, err := collection.InsertOne(ctx, newPuzzle)
	if err != nil {
		return nil, err
	}
	return &newPuzzle, nil
}

func (dal *MongoDAL) UpdatePuzzle(ctx context.Context, id string, puzzle gen.NewPuzzle) (*gen.Puzzle, error) {
	collection := dal.db.Collection("puzzles")
	_, err := collection.UpdateOne(ctx, bson.M{"id": id}, bson.D{{"$set", puzzle}})
	if err != nil {
		return nil, err
	}
	updatedPuzzle := gen.Puzzle{
		Id:          id,
		Name:        puzzle.Name,
		Description: puzzle.Description,
		Price:       puzzle.Price,
		Quantity:    puzzle.Quantity,
		Type:        puzzle.Type,
	}
	return &updatedPuzzle, nil
}

func (dal *MongoDAL) DeletePuzzle(ctx context.Context, id string) error {
	collection := dal.db.Collection("puzzles")
	_, err := collection.DeleteOne(ctx, bson.M{"id": id})
	return err
}
