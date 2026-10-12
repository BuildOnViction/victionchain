package tracers

import "testing"

// builtinTracers lists the tracers bundled under internal/tracers. Keeping an
// explicit copy here catches accidental additions or removals of trusted code.
var builtinTracers = []string{
	"4byteTracer",
	"bigramTracer",
	"callTracer",
	"evmdisTracer",
	"noopTracer",
	"opcountTracer",
	"prestateTracer",
	"trigramTracer",
	"unigramTracer",
}

func TestIsBuiltinAcceptsBundledTracers(t *testing.T) {
	for _, name := range builtinTracers {
		if !IsBuiltin(name) {
			t.Errorf("built-in tracer %q rejected", name)
		}
	}
	if len(all) != len(builtinTracers) {
		t.Errorf("bundled tracer count mismatch: have %d, want %d", len(all), len(builtinTracers))
	}
}

func TestIsBuiltinRejectsEverythingElse(t *testing.T) {
	rejected := []string{
		"",
		"unknownTracer",
		"CallTracer",
		"calltracer",
		" callTracer",
		"callTracer ",
		"callTracer\n",
		"callTracerX",
		"callTracer;Duktape.Pointer(this)",
		"{step:function(){},fault:function(){},result:function(){return Duktape.version}}",
		`{step:function(){},fault:function(){},result:function(){return Duktape["Pointer"](this)}}`,
		"{step:function(){},fault:function(){},result:function(){return Duktape[String.fromCharCode(80,111,105,110,116,101,114)](this)}}",
	}
	for _, name := range rejected {
		if IsBuiltin(name) {
			t.Errorf("non built-in tracer %q accepted", name)
		}
	}
}
