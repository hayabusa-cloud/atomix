# atomix

[![Go Reference](https://pkg.go.dev/badge/code.hybscloud.com/atomix.svg)](https://pkg.go.dev/code.hybscloud.com/atomix)
[![Go Report Card](https://goreportcard.com/badge/github.com/hayabusa-cloud/atomix)](https://goreportcard.com/report/github.com/hayabusa-cloud/atomix)
[![Codecov](https://codecov.io/gh/hayabusa-cloud/atomix/graph/badge.svg)](https://codecov.io/gh/hayabusa-cloud/atomix)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

**Langues:** [English](README.md) | [简体中文](README.zh-CN.md) | [日本語](README.ja.md) | [Español](README.es.md) | Français

Opérations atomiques avec ordonnancement mémoire explicite pour Go.

## Présentation

Le package `sync/atomic` de Go fournit des opérations atomiques avec cohérence séquentielle. Cette bibliothèque expose les ordonnancements du modèle mémoire C++11/C11 (Relaxed, Acquire, Release, AcqRel) via des implémentations spécifiques à l'architecture.

```go
import "code.hybscloud.com/atomix"

var counter atomix.Int64

// API basée sur méthodes avec suffixe d'ordonnancement
counter.AddRelaxed(1)    // Relaxed : aucune synchronisation
counter.Add(1)           // AcqRel : ordonnancement sûr par défaut

// API basée sur pointeurs pour mémoire brute
var flags int32
atomix.Relaxed.StoreInt32(&flags, 1)
val := atomix.Acquire.LoadInt32(&flags)
```

## Installation

```bash
go get code.hybscloud.com/atomix
```

**Prérequis :** Go 1.25+

## Ordonnancement Mémoire

La bibliothèque implémente quatre ordonnancements du modèle mémoire C++11 :

| Ordonnancement | Sémantique |
|----------------|------------|
| **Relaxed** | Atomicité uniquement. Aucune contrainte de synchronisation ou d'ordre. |
| **Acquire** | Les lectures et écritures ultérieures ne peuvent pas être réordonnées avant ce load. S'utilise en paire avec les stores Release. |
| **Release** | Les lectures et écritures antérieures ne peuvent pas être réordonnées après ce store. S'utilise en paire avec les loads Acquire. |
| **AcqRel** | Combine les sémantiques Acquire et Release. Pour les opérations lecture-modification-écriture. |

### Sélection de l'Ordonnancement

Les méthodes par défaut (sans suffixe d'ordonnancement) utilisent :
- Opérations Load : Relaxed
- Opérations Store : Relaxed
- Opérations lecture-modification-écriture : AcqRel

**Note :** sync/atomic utilise acquire pour Load et release pour Store (cohérence séquentielle sur x86). atomix utilise Relaxed par défaut pour des performances maximales sur les architectures faiblement ordonnées. Utilisez `LoadAcquire`/`StoreRelease` pour un ordonnancement équivalent à sync/atomic.

### Quand Utiliser Chaque Ordonnancement

| Cas d'Usage | Ordonnancement | Raison |
|-------------|----------------|--------|
| Compteurs de statistiques | Relaxed | Pas de synchronisation nécessaire ; cohérence éventuelle acceptable |
| Comptage de références | AcqRel | Assure la visibilité de l'état de l'objet avant désallocation |
| Flags producteur-consommateur | Release/Acquire | Le producteur libère les données, le consommateur acquiert |
| Acquisition de spinlock | Acquire | Les lectures de section critique doivent voir les écritures précédentes |
| Libération de spinlock | Release | Les écritures de section critique doivent se terminer avant le déverrouillage |
| Verrous de séquence | AcqRel | Les deux directions nécessitent un ordonnancement |

## Types

### Types Valeur

| Type | Taille | Description |
|------|--------|-------------|
| `Bool` | 4 octets | Booléen atomique (soutenu par uint32) |
| `Int32`, `Uint32` | 4 octets | Entiers 32 bits |
| `Int64`, `Uint64` | 8 octets | Entiers 64 bits |
| `Uintptr` | 8 octets | Entier taille pointeur |
| `Pointer[T]` | 8 octets | Pointeur atomique générique |
| `Int128`, `Uint128` | 16 octets | Entiers 128 bits (requiert alignement 16 octets) |

### Types avec Padding

Les variantes avec padding (`Int64Padded`, `Uint64Padded`, etc.) occupent une ligne de cache complète (64 octets) pour éviter le faux partage lorsque plusieurs variables atomiques sont accédées par différents cœurs CPU.

```go
// Sans padding : les variables peuvent partager une ligne de cache, causant de la contention
var a, b atomix.Int64  // Peuvent être adjacentes en mémoire

// Avec padding : chaque variable occupe sa propre ligne de cache
var a, b atomix.Int64Padded  // Séparation de 64 octets garantie
```

## Opérations

| Opération | Retourne | Description |
|-----------|----------|-------------|
| `Load` | valeur | Lecture atomique |
| `Store` | — | Écriture atomique |
| `Swap` | ancienne valeur | Échange atomique |
| `CompareAndSwap` | bool | Retourne true si l'échange a eu lieu |
| `CompareExchange` | ancienne valeur | Retourne la valeur précédente quel que soit le résultat |
| `Add`, `Sub` | nouvelle valeur | Arithmétique atomique |
| `Inc`, `Dec` | nouvelle valeur | Incrément/décrément atomique de 1 |
| `And`, `Or`, `Xor` | ancienne valeur | Opérations bit à bit atomiques |
| `Max`, `Min` | ancienne valeur | Maximum/minimum atomique |

**Sémantique des valeurs de retour :** Add/Sub/Inc/Dec retournent la **nouvelle** valeur (comme sync/atomic). Swap/And/Or/Xor/Max/Min retournent l'**ancienne** valeur.

### CompareAndSwap vs CompareExchange

```go
// CompareAndSwap : retourne succès/échec
if v.CompareAndSwap(old, new) {
    // Succès
}

// CompareExchange : retourne la valeur précédente (permet des boucles CAS sans Load séparé)
for {
    old := v.Load()
    new := transform(old)
    if v.CompareExchange(old, new) == old {
        break  // Succès
    }
}
```

## API Pointeurs

Pour l'interopération avec des régions mappées en mémoire, mémoire partagée, ou anneaux io_uring :

```go
var flags int32

atomix.Relaxed.StoreInt32(&flags, 1)
val := atomix.Acquire.LoadInt32(&flags)
atomix.Release.CompareAndSwapInt32(&flags, 0, 1)
```

L'API pointeurs opère sur des `*int32`, `*int64`, etc., bruts plutôt que des types wrapper. Utile quand les variables atomiques ne peuvent pas utiliser de types wrapper (ex., champs dans des structures partagées avec le noyau).

## Opérations 128 bits

Les opérations atomiques 128 bits requièrent un alignement de 16 octets. Utiliser les helpers de placement pour la mémoire partagée :

```go
buf := make([]byte, 32)
_, ptr := atomix.PlaceAlignedUint128(buf, 0)
ptr.Store(lo, hi)

var v atomix.Uint128  // Le type assure l'alignement
v.Store(lo, hi)
```

| Architecture | Implémentation 128 bits |
|--------------|-------------------------|
| amd64 | `LOCK CMPXCHG16B` |
| arm64 | `LDXP/STXP` (défaut) ou `CASP` (`-tags=lse2`) |
| riscv64, loong64 | Émulation spinlock (LL/SC sur les 64 bits bas) |

**Note :** Les atomiques 128 bits sont principalement utiles pour les patterns de CAS double mot (ex., structures de données lock-free avec compteurs de version).

## Implémentation par Architecture

### x86-64 (TSO)

x86-64 fournit Total Store Ordering (TSO), un modèle mémoire fort où :
- Toutes les charges ont une sémantique acquire implicite
- Tous les stores ont une sémantique release implicite
- L'ordonnancement store-load requiert une barrière explicite (MFENCE) ou instruction verrouillée

Par conséquent, toutes les variantes d'ordonnancement compilent vers un code machine identique sur x86-64. Le bénéfice principal de l'ordonnancement explicite sur x86-64 est la documentation et la portabilité.

| Opération | Instruction | Notes |
|-----------|-------------|-------|
| Load | `MOV` | Accès mémoire simple |
| Store | `MOV` | Accès mémoire simple |
| Add | `LOCK XADD` | Retourne ancienne valeur |
| Swap | `XCHG` | LOCK implicite |
| CAS | `LOCK CMPXCHG` | |
| And/Or/Xor | boucle `LOCK CMPXCHG` | Retourne l'ancienne valeur via boucle CAS |
| CAS128 | `LOCK CMPXCHG16B` | |

Load et Store sont implémentés en Go pur pour l'inlining du compilateur.

### ARM64 (Faiblement Ordonné)

ARM64 a un modèle mémoire faible nécessitant des instructions d'ordonnancement explicites. LSE (Large System Extensions) fournit des instructions atomiques avec suffixes d'ordonnancement :

**Signification des suffixes :** Sans suffixe = Relaxed, `A` = Acquire, `L` = Release, `AL` = Acquire-Release

| Opération | Relaxed | Acquire | Release | AcqRel |
|-----------|---------|---------|---------|--------|
| Load | `LDR` | `LDAR` | — | — |
| Store | `STR` | — | `STLR` | — |
| Add | `LDADD` | `LDADDA` | `LDADDL` | `LDADDAL` |
| CAS | `CAS` | `CASA` | `CASL` | `CASAL` |
| Swap | `SWP` | `SWPA` | `SWPL` | `SWPAL` |
| And | `LDCLR`† | `LDCLRA` | `LDCLRL` | `LDCLRAL` |
| Or | `LDSET` | `LDSETA` | `LDSETL` | `LDSETAL` |
| Xor | `LDEOR` | `LDEORA` | `LDEORL` | `LDEORAL` |

† `LDCLR` efface les bits (AND avec complément). Pour implémenter `And(mask)`, passer `~mask`.

Les load/store relâchés sont implémentés en Go pur pour l'inlining. Les autres ordonnancements utilisent de l'assembleur avec les instructions LSE.

#### Opérations 128 bits

| Tag de Build | Instructions | Matériel Cible |
|--------------|--------------|----------------|
| (défaut) | `LDXP/STXP` (boucle LL/SC) | Tout ARMv8+ |
| `-tags=lse2` | `CASP` (instruction unique) | ARMv8.4+ avec LSE2 |

LL/SC (Load-Link/Store-Conditional) réessaie en cas de contention. CASP fournit une atomicité en une seule instruction mais nécessite du matériel plus récent.

### RISC-V 64 bits

RISC-V RVWMO (Ordonnancement Mémoire Faible) utilise des instructions fence explicites :

| Opération | Implémentation |
|-----------|----------------|
| Load Relaxed | `LD` |
| Load Acquire | `LD` + `FENCE R,RW` |
| Store Relaxed | `SD` |
| Store Release | `FENCE RW,W` + `SD` |
| RMW | Instructions `AMO` avec modificateurs `.aq`/`.rl` |

Les opérations 128 bits utilisent une émulation basée sur spinlock.

### LoongArch 64 bits

LoongArch utilise des instructions DBAR (barrière de données) :

| Opération | Implémentation |
|-----------|----------------|
| Load Relaxed | `LD.D` |
| Load Acquire | `LD.D` + `DBAR` |
| Store Relaxed | `ST.D` |
| Store Release | `DBAR` + `ST.D` |
| RMW | Instructions `AM*_DB` |

Les opérations 128 bits utilisent une émulation basée sur spinlock.

### Fallback

Les architectures non supportées utilisent `sync/atomic`, qui fournit la cohérence séquentielle. Les opérations 128 bits sur les architectures fallback **ne sont pas atomiques** (deux opérations 64 bits séparées).

## Fondement de la Conception

### Pourquoi l'Ordonnancement Mémoire Explicite ?

1. **Performance sur architectures faibles** : ARM64/RISC-V peuvent utiliser des instructions plus faibles (plus rapides) quand l'ordonnancement complet n'est pas nécessaire
2. **Documentation** : Le suffixe d'ordonnancement documente l'intention de synchronisation
3. **Portabilité** : Le code spécifie explicitement les exigences plutôt que de dépendre de garanties spécifiques à l'architecture
4. **Correction** : Rend les décisions d'ordonnancement mémoire explicites et vérifiables

### Pourquoi Ne Pas Simplement Utiliser sync/atomic ?

sync/atomic fournit la cohérence séquentielle, qui est :
- **Suffisante** pour la plupart des cas d'usage
- **Portable** à travers toutes les architectures
- **Simple** à raisonner

Utiliser atomix quand :
- Construction de structures de données lock-free haute performance
- Interopération avec le noyau ou interfaces matérielles (io_uring, mémoire partagée)
- Portage de code C/C++ avec ordonnancement mémoire explicite
- Ciblage ARM64/RISC-V où l'ordonnancement plus faible apporte un bénéfice mesurable

## Support des Plateformes

| Plateforme | Implémentation |
|------------|----------------|
| linux/amd64 | Assembleur natif |
| linux/arm64 | Assembleur natif avec LSE |
| linux/riscv64 | Assembleur natif (128 bits émulé) |
| linux/loong64 | Assembleur natif (128 bits émulé) |
| darwin/amd64, darwin/arm64 | Assembleur natif |
| freebsd/amd64, freebsd/arm64 | Assembleur natif |
| Autres | Repli sur sync/atomic |

## Intrinsèques du Compilateur

Pour exploiter pleinement les performances, atomix peut être intégré au compilateur Go pour émettre des instructions atomiques en ligne, éliminant le surcoût des appels de fonction. Voir [intrinsics.md](./intrinsics.md) pour l'approche d'implémentation.

## Licence

MIT — voir [LICENSE](./LICENSE).

©2026 [Hayabusa Cloud Co., Ltd.](https://code.hybscloud.com/)
