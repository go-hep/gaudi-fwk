#include "c-gaudi/gaudi.h"
#include "GaudiKernel/IAlgTool.h"

/* IAlgTool */
CGaudi_StatusCode
CGaudi_IAlgTool_configure(CGaudi_IAlgTool self)
{
  StatusCode sc = ((IAlgTool*)self)->configure();
  return *(CGaudi_StatusCode*)(&sc);
}

CGaudi_StatusCode
CGaudi_IAlgTool_initialize(CGaudi_IAlgTool self)
{
  StatusCode sc = ((IAlgTool*)self)->initialize();
  return *(CGaudi_StatusCode*)(&sc);
}

CGaudi_StatusCode
CGaudi_IAlgTool_start(CGaudi_IAlgTool self)
{
  StatusCode sc = ((IAlgTool*)self)->start();
  return *(CGaudi_StatusCode*)(&sc);
}

CGaudi_StatusCode
CGaudi_IAlgTool_stop(CGaudi_IAlgTool self)
{
  StatusCode sc = ((IAlgTool*)self)->stop();
  return *(CGaudi_StatusCode*)(&sc);
}

CGaudi_StatusCode
CGaudi_IAlgTool_finalize(CGaudi_IAlgTool self)
{
  StatusCode sc = ((IAlgTool*)self)->finalize();
  return *(CGaudi_StatusCode*)(&sc);
}

CGaudi_StatusCode
CGaudi_IAlgTool_terminate(CGaudi_IAlgTool self)
{
  StatusCode sc = ((IAlgTool*)self)->terminate();
  return *(CGaudi_StatusCode*)(&sc);
}

CGaudi_StatusCode
CGaudi_IAlgTool_reinitialize(CGaudi_IAlgTool self)
{
  StatusCode sc = ((IAlgTool*)self)->reinitialize();
  return *(CGaudi_StatusCode*)(&sc);
}

CGaudi_StatusCode
CGaudi_IAlgTool_restart(CGaudi_IAlgTool self)
{
  StatusCode sc = ((IAlgTool*)self)->restart();
  return *(CGaudi_StatusCode*)(&sc);
}

CGaudi_StatusCode
CGaudi_IAlgTool_sysInitialize(CGaudi_IAlgTool self)
{
  StatusCode sc = ((IAlgTool*)self)->sysInitialize();
  return *(CGaudi_StatusCode*)(&sc);
}

CGaudi_StatusCode
CGaudi_IAlgTool_sysReinitialize(CGaudi_IAlgTool self)
{
  StatusCode sc = ((IAlgTool*)self)->sysReinitialize();
  return *(CGaudi_StatusCode*)(&sc);
}

CGaudi_StatusCode
CGaudi_IAlgTool_sysRestart(CGaudi_IAlgTool self)
{
  StatusCode sc = ((IAlgTool*)self)->sysRestart();
  return *(CGaudi_StatusCode*)(&sc);
}

CGaudi_StatusCode
CGaudi_IAlgTool_sysStop(CGaudi_IAlgTool self)
{
  StatusCode sc = ((IAlgTool*)self)->sysStop();
  return *(CGaudi_StatusCode*)(&sc);
}

CGaudi_StatusCode
CGaudi_IAlgTool_sysFinalize(CGaudi_IAlgTool self)
{
  StatusCode sc = ((IAlgTool*)self)->sysFinalize();
  return *(CGaudi_StatusCode*)(&sc);
}
