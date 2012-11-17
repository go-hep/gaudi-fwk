// test package 'pkg2'
package pkg2

import "time"

import "github.com/sbinet/go-gaudi/pkg/gaudi/kernel"

// --- alg1 ---

type alg1 struct {
	kernel.Algorithm
}

func (self *alg1) Initialize() kernel.Error {
	self.MsgInfo("== initialize ==\n")
	if !self.Algorithm.Initialize().IsSuccess() {
		self.MsgError("could not initialize base-class\n")
		return kernel.StatusCode(1)
	}
	return kernel.StatusCode(0)
}

func (self *alg1) Execute(ctx kernel.IEvtCtx) kernel.Error {
	self.MsgDebug("== execute == [%v]\n", ctx.Idx())
	return kernel.StatusCode(0)
}

func (self *alg1) Finalize() kernel.Error {
	self.MsgInfo("== finalize ==\n")
	return kernel.StatusCode(0)
}

// --- svc2 ---
type svc2 struct {
	kernel.Service
}

func (self *svc2) InitializeSvc() kernel.Error {
	self.MsgInfo("~~ initialize ~~\n")
	if !self.Service.InitializeSvc().IsSuccess() {
		self.MsgError("could not initialize base-class\n")
		return kernel.StatusCode(1)
	}
	return kernel.StatusCode(0)
}

func (self *svc2) FinalizeSvc() kernel.Error {
	self.MsgInfo("~~ finalize ~~\n")
	return kernel.StatusCode(0)
}

type simple_counter struct {
	Cnt int
}

// --- alg adder ---
type alg_adder struct {
	kernel.Algorithm
	val     float64
	cnt_key string
}

func (self *alg_adder) Initialize() kernel.Error {
	self.MsgInfo("== initialize ==\n")
	if !self.Algorithm.Initialize().IsSuccess() {
		self.MsgError("could not initialize base-class\n")
		return kernel.StatusCode(1)
	}

	self.MsgInfo("--> val: %v\n", self.val)
	self.val = self.GetProperty("Val").(float64)
	self.MsgInfo("--> val: %v\n", self.val)
	self.cnt_key = self.GetProperty("SimpleCounter").(string)
	self.MsgInfo("--> cnt: %v\n", self.cnt_key)
	return kernel.StatusCode(0)
}

func (self *alg_adder) Execute(ctx kernel.IEvtCtx) kernel.Error {
	self.MsgDebug("== execute == [%v]\n", ctx.Idx())

	njets := 1 + ctx.Idx()
	val := self.val + 1
	store := self.EvtStore(ctx)

	if store.Has("njets") {
		njets += store.Get("njets").(int)
	}
	store.Put("njets", njets)

	if store.Has("ptjets") {
		val += store.Get("ptjets").(float64)
	}
	store.Put("ptjets", val)

	cnt := &simple_counter{0}
	if store.Has(self.cnt_key) {
		cnt = store.Get(self.cnt_key).(*simple_counter)
	}
	cnt.Cnt += 1
	store.Put(self.cnt_key, cnt)
	return kernel.StatusCode(0)
}

func (self *alg_adder) Finalize() kernel.Error {
	self.MsgInfo("== finalize ==\n")
	return kernel.StatusCode(0)
}

// --- alg_dumper ---
type alg_dumper struct {
	kernel.Algorithm
	njets_key  string
	ptjets_key string
	cnt_key    string
	cnt_val    int
}

func (self *alg_dumper) Initialize() kernel.Error {
	self.MsgInfo("== initialize ==\n")
	if !self.Algorithm.Initialize().IsSuccess() {
		self.MsgError("could not initialize base-class\n")
		return kernel.StatusCode(1)
	}
	self.njets_key = self.GetProperty("NbrJets").(string)
	self.ptjets_key = self.GetProperty("PtJets").(string)
	self.cnt_key = self.GetProperty("SimpleCounter").(string)
	self.cnt_val = self.GetProperty("ExpectedValue").(int)
	return kernel.StatusCode(0)
}

func (self *alg_dumper) Execute(ctx kernel.IEvtCtx) kernel.Error {
	self.MsgDebug("== execute == [%v]\n", ctx.Idx())

	store := self.EvtStore(ctx)
	njets := store.Get(self.njets_key).(int)
	ptjets := store.Get(self.ptjets_key).(float64)
	cnt := store.Get(self.cnt_key).(*simple_counter)

	if self.cnt_val == cnt.Cnt {
		self.MsgDebug("[ctx:%03v] njets: %03v ptjets: %8.3v [OK]\n",
			ctx.Idx(), njets, ptjets)
	} else {
		self.MsgError("[ctx:%03v] njets: %03v ptjets: %8.3v (%v|%v) [ERR]\n",
			ctx.Idx(), njets, ptjets, self.cnt_val, cnt.Cnt)
		panic("race condition detected in evt-store")
	}

	// simulate a cpu-burning algorithm
	time.Sleep(3 * 10e6) // 3ms
	return kernel.StatusCode(0)
}

func (self *alg_dumper) Finalize() kernel.Error {
	self.MsgInfo("== finalize ==\n")
	return kernel.StatusCode(0)
}

// --- factory function ---
func New(t, n string) kernel.IComponent {
	switch t {
	case "alg1":
		self := &alg1{}
		_ = kernel.NewAlg(&self.Algorithm, t, n)
		kernel.RegisterComp(self)
		return self
	case "svc2":
		self := &svc2{}
		_ = kernel.NewSvc(&self.Service, t, n)
		kernel.RegisterComp(self)
		return self
	case "alg_adder":
		self := &alg_adder{}
		_ = kernel.NewAlg(&self.Algorithm, t, n)
		kernel.RegisterComp(self)

		self.DeclareProperty("Val", -99.)
		self.DeclareProperty("SimpleCounter", "cnt")
		return self

	case "alg_dumper":
		self := &alg_dumper{}
		_ = kernel.NewAlg(&self.Algorithm, t, n)
		kernel.RegisterComp(self)

		self.DeclareProperty("SimpleCounter", "cnt")
		self.DeclareProperty("ExpectedValue", -1)
		self.DeclareProperty("NbrJets", "njets")
		self.DeclareProperty("PtJets", "ptjets")
		return self

	default:
		err := "no such type [" + t + "]"
		panic(err)
	}
	return nil
}
/* EOF */
