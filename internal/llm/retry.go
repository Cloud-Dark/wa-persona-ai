package llm

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/rs/zerolog/log"
)

// RetryProvider wraps a Provider with retry logic and optional fallback.
type RetryProvider struct {
	primary     Provider
	fallback    Provider
	maxAttempts int
	initialMs   int
	maxMs       int
	multiplier  float64
}

// NewRetryProvider creates a retry wrapper around a provider.
func NewRetryProvider(primary, fallback Provider, maxAttempts, initialMs, maxMs int, multiplier float64) *RetryProvider {
	return &RetryProvider{
		primary:     primary,
		fallback:    fallback,
		maxAttempts: maxAttempts,
		initialMs:   initialMs,
		maxMs:       maxMs,
		multiplier:  multiplier,
	}
}

func (r *RetryProvider) Name() string {
	return r.primary.Name()
}

func (r *RetryProvider) Generate(ctx context.Context, req *Request) (*Response, error) {
	var lastErr error
	delay := float64(r.initialMs)

	for attempt := 1; attempt <= r.maxAttempts; attempt++ {
		resp, err := r.primary.Generate(ctx, req)
		if err == nil {
			return resp, nil
		}

		lastErr = err
		log.Warn().Err(err).Int("attempt", attempt).Str("provider", r.primary.Name()).Msg("LLM request failed")

		if attempt < r.maxAttempts {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(time.Duration(delay) * time.Millisecond):
			}
			delay = math.Min(delay*r.multiplier, float64(r.maxMs))
		}
	}

	// Try fallback if available
	if r.fallback != nil {
		log.Info().Str("fallback", r.fallback.Name()).Msg("switching to fallback LLM provider")
		resp, err := r.fallback.Generate(ctx, req)
		if err == nil {
			return resp, nil
		}
		return nil, fmt.Errorf("primary (%s): %w; fallback (%s): %v", r.primary.Name(), lastErr, r.fallback.Name(), err)
	}

	return nil, fmt.Errorf("all %d attempts failed: %w", r.maxAttempts, lastErr)
}
