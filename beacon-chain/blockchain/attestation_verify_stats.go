package blockchain

import "sync"

// AttestationVerifyStat store the statistics of attestation verification.
type AttestationVerifyStat struct {
	successfulCount uint
	failureCount    uint
	failureReasons  map[string]uint

	// mutex for successfulCount
	sms sync.Mutex

	// mutex for failureCount and failureReasons
	smf sync.Mutex
}

// Successful increments the successful count.
func (s *AttestationVerifyStat) Successful() {
	s.sms.Lock()
	defer s.sms.Unlock()

	s.successfulCount++
}

// Failure increments the failure count and stores the reason.
func (s *AttestationVerifyStat) Failure(reason string) {
	s.smf.Lock()
	defer s.smf.Unlock()

	s.failureCount++

	if s.failureReasons == nil {
		s.failureReasons = make(map[string]uint)
	}

	s.failureReasons[reason]++
}

// Stats returns the statistics.
func (s *AttestationVerifyStat) Stats() AttestationVerifyStats {
	s.sms.Lock()
	defer s.sms.Unlock()

	s.smf.Lock()
	defer s.smf.Unlock()

	return AttestationVerifyStats{
		SuccessfulCount: s.successfulCount,
		FailureCount:    s.failureCount,
		FailureReasons:  s.failureReasons,
	}
}

// AttestationVerifyStats the statistics of attestation verification.
type AttestationVerifyStats struct {
	SuccessfulCount uint
	FailureCount    uint
	FailureReasons  map[string]uint
}
