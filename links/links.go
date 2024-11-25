package links

import (
	"context"
	"fmt"
	"net/url"
)

type contextKey string

const (
	ctxProtocol   contextKey = "protocol"
	ctxHost       contextKey = "host"
	ctxPort       contextKey = "port"
	ctxUrlVersion contextKey = "urlVersion"
)

func URLBuild(ctx context.Context, oldURL string) (string, error) {
	parsedURL, err := url.Parse(oldURL)
	if err != nil {
		fmt.Println("Error parsing old URL:", err)
		return "", err
	}

	fmt.Printf("\nOld URL: %s\n", parsedURL.String())

	newProto := ctx.Value(ctxProtocol).(string)
	newHost := ctx.Value(ctxHost).(string)
	newPort := ctx.Value(ctxPort).(string)
	newUrlVersion := ctx.Value(ctxUrlVersion).(string)

	parsedURL.Scheme = newProto

	if newPort == "" {
		parsedURL.Host = newHost
	} else {
		parsedURL.Host = newHost + ":" + newPort
	}

	if newUrlVersion != "" {
		parsedURL.Path = newUrlVersion + parsedURL.Path
	}

	fmt.Printf("New URL: %s\n", parsedURL.String())

	return parsedURL.String(), nil
}
