package security

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func BenchmarkParseChannelWithOptions(b *testing.B) {
	in := "xm54Sj0srWlSEctra-yU6ZA6Z2e6pp7c/a/roman/is/da/best/?opt1=true&opt2=false"

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ParseChannel([]byte(in))
	}
}

func BenchmarkParseChannelStatic(b *testing.B) {
	in := "xm54Sj0srWlSEctra-yU6ZA6Z2e6pp7c/a/roman/is/da/best/"

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ParseChannel([]byte(in))
	}
}

func TestParseChannel(t *testing.T) {
	tests := []struct {
		k  string
		ch string
		o  []string
		t  uint8
	}{
		{k: "emitter", ch: "a/", t: ChannelStatic},
		{k: "emitter", ch: "a/b/c/", t: ChannelStatic},
		{k: "emitter", ch: "test-channel/", t: ChannelStatic},
		{k: "emitter", ch: "test-channel/+/and-more/", t: ChannelWildcard},
		{k: "emitter", ch: "a/-/x/", t: ChannelStatic},
		{k: "emitter", ch: "a/b/c/d/", t: ChannelStatic},
		{k: "emitter", ch: "a/b/c/+/", t: ChannelWildcard},
		{k: "emitter", ch: "a/+/c/+/", t: ChannelWildcard},
		{k: "emitter", ch: "b/+/", t: ChannelWildcard},
		{k: "0TJnt4yZPL73zt35h1UTIFsYBLetyD_g", ch: "emitter/", o: []string{"test=true", "something=7"}, t: ChannelStatic},
		{k: "emitter", ch: "a/b/c/d/", o: []string{"test=true", "something=7"}, t: ChannelStatic},

		// Invalid channels
		{t: ChannelInvalid},
		{k: "emitter", ch: "a/@/x/", t: ChannelInvalid},
		{k: "emitter", ch: "a", t: ChannelInvalid},
		{k: "emitter", ch: "a/b/c", t: ChannelInvalid},
		{k: "emitter", ch: "a//b/", t: ChannelInvalid},
		{k: "emitter", ch: "a//////b/c", t: ChannelInvalid},
		{k: "emitter", ch: "*", t: ChannelInvalid},
		{k: "emitter", ch: "+", t: ChannelInvalid},
		{k: "emitter", ch: "a/+", t: ChannelInvalid},
		{k: "emitter", ch: "b/+", t: ChannelInvalid},
		{k: "emitter", ch: "b/+a/", t: ChannelInvalid},
		{k: "emitter", ch: "", t: ChannelInvalid},
		{k: "emitter", ch: "/", t: ChannelInvalid},
		{k: "emitter", ch: "//", t: ChannelInvalid},
		{k: "emitter", ch: "a//", t: ChannelInvalid},
		{k: "emitter", ch: "a/b/c/d/", o: []string{"test=true", "something=7", "more=_"}, t: ChannelInvalid},
		{k: "emitter", ch: "a/b/c/d/", o: []string{"test==true"}, t: ChannelInvalid},
		{k: "emitter", ch: "a/", o: []string{"=true"}, t: ChannelInvalid},
		{k: "emitter", ch: "a/", o: []string{"test="}, t: ChannelInvalid},
		//		{k: "emitter", ch: "a/b/c/d", o: []string{"test=="}, err: true},
	}

	for _, tc := range tests {
		// First we need to build the input to parse
		in := tc.k + "/" + tc.ch
		if len(tc.o) > 0 {
			in += "?"
			in += strings.Join(tc.o, "&")
		}

		// Parse the channel now
		out := ParseChannel([]byte(in))
		assert.Equal(t, tc.t, out.ChannelType, "input: "+in)
		if tc.t != ChannelInvalid && out.ChannelType != ChannelInvalid {

			// Make sure this always ends with a trailing slash
			if !strings.HasSuffix(tc.ch, "/") {
				tc.ch += "/"
			}

			//assert.Equal(t, ChannelStatic, out.Type)
			assert.Equal(t, tc.k, string(out.Key), "input: "+in)
			assert.Equal(t, tc.ch, string(out.Channel), "input: "+in)

			// Check the options
			for _, opt := range tc.o {
				target := strings.Split(opt, "=")[0]

				found := false
				for _, kvp := range out.Options {
					if kvp.Key == target {
						found = true
						assert.Equal(t, strings.Split(opt, "=")[1], kvp.Value)
					}
				}

				assert.Equal(t, true, found, "unable to find key = "+target)
			}
		}
	}
}
