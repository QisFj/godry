package informer

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_informer(t *testing.T) {
	var list []Object[int]
	var lw ListAndWatch[int] = &PollListAndWatch[int]{
		ListFunc: func(ctx context.Context) ([]Object[int], error) {
			return list, nil
		},
		PollPeriod: 500 * time.Millisecond,
	}

	inf, err := New(lw)
	require.NoError(t, err)
	events := []string{}
	inf.AddEventandler(EventHandlerFunc[int](func(event Event[int]) {
		name := ""
		if event.OldObject != nil {
			name = event.OldObject.Name
		} else {
			name = event.NewObject.Name
		}
		t.Logf("EVENT: %s %s", event.Type, name)
		events = append(events, string(event.Type)+" "+name)
	}))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	list = []Object[int]{
		{ObjectMeta: ObjectMeta{Name: "1", ResourceVersion: "1"}, Content: 1},
		{ObjectMeta: ObjectMeta{Name: "2", ResourceVersion: "1"}, Content: 2},
		{ObjectMeta: ObjectMeta{Name: "3", ResourceVersion: "1"}, Content: 3},
	}
	go func() {
		WaitSynced(nil, inf)
		list = []Object[int]{
			{ObjectMeta: ObjectMeta{Name: "1", ResourceVersion: "2"}, Content: 1},
			{ObjectMeta: ObjectMeta{Name: "2", ResourceVersion: "1"}, Content: 22}, // this would not encourage an event
			{ObjectMeta: ObjectMeta{Name: "4", ResourceVersion: "1"}, Content: 4},
		}
		time.Sleep(500 * time.Millisecond) // wait a pool
		cancel()
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
