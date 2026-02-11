package geometrias

import (
	"sort"
	"sync"
	"time"

	pb "github.com/odysseia-greek/makedonia/eukleides/proto"
)

type GlobalKey struct{ Service, Word string }
type SessKey struct{ Session, Service, Word string }

type Counter struct {
	Count    int64
	LastUsed time.Time
}

type Store struct {
	mu      sync.RWMutex
	global  map[GlobalKey]*Counter
	session map[SessKey]*Counter
}

func NewStore() *Store {
	return &Store{
		global:  make(map[GlobalKey]*Counter, 1024),
		session: make(map[SessKey]*Counter, 1024),
	}
}

// Inc increments both global and per-session counters in a single critical section.
func (s *Store) Inc(sessionID, service, word string, ts time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// global
	gk := GlobalKey{Service: service, Word: word}
	if c := s.global[gk]; c != nil {
		c.Count++
		if ts.After(c.LastUsed) {
			c.LastUsed = ts
		}
	} else {
		s.global[gk] = &Counter{Count: 1, LastUsed: ts}
	}

	// per-session
	sk := SessKey{Session: sessionID, Service: service, Word: word}
	if c := s.session[sk]; c != nil {
		c.Count++
		if ts.After(c.LastUsed) {
			c.LastUsed = ts
		}
	} else {
		s.session[sk] = &Counter{Count: 1, LastUsed: ts}
	}
}

type row struct {
	service  string
	word     string
	count    int64
	lastUsed time.Time
}

func (s *Store) TopFiveGlobal() []*pb.TopFive {
	s.mu.RLock()
	out := make([]row, 0, len(s.global))
	for k, v := range s.global {
		out = append(out, row{service: k.Service, word: k.Word, count: v.Count, lastUsed: v.LastUsed})
	}
	s.mu.RUnlock()
	return top5ToProto(out)
}

// Top 5 within a service (global counters filtered by service)
func (s *Store) TopFiveByService(service string) []*pb.TopFive {
	s.mu.RLock()
	out := make([]row, 0, 32)
	for k, v := range s.global {
		if k.Service == service {
			out = append(out, row{service: k.Service, word: k.Word, count: v.Count, lastUsed: v.LastUsed})
		}
	}
	s.mu.RUnlock()
	return top5ToProto(out)
}

func (s *Store) TopFiveForSession(session string) []*pb.TopFive {
	s.mu.RLock()
	out := make([]row, 0, 32)
	for k, v := range s.session {
		if k.Session == session {
			out = append(out, row{service: k.Service, word: k.Word, count: v.Count, lastUsed: v.LastUsed})
		}
	}
	s.mu.RUnlock()
	return top5ToProto(out)
}

// Sort by count desc, then lastUsed desc; take 5; convert to proto.
func top5ToProto(rows []row) []*pb.TopFive {
	if len(rows) == 0 {
		return []*pb.TopFive{}
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].count != rows[j].count {
			return rows[i].count > rows[j].count
		}
		return rows[i].lastUsed.After(rows[j].lastUsed)
	})
	n := 5
	if len(rows) < n {
		n = len(rows)
	}
	out := make([]*pb.TopFive, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, &pb.TopFive{
			ServiceName: rows[i].service,
			Word:        rows[i].word,
			LastUsed:    rows[i].lastUsed.UTC().Format(time.RFC3339Nano),
			Count:       rows[i].count, // make sure TopFive.count is int64 in proto
		})
	}
	return out
}
