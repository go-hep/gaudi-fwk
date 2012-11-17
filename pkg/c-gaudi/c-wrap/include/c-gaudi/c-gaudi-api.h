#ifndef CGAUDI_CGAUDI_API_H
#define CGAUDI_CGAUDI_API_H 1

#include <stdint.h>
#include <stddef.h>
#include <stdbool.h>

#if __GNUC__ >= 4
#  define CGAUDI_HASCLASSVISIBILITY
#endif

#if defined(CGAUDI_HASCLASSVISIBILITY)
#  define CGAUDI_IMPORT __attribute__((visibility("default")))
#  define CGAUDI_EXPORT __attribute__((visibility("default")))
#  define CGAUDI_LOCAL  __attribute__((visibility("hidden")))
#else
#  define CGAUDI_IMPORT
#  define CGAUDI_EXPORT
#  define CGAUDI_LOCAL
#endif

#define CGAUDI_API CGAUDI_EXPORT

#endif /* !CGAUDI_CGAUDI_API_H */

