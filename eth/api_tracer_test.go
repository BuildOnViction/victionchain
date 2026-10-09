package eth

import (
	"context"
	"testing"

	"github.com/tomochain/tomochain/core/types"
	"github.com/tomochain/tomochain/core/vm"
)

func TestValidateTracerConfig(t *testing.T) {
	str := func(s string) *string { return &s }

	tests := []struct {
		name   string
		config *TraceConfig
		err    error
	}{
		{"nil config", nil, nil},
		{"nil tracer", &TraceConfig{}, nil},
		{"call tracer", &TraceConfig{Tracer: str("callTracer")}, nil},
		{"prestate tracer", &TraceConfig{Tracer: str("prestateTracer")}, nil},
		{"empty name", &TraceConfig{Tracer: str("")}, errCustomTracer},
		{"wrong case", &TraceConfig{Tracer: str("CallTracer")}, errCustomTracer},
		{"whitespace", &TraceConfig{Tracer: str(" callTracer")}, errCustomTracer},
		{"suffix", &TraceConfig{Tracer: str("callTracer;Duktape.Pointer(this)")}, errCustomTracer},
		{"raw js", &TraceConfig{Tracer: str("{step:function(){},fault:function(){},result:function(){return Duktape.version}}")}, errCustomTracer},
		{"bracket access", &TraceConfig{Tracer: str(`{step:function(){},fault:function(){},result:function(){return Duktape["Pointer"](this)}}`)}, errCustomTracer},
	}
	for _, tt := range tests {
		if err := validateTracerConfig(tt.config); err != tt.err {
			t.Errorf("%s: have %v, want %v", tt.name, err, tt.err)
		}
	}
}

// TestTraceEntryPointsRejectCustomTracer ensures the API entry points reject a
// custom tracer before touching the backend. The API has no backend wired in,
// so reaching any chain access would panic.
func TestTraceEntryPointsRejectCustomTracer(t *testing.T) {
	js := "{step:function(){},fault:function(){},result:function(){return Duktape.version}}"
	config := &TraceConfig{Tracer: &js}
	api := &PrivateDebugAPI{}

	if _, err := api.TraceChain(context.Background(), 0, 1, config); err != errCustomTracer {
		t.Errorf("TraceChain: have %v, want %v", err, errCustomTracer)
	}
	block := types.NewBlockWithHeader(&types.Header{})
	if _, err := api.traceBlock(context.Background(), block, config); err != errCustomTracer {
		t.Errorf("traceBlock: have %v, want %v", err, errCustomTracer)
	}
	if _, err := api.traceTx(context.Background(), nil, vm.Context{}, nil, config); err != errCustomTracer {
		t.Errorf("traceTx: have %v, want %v", err, errCustomTracer)
	}
}
