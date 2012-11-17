package evtproc

import "time"
//import "fmt"

import "github.com/sbinet/go-gaudi/pkg/gaudi/kernel"

// --- evt state ---
type evtstate struct {
	idx  int
	sc   kernel.Error
	data kernel.DataStore
}

func new_evtstate(idx int) *evtstate {
	self := &evtstate{}
	self.idx = idx
	self.sc = kernel.StatusCode(0)
	self.data = make(kernel.DataStore)
	return self
}

func (self *evtstate) Idx() int {
	return self.idx
}
func (self *evtstate) Store() kernel.DataStore {
	return self.data
}

// --- evt processor ---
type evtproc struct {
	kernel.Service
	algs []kernel.IAlgorithm
	nworkers int
}

func (self *evtproc) InitializeSvc() kernel.Error {

	sc := self.Service.InitializeSvc()
	if !sc.IsSuccess() {
		return sc
	}
	self.nworkers = self.GetProperty("NbrWorkers").(int)
	//self.nworkers = 2
	self.MsgInfo("n-workers: %v\n", self.nworkers)
	svcloc := kernel.GetSvcLocator()
	if svcloc == nil {
		self.MsgError("could not retrieve ISvcLocator !\n")
		return kernel.StatusCode(1)
	}
	appmgr := svcloc.(kernel.IComponentMgr).GetComp("app-mgr")
	if appmgr == nil {
		self.MsgError("could not retrieve 'app-mgr'\n")
	}
	propmgr := appmgr.(kernel.IProperty)
	alg_names := propmgr.GetProperty("Algs").([]string)
	self.MsgInfo("got alg-names: %v\n", len(alg_names))

	if len(alg_names)>0 {
		comp_mgr := appmgr.(kernel.IComponentMgr)
		self.algs = make([]kernel.IAlgorithm, len(alg_names))
		for i,alg_name := range alg_names {
			ialg,isalg := comp_mgr.GetComp(alg_name).(kernel.IAlgorithm)
			if isalg {
				self.algs[i] = ialg
			}
		}
	}
	self.MsgInfo("got alg-list: %v\n", len(self.algs))

	return kernel.StatusCode(0)
}

func (self *evtproc) ExecuteEvent(ictx kernel.IEvtCtx) kernel.Error {
	ctx := ictx.Idx()
	self.MsgDebug("executing event [%v]... (#algs: %v)\n", ctx, len(self.algs))
	for i,alg := range self.algs {
		self.MsgDebug("-- ctx:%03v --> [%s]...\n", ctx, alg.CompName())
		if !alg.Execute(ictx).IsSuccess() {
			self.MsgError("pb executing alg #%v (%s) for ctx:%v\n",
				i,alg.CompName(), ictx.Idx())
			return kernel.StatusCode(1)
		}
	}
	self.MsgDebug("data: %v\n",ictx.Store())
	return kernel.StatusCode(0)
}

func (self *evtproc) ExecuteRun(evtmax int) kernel.Error {
	self.MsgInfo("execute-run [%v]\n", evtmax)
	sc := self.NextEvent(evtmax)
	return sc
}

func (self *evtproc) NextEvent(evtmax int) kernel.Error {

	if self.nworkers > 1 {
		return self.mp_NextEvent(evtmax)
	}
	return self.seq_NextEvent(evtmax)
}

func (self *evtproc) seq_NextEvent(evtmax int) kernel.Error {
	
	self.MsgInfo("nextEvent[%v]...\n", evtmax)
	for i:=0; i<evtmax; i++ {
		ctx := new_evtstate(i)
		if !self.ExecuteEvent(ctx).IsSuccess() {
			self.MsgError("failed to execute evt idx %03v\n", i)
			return kernel.StatusCode(1)
		}
	}
	return kernel.StatusCode(0)
}
func (self *evtproc) mp_NextEvent(evtmax int) kernel.Error {

	handle := func(evt *evtstate, out_queue chan <- *evtstate) {
		self.MsgInfo("nextEvent[%v]...\n", evt.idx)
		evt.sc = self.ExecuteEvent(evt)
		out_queue <- evt
	}

	serve_evts := func(in_evt_queue <- chan *evtstate, out_evt_queue chan <- *evtstate, quit <- chan bool) {
		for {
			select {
			case ievt := <-in_evt_queue:
				go handle(ievt, out_evt_queue)
			case <-quit:
				//println("quit requested !")
				return
			}
		}
	}

	start_evt_server := func(nworkers int) (in_evt_queue,
		                                    out_evt_queue chan *evtstate,
		                                    quit chan bool) {
		in_evt_queue = make(chan *evtstate, nworkers)
		out_evt_queue = make(chan *evtstate)
		quit = make(chan bool)
		go serve_evts(in_evt_queue, out_evt_queue, quit)
		return in_evt_queue, out_evt_queue, quit
	}

	in_evt_queue, out_evt_queue, quit := start_evt_server(self.nworkers)
	//self.MsgInfo("sending requests...\n")
	for i:=0; i<evtmax; i++ {
		in_evt_queue <- new_evtstate(i)
	}
	//self.MsgInfo("sending requests... [done]\n")
	n_fails := 0
	n_processed := 0
	for evt := range out_evt_queue {
		//self.MsgDebug("out-evt-queue: %v %v\n",evt.idx, evt.sc)
		if !evt.sc.IsSuccess() {
			n_fails++
		}
		n_processed++
		if n_processed == evtmax {
			quit <- true
			close(out_evt_queue)
			//self.MsgDebug("closing evt server...\n")
			break
		}
	}
	if n_fails != 0 {
		return kernel.StatusCode(1)
	}
	return kernel.StatusCode(0)
}

func (self *evtproc) StopRun() kernel.Error {
	self.MsgInfo("stopping run...\n")
	return kernel.StatusCode(0)
}

// ---

func (self *evtproc) test_0() {
	handle := func(queue <-chan int) kernel.Error {
		sc := kernel.StatusCode(0)
		for i := range queue {
			ctx := new_evtstate(i)
			self.MsgInfo("   --> handling [%i]...\n",i)
			sc = self.ExecuteEvent(ctx)
		}
		return sc
	}

	max_in_flight := 4
	serve := func(queue <-chan int, quit <-chan bool) kernel.Error {
		for i := 0; i < max_in_flight; i++ {
			go handle(queue)
		}
		<-quit // wait to be told to exit
		return kernel.StatusCode(0)
	}

	quit := make(chan bool)

	self.MsgInfo("-- filling the event queue...\n")
	queue := make(chan int)
	go func() {
		for i := 0; i < 20; i++ {
			queue <- i
		}
	}()
	self.MsgInfo("-- starting to serve 20 events...\n")
	go serve(queue, quit)
	self.MsgInfo("-- requests sent...\n")
	time.Sleep(2000000000)
	quit <- true
	self.MsgInfo("-- done.\n")
}

// check implementations match interfaces
var _ kernel.IEvtCtx = (*evtstate)(nil)

var _ kernel.IComponent = (*evtproc)(nil)
var _ kernel.IEvtProcessor = (*evtproc)(nil)
var _ kernel.IProperty = (*evtproc)(nil)
var _ kernel.IService = (*evtproc)(nil)

// --- factory function ---
func New(t,n string) kernel.IComponent {
	switch t {
	case "evtproc":
		self := &evtproc{}
	
		//self.properties.props = make(map[string]interface{})
		//self.name = name
		self.algs = []kernel.IAlgorithm{}
		_ = kernel.NewSvc(&self.Service, t, n)
		kernel.RegisterComp(self)
		self.DeclareProperty("NbrWorkers", 2)
		return self

	default:
		err := "no such type ["+t+"]"
		panic(err)
	}
	return nil
}

/* EOF */
