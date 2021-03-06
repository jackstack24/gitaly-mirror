package featureflag

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc/metadata"
)

// OutgoingCtxWithFeatureFlag is used to enable a feature flag in the outgoing
// context metadata. The returned context is meant to be used in a client where
// the outcoming context is transferred to an incoming context.
func OutgoingCtxWithFeatureFlag(ctx context.Context, flag string) context.Context {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.New(map[string]string{})
	}
	md.Set(HeaderKey(flag), "true")
	return metadata.NewOutgoingContext(ctx, md)
}

// IncomingCtxWithFeatureFlag is used to enable a feature flag in the incoming
// context. This is NOT meant for use in clients that transfer the context
// across process boundaries.
func IncomingCtxWithFeatureFlag(ctx context.Context, flag string) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(map[string]string{})
	}
	md.Set(HeaderKey(flag), "true")
	return metadata.NewIncomingContext(ctx, md)
}

func OutgoingCtxWithRubyFeatureFlags(ctx context.Context, flags ...string) context.Context {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.New(map[string]string{})
	}

	for _, flag := range flags {
		md.Set(rubyHeaderKey(flag), "true")
	}

	return metadata.NewOutgoingContext(ctx, md)
}

func rubyHeaderKey(flag string) string {
	return fmt.Sprintf("gitaly-feature-ruby-%s", strings.ReplaceAll(flag, "_", "-"))
}
