package db_test

import (
	"context"
	"testing"

	"github.com/dai65527/microservice-handson/platform/db/db"
	"github.com/dai65527/microservice-handson/platform/db/model"
	"github.com/stretchr/testify/assert"
)

func TestDB(t *testing.T) {
	d := db.New()

	// Customer
	c, err := d.GetCustomer(context.TODO(), "7c0cde05-4df0-47f4-94c4-978dd9f56e5c")
	assert.NoError(t, err)
	assert.Equal(t, "7c0cde05-4df0-47f4-94c4-978dd9f56e5c", c.ID)
	assert.Equal(t, "goldie", c.Name)

	_, err = d.GetCustomer(context.TODO(), "hogehoge")
	assert.Equal(t, err, db.ErrNotFound)

	c, err = d.CreateCustomer(context.TODO(), "gopher")
	assert.NoError(t, err)
	assert.NotEmpty(t, c.ID)
	assert.Equal(t, "gopher", c.Name)
	id := c.ID

	_, err = d.CreateCustomer(context.TODO(), "gopher")
	assert.Equal(t, err, db.ErrAlreadyExists)

	c, err = d.GetCustomer(context.TODO(), id)
	assert.NoError(t, err)
	assert.Equal(t, id, c.ID)
	assert.Equal(t, "gopher", c.Name)

	c, err = d.GetCustomerByName(context.TODO(), "gopher")
	assert.NoError(t, err)
	assert.Equal(t, id, c.ID)
	assert.Equal(t, "gopher", c.Name)

	_, err = d.GetCustomerByName(context.TODO(), "hogehoge")
	assert.Equal(t, err, db.ErrNotFound)

	// Item
	newItem := &model.Item{
		Title:      "iPhone 13",
		Price:      42000,
		CustomerID: c.ID,
	}
	item, err := d.CreateItem(context.TODO(), newItem)
	assert.NoError(t, err)
	assert.NotEmpty(t, item.ID)
	assert.Equal(t, "iPhone 13", item.Title)
	assert.Equal(t, uint64(42000), item.Price)
	assert.Equal(t, c.ID, item.CustomerID)
	id = item.ID

	item, err = d.GetItem(context.TODO(), id)
	assert.NoError(t, err)
	assert.NotEmpty(t, item.ID)
	assert.Equal(t, "iPhone 13", item.Title)
	assert.Equal(t, uint64(42000), item.Price)
	assert.Equal(t, c.ID, item.CustomerID)

	_, err = d.GetItem(context.TODO(), "hogehoge")
	assert.Equal(t, err, db.ErrNotFound)

	items, err := d.GetAllItems(context.TODO())
	assert.NoError(t, err)
	assert.Len(t, items, 3) // includes initial dataset
}
