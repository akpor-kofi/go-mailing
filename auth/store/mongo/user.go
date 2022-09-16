package mongo

import (
	"context"
	"github.com/akpor-kofi/auth/models"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var ctx = context.Background()
var ms = 5 * time.Second

type UserStore struct {
	coll *mgm.Collection
}

func NewUserStore() *UserStore {
	user := &models.User{}

	coll := mgm.Coll(user)

	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "email", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	}

	coll.Indexes().CreateOne(ctx, indexModel)

	return &UserStore{coll}
}

func (u *UserStore) Create(user *models.User) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, ms)
	defer cancel()

	err := u.coll.CreateWithCtx(ctxWithTimeout, user)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserStore) Get(id string) (*models.User, error) {
	user := new(models.User)

	err := u.coll.FindByID(id, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserStore) List() ([]*models.User, error) {
	var users []*models.User

	err := u.coll.SimpleFind(&users, bson.M{})
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *UserStore) Update(id string, user *models.User) (*models.User, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, ms)
	defer cancel()

	_, err := u.coll.UpdateByID(ctxWithTimeout, id, user)
	if err != nil {
		return nil, err
	}

	updatedUser := new(models.User)

	err = u.coll.FindByID(id, updatedUser)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (u *UserStore) Delete(id string) error {
	//TODO implement me
	panic("implement me")
}

func (u *UserStore) GetByEmail(email string) (*models.User, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, ms)
	defer cancel()

	user := new(models.User)

	sr := u.coll.FindOne(ctxWithTimeout, bson.M{"email": email})

	if sr.Err() != nil {
		return nil, sr.Err()
	}

	sr.Decode(user)

	return user, nil
}
