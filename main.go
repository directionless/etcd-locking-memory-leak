package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"github.com/kolide/kit/stringutil"
	"github.com/pkg/errors"
)

func lockOne(etcd *clientv3.Client) error {
	ctx := context.TODO()

	session, err := concurrency.NewSession(etcd)
	if err != nil {
		return errors.Wrap(err, "error creating etcd session")
	}
	defer session.Close()

	// This formatting matches our use case. It's a bit out of place here.
	lockName := path.Join(
		"/locks",
		stringutil.RandomString(32),                                                          //ulid
		fmt.Sprintf("pack:_%s:%s", stringutil.RandomString(20), stringutil.RandomString(20)), //query
		fmt.Sprintf("%s-%s", stringutil.RandomString(16), stringutil.RandomString(16)),       //host
	)
	lock := concurrency.NewMutex(session, lockName)
	if err := lock.Lock(ctx); err != nil {
		return errors.Wrap(err, "error acquiring lock")
	}
	defer lock.Unlock(ctx)

	// Locked! Let's do some work.
	// Not really, this is a lock benchmark...
	fmt.Printf("locked %s\n", lockName)
	timer := time.NewTimer(time.Millisecond * time.Duration(rand.Intn(1000)))
	<-timer.C
	return nil
}

func locker(etcd *clientv3.Client) {
	for {
		err := lockOne(etcd)
		if err != nil {
			fmt.Printf("GOT ERR: %v\n", err)
		}
	}
}

func main() {
	ctx := context.Background()

	// Set endpoints, and allow ENV to override
	endpoints := []string{"http://localhost:2379"}
	if env, ok := os.LookupEnv("ENDPOINTS"); ok {
		endpoints = strings.Split(env, ",")
	}

	etcd, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 2 * time.Second,
	})
	if err != nil {
		panic(errors.Wrap(err, "error creating etcd client"))
	}

	fmt.Printf("Starting lock actors. Talking to %s\n", endpoints)
	for i := 1; i <= 10; i++ {
		go locker(etcd)
	}

	// Wait for control-c
	<-ctx.Done()

}
