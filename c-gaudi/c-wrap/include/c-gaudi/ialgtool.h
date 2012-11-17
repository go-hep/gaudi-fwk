#ifndef CGAUDI_IALGTOOL_H
#define CGAUDI_IALGTOOL_H 1

#include "c-gaudi/c-gaudi-fwd.h"
#include "c-gaudi/c-gaudi-api.h"

#ifdef __cplusplus
extern "C" {
#endif

/* IAlgTool */
CGAUDI_API
CGaudi_StatusCode
CGaudi_IAlgTool_configure(CGaudi_IAlgTool self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IAlgTool_initialize(CGaudi_IAlgTool self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IAlgTool_start(CGaudi_IAlgTool self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IAlgTool_stop(CGaudi_IAlgTool self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IAlgTool_finalized(CGaudi_IAlgTool self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IAlgTool_terminate(CGaudi_IAlgTool self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IAlgTool_reinitialize(CGaudi_IAlgTool self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IAlgTool_restart(CGaudi_IAlgTool self);


CGAUDI_API
CGaudi_StatusCode
CGaudi_IAlgTool_sysInitialize(CGaudi_IAlgTool self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IAlgTool_sysReinitialize(CGaudi_IAlgTool self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IAlgTool_sysRestart(CGaudi_IAlgTool self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IAlgTool_sysStop(CGaudi_IAlgTool self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IAlgTool_sysFinalize(CGaudi_IAlgTool self);

#ifdef __cplusplus
} /* !extern "C" */
#endif

#endif /* !CGAUDI_IALGTOOL_H */
