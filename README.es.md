# atomix

[![Go Reference](https://pkg.go.dev/badge/code.hybscloud.com/atomix.svg)](https://pkg.go.dev/code.hybscloud.com/atomix)
[![Go Report Card](https://goreportcard.com/badge/github.com/hayabusa-cloud/atomix)](https://goreportcard.com/report/github.com/hayabusa-cloud/atomix)
[![Codecov](https://codecov.io/gh/hayabusa-cloud/atomix/graph/badge.svg)](https://codecov.io/gh/hayabusa-cloud/atomix)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

**Idiomas:** [English](README.md) | [简体中文](README.zh-CN.md) | [日本語](README.ja.md) | Español | [Français](README.fr.md)

Operaciones atómicas con ordenamiento de memoria explícito para Go.

## Descripción

El paquete `sync/atomic` de Go proporciona operaciones atómicas con consistencia secuencial. Esta biblioteca expone los ordenamientos del modelo de memoria C++11/C11 (Relaxed, Acquire, Release, AcqRel) mediante implementaciones específicas de arquitectura.

```go
import "code.hybscloud.com/atomix"

var counter atomix.Int64

// API basada en métodos con sufijo de ordenamiento
counter.AddRelaxed(1)    // Relaxed: sin sincronización
counter.Add(1)           // AcqRel: ordenamiento seguro por defecto

// API basada en punteros para memoria cruda
var flags int32
atomix.Relaxed.StoreInt32(&flags, 1)
val := atomix.Acquire.LoadInt32(&flags)
```

## Instalación

```bash
go get code.hybscloud.com/atomix
```

**Requisitos:** Go 1.25+

## Ordenamiento de Memoria

La biblioteca implementa cuatro ordenamientos del modelo de memoria C++11:

| Ordenamiento | Semántica |
|--------------|-----------|
| **Relaxed** | Solo atomicidad. Sin restricciones de sincronización u orden. |
| **Acquire** | Las lecturas y escrituras posteriores no pueden reordenarse antes de esta carga. Se empareja con stores Release. |
| **Release** | Las lecturas y escrituras anteriores no pueden reordenarse después de este store. Se empareja con loads Acquire. |
| **AcqRel** | Combina semánticas Acquire y Release. Para operaciones lectura-modificación-escritura. |

### Selección de Ordenamiento

Los métodos por defecto (sin sufijo de ordenamiento) usan:
- Operaciones Load: Relaxed
- Operaciones Store: Relaxed
- Operaciones lectura-modificación-escritura: AcqRel

**Nota:** sync/atomic usa acquire para Load y release para Store (consistencia secuencial en x86). atomix usa Relaxed por defecto, que se mapea a instrucciones distintas en arquitecturas débilmente ordenadas (ej. LDR vs LDAR en ARM64). Use `LoadAcquire`/`StoreRelease` cuando requiera ordenamiento equivalente a sync/atomic.

### Cuándo Usar Cada Ordenamiento

| Caso de Uso | Ordenamiento | Razón |
|-------------|--------------|-------|
| Contadores de estadísticas | Relaxed | Sin necesidad de sincronización; consistencia eventual aceptable |
| Conteo de referencias | AcqRel | Asegura visibilidad del estado del objeto antes de desasignación |
| Flags productor-consumidor | Release/Acquire | El productor libera datos, el consumidor adquiere |
| Adquisición de spinlock | Acquire | Lecturas de sección crítica deben ver escrituras previas |
| Liberación de spinlock | Release | Escrituras de sección crítica deben completarse antes de desbloquear |
| Sequence locks | AcqRel | Ambas direcciones necesitan ordenamiento |

## Tipos

### Tipos de Valor

| Tipo | Tamaño | Descripción |
|------|--------|-------------|
| `Bool` | 4 bytes | Booleano atómico (respaldado por uint32) |
| `Int32`, `Uint32` | 4 bytes | Enteros de 32 bits |
| `Int64`, `Uint64` | 8 bytes | Enteros de 64 bits |
| `Uintptr` | 8 bytes | Entero tamaño de puntero |
| `Pointer[T]` | 8 bytes | Puntero atómico genérico |
| `Int128`, `Uint128` | 16 bytes | Enteros de 128 bits (requiere alineación de 16 bytes) |

### Tipos con Padding

Las variantes con padding (`Int64Padded`, `Uint64Padded`, etc.) ocupan una línea de caché completa (64 bytes) para prevenir compartición falsa cuando múltiples variables atómicas son accedidas por diferentes núcleos de CPU.

```go
// Sin padding: las variables pueden compartir línea de caché, causando contención
var a, b atomix.Int64  // Pueden ser adyacentes en memoria

// Con padding: cada variable ocupa su propia línea de caché
var a, b atomix.Int64Padded  // Separación de 64 bytes garantizada
```

## Operaciones

| Operación | Retorna | Descripción |
|-----------|---------|-------------|
| `Load` | valor | Lectura atómica |
| `Store` | — | Escritura atómica |
| `Swap` | valor antiguo | Intercambio atómico |
| `CompareAndSwap` | bool | Retorna true si el intercambio ocurrió |
| `CompareExchange` | valor antiguo | Retorna valor previo independientemente del resultado |
| `Add`, `Sub` | valor nuevo | Aritmética atómica |
| `Inc`, `Dec` | valor nuevo | Incremento/decremento atómico de 1 |
| `And`, `Or`, `Xor` | valor antiguo | Operaciones bit a bit atómicas |
| `Max`, `Min` | valor antiguo | Máximo/mínimo atómico |

**Semántica de valores de retorno:** Add/Sub/Inc/Dec retornan el valor **nuevo** (como sync/atomic). Swap/And/Or/Xor/Max/Min retornan el valor **antiguo**.

### CompareAndSwap vs CompareExchange

```go
// CompareAndSwap: retorna éxito/fallo
if v.CompareAndSwap(old, new) {
    // Éxito
}

// CompareExchange: retorna valor previo (permite bucles CAS sin Load separado)
for {
    old := v.Load()
    new := transform(old)
    if v.CompareExchange(old, new) == old {
        break  // Éxito
    }
}
```

## API de Punteros

Para interoperación con regiones mapeadas en memoria, memoria compartida, o anillos io_uring:

```go
var flags int32

atomix.Relaxed.StoreInt32(&flags, 1)
val := atomix.Acquire.LoadInt32(&flags)
atomix.Release.CompareAndSwapInt32(&flags, 0, 1)
```

La API de punteros opera sobre `*int32`, `*int64`, etc., crudos en lugar de tipos wrapper. Es útil cuando las variables atómicas no pueden usar tipos wrapper (ej., campos en estructuras compartidas con el kernel).

## Operaciones de 128 bits

Las operaciones atómicas de 128 bits requieren alineación de 16 bytes. Usar helpers de colocación para memoria compartida:

```go
buf := make([]byte, 32)
_, ptr := atomix.PlaceAlignedUint128(buf, 0)
ptr.Store(lo, hi)

var v atomix.Uint128  // El tipo asegura alineación
v.Store(lo, hi)
```

| Arquitectura | Implementación 128 bits |
|--------------|-------------------------|
| amd64 | `LOCK CMPXCHG16B` |
| arm64 | `LDXP/STXP` (defecto) o `CASP` (`-tags=lse2`) |
| riscv64, loong64 | Emulación spinlock (LL/SC en los 64 bits bajos) |

**Nota:** Los atómicos de 128 bits son principalmente útiles para patrones de CAS de doble palabra (ej., estructuras de datos lock-free con contadores de versión).

## Implementación por Arquitectura

### x86-64 (TSO)

x86-64 proporciona Total Store Ordering (TSO), un modelo de memoria fuerte donde:
- Todas las cargas tienen semántica acquire implícita
- Todos los stores tienen semántica release implícita
- El ordenamiento store-load requiere barrera explícita (MFENCE) o instrucción bloqueada

Consecuentemente, todas las variantes de ordenamiento compilan a código máquina idéntico en x86-64. El beneficio principal del ordenamiento explícito en x86-64 es documentación y portabilidad.

| Operación | Instrucción | Notas |
|-----------|-------------|-------|
| Load | `MOV` | Acceso a memoria plano |
| Store | `MOV` | Acceso a memoria plano |
| Add | `LOCK XADD` | Retorna valor antiguo |
| Swap | `XCHG` | LOCK implícito |
| CAS | `LOCK CMPXCHG` | |
| And/Or/Xor | bucle `LOCK CMPXCHG` | Retorna valor antiguo via bucle CAS |
| CAS128 | `LOCK CMPXCHG16B` | |

Load y Store están implementados en Go puro para inlining del compilador.

### ARM64 (Débilmente Ordenado)

ARM64 tiene un modelo de memoria débil que requiere instrucciones de ordenamiento explícitas. LSE (Large System Extensions) proporciona instrucciones atómicas con sufijos de ordenamiento:

**Significado de sufijos:** Sin sufijo = Relaxed, `A` = Acquire, `L` = Release, `AL` = Acquire-Release

| Operación | Relaxed | Acquire | Release | AcqRel |
|-----------|---------|---------|---------|--------|
| Load | `LDR` | `LDAR` | — | — |
| Store | `STR` | — | `STLR` | — |
| Add | `LDADD` | `LDADDA` | `LDADDL` | `LDADDAL` |
| CAS | `CAS` | `CASA` | `CASL` | `CASAL` |
| Swap | `SWP` | `SWPA` | `SWPL` | `SWPAL` |
| And | `LDCLR`† | `LDCLRA` | `LDCLRL` | `LDCLRAL` |
| Or | `LDSET` | `LDSETA` | `LDSETL` | `LDSETAL` |
| Xor | `LDEOR` | `LDEORA` | `LDEORL` | `LDEORAL` |

† `LDCLR` limpia bits (AND con complemento). Para implementar `And(mask)`, pasar `~mask`.

Load/store relajados están implementados en Go puro para inlining. Otros ordenamientos usan ensamblador con instrucciones LSE.

#### Operaciones de 128 bits

| Tag de Build | Instrucciones | Hardware Objetivo |
|--------------|---------------|-------------------|
| (defecto) | `LDXP/STXP` (bucle LL/SC) | Todo ARMv8+ |
| `-tags=lse2` | `CASP` (instrucción única) | ARMv8.4+ con LSE2 |

LL/SC (Load-Link/Store-Conditional) reintenta en contención. CASP proporciona atomicidad de instrucción única pero requiere hardware más nuevo.

### RISC-V 64 bits

RISC-V RVWMO (Ordenamiento de Memoria Débil) usa instrucciones fence explícitas:

| Operación | Implementación |
|-----------|----------------|
| Load Relaxed | `LD` |
| Load Acquire | `LD` + `FENCE R,RW` |
| Store Relaxed | `SD` |
| Store Release | `FENCE RW,W` + `SD` |
| RMW | Instrucciones `AMO` con modificadores `.aq`/`.rl` |

Las operaciones de 128 bits usan emulación basada en spinlock.

### LoongArch 64 bits

LoongArch usa instrucciones DBAR (barrera de datos):

| Operación | Implementación |
|-----------|----------------|
| Load Relaxed | `LD.D` |
| Load Acquire | `LD.D` + `DBAR` |
| Store Relaxed | `ST.D` |
| Store Release | `DBAR` + `ST.D` |
| RMW | Instrucciones `AM*_DB` |

Las operaciones de 128 bits usan emulación basada en spinlock.

### Fallback

Las arquitecturas no soportadas usan `sync/atomic`, que proporciona consistencia secuencial. Las operaciones de 128 bits en arquitecturas fallback **no son atómicas** (dos operaciones de 64 bits separadas).

## Fundamento del Diseño

### Ordenamiento de Memoria Explícito

1. **Selección de instrucciones en arquitecturas débiles**: ARM64/RISC-V seleccionan instrucciones diferentes según los requisitos de ordenamiento
2. **Documentación**: El sufijo de ordenamiento documenta la intención de sincronización
3. **Portabilidad**: El código especifica explícitamente requisitos en lugar de depender de garantías específicas de arquitectura
4. **Corrección**: Hace las decisiones de ordenamiento de memoria explícitas y revisables

### Comparación con sync/atomic

sync/atomic proporciona consistencia secuencial, que es:
- **Suficiente** para la mayoría de casos de uso
- **Portable** a través de todas las arquitecturas
- **Simple** de razonar

Usar atomix cuando:
- Construir estructuras de datos lock-free
- Interoperar con kernel o interfaces de hardware (io_uring, memoria compartida)
- Portar código C/C++ con ordenamiento de memoria explícito
- Apuntar a ARM64/RISC-V donde el ordenamiento explícito controla la selección de instrucciones

## Soporte de Plataformas

| Plataforma | Implementación |
|------------|----------------|
| linux/amd64 | Ensamblador nativo |
| linux/arm64 | Ensamblador nativo con LSE |
| linux/riscv64 | Ensamblador nativo (128 bits emulado) |
| linux/loong64 | Ensamblador nativo (128 bits emulado) |
| darwin/amd64, darwin/arm64 | Ensamblador nativo |
| freebsd/amd64, freebsd/arm64 | Ensamblador nativo |
| Otros | Fallback a sync/atomic |

## Intrínsecos del Compilador

atomix proporciona un compilador Go personalizado que emite instrucciones atómicas inline en lugar de llamadas a funciones. Esto transforma las llamadas a funciones en instrucciones CPU únicas, eliminando la sobrecarga de llamadas.

### Inicio Rápido

```bash
# Instalar el compilador con intrínsecos personalizados
make install-compiler

# Compilar con intrínsecos
make build

# Probar con intrínsecos
make test

# Verificar que los intrínsecos se aplican
make verify
```

### Qué Hace el Compilador

El compilador personalizado añade operaciones SSA para intrínsecos de atomix:

| Operación | x86-64 | ARM64 |
|-----------|--------|-------|
| Load (Relaxed) | `MOV` | `LDR` |
| Load (Acquire) | `MOV` | `LDAR` |
| Store (Relaxed) | `MOV` | `STR` |
| Store (Release) | `MOV` | `STLR` |
| Add (AcqRel) | `LOCK XADD` | `LDADDAL` |
| CAS | `LOCK CMPXCHG` | `CASAL` |

**Optimización x86-64 TSO:** Los stores Release usan `MOV` plano en lugar de `XCHG`, aprovechando el Total Store Ordering de x86-64 que proporciona semántica release implícita para todos los stores.

### Configuración Manual del Compilador

Si prefieres configuración manual sobre el Makefile:

```bash
# Clonar el compilador con intrínsecos
git clone --branch atomix https://github.com/hayabusa-cloud/go.git ~/github.com/go

# Compilar el compilador
cd ~/github.com/go/src && ./make.bash

# Usar para atomix
GOROOT=~/github.com/go ~/github.com/go/bin/go build ./...
```

Ver [intrinsics.md](./intrinsics.md) para documentación detallada de implementación.

## Licencia

MIT — ver [LICENSE](./LICENSE).

©2026 [Hayabusa Cloud Co., Ltd.](https://code.hybscloud.com/)
