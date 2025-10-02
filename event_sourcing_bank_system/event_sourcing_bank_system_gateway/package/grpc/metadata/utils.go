package metadata

import (
	"context"
	"errors"
	"strconv"

	"google.golang.org/grpc/metadata"
)

func GetUserLoginID(ctx context.Context) (int64, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, errors.New("metadata in ctx not exists")
	}
	mdUID := md.Get("uid")
	if len(mdUID) == 0 {
		return 0, errors.New("metadata uid in ctx not exists")
	}
	uid, err := strconv.ParseInt(mdUID[0], 10, 64) // metadata.Get returns an array of values for the key
	if err != nil {
		return 0, err
	}
	if uid < 1 {
		return 0, errors.New("uid invalid")
	}
	return uid, nil
}

func GetToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("token in ctx not exists")
	}
	token := md.Get("token")
	if len(token) == 0 {
		return "", errors.New("token in ctx not exists")
	}
	return token[0], nil
}

func GetIPAddress(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("ip in ctx not exists")
	}
	ip := md.Get("ip_address")
	if len(ip) == 0 {
		return "", errors.New("ip in ctx not exists")
	}

	return ip[0], nil
}

func GetLanguage(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "vi", nil
	}
	language := md.Get("language")
	if len(language) == 0 {
		return "vi", nil
	}
	return language[0], nil
}
