package kernel

import "fmt"

type comps_db map[string]IComponent

/// the main entry point to the gaudi framework: the service locator
var g_isvcloc ISvcLocator = nil
/// the central repository of all gaudi components
var g_compsdb comps_db

/// the interface to assess if everything went right and report why otherwise
type Error interface {
	/// embed 'error' builtin interface
	error
	/// return the error code (0 is success)
	Code() int
	/// the error reason
	Err() error
	IsSuccess() bool
	IsFailure() bool
	IsRecoverable() bool
}

type statuscode struct {
	code int
	err  error
}
// check statuscode implements kernel.Error
var _ = Error(&statuscode{})
//

func StatusCode(code int) Error {
	return &statuscode{code: code, err: nil}
}

func StatusCodeWithErr(code int, err error) Error {
	return &statuscode{code: code, err: err}
}

func (sc *statuscode) Code() int {
	return sc.code
}

func (sc *statuscode) Err() error {
	return sc.err
}

func (sc *statuscode) Error() string {
	return fmt.Sprintf("code:%v err:%v", sc.code, sc.err)
}

func (sc *statuscode) IsSuccess() bool {
	return sc.code == 0
}

func (sc *statuscode) IsFailure() bool {
	return !sc.IsSuccess()
}

func (sc *statuscode) IsRecoverable() bool {
	return sc.code == 2
}

func GetSvcLocator() ISvcLocator {
	return g_isvcloc
}

func ComponentMgr() IComponentMgr {
	isvcloc := GetSvcLocator()
	imgr, ok := isvcloc.(IComponentMgr)
	if ok {
		return imgr
	}
	return nil
}

func RegisterComp(c IComponent) bool {
	if c == nil {
		return false
	}
	n := c.CompName()
	oldcomp, already_there := g_compsdb[n]
	if already_there {
		if oldcomp == c {
			// double registration of the same component...
			// silly but harmless.
			return true
		}
		// already existing component with that same name !
		err := fmt.Sprintf("a component with name [%s] was already registered ! (old-type: %T, new-type: %T)",
			n, oldcomp, c)
		panic(err)
	}
	//fmt.Printf("--> registering [%T/%s]...\n", c, n)
	g_compsdb[n] = c
	//fmt.Printf("--> registering [%T/%s]... [done]\n", c, n)
	return true
}

/// the central interface to the Gaudi component object model
/// a component has a type and an instance name
type IComponent interface {
	CompName() string
	CompType() string
}

type IComponentMgr interface {
	GetComp(n string) IComponent
	GetComps() []IComponent
	HasComp(n string) bool
}

type Property struct {
	Name  string
	Value interface{}
}
type IProperty interface {
	/// declare a property by name and default value
	DeclareProperty(name string, value interface{})
	/// set the property value
	SetProperty(name string, value interface{}) Error
	/// get the property value by name
	GetProperty(name string) interface{}
	/// get the list of properties
	GetProperties() []Property
}

type IService interface {
	IComponent
	//SysInitializeSvc() Error
	//SysFinalizeSvc() Error

	InitializeSvc() Error
	FinalizeSvc() Error
}

type IAlgorithm interface {
	IComponent
	//SysInitialize() Error
	//SysExecute(evtctx IEvtCtx) Error
	//SysFinalize() Error
	Initialize() Error
	Execute(evtctx IEvtCtx) Error
	Finalize() Error
}

type IAlgTool interface {
	IComponent
	//SysInitializeTool() Error
	//SysFinalizeTool() Error

	InitializeTool() Error
	FinalizeTool() Error
}

type DataStore map[string]interface{}

type IEvtCtx interface {
	Idx() int
	Store() DataStore
	//Id() int
}

type IEvtProcessor interface {
	IComponent
	ExecuteEvent(evtctx IEvtCtx) Error
	ExecuteRun(maxevt int) Error
	NextEvent(maxevt int) Error
	StopRun() Error
}

type IEvtSelector interface {
	IComponent
	CreateContext(ctx *IEvtCtx) Error
	Next(ctx *IEvtCtx, jump int) Error
	Previous(ctx *IEvtCtx, jump int) Error
	Last(ctx *IEvtCtx) Error
	Rewind(ctx *IEvtCtx) Error
}

type IAppMgr interface {
	IComponent
	Configure() Error
	Initialize() Error
	Start() Error
	/// Run the complete job (from Initialize to Terminate)
	Run() Error
	Stop() Error
	Finalize() Error
	Terminate() Error
}

type IAlgMgr interface {
	//IComponent
	AddAlgorithm(alg IAlgorithm) Error
	RemoveAlgorithm(alg IAlgorithm) Error
	HasAlgorithm(algname string) bool
}

type ISvcMgr interface {
	//IComponent
	AddService(svc string) Error
	RemoveService(svc string) Error
	HasService(svc string) Error
}

type ISvcLocator interface {
	//IComponent
	GetService(svc string) IService
	GetServices() []IService
	ExistsService(svc string) bool
}

type IDataStore interface {
	//IComponent
	Get(key string) interface{}
	Put(key string, value interface{})
	Has(key string) bool
	//Keys() []string // ??
}

type IDataStoreMgr interface {
	IComponent
	Store(ctx IEvtCtx) IDataStore
}

type IMessager interface {
	Msg(lvl OutputLevel, format string, a ...interface{}) (int, error)
	MsgVerbose(format string, a ...interface{}) (int, error)
	MsgDebug(format string, a ...interface{}) (int, error)
	MsgInfo(format string, a ...interface{}) (int, error)
	MsgWarning(format string, a ...interface{}) (int, error)
	MsgError(format string, a ...interface{}) (int, error)
	MsgFatal(format string, a ...interface{}) (int, error)
	MsgAlways(format string, a ...interface{}) (int, error)
}

/// handle to a concurrent output stream
type IOutputStream interface {
	/// writes (and possibly commit) data to the stream
	Write(data interface{}) Error
	/// closes and flushes the output stream
	Close() Error
	/// returns the name of the output stream (ie: URI)
	Name() string
	/// returns the file-descriptor associated to that output stream
	Fd() int
}

/// interface to a concurrent output stream server
type IOutputStreamSvc interface {
	//IComponent

	/// returns a new output stream
	NewOutputStream(stream_name string) IOutputStream
}

/* EOF */
