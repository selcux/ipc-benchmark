package tcp

import "testing"

func TestPingPong(t *testing.T) {
	tcp := NewTcp(1512, 2)
	_, err := tcp.PingPong()
	if err != nil {
		t.Error(err)
	}
}
