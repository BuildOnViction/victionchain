// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package tracers

import (
	"encoding/json"
	"errors"
	"github.com/tomochain/tomochain/common"
	"github.com/tomochain/tomochain/core/state"
	"github.com/tomochain/tomochain/core/vm"
	"github.com/tomochain/tomochain/params"
	"math/big"
	"strings"
	"testing"
)

type account struct{}

func (account) SubBalance(amount *big.Int)                          {}
func (account) AddBalance(amount *big.Int)                          {}
func (account) SetAddress(common.Address)                           {}
func (account) Value() *big.Int                                     { return nil }
func (account) SetBalance(*big.Int)                                 {}
func (account) SetNonce(uint64)                                     {}
func (account) Balance() *big.Int                                   { return nil }
func (account) Address() common.Address                             { return common.Address{} }
func (account) ReturnGas(*big.Int)                                  {}
func (account) SetCode(common.Hash, []byte)                         {}
func (account) ForEachStorage(cb func(key, value common.Hash) bool) {}

type dummyStatedb struct {
	state.StateDB
}

func (*dummyStatedb) GetRefund() uint64 { return 1337 }

func runTrace(tracer Tracer) (json.RawMessage, error) {
	var (
		env             = vm.NewEVM(vm.Context{BlockNumber: big.NewInt(1)}, &dummyStatedb{}, nil, params.TestChainConfig, vm.Config{Debug: true, Tracer: tracer})
		gasLimit uint64 = 31000
		startGas uint64 = 10000
		value           = big.NewInt(0)
		contract        = vm.NewContract(account{}, account{}, value, startGas)
	)
	contract.Code = []byte{byte(vm.PUSH1), 0x1, byte(vm.PUSH1), 0x1, 0x0}

	tracer.CaptureTxStart(gasLimit)
	tracer.CaptureStart(env, contract.Caller(), contract.Address(), false, []byte{}, startGas, value)
	ret, err := env.Interpreter().Run(contract, []byte{}, false)
	tracer.CaptureEnd(ret, startGas-contract.Gas, err)
	// Rest gas assumes no refund
	tracer.CaptureTxEnd(contract.Gas)
	if err != nil {
		return nil, err
	}
	return tracer.GetResult()
}

func TestHaltBetweenSteps(t *testing.T) {
	tracer, err := New("{step: function() {}, fault: function() {}, result: function() { return null; }}")
	if err != nil {
		t.Fatal(err)
	}
	env := vm.NewEVM(vm.Context{BlockNumber: big.NewInt(1), GasPrice: big.NewInt(1)}, &dummyStatedb{}, nil, params.TestChainConfig, vm.Config{Tracer: tracer})
	scope := &vm.CallCtx{
		Contract: vm.NewContract(&account{}, &account{}, big.NewInt(0), 0),
	}
	tracer.CaptureStart(env, common.Address{}, common.Address{}, false, []byte{}, 0, big.NewInt(0))
	tracer.CaptureState(0, 0, 0, 0, scope, nil, 0, nil)
	timeout := errors.New("stahp")
	tracer.Stop(timeout)
	tracer.CaptureState(0, 0, 0, 0, scope, nil, 0, nil)

	if _, err := tracer.GetResult(); !strings.Contains(err.Error(), timeout.Error()) {
		t.Errorf("Expected timeout error, got %v", err)
	}
}

type tracerCtor = func(string) (Tracer, error)

func TestDuktapeTracer(t *testing.T) {
	testTracer(t, New)
}

func testTracer(t *testing.T, newTracer tracerCtor) {
	execTracer := func(code string) ([]byte, string) {
		t.Helper()
		tracer, err := newTracer(code)
		if err != nil {
			t.Fatal(err)
		}
		ret, err := runTrace(tracer)
		if err != nil {
			return nil, err.Error() // Stringify to allow comparison without nil checks
		}
		return ret, ""
	}
	for i, tt := range []struct {
		code string
		want string
		fail string
	}{
		{ // tests that we don't panic on bad arguments to memory access
			code: "{depths: [], step: function(log) { this.depths.push(log.memory.slice(-1,-2)); }, fault: function() {}, result: function() { return this.depths; }}",
			want: `[{},{},{}]`,
		}, { // tests that we don't panic on bad arguments to stack peeks
			code: "{depths: [], step: function(log) { this.depths.push(log.stack.peek(-1)); }, fault: function() {}, result: function() { return this.depths; }}",
			want: `["0","0","0"]`,
		}, { //  tests that we don't panic on bad arguments to memory getUint
			code: "{ depths: [], step: function(log, db) { this.depths.push(log.memory.getUint(-64));}, fault: function() {}, result: function() { return this.depths; }}",
			want: `["0","0","0"]`,
		}, { // tests some general counting
			code: "{count: 0, step: function() { this.count += 1; }, fault: function() {}, result: function() { return this.count; }}",
			want: `3`,
		}, { // tests that depth is reported correctly
			code: "{depths: [], step: function(log) { this.depths.push(log.stack.length()); }, fault: function() {}, result: function() { return this.depths; }}",
			want: `[0,1,2]`,
		}, { // tests memory length
			code: "{lengths: [], step: function(log) { this.lengths.push(log.memory.length()); }, fault: function() {}, result: function() { return this.lengths; }}",
			want: `[0,0,0]`,
		}, { // tests to-string of opcodes
			code: "{opcodes: [], step: function(log) { this.opcodes.push(log.op.toString()); }, fault: function() {}, result: function() { return this.opcodes; }}",
			want: `["PUSH1","PUSH1","STOP"]`,
		}, { // tests intrinsic gas
			code: "{depths: [], step: function() {}, fault: function() {}, result: function(ctx) { return ctx.gasPrice+'.'+ctx.gasUsed+'.'+ctx.intrinsicGas; }}",
			want: `"100000.6.21000"`,
		}, {
			code: "{res: null, step: function(log) {}, fault: function() {}, result: function() { return toWord('0xffaa') }}",
			want: `{"0":0,"1":0,"2":0,"3":0,"4":0,"5":0,"6":0,"7":0,"8":0,"9":0,"10":0,"11":0,"12":0,"13":0,"14":0,"15":0,"16":0,"17":0,"18":0,"19":0,"20":0,"21":0,"22":0,"23":0,"24":0,"25":0,"26":0,"27":0,"28":0,"29":0,"30":255,"31":170}`,
		}, { // test feeding a buffer back into go
			code: "{res: null, step: function(log) { var address = log.contract.getAddress(); this.res = toAddress(address); }, fault: function() {}, result: function() { return this.res }}",
			want: `{"0":0,"1":0,"2":0,"3":0,"4":0,"5":0,"6":0,"7":0,"8":0,"9":0,"10":0,"11":0,"12":0,"13":0,"14":0,"15":0,"16":0,"17":0,"18":0,"19":0}`,
		}, {
			code: "{res: null, step: function(log) { var address = '0x0000000000000000000000000000000000000000'; this.res = toAddress(address); }, fault: function() {}, result: function() { return this.res }}",
			want: `{"0":0,"1":0,"2":0,"3":0,"4":0,"5":0,"6":0,"7":0,"8":0,"9":0,"10":0,"11":0,"12":0,"13":0,"14":0,"15":0,"16":0,"17":0,"18":0,"19":0}`,
		}, {
			code: "{res: null, step: function(log) { var address = Array.prototype.slice.call(log.contract.getAddress()); this.res = toAddress(address); }, fault: function() {}, result: function() { return this.res }}",
			want: `{"0":0,"1":0,"2":0,"3":0,"4":0,"5":0,"6":0,"7":0,"8":0,"9":0,"10":0,"11":0,"12":0,"13":0,"14":0,"15":0,"16":0,"17":0,"18":0,"19":0}`,
		},
	} {
		if have, err := execTracer(tt.code); tt.want != string(have) || tt.fail != err {
			t.Errorf("testcase %d: expected return value to be '%s' got '%s', error to be '%s' got '%s'\n\tcode: %v", i, tt.want, string(have), tt.fail, err, tt.code)
		}
	}
}
