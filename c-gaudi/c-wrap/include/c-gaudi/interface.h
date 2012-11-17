#ifndef CGAUDI_INTERFACE_H
#define CGAUDI_INTERFACE_H 1

#include "c-gaudi/c-gaudi-api.h"
#include "c-gaudi/c-gaudi-fwd.h"

#ifdef __cplusplus
extern "C" {
#endif

/* InterfaceID */

CGAUDI_API
int
CGaudi_InterfaceID_versionMatch(CGaudi_InterfaceID self, CGaudi_InterfaceID other);

CGAUDI_API
int
CGaudi_InterfaceID_fullMatch(CGaudi_InterfaceID self, CGaudi_InterfaceID other);

/* IInterface */

CGAUDI_API
CGaudi_InterfaceID
CGaudi_IInterface_InterfaceID(CGaudi_IInterface self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IInterface_queryInterface(CGaudi_IInterface self, CGaudi_InterfaceID iid, void **p);

CGAUDI_API
unsigned long
CGaudi_IInterface_addRef(CGaudi_IInterface self);

CGAUDI_API
unsigned long
CGaudi_IInterface_release(CGaudi_IInterface self);

CGAUDI_API
unsigned long
CGaudi_IInterface_refCount(CGaudi_IInterface self);

/* INamedInterface */
CGAUDI_API
const char*
CGaudi_INamedInterface_name(CGaudi_INamedInterface self);

#ifdef __cplusplus
} /* !extern "C" */
#endif

#endif /* !CGAUDI_INTERFACE_H */

