// reference implementation of an IDataStore
package datastore

import "github.com/sbinet/go-gaudi/pkg/gaudi/kernel"

// --- datastore helper ---

type datastore struct {
	store kernel.DataStore //map[string]interface{}
	//store map[string]chan interface{}
}

func (self *datastore) Put(key string, value interface{}) {
	// _, ok := self.store[key]
	// if !ok {
	//  	self.store[key] = make(chan interface{}, 1)
	// }
	
	// self.store[key] <- value
	self.store[key] = value
}

func (self *datastore) Get(key string) interface{} {
	value, ok := self.store[key]
	// if !ok {
	//  	self.store[key] = make(chan interface{}, 1)
	// }
	if ok {
		//v := <-value
		//self.store[key] <- v
		return value
	}
	return nil
}

func (self *datastore) Has(key string) bool {
	_, ok := self.store[key]
	// if !ok {
	// 	delete(self.store, key)
	// }
	return ok
}

func (self *datastore) ClearStore() kernel.Error {
	for k, _ := range self.store {
		delete(self.store, k)
	}
	return kernel.StatusCode(0)
}

// --- datastore service ---

type datastoresvc struct {
	kernel.Service
}

func (self *datastoresvc) InitializeSvc() kernel.Error {
	self.MsgInfo("~~ initialize [datastore svc] ~~\n")
	return kernel.StatusCode(0)
}

func (self *datastoresvc) FinalizeSvc() kernel.Error {
	self.MsgInfo("~~ finalize [datastore svc] ~~\n")
	return kernel.StatusCode(0)
}

func (self *datastoresvc) Store(ctx kernel.IEvtCtx) kernel.IDataStore {
	store := ctx.Store()
	n := self.CompName()
	if _, ok := store[n]; !ok {
		dstore := make(kernel.DataStore)
		//dstore := make(map[string]chan interface{})
		store[n] = dstore
	}

	dstore := store[n].(kernel.DataStore)
	//dstore := store[n].(map[string]chan interface{})
	if dstore != nil {
		return &datastore{dstore}
	}

	return nil
}

// check matching interfaces
var _ kernel.IDataStore = (*datastore)(nil)
var _ kernel.IComponent = (*datastoresvc)(nil)
var _ kernel.IService = (*datastoresvc)(nil)
var _ kernel.IProperty = (*datastoresvc)(nil)

// --- factory function ---
func New(t, n string) kernel.IComponent {
	switch t {
	case "datastoresvc":
		self := &datastoresvc{}
		//self.stores = make([]datastore, 1)
		_ = kernel.NewSvc(&self.Service, t, n)
		kernel.RegisterComp(self)
		return self
	default:
		err := "no such type [" + t + "]"
		panic(err)
	}
	return nil
}

/* EOF */
