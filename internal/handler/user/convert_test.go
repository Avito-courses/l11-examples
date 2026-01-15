package user

import (
	"reflect"
	"testing"
	"time"

	"github.com/Avito-courses/l11-examples/internal/model/user"
)

func TestModelToResponse(t *testing.T) {
	type args struct {
		user user.User
	}
	tests := []struct {
		name string
		args args
		want User
	}{
		{
			name: "ok",
			args: args{
				user: user.User{
					ID:        0,
					Name:      "",
					Phone:     "",
					Rating:    0,
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
				},
			},
			want: User{
				ID:     0,
				Name:   "",
				Phone:  "",
				Rating: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ModelToResponse(tt.args.user); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ModelToResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}
