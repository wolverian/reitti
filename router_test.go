package reitti

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExampleRouter_simple() {
	r := &Router{}
	r.Add("repos/{owner}/{repo}/issues", func(ctx context.Context, owner, repo string) (any, error) {
		return fmt.Sprintf("owner=%s, repo=%s", owner, repo), nil
	})

	handler, _ := r.Match("repos/wolverian/reitti/issues")
	result, _ := handler(context.Background())
	fmt.Printf("result: %s\n", result)
	_, err := r.Match("foobar")
	fmt.Printf("error: %s\n", err)

	// Output:
	// result: owner=wolverian, repo=reitti
	// error: no handler found for route: "foobar"
}

func TestRouter(t *testing.T) {
	tests := []struct {
		name         string
		routes       []string
		path         string
		wantMatchErr string
		want         []string
		wantErr      error
	}{
		{
			name:         "empty router",
			routes:       []string{},
			path:         "repos/wolverian/reitti/issues",
			wantMatchErr: `no handler found for route: "repos/wolverian/reitti/issues"`,
		},
		{
			name:         "empty path",
			routes:       []string{"repos/{owner}/{repo}/issues"},
			wantMatchErr: `no handler found for route: ""`,
		},
		{
			name:         "no matching handler",
			routes:       []string{"repos/{owner}/{repo}/issues"},
			path:         "repos/wolverian/reitti",
			wantMatchErr: `no handler found for route: "repos/wolverian/reitti"`,
		},
		{
			name:   "github issues",
			routes: []string{"repos/{owner}/{repo}/issues"},
			path:   "repos/wolverian/reitti/issues",
			want:   []string{"wolverian", "reitti"},
		},
		{
			name: "multiple routes",
			routes: []string{
				"repos/{owner}",
				"repos/{owner}/{repo}",
				"repos/{owner}/{repo}/issues/{issue}",
				"repos/{owner}/{repo}/issues",
				"repos/{owner}/{repo}/issues/{issue}",
			},
			path: "repos/wolverian/reitti/issues",
			want: []string{"wolverian", "reitti"},
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
			if tt.wantMatchErr != "" {
				assert.EqualError(t, err, tt.wantMatchErr, "we get the expected error when no matching handler is found")
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
