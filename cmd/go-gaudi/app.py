#!/usr/bin/env python

"""\
Driver for the ng-gaudi 'framework'.

This will read a python jobo file (or a set of those) and create a go-binary from the assembled set of go-packages and ng-components.
"""

### imports ------------------------------------------------------------------
import sys
import os
import subprocess
import tempfile
if not tempfile.tempdir:
    tempfile.tempdir = "/tmp"
import shutil

### functions ----------------------------------------------------------------
def _make_configurable(pkg_name, name, **kwds):
    """create a configurable for a component from package `pkg_name` with
    name `name`. the keyword arguments are its properties
    """
    kwds['name']=name
    return Configurable(pkg_name=pkg_name, **kwds)
Alg = _make_configurable
Svc = _make_configurable
Tool= _make_configurable

def _save_formatted_go_src(fname, src, show=False):
    if os.path.exists(fname):
        os.remove(fname)
    with open(fname, "w") as f:
        if isinstance(src, (list,tuple)):
            print >> f, "\n".join(src)
        else:
            print >> f, src
    subprocess.check_call(["gofmt", "-w",fname])

    if show:
        print ":"*80
        print "::: [%s]" % (fname,)
        with open(fname, 'r') as f:
            print "".join(f.readlines())
        print ":"*80

    return

### classes ------------------------------------------------------------------
class Lvl:
    VERBOSE = 0
    DEBUG = 1
    INFO = 2
    WARNING = 3
    ERROR = 4
    FATAL = 5
    ALWAYS = 6
    
class Configurable(object):
    def __init__(self, pkg_name, **kwds):
        self.pkg_name, self.comp_type = pkg_name.split(":")
        self.name = kwds['name']
        self.props = dict(kwds)
        del self.props['name']

    def __repr__(self):
        return "<Configurable(pkg='%s',name='%s')>" % (
            self.pkg_name, self.name,
            )
    
class cfglist(list):
    def __init__(self, *args, **kw):
        list.__init__(self, *args, **kw)
        self._db = {}
    def __iadd__(self, o):
        if isinstance(o, Configurable):
            self._db[o.name] = o
            o = [o]
        return list.__iadd__(self, o)
    def __getitem__(self, idx):
        if isinstance(idx, basestring):
            return self._db[idx]
        return list.__getitem__(idx)
    def __delitem__(self, idx):
        # FIXME: keep dict and list in-sync !!
        if isinstance(idx, basestring):
            del self._db[idx]
            return
        return list.__getitem__(idx)
    
class attrdict(dict):
    def __setattr__(self, k, v):
        self[k] = v
    def __getattr__(self, k):
        return self[k]
    
class AppMgr(object):
    def __init__(self):
        self.algs = cfglist()
        self.svcs = cfglist()
        self.toolsvc = cfglist()

        self.props = attrdict(
            EvtMax=10,
            OutputLevel=1,
            )
        
        self._workdir = tempfile.mkdtemp(prefix='ng-gaudi-go-')
        #self._workdir = "/tmp/%s/ng-gaudi/gaudi-jobopt" % (os.getenv('LOGNAME'),)'
        return

    def configure(self, jobopts):
        if not isinstance(jobopts, (list,tuple)):
            raise TypeError("jobopts should be a sequence")
        self.jobopts = jobopts[:]
        for i,jobo in enumerate(self.jobopts):
            jobo = os.path.expandvars(os.path.expanduser(jobo))
            self.jobopts[i] = jobo
            if not os.path.exists(jobo):
                raise RuntimeError("no such file [%s]" % jobo)
            print "::: including [%s]..." % (jobo,)
            execfile(jobo)
            print "::: including [%s]... [ok]" % (jobo,)
            
    def run(self):
        print "::: algs:", len(self.algs)
        print "::: svcs:", len(self.svcs)
        print "::: tool:", len(self.toolsvc)

        fname="gaudi_jobopt.go"
        self._gen_golang_pkg(fname)
        self._compile_golang_pkg(fname)

        exitcode = 0
        
        orig_dir = os.getcwd()
        try:
            os.chdir(self._workdir)
            #print os.listdir('.')
            cmd = ["./gaudi-main",]
            #cmd = ["/usr/bin/time", "-o", "%s/gaudi.profile.out" % orig_dir, "./gaudi-main",]
            subprocess.check_call(cmd)
        except Exception:
            import traceback
            traceback.print_exc()
            exitcode = 1
        finally:
            os.chdir(orig_dir)
            if exitcode == 0:
                shutil.rmtree(self._workdir)
        return exitcode

    def _gen_golang_pkg(self, fname=None):
        if fname is None:
            fname = "gaudi_jobopt.go"

        fname = os.path.join(self._workdir, os.path.basename(fname))
        go_pkg = ["package gaudi_jobopt",
                  "",
                  'import "bitbucket.org/binet/ng-go-gaudi/pkg/gaudi/kernel"',
                  ]
        # a symbol table for pkgs (we may need to mangle them...)
        gaudi_pkgs = {}
        _pkg_idx = 0
        for comp in self.algs+self.svcs+self.toolsvc:
            if not comp.pkg_name in gaudi_pkgs:
                gaudi_pkgs[comp.pkg_name] = "_gaudi_pkg_%08i" % (_pkg_idx,)
                _pkg_idx += 1
        for pkg,idx in gaudi_pkgs.iteritems():
            go_pkg += ['import %s "%s"' % (idx,pkg,)]

        go_pkg += ["",
                   "type CompCfg struct {",
                   " Instance kernel.IComponent",
                   " Name string",
                   "}",
                   "var Algs []CompCfg",
                   "var Svcs []CompCfg",
                   "var Tools []CompCfg",
                   "var AppMgr kernel.IAppMgr",
                   ""
                   ]
        go_pkg += ["func init() {"]

        go_pkg += [
            "AppMgr = kernel.NewAppMgr()",
            "{",
            "c := AppMgr.(kernel.IProperty)",
            ]
        try:
            import multiprocessing as _mp
            go_pkg += [
                "c.SetProperty(\"NbrProcs\", %r)" % (_mp.cpu_count(),)
                ]
        except ImportError:
            pass

        def to_gorepr(v):
            #print "-- got [%r] type:[%s]..." % (v,type(v).__name__)
            if isinstance(v, basestring):
                return '\"%s\"' % (v,)
            elif isinstance(v, list):
                cls = type(v[0])
                cls_str = 'string' if isinstance(v[0], basestring) else cls.__name__
                return '[]%s{%s}' % (cls_str,
                                     ",".join(map(to_gorepr, v)))
            else:
                return v
            
        for k,v in self.props.iteritems():
            go_pkg += [
                "c.SetProperty(\"%s\", %r)" % (k, to_gorepr(v))
                ]
        go_pkg += [
            "}",
            ]
        
        for comps,name in [(self.algs, "Algs"),
                           (self.svcs, "Svcs"),
                           (self.toolsvc, "Tools")]:
            go_pkg += ["%s = []CompCfg{" % (name,)]
            for i,comp in enumerate(comps):
                pkg_ident = gaudi_pkgs[comp.pkg_name]
                comma = "," if i+1 != len(comps) else "}"
                go_pkg += [ 'CompCfg{Instance:%s.New(\"%s\",\"%s\"), Name:"%s"}%s' %
                            (pkg_ident,
                             comp.comp_type,
                             comp.name,
                             comp.name,
                             comma)]
                pass
            for i,comp in enumerate(comps):
                if len(comp.props) <= 0:
                    continue
                go_pkg += [
                    "{",
                    "// %s" % (comp.name,),
                    "c,ok := %s[%i].Instance.(kernel.IProperty)" % (name, i),
                    "if ok {"
                    ]
                for k, val in comp.props.iteritems():
                    fmt = "c.SetProperty(\"%s\", %s)"
                    val = fmt % (k, to_gorepr(val))
                    go_pkg += [ val ]
                go_pkg += ["}","}"]

        go_pkg += ["}"]

        go_pkg += ["", "/* EOF */", ""]
        _save_formatted_go_src(fname=fname, src=go_pkg, show=False)

        return

    def _compile_golang_pkg(self, fname):
        goarch = os.getenv('GOARCH')
        arch = '32bit' if '32' in goarch else '64bit'
        
        compiler = {
            '64bit': '6g',
            '32bit': '8g',
            }[arch]
        linker = {
            '64bit': '6l',
            '32bit': '8l',
            }[arch]
        orig_dir = os.getcwd()
        work_dir = self._workdir #os.path.dirname(fname)
        try:
            os.chdir(work_dir)

            # obj_fname = os.path.basename(fname)
            # obj_fname = os.splitext(obj_fname)[0] + '.o'
            cmd = [compiler, fname]
            #print "::: curdir: %s" % (os.getcwd(),)
            print "::: compiling 'gaudi_jobopt'..."
            subprocess.check_call(cmd)
            print "::: compiling 'gaudi_jobopt'... [ok]"

            go_main = '''\n
package main

import "fmt"
//import "os"

import "bitbucket.org/binet/ng-go-gaudi/pkg/gaudi/kernel"
import "./gaudi_jobopt"

func main() {
   fmt.Printf("::: gaudi\\n")
   app := gaudi_jobopt.AppMgr
   msg := app.(kernel.IMessager)
   msg.MsgInfo(" -> created [%s/%s]\\n", app.CompType(), app.CompName())
   msg.MsgInfo("::: configure...\\n")
   app_prop,_ := app.(kernel.IProperty)
   {
     mgr,_ := app.(kernel.IAlgMgr)
     algs := make([]string, len(gaudi_jobopt.Algs))
     for i,alg := range gaudi_jobopt.Algs {
     ialg,ok := alg.Instance.(kernel.IAlgorithm)
       if ok {
          ai := alg.Instance
          mgr.AddAlgorithm(ialg)
          algs[i] = ialg.CompName()
          msg.MsgDebug("algorithm [%T/%T/%s] registered\\n", ai, ialg, ialg.CompName())
          iprop, ok := alg.Instance.(kernel.IProperty)
          if ok {
           if iprop != nil {
              msg.MsgDebug("alg [%s] implements kernel.IProperty\\n", ialg.CompName())
           }
          }
       }
     }
     app_prop.SetProperty("Algs", algs)
   }
   {
     mgr,_ := app.(kernel.ISvcMgr)
     svcs := make([]string, len(gaudi_jobopt.Svcs))
     for i,svc := range gaudi_jobopt.Svcs {
       isvc,ok := svc.Instance.(kernel.IService)
       if ok && isvc.CompName() != "" {
          msg.MsgDebug("adding service [%s]...\\n", isvc.CompName())
          if !mgr.AddService(isvc.CompName()).IsSuccess() {
             msg.MsgError("pb adding svc [%s]\\n", isvc.CompName())
          }
          svcs[i] = isvc.CompName()
          msg.MsgDebug("adding service [%s]... [done]\\n", isvc.CompName())
       }
     }
     app_prop.SetProperty("Svcs", svcs)
   }
   sc := app.Configure()
   msg.MsgInfo("::: configure... [%d]\\n", sc.Code())
   msg.MsgInfo("::: run...\\n")
   sc = app.Run()
   msg.MsgInfo("::: run... [%d]\\n", sc.Code())
   msg.MsgInfo("::: bye.\\n")
}

/* EOF */
'''
            go_main_fname = os.path.join(self._workdir, "gaudi_main.go")
            _save_formatted_go_src(fname=go_main_fname, src=go_main,
                                   show=False)
            cmd = [compiler, "-o", "gaudi_main.o", go_main_fname]
            print "::: compiling 'gaudi_main'..."
            subprocess.check_call(cmd)
            print "::: compiling 'gaudi_main'... [ok]"
            print "::: linking 'gaudi_main'..."
            cmd = [linker, "-o", "gaudi-main", "gaudi_main.o"]
            subprocess.check_call(cmd)
            print "::: linking 'gaudi_main'... [ok]"
        finally:
            os.chdir(orig_dir)
        return
    
app = AppMgr()

if __name__ == "__main__":
    print ":"*80
    print "::: welcome to ng-go-gaudi"
    print ":"*80
    jobopts = []
    if len(sys.argv)>1:
        for arg in sys.argv[1:]:
            if arg.endswith(".py"):
                jobopts += [arg]
    if len(jobopts) == 0:
        jobopts = ["jobopt.py"]
    app.configure(jobopts)
    sys.exit(app.run())
