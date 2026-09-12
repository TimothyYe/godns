package lib

import (
	"net"
	"testing"
	"time"
)

func addrList(cidrs ...string) []net.Addr {
	var out []net.Addr
	for _, c := range cidrs {
		ip, n, err := net.ParseCIDR(c)
		if err != nil {
			panic(err)
		}
		out = append(out, &net.IPNet{IP: ip, Mask: n.Mask})
	}
	return out
}

func TestPickLocalAddr(t *testing.T) {
	addrs := addrList(
		"127.0.0.1/8",      // loopback, never picked
		"fe80::1/64",       // link-local, never picked
		"169.254.10.10/16", // IPv4 link-local, never picked
		"2001:db8::10/64",  // global IPv6
		"192.168.1.10/24",  // private IPv4, allowed for binding
	)

	cases := []struct {
		network string
		want    string
	}{
		{"tcp4", "192.168.1.10"},
		{"tcp6", "2001:db8::10"},
		{"tcp", "192.168.1.10"}, // unconstrained prefers IPv4
	}
	for _, c := range cases {
		got := pickLocalAddr(addrs, c.network)
		if got == nil || got.String() != c.want {
			t.Errorf("pickLocalAddr(%s) = %v, want %s", c.network, got, c.want)
		}
	}

	// Unconstrained falls back to IPv6 when no IPv4 is present.
	v6only := addrList("fe80::1/64", "2001:db8::10/64")
	if got := pickLocalAddr(v6only, "tcp"); got == nil || got.String() != "2001:db8::10" {
		t.Errorf("pickLocalAddr(tcp, v6 only) = %v, want 2001:db8::10", got)
	}

	// Nothing usable.
	if got := pickLocalAddr(addrList("127.0.0.1/8", "fe80::1/64"), "tcp4"); got != nil {
		t.Errorf("pickLocalAddr on loopback/link-local = %v, want nil", got)
	}
	if got := pickLocalAddr(addrList("192.168.1.10/24"), "tcp6"); got != nil {
		t.Errorf("pickLocalAddr(tcp6) with only IPv4 = %v, want nil", got)
	}
}

func TestInterfaceDialerErrors(t *testing.T) {
	if _, err := interfaceDialer("godns-no-such-if0", "tcp4", time.Second); err == nil {
		t.Error("expected error for unknown interface")
	}

	// Loopback exists everywhere but only carries non-global addresses.
	lo := loopbackInterface(t)
	if _, err := interfaceDialer(lo.Name, "tcp4", time.Second); err == nil {
		t.Errorf("expected error binding to loopback %s, which has no global-unicast address", lo.Name)
	}
}

// TestBindToInterfaceLoopback exercises the platform socket option by
// pinning a dial to the loopback device and connecting to a local listener.
func TestBindToInterfaceLoopback(t *testing.T) {
	if !deviceBindSupported {
		t.Skip("device binding not supported on this platform")
	}

	ln, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	lo := loopbackInterface(t)
	d := &net.Dialer{Timeout: 2 * time.Second, Control: bindToInterface(lo)}
	conn, err := d.Dial("tcp4", ln.Addr().String())
	if err != nil {
		t.Fatalf("dial bound to %s failed: %v", lo.Name, err)
	}
	conn.Close()
}

func loopbackInterface(t *testing.T) *net.Interface {
	t.Helper()
	ifaces, err := net.Interfaces()
	if err != nil {
		t.Fatal(err)
	}
	for i := range ifaces {
		if ifaces[i].Flags&net.FlagLoopback != 0 {
			return &ifaces[i]
		}
	}
	t.Skip("no loopback interface found")
	return nil
}
