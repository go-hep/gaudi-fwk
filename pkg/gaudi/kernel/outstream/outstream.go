// package containing components to implement gaudi output stream algorithms
package outstream

//import "io"
import "os"
import "encoding/gob"
import "encoding/json"

import "github.com/sbinet/go-gaudi/pkg/gaudi/kernel"

var g_keys = []string{
	"njets",
	"ptjets",
	"cnt",
}

// simple interface to gather all encoders
type iwriter interface {
	Encode(v interface{}) error
}

// --- gob_outstream ---
type gob_outstream struct {
	kernel.Algorithm
	w          *os.File
	enc        *gob.Encoder
	item_names []string
}

func (self *gob_outstream) Initialize() kernel.Error {
	self.MsgDebug("== initialize ==\n")
	if !self.Algorithm.Initialize().IsSuccess() {
		self.MsgError("could not initialize base-class\n")
		return kernel.StatusCode(1)
	}

	self.item_names = self.GetProperty("Items").([]string)
	fname := self.GetProperty("Output").(string)

	self.MsgInfo("output file: [%v]\n", fname)
	self.MsgInfo("items: %v\n", self.item_names)

	w, err := os.Create(fname)
	if err != nil {
		self.MsgError("problem while opening file [%v]: %v\n", fname, err)
		return kernel.StatusCode(1)
	}
	self.w = w
	self.enc = gob.NewEncoder(self.w)
	return kernel.StatusCode(0)
}

func (self *gob_outstream) Execute(ctx kernel.IEvtCtx) kernel.Error {
	self.MsgDebug("== execute ==\n")
	store := self.EvtStore(ctx)
	if store == nil {
		self.MsgError("could not retrieve evt-store\n")
	}
	var err error
	allgood := true
	for _, k := range self.item_names {
		err = self.enc.Encode(store.Get(k))
		if err != nil {
			self.MsgError("error while writing store content at [%v]: %v\n",
				k, err)
			allgood = false
		}
	}

	/*
		hdr_offset := 1
		val := make([]interface{}, len(self.item_names)+hdr_offset)
		val[0] = ctx.Idx()

		for i,k := range self.item_names {
			val[i+hdr_offset] = store.Get(k)
		}
		for idx,v := range val {
			err = self.enc.Encode(val)
			if err != nil {
				self.MsgError("error while encoding data [%v|%v]: %v\n", 
					idx, v, err)
				allgood = false
			}
		}
	*/

	if !allgood {
		return kernel.StatusCode(1)
	}
	return kernel.StatusCode(0)
}

func (self *gob_outstream) Finalize() kernel.Error {
	self.MsgDebug("== finalize ==\n")
	//self.w.Close()

	return kernel.StatusCode(0)
}

// --- json_outstream ---
type json_outstream struct {
	kernel.Algorithm
	item_names []string
	handle     kernel.IOutputStream
}

func (self *json_outstream) Initialize() kernel.Error {
	self.MsgDebug("== initialize ==\n")
	if !self.Algorithm.Initialize().IsSuccess() {
		self.MsgError("could not initialize base-class\n")
		return kernel.StatusCode(1)
	}

	self.item_names = self.GetProperty("Items").([]string)

	fname := self.GetProperty("Output").(string)
	self.MsgInfo("output file: [%v]\n", fname)
	self.MsgInfo("items: %v\n", self.item_names)

	svcloc := kernel.GetSvcLocator()
	if svcloc == nil {
		self.MsgError("could not retrieve service locator !\n")
		return kernel.StatusCode(1)
	}
	svc := svcloc.GetService("outstreamsvc").(kernel.IOutputStreamSvc)
	if svc == nil {
		self.MsgError("could not retrieve [outstreamsvc] !\n")
		return kernel.StatusCode(1)
	}

	self.handle = svc.NewOutputStream(fname)
	if self.handle == nil {
		self.MsgError("could not retrieve json output stream [%s] !\n", fname)
		return kernel.StatusCode(1)
	}

	return kernel.StatusCode(0)
}

func (self *json_outstream) Execute(ctx kernel.IEvtCtx) kernel.Error {
	self.MsgDebug("== execute ==\n")
	store := self.EvtStore(ctx)
	if store == nil {
		self.MsgError("could not retrieve evt-store\n")
	}

	hdr_offset := 1

	val := make([]interface{}, len(self.item_names)+hdr_offset)
	val[0] = ctx.Idx()

	for i, k := range self.item_names {
		val[i+hdr_offset] = store.Get(k)
	}

	return self.handle.Write(val)
}

func (self *json_outstream) Finalize() kernel.Error {
	self.MsgDebug("== finalize ==\n")

	return kernel.StatusCode(0)
}

/// an output stream using JSON as a format
type json_outstream_handle struct {
	svc  kernel.IService
	w    *os.File
	enc  *json.Encoder
	data chan interface{}
	errs chan error
	quit chan bool
}

func (self *json_outstream_handle) Write(data interface{}) kernel.Error {
	/* // FIXME: how to get the error back ??
	err := self.enc.Encode(data)
	if err != nil {
		return kernel.StatusCodeWithErr(1, err)
	}
	*/
	self.data <- data
	select {
	case err := <-self.errs:
		if err != nil {
			msg := self.svc.(kernel.IMessager)
			msg.MsgError("--> write got: %v\n", err)
			return kernel.StatusCodeWithErr(1, err)
		}
	default:
		return kernel.StatusCode(0)
	}
	return kernel.StatusCode(0)
}

func (self *json_outstream_handle) Close() kernel.Error {

	self.quit <- true

	msg := self.svc.(kernel.IMessager)
	msg.MsgDebug("--> closing json-handle [%v]\n", self.w.Name())

	// _,ok = <-self.quit
	// if ok {
	// 	self.quit <- true
	// }
	close(self.quit)
	close(self.errs)
	close(self.data)

	fd := self.w.Fd()
	if fd >= 0 {
		fname := self.w.Name()
		err := self.w.Close()
		if err != nil {
			msg.MsgError("closing fd: [%v] name [%v]. err: %v\n",
				fd, fname, err)
			return kernel.StatusCodeWithErr(1, err)
		}
	}
	return kernel.StatusCode(0)
}

func (self *json_outstream_handle) Name() string {
	if self.w != nil {
		return self.w.Name()
	}
	return "<N/A>"
}

func (self *json_outstream_handle) Fd() int {
	if self.w != nil {
		return int(self.w.Fd())
	}
	return -1
}

// --- json outputstream srv
type json_outstream_svc struct {
	kernel.Service
	streams map[string]kernel.IOutputStream
}

func (self *json_outstream_svc) InitializeSvc() kernel.Error {

	self.MsgDebug("== initialize ==\n")
	if !self.Service.InitializeSvc().IsSuccess() {
		self.MsgError("could not initialize base-class\n")
		return kernel.StatusCode(1)
	}

	return kernel.StatusCode(0)
}

func (self *json_outstream_svc) FinalizeSvc() kernel.Error {

	self.MsgDebug("== finalize ==\n")
	if !self.Service.FinalizeSvc().IsSuccess() {
		self.MsgError("could not finalize base-class\n")
		return kernel.StatusCode(1)
	}

	self.MsgDebug("== closing output streams...\n")
	allgood := true
	for n, stream := range self.streams {
		self.MsgDebug("-- closing [%s]...\n", n)
		sc := stream.Close()
		if !sc.IsSuccess() {
			self.MsgError("problem closing stream [%s]: %v\n", n, sc)
			allgood = false
		}
	}
	if !allgood {
		return kernel.StatusCode(1)
	}

	return kernel.StatusCode(0)
}

func (self *json_outstream_svc) NewOutputStream(n string) kernel.IOutputStream {
	_, ok := self.streams[n]
	if ok {
		self.MsgError("a json output stream with name [%s] "+
			"has already been created !\n", n)
	}

	stream := &json_outstream_handle{}
	stream.svc = self
	self.streams[n] = stream

	make_chans := func(w iwriter) (chan interface{},
		chan error,
		chan bool) {
		in := make(chan interface{})
		errs := make(chan error)
		quit := make(chan bool)
		go func() {
			for {
				select {
				case data := <-in:
					err := w.Encode(data)
					//errs <- err // FIXME: how to pass back the errors ?!?
					if err != nil {
						println("** error **", err)
					}
				case <-quit:
					return
				}

			}
		}()
		return in, errs, quit
	}

	w, err := os.Create(n)
	if err != nil {
		self.MsgError("problem opening stream [%v]: %v\n", n, err)
		return nil
	}

	stream.w = w
	stream.enc = json.NewEncoder(stream.w)
	stream.data, stream.errs, stream.quit = make_chans(stream.enc)

	return stream
}

// check implementations match interfaces
var _ kernel.IComponent = (*gob_outstream)(nil)
var _ kernel.IAlgorithm = (*gob_outstream)(nil)

var _ kernel.IComponent = (*json_outstream)(nil)
var _ kernel.IAlgorithm = (*json_outstream)(nil)

var _ kernel.IComponent = (*json_outstream_svc)(nil)
var _ kernel.IService = (*json_outstream_svc)(nil)
var _ kernel.IOutputStreamSvc = (*json_outstream_svc)(nil)

var _ kernel.IOutputStream = (*json_outstream_handle)(nil)

// --- factory ---
func New(t, n string) kernel.IComponent {
	switch t {
	case "gob_outstream":
		self := &gob_outstream{}
		kernel.NewAlg(&self.Algorithm, t, n)
		kernel.RegisterComp(self)

		// properties
		self.DeclareProperty("Output", "foo.gob")
		self.DeclareProperty("Items", g_keys)
		return self

	case "json_outstream":
		self := &json_outstream{}
		kernel.NewAlg(&self.Algorithm, t, n)
		kernel.RegisterComp(self)

		// properties
		self.DeclareProperty("Output", "foo.json")
		self.DeclareProperty("Items", g_keys)
		return self

	case "json_outstream_svc":
		self := &json_outstream_svc{}
		kernel.NewSvc(&self.Service, t, n)
		kernel.RegisterComp(self)

		self.streams = make(map[string]kernel.IOutputStream)
		// properties

		return self

	default:
		err := "no such type [" + t + "]"
		panic(err)
	}
	return nil
}
/* EOF */
