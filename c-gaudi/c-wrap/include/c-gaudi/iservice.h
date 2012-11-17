#ifndef CGAUDI_ISERVICE_H
#define CGAUDI_ISERVICE_H 1

#include "c-gaudi/c-gaudi-fwd.h"
#include "c-gaudi/c-gaudi-api.h"

#ifdef __cplusplus
extern "C" {
#endif

/* IService */
CGAUDI_API
CGaudi_StatusCode
CGaudi_IService_configure(CGaudi_IService self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IService_initialize(CGaudi_IService self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IService_start(CGaudi_IService self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IService_stop(CGaudi_IService self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IService_finalized(CGaudi_IService self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IService_terminate(CGaudi_IService self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IService_reinitialize(CGaudi_IService self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IService_restart(CGaudi_IService self);


  /** Get the current state.
   */
  /* virtual Gaudi::StateMachine::State FSMState() const = 0; */

CGAUDI_API
CGaudi_StatusCode
CGaudi_IService_sysInitialize(CGaudi_IService self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IService_sysReinitialize(CGaudi_IService self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IService_sysRestart(CGaudi_IService self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IService_sysStop(CGaudi_IService self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IService_sysFinalize(CGaudi_IService self);

#ifdef __cplusplus
} /* !extern "C" */
#endif

#endif /* !CGAUDI_ISERVICE_H */
