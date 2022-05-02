package informer

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type impl[ObjectContent any] struct {
	Interface[ObjectContent]

	ListFunc func() []Object[ObjectContent]
}

func (impl impl[ObjectContent]) List(ctx context.Context) ([]Object[ObjectContent], error) {
	return impl.ListFunc(), nil
}

func Test_informer(t *testing.T) {
	var list []Object[int]
	var intf Interface[int] = &impl[int]{
		ListFunc: func() []Object[int] {
			t.Log("LIST")
			return list
		},
	}

	inf, err := New(Config{
		ResyncPeriod: 500 * time.Millisecond,
	}, intf)
	require.NoError(t, err)
	events := []string{}
	inf.AddHandler(EventHandlerFunc[int](func(event Event[int]) {
		name := ""
		if event.OldObject != nil {
			name = event.OldObject.Name
		} else {
			name = event.NewObject.Name
		}
		t.Logf("EVENT: %s %s", event.Type, name)
		events = append(events, string(event.Type)+" "+name)
	}))
	ctx, cancel := context.WithTimeout(context.Background(), 700*time.Millisecond)
	defer cancel()

	list = []Object[int]{
		{ObjectMeta: ObjectMeta{Name: "1", ResourceVersion: "1"}, Content: 1},
		{ObjectMeta: ObjectMeta{Name: "2", ResourceVersion: "1"}, Content: 2},
		{ObjectMeta: ObjectMeta{Name: "3", ResourceVersion: "1"}, Content: 3},
	}
	go func() {
		time.Sleep(300 * time.Millisecond)
		list = []Object[int]{
			{ObjectMeta: ObjectMeta{Name: "1", ResourceVersion: "2"}, Content: 1},
			{ObjectMeta: ObjectMeta{Name: "2", ResourceVersion: "1"}, Content: 22}, // this would not encourage an event
			{ObjectMeta: ObjectMeta{Name: "4", ResourceVersion: "1"}, Content: 4},
		}
	}()
	inf.Run(ctx)
	requireEqualAfterSort := func(exp, act []string) {
		sort.Strings(exp)
		sort.Strings(act)
		require.Equal(t, exp, act)
	}
	requireEqualAfterSort([]string{
		"CREATE 1",
		"CREATE 2",
		"CREATE 3",
		"UPDATE 1",
		"DELETE 3",
		"CREATE 4",
	}, events)
}
