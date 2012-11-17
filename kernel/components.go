package kernel

import "fmt"

type Component struct {
	comp_name string
	comp_type string
}

func (self *Component) CompName() string {
	return self.comp_name
}

func (self *Component) CompType() string {
	return self.comp_type
}

func NewComponent(t, n string) IComponent {
	self := &Component{comp_name: n, comp_type: t}
	g_compsdb[n] = self
	return self
}

// --- properties ---
type properties struct {
	props map[string]interface{}
}

func (self *properties) DeclareProperty(n string, v interface{}) {
	self.SetProperty(n, v)
}
func (self *properties) SetProperty(n string, v interface{}) Error {
	self.props[n] = v
	return StatusCode(0)
}
func (self *properties) GetProperty(n string) interface{} {
	v, ok := self.props[n]
	if ok {
		return v
	}
	return nil
}
func (self *properties) GetProperties() []Property {
	props := make([]Property, len(self.props))
	i := 0
	for k, v := range self.props {
		props[i] = Property{Name: k, Value: v}
		i++
	}
	return props
}

// --- output level ---
type OutputLevel int

const (
	LVL_VERBOSE OutputLevel = iota
	LVL_DEBUG
	LVL_INFO
	LVL_WARNING
	LVL_ERROR
	LVL_FATAL
	LVL_ALWAYS
)

func (self OutputLevel) String() string {
	switch self {
	case LVL_VERBOSE:
		return "VERBOSE"
	case LVL_DEBUG:
		return "DEBUG"
	case LVL_INFO:
		return "INFO"
	case LVL_WARNING:
		return "WARNING"
	case LVL_ERROR:
		return "ERROR"
	case LVL_FATAL:
		return "FATAL"
	case LVL_ALWAYS:
		return "ALWAYS"
	default:
		return "???"
	}
	return "???"
}

// msgstream
type msgstream struct {
	name  string
	level OutputLevel
}

func (self *msgstream) SetOutputLevel(lvl OutputLevel) {
	self.level = lvl
}
func (self *msgstream) OutputLevel() OutputLevel {
	return self.level
}
func (self *msgstream) Msg(lvl OutputLevel, format string, a ...interface{}) (int, error) {
	//return 0, nil
	if self.level <= lvl {
		s := fmt.Sprintf(format, a...)
		return fmt.Printf("%-15s %6s %s", self.name, lvl, s)
	}
	return 0, nil
}
func (self *msgstream) MsgVerbose(format string, a ...interface{}) (int, error) {
	return self.Msg(LVL_VERBOSE, format, a...)
}
func (self *msgstream) MsgDebug(format string, a ...interface{}) (int, error) {
	return self.Msg(LVL_DEBUG, format, a...)
}
func (self *msgstream) MsgInfo(format string, a ...interface{}) (int, error) {
	return self.Msg(LVL_INFO, format, a...)
}
func (self *msgstream) MsgWarning(format string, a ...interface{}) (int, error) {
	return self.Msg(LVL_WARNING, format, a...)
}
func (self *msgstream) MsgError(format string, a ...interface{}) (int, error) {
	return self.Msg(LVL_ERROR, format, a...)
}
func (self *msgstream) MsgFatal(format string, a ...interface{}) (int, error) {
	return self.Msg(LVL_FATAL, format, a...)
}
func (self *msgstream) MsgAlways(format string, a ...interface{}) (int, error) {
	return self.Msg(LVL_ALWAYS, format, a...)
}

// algorithm
type Algorithm struct {
	Component
	properties
	msgstream
	evtstore IDataStoreMgr
	detstore IDataStoreMgr
	//stores map[string]IDataStoreMgr
}

// convenience function
func (self *Algorithm) EvtStore(ctx IEvtCtx) IDataStore {
	return self.evtstore.Store(ctx)
}
func (self *Algorithm) DetStore(ctx IEvtCtx) IDataStore {
	return self.detstore.Store(ctx)
}

//func (self *Algorithm) SysInitialize() Error {
//	return self.Initialize()
//}

//func (self *Algorithm) SysExecute(ctx IEvtCtx) Error {
//	self.MsgInfo("sys-execute... [%v]\n", ctx)
//	println("==>",self.CompName(),"sys-execute [",ctx,"]")
//	return self.Execute(ctx)
//}

//func (self *Algorithm) SysFinalize() Error {
//	return self.Finalize()
//}

func (self *Algorithm) Initialize() Error {
	lvl := self.GetProperty("OutputLevel").(int)
	self.SetOutputLevel(OutputLevel(lvl))

	self.MsgInfo("initialize...\n")
	svcloc := GetSvcLocator()
	if svcloc == nil {
		self.MsgError("could not retrieve svclocator\n")
		return StatusCode(1)
	}

	self.evtstore = svcloc.GetService("evt-store").(IDataStoreMgr)
	self.detstore = svcloc.GetService("det-store").(IDataStoreMgr)
	return StatusCode(0)
}

func (self *Algorithm) Execute(ctx IEvtCtx) Error {
	self.MsgDebug("execute... [%v]\n", ctx)
	println("==>", self.CompName(), "execute [", ctx, "]")
	return StatusCode(0)
}

func (self *Algorithm) Finalize() Error {
	self.MsgInfo("finalize...\n")
	return StatusCode(0)
}

/*
func (self *Algorithm) Store(ctx IEvtCtx, n string) IDataStore {
	store, ok := self.stores[n]
	if !ok {
        delete(self.stores, n)
		return nil
	}
	return store.Store(ctx)
}
*/

func NewAlg(comp IComponent, t, n string) IAlgorithm {
	self := comp.(*Algorithm)
	self.Component.comp_name = n
	self.Component.comp_type = t
	self.properties.props = make(map[string]interface{})
	self.msgstream.name = n
	self.msgstream.level = LVL_INFO

	self.evtstore = nil
	self.detstore = nil

	self.DeclareProperty("OutputLevel", int(LVL_INFO))

	//g_compsdb[n] = self

	return self
}

// service
type Service struct {
	Component
	properties
	msgstream
}

//func (self *Service) SysInitializeSvc() Error {
//	return self.InitializeSvc()
//}

//func (self *Service) SysFinalizeSvc() Error {
//	return self.FinalizeSvc()
//}

func (self *Service) InitializeSvc() Error {
	self.MsgInfo("initialize...\n")
	return StatusCode(0)
}

func (self *Service) FinalizeSvc() Error {
	self.MsgInfo("finalize...\n")
	return StatusCode(0)
}

func NewSvc(comp IComponent, t, n string) IService {
	self := comp.(*Service)
	self.Component.comp_name = n
	self.Component.comp_type = t
	self.properties.props = make(map[string]interface{})

	self.msgstream.name = n
	self.msgstream.level = LVL_INFO

	//g_compsdb[n] = self
	return self
}

// algtool
type AlgTool struct {
	Component
	properties
	msgstream
	parent IComponent
}

func (self *AlgTool) CompName() string {
	// FIXME: implement toolsvc !
	if self.parent != nil {
		return self.parent.CompName() + "." + self.Component.CompName()
	}
	return "ToolSvc." + self.Component.CompName()
}

//func (self *AlgTool) SysInitializeTool() Error {
//	return self.InitializeTool()
//}

//func (self *AlgTool) SysFinalizeTool() Error {
//	return self.FinalizeTool()
//}

func (self *AlgTool) InitializeTool() Error {
	self.MsgInfo("initialize...\n")
	return StatusCode(0)
}

func (self *AlgTool) FinalizeTool() Error {
	self.MsgInfo("finalize...\n")
	return StatusCode(0)
}

func NewTool(comp IComponent, t, n string, parent IComponent) IAlgTool {
	self := comp.(*AlgTool)
	self.Component.comp_name = n
	self.Component.comp_type = t
	self.properties.props = make(map[string]interface{})
	self.msgstream = msgstream{name: self.CompName(), level: LVL_INFO}
	self.parent = parent

	//g_compsdb[n] = self
	return self
}

func init() {
	g_compsdb = make(comps_db)
	//fmt.Printf("--> components: %v\n", g_compsdb)
}

// checking implementations match interfaces
var _ IAlgorithm = (*Algorithm)(nil)
var _ IComponent = (*Algorithm)(nil)
var _ IProperty = (*Algorithm)(nil)

var _ IAlgTool = (*AlgTool)(nil)
var _ IComponent = (*AlgTool)(nil)
var _ IProperty = (*AlgTool)(nil)

var _ IService = (*Service)(nil)
var _ IComponent = (*Service)(nil)
var _ IProperty = (*Service)(nil)

/* EOF */
