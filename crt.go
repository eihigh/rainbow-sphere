//go:build ignore

package main

//kage:unit pixels

// Based on https://github.com/XorDev/1PassBlur
// https://www.shadertoy.com/view/DtscRf

// Please adjust the clarity and brightness with these consts
const base = 0.5
const glow = 1.5
const radius = 16.0
const samples = 32.0

func Fragment(dstPos vec4, srcPos vec2, color vec4) vec4 {
	blur := vec4(0)
	weights := 0.0
	scale := radius / sqrt(samples)
	offset := scale * normalize(fract(cos(srcPos*mat2(195, 174, 286, 183))*742)-0.5) // random value
	rot := mat2(-0.7373688, -0.6754904, 0.6754904, -0.7373688)
	for i := 0.0; i < samples; i += 1 {
		// rotate by golden angle
		offset *= rot
		dist := sqrt(i)
		pos := srcPos + offset*dist
		color := imageSrc0At(pos)

		weight := 1.0 / (1 + dist)
		blur += color * weight
		weights += weight
	}
	blur /= weights
	clr := mix(blur*glow, imageSrc0At(srcPos), base)

	rgb := clr.rgb
	rgb = clamp(mix(rgb, rgb*rgb, 0.4), 0, 1)

	// vignette
	uv := (srcPos - imageSrc0Origin()) / imageSrc0Size()
	vig := 40 * uv.x * uv.y * (1 - uv.x) * (1 - uv.y)
	rgb *= vec3(pow(vig, 0.3))
	rgb *= vec3(0.95, 1.05, 0.95)

	// scan lines
	n := floor(imageDstSize().y/480) + 1
	rgb *= 1.0 - fract(dstPos.y/n)*0.3

	rgb *= 1.4
	return vec4(rgb, clr.a)
}
