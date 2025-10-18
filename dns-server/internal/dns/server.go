package dns

import (
	"context"
	"dns-server/internal/config"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	miekg_dns "github.com/miekg/dns"
)

type cacheEntry struct {
	msg *miekg_dns.Msg
	expiry time.Time
}

type Server struct {
	cfg           *config.Config
	client        *miekg_dns.Client
	upstreamAddrs []string

	cache map[string]cacheEntry
	mu sync.RWMutex
}

func NewServer(cfg *config.Config) (*Server, error) {
	s := &Server{
		cfg:    cfg,
		client: &miekg_dns.Client{Net: "udp", Timeout: 3 * time.Second},
		cache: make(map[string]cacheEntry),
	}

	if len(cfg.Upstream) > 0 {
		s.upstreamAddrs = cfg.Upstream
	} else {
		s.upstreamAddrs = []string{"8.8.8.8:53", "1.1.1.1:53"}
	}

	s.upstreamAddrs = sanitizeUpstreams(cfg.Listen, s.upstreamAddrs)
	return s, nil
}

func (s *Server) ServeDNS(w miekg_dns.ResponseWriter, r *miekg_dns.Msg) {
	if len(r.Question) != 1 {
		s.forward(w, r)
		return
	}
	q := r.Question[0]

	name := strings.ToLower(miekg_dns.Fqdn(q.Name))

	switch q.Qtype {
	case miekg_dns.TypeA:
		ip, ok := s.lookupA(name)
		if ok {
			msg := new(miekg_dns.Msg)
			msg.SetReply(r)
			msg.Authoritative = true

			rr := &miekg_dns.A{
				Hdr: miekg_dns.RR_Header{
					Name:   name,
					Rrtype: miekg_dns.TypeA,
					Class:  miekg_dns.ClassINET,
					Ttl:    s.cfg.TTL,
				},
				A: net.ParseIP(ip),
			}

			msg.Answer = append(msg.Answer, rr)
			_ = w.WriteMsg(msg)
			return
		}

		s.forward(w, r)

	default:
		s.forward(w, r)
	}
}

func (s *Server) lookupA(name string) (string, bool) {
	// Ищем имя в конфиге, поддерживает и с точкой и без на конце
	n := strings.TrimSuffix(name, ".")
	if ip, ok := s.cfg.Records[n]; ok {
		return ip, true
	}
	if ip, ok := s.cfg.Records[name]; ok {
		return ip, true
	}
	return "", false
}

func (s *Server) cacheKey(r *miekg_dns.Msg) string {
	if len(r.Question) == 0 {
		return ""
	}
	q := r.Question[0]
	return strings.ToLower(miekg_dns.Fqdn(q.Name)) + ":" + miekg_dns.TypeToString[q.Qtype]
}

func (s *Server) forward(w miekg_dns.ResponseWriter, r *miekg_dns.Msg) {
	key := s.cacheKey(r)

	s.mu.RLock()
	entry, ok := s.cache[key]
	s.mu.RUnlock()

	if ok && time.Now().Before(entry.expiry) {
		log.Printf("used cache: %s", key)
		cached := entry.msg.Copy()
		cached.SetReply(r)
		_ = w.WriteMsg(cached)
		return
	}

	// Ищем в апстрим
	for _, ns := range s.upstreamAddrs {
		resp, _, err := s.client.Exchange(r, ns)
		if err == nil && resp != nil {
			// Кешируем
			s.mu.Lock()
			s.cache[key] = cacheEntry{
				msg: resp.Copy(),
				expiry: time.Now().Add(time.Duration(s.cfg.TTL) * time.Second),
			}
			s.mu.Unlock()

			_ = w.WriteMsg(resp)
			return
		}
	}

	// SERVFAIL
	m := new(miekg_dns.Msg)
	m.SetReply(r)
	m.Rcode = miekg_dns.RcodeServerFailure
	_ = w.WriteMsg(m)
}

func sanitizeUpstreams(listen string, ns []string) []string {
	/// Защита от петли
	listenHost, _, _ := net.SplitHostPort(listen)
	if listenHost == "" {
		listenHost = "0.0.0.0"
	}

	bad := map[string]bool{
		"127.0.0.1": true, "::1": true, "127.0.0.53": true,
		listenHost: true,
	}

	var out []string
	for _, s := range ns {
		host, port, _ := net.SplitHostPort(s)
		if host == "" {
			host = s
			port = "53"
		}
		if !bad[host] {
			out = append(out, net.JoinHostPort(host, port))
		}
	}
	return out
}

func (s *Server) startCacheCleaner(ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <- ctx.Done():
			return

		case <- ticker.C:
		now := time.Now()
		s.mu.Lock()
			for k, v := range s.cache {
				if now.After(v.expiry) {
					delete(s.cache, k)
				}
			}
		s.mu.Unlock()
		}
	}
}

func (s *Server) Run(ctx context.Context) error {
	go s.startCacheCleaner(ctx)

	udp := &miekg_dns.Server{Addr: s.cfg.Listen, Net: "udp", Handler: miekg_dns.HandlerFunc(s.ServeDNS)}
	tcp := &miekg_dns.Server{Addr: s.cfg.Listen, Net: "tcp", Handler: miekg_dns.HandlerFunc(s.ServeDNS)}

	errCh := make(chan error, 2)

	go func() { errCh <- udp.ListenAndServe() }()
	go func() { errCh <- tcp.ListenAndServe() }()

	go func() {
		<-ctx.Done()
		_ = udp.Shutdown()
		_ = tcp.Shutdown()
	}()

	return <-errCh
}
