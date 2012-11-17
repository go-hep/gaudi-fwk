#ifndef CGAUDI_IALGORITHM_H
#define CGAUDI_IALGORITHM_H 1

#include "c-gaudi/c-gaudi-api.h"
#include "c-gaudi/c-gaudi-fwd.h"

#ifdef __cplusplus
extern "C" {
#endif

/* IAlgorithm */
CGAUDI_API
CGaudi_StatusCode
CGaudi_IAlgorithm_execute(CGaudi_IAlgorithm self);

CGAUDI_API
int 
CGaudi_IAlgorithm_isInitialized(CGaudi_IAlgorithm self);

CGAUDI_API
int
CGaudi_IAlgorithm_isFinalized(CGaudi_IAlgorithm self);

CGAUDI_API
int
CGaudi_IAlgorithm_isExecuted(CGaudi_IAlgorithm self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IAlgorithm_configure(CGaudi_IAlgorithm self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IAlgorithm_initialize(CGaudi_IAlgorithm self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IAlgorithm_start(CGaudi_IAlgorithm self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IAlgorithm_stop(CGaudi_IAlgorithm self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IAlgorithm_finalized(CGaudi_IAlgorithm self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IAlgorithm_terminate(CGaudi_IAlgorithm self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IAlgorithm_reinitialize(CGaudi_IAlgorithm self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IAlgorithm_restart(CGaudi_IAlgorithm self);


  /** Get the current state.
   */
  /* virtual Gaudi::StateMachine::State FSMState() const = 0; */

CGAUDI_API
CGaudi_StatusCode
CGaudi_IAlgorithm_sysInitialize(CGaudi_IAlgorithm self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IAlgorithm_sysReinitialize(CGaudi_IAlgorithm self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IAlgorithm_sysRestart(CGaudi_IAlgorithm self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IAlgorithm_sysExecute(CGaudi_IAlgorithm self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IAlgorithm_sysStop(CGaudi_IAlgorithm self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IAlgorithm_sysFinalize(CGaudi_IAlgorithm self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IAlgorithm_sysBeginRun(CGaudi_IAlgorithm self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IAlgorithm_sysEndRun(CGaudi_IAlgorithm self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IAlgorithm_beginRun(CGaudi_IAlgorithm self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IAlgorithm_endRun(CGaudi_IAlgorithm self);

CGAUDI_API
void
CGaudi_IAlgorithm_resetExecuted(CGaudi_IAlgorithm self);

CGAUDI_API
void
CGaudi_IAlgorithm_setExecuted(CGaudi_IAlgorithm self, int state);

CGAUDI_API
int
CGaudi_IAlgorithm_isEnabled(CGaudi_IAlgorithm self);

CGAUDI_API
int
CGaudi_IAlgorithm_filterPassed(CGaudi_IAlgorithm self);

CGAUDI_API
void
CGaudi_IAlgorithm_setFilterPassed(CGaudi_IAlgorithm self, int state);

#ifdef __cplusplus
} /* !extern "C" */
#endif

#endif /* !CGAUDI_IALGORITHM_H */
