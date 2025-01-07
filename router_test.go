package reitti

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouter(t *testing.T) {
	tests := []struct {
		name         string
		routes       []string
		path         string
		wantMatchErr error
		want         []string
		wantErr      error
	}{
		{"empty router", []string{}, "repos/wolverian/reitti/issues", ErrNoRoute, nil, nil},
		{"empty path", []string{"repos/{owner}/{repo}/issues"}, "", ErrNoRoute, nil, nil},
		{"no match", []string{"repos/{owner}/{repo}/issues"}, "repos/wolverian/reitti", ErrNoRoute, nil, nil},
		{
			name:   "github issues",
			routes: []string{"repos/{owner}/{repo}/issues"},
			path:   "repos/wolverian/reitti/issues",
			want:   []string{"wolverian", "reitti"},
		},
		{
			name:   "multiple routes",
			routes: []string{"repos/{owner}", "repos/{owner}/{repo}", "repos/{owner}/{repo}/issues/{issue}", "repos/{owner}/{repo}/issues", "repos/{owner}/{repo}/issues/{issue}"},
			path:   "repos/wolverian/reitti/issues",
			want:   []string{"wolverian", "reitti"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			t.Cleanup(cancel)

			r := &Router{}
			for _, route := range tt.routes {
				r.Add(route, func(ctx context.Context, args ...string) (any, error) {
					return args, nil
				})
			}
			handler, err := r.Match(tt.path)
			if tt.wantMatchErr != nil {
				assert.EqualError(t, err, tt.wantMatchErr.Error(), "we get the expected error when no matching handler is found")
				return
			}
			assert.NoError(t, err, "the handler is found")
			result, err := handler(ctx)
			if tt.wantErr != nil {
				assert.EqualError(t, err, tt.wantErr.Error(), "the handler returns the expected error")
				return
			}
			assert.NoError(t, err, "the handler does not return an error")
			assert.Equal(t, tt.want, result, "the handler returns the expected result")
		})
	}
}
