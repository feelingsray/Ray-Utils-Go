package driver

import (
	"bytes"
	"context"
	"fmt"
	
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewClient(host string, port int, username, password string, ctx context.Context) (*mongo.Client, error) {
	var (
		buf     bytes.Buffer
		charMap = map[int32]string{'@': "%40", ':': "%3A", '/': "%2F", '%': "%25"}
	)
	for _, c := range password {
		if v, ok := charMap[c]; ok {
			for _, i := range v {
				buf.WriteRune(i)
			}
		} else {
			buf.WriteRune(c)
		}
	}
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d", username, buf.String(), host, port)
	connect, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	return connect, nil
}

func NewClientNoAuth(host string, port int, ctx context.Context) (*mongo.Client, error) {
	uri := fmt.Sprintf("mongodb://%s:%d", host, port)
	connect, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	return connect, nil
}
