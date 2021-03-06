package featureflag

import (
	"context"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc/metadata"
)

var (
	flagChecks = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gitaly_feature_flag_checks_total",
			Help: "Number of enabled/disabled checks for Gitaly server side feature flags",
		},
		[]string{"flag", "enabled"},
	)
)

func init() {
	prometheus.MustRegister(flagChecks)
}

// IsEnabled checks if the feature flag is enabled for the passed context.
// Only return true if the metadata for the feature flag is set to "true"
func IsEnabled(ctx context.Context, flag string) bool {
	enabled := isEnabled(ctx, flag)
	flagChecks.WithLabelValues(flag, strconv.FormatBool(enabled)).Inc()
	return enabled
}

func isEnabled(ctx context.Context, flag string) bool {
	if flag == "" {
		return false
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return false
	}

	val, ok := md[HeaderKey(flag)]
	if !ok {
		return false
	}

	return len(val) > 0 && val[0] == "true"
}

// IsDisabled is the inverse operation of IsEnabled
func IsDisabled(ctx context.Context, flag string) bool {
	return !IsEnabled(ctx, flag)
}

const ffPrefix = "gitaly-feature-"

// HeaderKey returns the feature flag key to be used in the metadata map
func HeaderKey(flag string) string {
	return ffPrefix + strings.ReplaceAll(flag, "_", "-")
}

// AllEnabledFlags returns all feature flags that use the Gitaly metadata
// prefix and are enabled. Note: results will not be sorted.
func AllEnabledFlags(ctx context.Context) []string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil
	}

	ffs := make([]string, 0, len(md))

	for k, v := range md {
		if !strings.HasPrefix(k, ffPrefix) {
			continue
		}
		if len(v) > 0 && v[0] == "true" {
			ffs = append(ffs, strings.TrimPrefix(k, ffPrefix))
		}
	}

	return ffs
}
