#ifndef CGAUDI_IAPPMGR_H
#define CGAUDI_IAPPMGR_H 1

#include "c-gaudi/c-gaudi-api.h"
#include "c-gaudi/c-gaudi-fwd.h"

#ifdef __cplusplus
extern "C" {
#endif

CGAUDI_API
CGaudi_StatusCode
CGaudi_IAppMgr_run(CGaudi_IAppMgr self);

#ifdef __cplusplus
} /* !extern "C" */
#endif

#endif /* !CGAUDI_IAPPMGR_H */
