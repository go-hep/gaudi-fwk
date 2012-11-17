#ifndef CGAUDI_ISVCLOCATOR_H
#define CGAUDI_ISVCLOCATOR_H 1

#include "c-gaudi/c-gaudi-api.h"
#include "c-gaudi/c-gaudi-fwd.h"

#ifdef __cplusplus
extern "C" {
#endif

CGAUDI_API
CGaudi_StatusCode
CGaudi_ISvcLocator_getService(CGaudi_ISvcLocator self,
                              const char *type_name,
                              CGaudi_IService *svc,
                              int createif);

CGAUDI_API
int
CGaudi_ISvcLocator_existsService(CGaudi_ISvcLocator self,
                                 const char *name);

  /*
   /// Return the list of Services
  virtual const std::list<IService*> &getServices() const = 0;
  */

#ifdef __cplusplus
} /* !extern "C" */
#endif

#endif /* !CGAUDI_ISVCLOCATOR_H */
