package mongodb

import (
	"auth-sso/internal/domain/models"
	"auth-sso/internal/storage"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

const (
	duplicateKeyError = 11000
)

type Storage struct {
	client   *mongo.Client
	database string
}

// New creates a new instance of the MongoDB storage.
func New(uri string, database string) (*Storage, error) {
	const op = "storage.mongodb.New"

	clientOptions := options.Client().ApplyURI(uri)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{
		client:   client,
		database: database,
	}, nil
}

// Close closes instated mongodb connection.
func (s *Storage) Close(ctx context.Context) error {
	if s.client != nil {
		return s.client.Disconnect(ctx)
	}
	return nil
}

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (uid string, err error) {
	const op = "storage.mongodb.SaveUser"

	uid = uuid.New().String()

	collection := s.client.Database(s.database).Collection("users")
	document := bson.D{
		{"uniqueId", uid},
		{"email", email},
		{"passwordHash", passHash},
	}

	result, err := collection.InsertOne(ctx, document)
	if err != nil {
		var writeErr mongo.WriteException
		if errors.As(err, &writeErr) {
			for _, e := range writeErr.WriteErrors {
				if e.Code == duplicateKeyError {
					return "", fmt.Errorf("%s: %w", op, storage.ErrorUserExists)
				}
			}
		}

		return "", fmt.Errorf("%s: %w", op, err)
	}

	if _, ok := result.InsertedID.(primitive.ObjectID); ok {
		return uid, nil
	}

	return "", fmt.Errorf("%s: failed to get inserted AppID", op)
}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "storage.mongodb.User"

	collection := s.client.Database(s.database).Collection("users")
	filter := bson.M{"email": email}

	var user models.User

	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrorUserNotFound)
		}

		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) App(ctx context.Context, appID int) (models.App, error) {
	const op = "storage.mongodb.App"

	collection := s.client.Database(s.database).Collection("apps")
	filter := bson.M{"appID": appID}

	var user models.App

	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.App{}, fmt.Errorf("%s: %w", op, storage.ErrorAppNotFound)
		}

		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) UserById(ctx context.Context, id string) (models.User, error) {
	const op = "storage.mongodb.UserById"

	collection := s.client.Database(s.database).Collection("users")
	filter := bson.M{"uniqueId": id}

	var user models.User

	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrorUserNotFound)
		}

		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) CreateNewValidation(ctx context.Context,
	user models.User,
	documentType string,
) (string, error) {
	const op = "storage.mongodb.CreateNewValidation"

	uid := uuid.New().String()
	now := time.Now()

	collection := s.client.Database(s.database).Collection("validations")
	document := bson.D{
		{"verificationId", uid},
		{"user", user},
		{"documentType", documentType},
		{"status", "pending"},
		{"createdAt", now},
		{"updatedAt", now},
	}

	_, err := collection.InsertOne(ctx, document)
	if err != nil {
		var writeErr mongo.WriteException
		if errors.As(err, &writeErr) {
			for _, e := range writeErr.WriteErrors {
				if e.Code == duplicateKeyError {
					return "", fmt.Errorf("%s: %w", op, storage.ErrorUserExists)
				}
			}
		}

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return uid, nil
}
