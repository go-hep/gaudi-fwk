#include "c-gaudi/gaudi.h"
#include "GaudiKernel/IAppMgrUI.h"

CGaudi_StatusCode
CGaudi_IAppMgr_run(CGaudi_IAppMgr self)
{
  StatusCode sc = ((IAppMgrUI*)self)->run();
  return *(CGaudi_StatusCode*)(&sc);
}

