package core

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DefaultFieldHook defines the interface to change default fields by hook
type DefaultFieldHook interface {
	DefaultUpdateAt()
	DefaultCreateAt()
	DefaultId()
}

// MgoModel defines the default fields to handle when operation happens
// import the MgoModel in document struct to make it working
type MgoModel struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	CreateAt time.Time          `json:"create_at" bson:"createAt"`
	UpdateAt time.Time          `json:"update_at" bson:"updateAt"`
}

// DefaultUpdateAt changes the default updateAt field
func (df *MgoModel) DefaultUpdateAt() {
	df.UpdateAt = time.Now().Local()
}

// DefaultCreateAt changes the default createAt field
func (df *MgoModel) DefaultCreateAt() {
	if df.CreateAt.IsZero() {
		df.CreateAt = time.Now().Local()
	}
}

// DefaultId changes the default _id field
func (df *MgoModel) DefaultId() {
	if df.ID.IsZero() {
		df.ID = primitive.NewObjectID()
	}
}
