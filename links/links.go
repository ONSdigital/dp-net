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
	ctxPathPrefix contextKey = "pathPrefix"
)

func URLBuild(ctx context.Context, oldURL string) (string, error) {
	parsedURL, err := url.Parse(oldURL)
	if err != nil {
		return "", fmt.Errorf("error parsing old URL: %v", err)
	}

	fmt.Printf("\nOld URL: %s\n", parsedURL.String())

	newProto := ctx.Value(ctxProtocol).(string)
	newHost := ctx.Value(ctxHost).(string)
	newPort := ctx.Value(ctxPort).(string)
	newPathPrefix := ctx.Value(ctxPathPrefix).(string)

	parsedURL.Scheme = newProto

	if newPort == "" {
		parsedURL.Host = newHost
	} else {
		parsedURL.Host = newHost + ":" + newPort
	}

	if newPathPrefix != "" {
		parsedURL.Path, err = url.JoinPath(newPathPrefix, parsedURL.Path)
		if err != nil {
			return "", fmt.Errorf("error joining paths: %v", err)
		}
	}

	fmt.Printf("New URL: %s\n", parsedURL.String())

	return parsedURL.String(), nil
}
