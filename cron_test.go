package timerqueue_test

import (
	"context"
	"fmt"
	"time"

	"github.com/gorhill/cronexpr"
	"github.com/treywelsh/timerqueue"
)

type Event struct {
	id   int
	Expr *cronexpr.Expression
}

func (e *Event) OnTimer(ctx context.Context, t time.Time) {
	fmt.Printf("  Event %d executed at %v\n", e.id, t)
}

func (e *Event) Next(tm time.Time) time.Time {
	return e.Expr.Next(tm)
}

// Schedule several events with a timerqueue, and dispatch
// them by calling Advance.
func ExampleCronQueue() {
	queue := timerqueue.New()

	// Start date
	tm := time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)

	for i := 1; i <= 7; i++ {

		// Occurences computation based on cron format
		cronStr := fmt.Sprintf("0 %d * * *", i)
		expr, err := cronexpr.Parse(cronStr)
		if err != nil {
			fmt.Printf("can't parse: %s: %s\n", cronStr, err)
			return
		}

		e := &Event{id: i, Expr: expr}

		queue.Schedule(e, e.Next(tm))
	}

	fmt.Println("Advancing to Jan 2...")
	queue.Advance(context.Background(), time.Date(2015, 1, 2, 0, 0, 0, 0, time.UTC))

	fmt.Println("Advancing to Jan 3...")
	queue.Advance(context.Background(), time.Date(2015, 1, 3, 0, 0, 0, 0, time.UTC))

	// Output:
	// Advancing to Jan 2...
	//   Event 1 executed at 2015-01-01 01:00:00 +0000 UTC
	//   Event 2 executed at 2015-01-01 02:00:00 +0000 UTC
	//   Event 3 executed at 2015-01-01 03:00:00 +0000 UTC
	//   Event 4 executed at 2015-01-01 04:00:00 +0000 UTC
	//   Event 5 executed at 2015-01-01 05:00:00 +0000 UTC
	//   Event 6 executed at 2015-01-01 06:00:00 +0000 UTC
	//   Event 7 executed at 2015-01-01 07:00:00 +0000 UTC
	// Advancing to Jan 3...
	//   Event 1 executed at 2015-01-02 01:00:00 +0000 UTC
	//   Event 2 executed at 2015-01-02 02:00:00 +0000 UTC
	//   Event 3 executed at 2015-01-02 03:00:00 +0000 UTC
	//   Event 4 executed at 2015-01-02 04:00:00 +0000 UTC
	//   Event 5 executed at 2015-01-02 05:00:00 +0000 UTC
	//   Event 6 executed at 2015-01-02 06:00:00 +0000 UTC
	//   Event 7 executed at 2015-01-02 07:00:00 +0000 UTC
}
