package timerqueue

import (
	"context"
	"math/rand"
	"testing"
	"time"
)

type object struct {
	value int
}

var executed int

func (o *object) OnTimer(ctx context.Context, t time.Time) {
	executed++
}

func (o *object) Next(t time.Time) time.Time {
	return t.Add(time.Hour)
}

func populateQueue(t *testing.T, now time.Time) *Queue {
	q := New()

	count := 200
	objects := make([]*object, count)

	// Add a bunch of objects to the queue in random order.
	for i, j := range rand.Perm(count) {
		tm := now.Add(time.Duration(i+1) * time.Hour)
		objects[j] = &object{j}
		q.Schedule(objects[j], tm)
	}

	if q.Len() != count {
		t.Error("invalid queue length:", q.Len())
	}

	return q
}

func TestQueue(t *testing.T) {
	for iter := 0; iter < 100; iter++ {
		now := time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)
		queue := populateQueue(t, now)

		// Make sure objects are removed from the queue in order.
		for prev := now; queue.Len() > 0; {
			_, tm := queue.PopFirst()
			if tm.Sub(prev) != time.Hour {
				t.Errorf("Invalid queue ordering.\n"+
					"     Got: %v\n"+
					"Expected: %v\n", tm, prev.Add(time.Hour))
			}
			prev = tm
		}
	}
}

func TestAdvanceOnce(t *testing.T) {
	for iter := 0; iter < 100; iter++ {
		now := time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)
		queue := populateQueue(t, now)

		executed = 0
		count := queue.Len()
		lastTime := now.Add(time.Duration(count) * time.Hour)

		for adv := 0; adv < 5; adv++ {
			queue.AdvanceOnce(context.Background(), lastTime)
			if executed != count {
				t.Errorf("AdvanceOnce failed.\n"+
					"Should have executed %d times.\n"+
					"Only executed %d times.\n", count, executed)
			}
		}
	}
}

func TestAdvance(t *testing.T) {
	for iter := 0; iter < 100; iter++ {
		now := time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)
		queue := populateQueue(t, now)

		executed = 0
		count := queue.Len()
		lastTime := now.Add(time.Duration(count) * time.Hour)

		for adv := 0; adv < 5; adv++ {
			queue.Advance(context.Background(), lastTime)
			if executed != count {
				t.Errorf("Advance failed.\n"+
					"Should have executed %d times.\n"+
					"Only executed %d times.\n", count, executed)
			}
			if queue.Len() != count {
				t.Errorf("Advance failed.\n"+
					"Queue len should be %d.\n"+
					"But is %d.\n", count, queue.Len())
			}
		}
	}
}
