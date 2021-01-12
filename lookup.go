package feeds

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethersphere/feeds/lookup"
	"github.com/ethersphere/swarm/storage"
	"k8s.io/helm/log"
)

var retrieveTimeout = 10 * time.Second

// Lookup retrieves a specific or latest feed update
// Lookup works differently depending on the configuration of `query`
// See the `query` documentation and helper functions:
// `NewQueryLatest` and `NewQuery`
func Lookup(ctx context.Context, ls LoadSaver, user common.Address, topic []byte) (interface{}, error) {

	timeLimit := time.Now()
	if timeLimit == 0 { // if time limit is set to zero, the user wants to get the latest update
		timeLimit = TimestampProvider.Now().Time
	}

	var readCount int32

	// Invoke the lookup engine.
	// The callback will be called every time the lookup algorithm needs to guess
	requestPtr, err := lookup.FluzCapacitorAlgorithm(ctx, timeLimit, lookup.NoClue, func(ctx context.Context, epoch lookup.Epoch, now uint64) (interface{}, error) {
		atomic.AddInt32(&readCount, 1)
		id := ID{
			Feed:  query.Feed,
			Epoch: epoch,
		}
		ctx, cancel := context.WithTimeout(ctx, retrieveTimeout)
		defer cancel()

		r := storage.NewRequest(id.Addr())

		b, err := ls.Load(ctx, id.Addr())

		//ch, err := h.chunkStore.Get(ctx, chunk.ModeGetLookup, r)
		//if err != nil {
		//if err == context.DeadlineExceeded || err == storage.ErrNoSuitablePeer { // chunk not found
		//return nil, nil
		//}
		//return nil, err
		//}

		var request Request
		if err := request.fromChunk(ch); err != nil {
			return nil, nil
		}
		if request.Time <= timeLimit {
			return &request, nil
		}
		return nil, nil
	})
	if err != nil {
		return nil, err
	}

	log.Info(fmt.Sprintf("Feed lookup finished in %d lookups", readCount))

	request, _ := requestPtr.(*Request)
	if request == nil {
		return nil, NewError(ErrNotFound, "no feed updates found")
	}
	return h.updateCache(request)

}
