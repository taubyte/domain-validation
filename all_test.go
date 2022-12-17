package lib

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/ipfs/go-cid"
	"github.com/miekg/dns"

	_ "embed"
)

//go:embed fixtures/test_wrong_public.pem
var fakePublicKeyData []byte

//go:embed fixtures/test_good_private.key
var testPrivateKeyData []byte

//go:embed fixtures/test_good_public.pem
var testPublicKeyData []byte

var (
	fqdn           = "example.com"
	project        cid.Cid
	projectString  = "QmNTj7BY3YNySGaUrfuWjkomCME8fTcynvgHiNdX1hyLME"
	port           = 9090
	defaultUdpPort = fmt.Sprintf("127.0.0.1:%d", port)
)

func init() {
	var err error
	project, err = cid.Decode(projectString)
	if err != nil {
		panic(fmt.Errorf("Failed to decode cid with: %w", err))
	}

	UseResolver(createResolver("udp", defaultUdpPort))
}

func TestE2E(t *testing.T) {
	claim, err := New(FQDN(fqdn), Project(project), PrivateKey(testPrivateKeyData), PublicKey(testPublicKeyData))
	if err != nil {
		t.Error(err)
		return
	}

	token, err := claim.Sign()
	if err != nil {
		t.Error(err)
		return
	}

	claimToCheck, err := FromToken(token, PublicKey(testPublicKeyData))
	if err != nil {
		t.Error(err)
		return
	}

	if claimToCheck.Address != claim.Address {
		t.Error("Decoded token contains wrong data", claimToCheck, claim)
		return
	}

	corruptedToken := token + "CoRruPteD"
	claimToCheck, err = FromToken(corruptedToken, PublicKey(testPublicKeyData))
	if err == nil {
		t.Error("did not fail for a corrupted token")
		return
	}

	_, err = FromToken(token, PublicKey(fakePublicKeyData))
	if err == nil {
		t.Error("Should have failed to validate signature with wrong key")
		return
	}

}

func TestDNSVerify(t *testing.T) {
	ctx, ctxC := context.WithCancel(context.Background())
	defer ctxC()

	projectsDomains := map[string]string{
		project.String(): fqdn,
	}
	zoneTXT := make(map[string]string)

	for prj, dom := range projectsDomains {
		prjCid, err := cid.Decode(prj)
		if err != nil {
			t.Error("Failed to decode cid with:", err)
			continue
		}
		claims, err := New(FQDN(dom), Project(prjCid), PrivateKey(testPrivateKeyData), PublicKey(testPublicKeyData))
		if err != nil {
			t.Error(err)
			continue
		}
		token, err := claims.Sign()
		if err != nil {
			t.Error(err)
			continue
		}

		entry := fmt.Sprintf("%s.%s.", claims.project[:8], claims.fqdn)
		zoneTXT[entry] = string(token)
	}

	dnsServer := &dns.Server{
		Addr: ":" + strconv.Itoa(port), Net: "udp",
		Handler: dns.HandlerFunc(func(w dns.ResponseWriter, r *dns.Msg) {
			msg := dns.Msg{}
			msg.SetReply(r)

			msg.Authoritative = true

			switch r.Question[0].Qtype {
			case dns.TypeTXT:
				reqDomain := r.Question[0].Name
				token, ok := zoneTXT[reqDomain]
				if ok == true {
					msg.Answer = append(msg.Answer, &dns.TXT{
						Hdr: dns.RR_Header{Name: reqDomain, Rrtype: dns.TypeTXT, Class: dns.ClassINET, Ttl: 60},
						Txt: []string{token},
					})
				}
			}
			err := w.WriteMsg(&msg)
			if err != nil {
				fmt.Println("Failed write msg with error ", err)
				return
			}
		}),
	}

	go func() error {
		if err := dnsServer.ListenAndServe(); err != nil {
			return fmt.Errorf("failed starting UPD Server error: %v", err)
		}
		return nil
	}()
	_, err := resolver.LookupTXT(ctx, projectString[:8]+"."+fqdn)
	if err != nil {
		t.Error("Failed tcp lookupTXT error: ", err)
		return
	}

	err = FromDNS(ctx, &project, fqdn, PublicKey(testPublicKeyData))
	if err != nil {
		t.Error(err)
	}

	go func() {
		select {
		case <-ctx.Done():
			dnsServer.Shutdown()
		}
	}()

}

func createResolver(netType string, port string) *net.Resolver {
	newResolver := &net.Resolver{
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: 3 * time.Second,
			}
			return d.DialContext(ctx, netType, port)
		},
	}
	return newResolver
}
