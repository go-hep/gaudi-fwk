#include "c-gaudi/gaudi.h"
#include "GaudiKernel/IService.h"

/* IService */
CGaudi_StatusCode
CGaudi_IService_configure(CGaudi_IService self)
{
  StatusCode sc = ((IService*)self)->configure();
  return *(CGaudi_StatusCode*)(&sc);
}

CGaudi_StatusCode
CGaudi_IService_initialize(CGaudi_IService self)
{
  StatusCode sc = ((IService*)self)->initialize();
  return *(CGaudi_StatusCode*)(&sc);
}

CGaudi_StatusCode
CGaudi_IService_start(CGaudi_IService self)
{
  StatusCode sc = ((IService*)self)->start();
  return *(CGaudi_StatusCode*)(&sc);
}

CGaudi_StatusCode
CGaudi_IService_stop(CGaudi_IService self)
{
  StatusCode sc = ((IService*)self)->stop();
  return *(CGaudi_StatusCode*)(&sc);
}

CGaudi_StatusCode
CGaudi_IService_finalize(CGaudi_IService self)
{
  StatusCode sc = ((IService*)self)->finalize();
  return *(CGaudi_StatusCode*)(&sc);
}

CGaudi_StatusCode
CGaudi_IService_terminate(CGaudi_IService self)
{
  StatusCode sc = ((IService*)self)->terminate();
  return *(CGaudi_StatusCode*)(&sc);
}

CGaudi_StatusCode
CGaudi_IService_reinitialize(CGaudi_IService self)
{
  StatusCode sc = ((IService*)self)->reinitialize();
  return *(CGaudi_StatusCode*)(&sc);
}

CGaudi_StatusCode
CGaudi_IService_restart(CGaudi_IService self)
{
  StatusCode sc = ((IService*)self)->restart();
  return *(CGaudi_StatusCode*)(&sc);
}

CGaudi_StatusCode
CGaudi_IService_sysInitialize(CGaudi_IService self)
{
  StatusCode sc = ((IService*)self)->sysInitialize();
  return *(CGaudi_StatusCode*)(&sc);
}

CGaudi_StatusCode
CGaudi_IService_sysReinitialize(CGaudi_IService self)
{
  StatusCode sc = ((IService*)self)->sysReinitialize();
  return *(CGaudi_StatusCode*)(&sc);
}

CGaudi_StatusCode
CGaudi_IService_sysRestart(CGaudi_IService self)
{
  StatusCode sc = ((IService*)self)->sysRestart();
  return *(CGaudi_StatusCode*)(&sc);
}

CGaudi_StatusCode
CGaudi_IService_sysStop(CGaudi_IService self)
{
  StatusCode sc = ((IService*)self)->sysStop();
  return *(CGaudi_StatusCode*)(&sc);
}

CGaudi_StatusCode
CGaudi_IService_sysFinalize(CGaudi_IService self)
{
  StatusCode sc = ((IService*)self)->sysFinalize();
  return *(CGaudi_StatusCode*)(&sc);
}
