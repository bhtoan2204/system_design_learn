package main

import (
	"context"
	"fmt"

	"event_sourcing_golang/eventstore"
	"event_sourcing_golang/eventstore/repos"
	"event_sourcing_golang/pkg/eventsourcing"
)

// Demo domain: BankAccount aggregate (no projections)
// Events
type AccountOpened struct{ Owner string }
type MoneyDeposited struct{ Amount int }
type MoneyWithdrawn struct{ Amount int }

type BankAccount struct {
	eventsourcing.AggregateRoot
	Owner   string
	Balance int
}

func (a *BankAccount) RegisterEvents(reg eventsourcing.RegisterEventsFunc) error {
	return reg(&AccountOpened{}, &MoneyDeposited{}, &MoneyWithdrawn{})
}

func (a *BankAccount) Transition(e eventsourcing.Event) error {
	switch v := e.Data.(type) {
	case *AccountOpened:
		a.Owner = v.Owner
	case *MoneyDeposited:
		a.Balance += v.Amount
	case *MoneyWithdrawn:
		a.Balance -= v.Amount
	}
	return nil
}

func main() {
	ctx := context.Background()

	// Serializer with aggregate registration
	s := eventsourcing.NewSerializer()
	_ = s.RegisterAggregate(&BankAccount{})

	// Use in-memory event store
	r := repos.NewInMemory(s)
	as := eventstore.NewAggregateStore(r)

	// Open a new account
	acc := &BankAccount{}
	_ = acc.SetID("acc-1")

	// Perform multiple transactions to trigger snapshots (every 10 events)
	_ = acc.ApplyChange(acc, &AccountOpened{Owner: "Alice"})
	for i := 0; i < 1000; i++ {
		_ = acc.ApplyChange(acc, &MoneyDeposited{Amount: 10})
		_ = acc.ApplyChange(acc, &MoneyWithdrawn{Amount: 5})
	}

	// Persist events
	if err := as.Save(ctx, acc); err != nil {
		panic(err)
	}

	// Load fresh aggregate
	loaded := &BankAccount{}
	if err := as.Get(ctx, "acc-1", loaded); err != nil {
		panic(err)
	}

	fmt.Printf("Owner=%s Balance=%d Version=%d\n", loaded.Owner, loaded.Balance, loaded.Root().Version())

	// List and print all events
	// Note: we need an instance with aggregate type set for deserialization
	tmp := &BankAccount{}
	_ = tmp.SetID("acc-1")
	events, err := r.EventStore().List(ctx, "acc-1", tmp)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Total events: %d\n", len(events))
	for _, e := range events[:5] { // print first few
		fmt.Printf("v=%d type=%s data=%+v\n", e.Version, e.EventType, e.Data)
	}
	if len(events) > 5 {
		last := events[len(events)-1]
		fmt.Printf("... last v=%d type=%s\n", last.Version, last.EventType)
	}

	// Snapshot info
	if ver, ok := r.EventStore().SnapshotVersion(ctx, "acc-1"); ok {
		fmt.Printf("Snapshot exists at version=%d\n", ver)
	} else {
		fmt.Println("No snapshot yet")
	}
}
