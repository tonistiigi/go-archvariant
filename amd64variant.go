//go:build amd64
// +build amd64

package archvariant

import (
	"fmt"
)

func cpuid(ax, cx uint32) (eax, ebx, ecx, edx uint32)
func xgetbv() (eax uint32)

const (
	sse3    = 0
	ssse3   = 9
	cx16    = 13
	sse4_1  = 19
	sse4_2  = 20
	popcnt  = 23
	fma     = 12
	movbe   = 22
	xsave   = 26
	osxsave = 27
	avx     = 28
	f16c    = 29

	v2Features = 1<<sse3 | 1<<ssse3 | 1<<cx16 | 1<<sse4_1 | 1<<sse4_2 | 1<<popcnt
	v3Features = v2Features | 1<<fma | 1<<movbe | 1<<xsave | 1<<osxsave | 1<<avx | 1<<f16c
)

const (
	bmi1      = 3
	avx2      = 5
	bmi2      = 8
	avx512_f  = 16
	avx512_dq = 17
	avx512_cd = 28
	avx512_bw = 30
	avx512_vl = 31

	v3ExtFeaturesBX = 1<<bmi1 | 1<<avx2 | 1<<bmi2
	v4ExtFeaturesBX = v3ExtFeaturesBX | 1<<avx512_f | 1<<avx512_dq | 1<<avx512_cd | 1<<avx512_bw | 1<<avx512_vl
)

const (
	lahf_lm = 0
	abm     = 5

	v2ExtFeatureCX = 1 << lahf_lm
	v3ExtFeatureCX = v2ExtFeatureCX | 1<<abm
)

const (
	// XCR0
	xmm       = 1
	ymm       = 2
	opmask    = 5
	zmm_hi16  = 6
	zmm_hi256 = 7

	v3OSSupport = 1<<xmm | 1<<ymm
	v4OSSupport = v3OSSupport | 1<<opmask | 1<<zmm_hi16 | 1<<zmm_hi256
)

func detectVersion() int {
	// highest basic calling parameter
	ax, _, _, _ := cpuid(0, 0)
	if ax < 7 {
		return 1
	}

	// highest extended calling parameter
	ax, _, _, _ = cpuid(0x80000000, 0)
	if ax < 0x80000001 {
		return 1
	}

	version := 4 // initialize to highest version

	// feature bits
	_, _, cx, _ := cpuid(1, 0)
	if cx&v3Features != v3Features {
		version = 2
		if cx&v2Features != v2Features {
			return 1
		}
	}

	// extended features
	_, bx, _, _ := cpuid(7, 0)
	if version == 4 {
		if bx&v4ExtFeaturesBX != v4ExtFeaturesBX {
			version = 3
		}
	}
	if version == 3 {
		if bx&v3ExtFeaturesBX != v3ExtFeaturesBX {
			version = 2
		}
	}

	// extended processor info and feature bits
	_, _, cx, _ = cpuid(0x80000001, 0)
	if version >= 3 {
		if cx&v3ExtFeatureCX != v3ExtFeatureCX {
			version = 2
		}
	}

	if version == 2 {
		if cx&v2ExtFeatureCX != v2ExtFeatureCX {
			return 1
		}
	}

	if version >= 3 {
		ax = xgetbv()
		if version == 4 {
			if !osAVX512Supported(ax) {
				version = 3
			}
		}
		if ax&v3OSSupport != v3OSSupport {
			version = 2
		}
	}

	return version
}

func AMD64Variant() string {
	return "v" + fmt.Sprintf("%d", detectVersion())
}