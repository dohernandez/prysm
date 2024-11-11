package blockchain_test

import (
	"sync"
	"testing"

	"github.com/prysmaticlabs/prysm/v5/beacon-chain/blockchain"
	"github.com/prysmaticlabs/prysm/v5/testing/require"
)

func TestAttestationVerifyStat_Stats_all_successfully(t *testing.T) {
	t.Parallel()

	s := &blockchain.AttestationVerifyStat{}

	for range 10 {
		s.Successful()
	}

	want := blockchain.AttestationVerifyStats{
		SuccessfulCount: 10,
		FailureCount:    0,
		FailureReasons:  nil,
	}

	got := s.Stats()

	require.Equal(t, want, got)
}

func TestAttestationVerifyStat_Stats_all_failure_same_reason(t *testing.T) {
	t.Parallel()

	s := &blockchain.AttestationVerifyStat{}

	for range 10 {
		s.Failure("reason")
	}

	want := blockchain.AttestationVerifyStats{
		SuccessfulCount: 0,
		FailureCount:    10,
		FailureReasons:  map[string]uint{"reason": 10},
	}

	got := s.Stats()

	require.Equal(t, want, got)
}

func TestAttestationVerifyStat_Stats(t *testing.T) {
	t.Parallel()

	s := &blockchain.AttestationVerifyStat{}

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()

		for range 10 {
			s.Successful()
		}
	}()

	var (
		reason1 = "reason1"
		reason2 = "reason2"
	)

	wg.Add(1)

	go func() {
		defer wg.Done()

		for i := range 5 {
			if i%2 == 0 {
				s.Failure(reason1)

				continue
			}

			s.Failure(reason2)
		}
	}()

	wg.Wait()

	want := blockchain.AttestationVerifyStats{
		SuccessfulCount: 10,
		FailureCount:    5,
		FailureReasons:  map[string]uint{"reason1": 3, "reason2": 2},
	}

	got := s.Stats()

	require.Equal(t, want, got)
}
