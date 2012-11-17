// test package 'pkg1'
package pkg1

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
	self.MsgDebug("== execute == [ctx:%v]\n", ctx.Idx())
	detstore := self.DetStore(ctx)
	self.MsgDebug("det-store: %v\n", detstore)
	return kernel.StatusCode(0)
}

func (self *alg1) Finalize() kernel.Error {
	self.MsgInfo("== finalize ==\n")
	return kernel.StatusCode(0)
}

// --- alg2 ---

type alg2 struct {
	kernel.Algorithm
}

func (self *alg2) Initialize() kernel.Error {
	self.MsgInfo("~~ initialize ~~\n")
	if !self.Algorithm.Initialize().IsSuccess() {
		self.MsgError("could not initialize base-class\n")
		return kernel.StatusCode(1)
	}
	return kernel.StatusCode(0)
}

func (self *alg2) Execute(ctx kernel.IEvtCtx) kernel.Error {
	self.MsgDebug("~~ execute ~~ [ctx:%v]\n", ctx.Idx())
	return kernel.StatusCode(0)
}

func (self *alg2) Finalize() kernel.Error {
	self.MsgInfo("~~ finalize ~~\n")
	return kernel.StatusCode(0)
}

// --- svc1 ---
type svc1 struct {
	kernel.Service
}

func (self *svc1) InitializeSvc() kernel.Error {
	self.MsgInfo("~~ initialize ~~\n")
	if !self.Service.InitializeSvc().IsSuccess() {
		self.MsgError("could not initialize base-class\n")
		return kernel.StatusCode(1)
	}
	return kernel.StatusCode(0)
}

func (self *svc1) FinalizeSvc() kernel.Error {
	self.MsgInfo("~~ finalize ~~\n")
	return kernel.StatusCode(0)
}

// --- tool1 ---
type tool1 struct {
	kernel.AlgTool
}

func (self *tool1) InitializeTool() kernel.Error {
	self.MsgInfo("~~ initialize ~~\n")
	if !self.AlgTool.InitializeTool().IsSuccess() {
		self.MsgError("could not initialize base-class\n")
		return kernel.StatusCode(1)
	}
	return kernel.StatusCode(0)
}

func (self *tool1) FinalizeTool() kernel.Error {
	self.MsgInfo("~~ finalize ~~\n")
	return kernel.StatusCode(0)
}

// check matching interfaces
var _ kernel.IComponent = (*alg1)(nil)
var _ kernel.IAlgorithm = (*alg1)(nil)
//var _ kernel.Algorithm = (*alg1)(nil)

var _ kernel.IComponent = (*alg2)(nil)
var _ kernel.IAlgorithm = (*alg2)(nil)

var _ kernel.IComponent = (*svc1)(nil)
var _ kernel.IService = (*svc1)(nil)

var _ kernel.IComponent = (*tool1)(nil)
var _ kernel.IAlgTool = (*tool1)(nil)

// --- factory function ---
func New(t,n string) kernel.IComponent {
	switch t {
	case "alg1":
		self := &alg1{}
		_ = kernel.NewAlg(&self.Algorithm,t,n)
		kernel.RegisterComp(self)
		return self
	case "alg2":
		self := &alg2{}
		_ = kernel.NewAlg(&self.Algorithm,t,n)
		kernel.RegisterComp(self)
		return self
	case "svc1":
		self := &svc1{}
		_ = kernel.NewSvc(&self.Service,t,n)
		kernel.RegisterComp(self)
		return self
	case "tool1":
		self := &tool1{}
		_ = kernel.NewTool(&self.AlgTool,t,n, nil)
		kernel.RegisterComp(self)
		return self
	default:
		err := "no such type ["+t+"]"
		panic(err)
	}
	return nil
}
/* EOF */
