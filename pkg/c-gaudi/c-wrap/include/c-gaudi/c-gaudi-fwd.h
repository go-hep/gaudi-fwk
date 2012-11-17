#ifndef CGAUDI_CGAUDI_FWD_H
#define CGAUDI_CGAUDI_FWD_H 1

#ifdef __cplusplus
extern "C" {
#endif


/* StatusCode */
struct CGaudi_StatusCode {
  unsigned long   code;      ///< The status code
};

typedef void* CGaudi_IInterface;
typedef void* CGaudi_INamedInterface;
typedef void* CGaudi_IAlgorithm;
typedef void* CGaudi_IService;
typedef void* CGaudi_IAlgTool;

typedef void* CGaudi_ISvcLocator;
typedef void* CGaudi_IAppMgr;

struct CGaudi_InterfaceID {
  unsigned long id;
  unsigned long major_ver;
  unsigned long minor_ver;
};

#ifdef __cplusplus
} /* !extern "C" */
#endif

#endif /* !CGAUDI_CGAUDI_FWD_H */
