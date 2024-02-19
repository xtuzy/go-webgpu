package wgpu

/*

#include <stdlib.h>
#include "./lib/wgpu.h"

static inline WGPUSurfaceTexture gowebgpu_surface_get_current_texture(WGPUSurface surface) {
	WGPUSurfaceTexture ref;
	wgpuSurfaceGetCurrentTexture(surface, &ref);
	return ref;
}

static inline void gowebgpu_surface_configure(WGPUSurface surface, WGPUSurfaceConfiguration descriptor) {
	wgpuSurfaceConfigure(surface, &descriptor);
}

*/
import "C"
import (
	"fmt"
	"unsafe"
)

type Surface struct {
	ref       C.WGPUSurface
	deviceRef C.WGPUDevice
}

func (p *Surface) GetPreferredFormat(adapter *Adapter) TextureFormat {
	format := C.wgpuSurfaceGetPreferredFormat(p.ref, adapter.ref)
	return TextureFormat(format)
}

type SurfaceCapabilities struct {
	Formats      []TextureFormat
	PresentModes []PresentMode
	AlphaModes   []CompositeAlphaMode
}

func (p *Surface) GetCurrentTexture() (*Texture, error) {
	ref := C.gowebgpu_surface_get_current_texture(p.ref)
	if ref.status != C.WGPUSurfaceGetCurrentTextureStatus_Success {
		return nil, &Error{
			Type:    ErrorType_Validation,
			Message: fmt.Sprintf("failed to get current texture: %d", ref.status),
		}
	}
	return &Texture{
		ref:       ref.texture,
		deviceRef: p.deviceRef,
	}, nil
}

func (p *Surface) Configure(device *Device, config *SwapChainDescriptor) {
	p.deviceRef = device.ref
	var cfg C.WGPUSurfaceConfiguration
	cfg.device = device.ref
	cfg.format = C.WGPUTextureFormat(config.Format)
	cfg.usage = C.uint32_t(config.Usage)
	cfg.width = C.uint32_t(config.Width)
	cfg.height = C.uint32_t(config.Height)
	cfg.presentMode = C.WGPUPresentMode(config.PresentMode)

	viewFormatCount := len(config.ViewFormats)
	if viewFormatCount > 0 {
		viewFormats := C.malloc(C.size_t(unsafe.Sizeof(C.WGPUTextureFormat(0))) * C.size_t(viewFormatCount))
		defer C.free(viewFormats)

		viewFormatsSlice := unsafe.Slice((*TextureFormat)(viewFormats), viewFormatCount)
		copy(viewFormatsSlice, config.ViewFormats)

		cfg.viewFormatCount = C.size_t(viewFormatCount)
		cfg.viewFormats = (*C.WGPUTextureFormat)(viewFormats)
	} else {
		cfg.viewFormatCount = 0
		cfg.viewFormats = nil
	}
	cfg.alphaMode = C.WGPUCompositeAlphaMode(config.AlphaMode)

	C.gowebgpu_surface_configure(p.ref, cfg)

}

func (p *Surface) GetCapabilities(adapter *Adapter) (ret SurfaceCapabilities) {
	var caps C.WGPUSurfaceCapabilities
	C.wgpuSurfaceGetCapabilities(p.ref, adapter.ref, &caps)

	if caps.alphaModeCount == 0 && caps.formatCount == 0 && caps.presentModeCount == 0 {
		return
	}
	if caps.formatCount > 0 {
		caps.formats = (*C.WGPUTextureFormat)(C.malloc(C.size_t(unsafe.Sizeof(C.WGPUTextureFormat(0))) * caps.formatCount))
		defer C.free(unsafe.Pointer(caps.formats))
	}
	if caps.presentModeCount > 0 {
		caps.presentModes = (*C.WGPUPresentMode)(C.malloc(C.size_t(unsafe.Sizeof(C.WGPUPresentMode(0))) * caps.presentModeCount))
		defer C.free(unsafe.Pointer(caps.presentModes))
	}
	if caps.alphaModeCount > 0 {
		caps.alphaModes = (*C.WGPUCompositeAlphaMode)(C.malloc(C.size_t(unsafe.Sizeof(C.WGPUCompositeAlphaMode(0))) * caps.alphaModeCount))
		defer C.free(unsafe.Pointer(caps.alphaModes))
	}

	C.wgpuSurfaceGetCapabilities(p.ref, adapter.ref, &caps)

	if caps.formatCount > 0 {
		formatsTmp := unsafe.Slice((*TextureFormat)(caps.formats), caps.formatCount)
		ret.Formats = make([]TextureFormat, caps.formatCount)
		copy(ret.Formats, formatsTmp)
	}
	if caps.presentModeCount > 0 {
		presentModesTmp := unsafe.Slice((*PresentMode)(caps.presentModes), caps.presentModeCount)
		ret.PresentModes = make([]PresentMode, caps.presentModeCount)
		copy(ret.PresentModes, presentModesTmp)
	}
	if caps.alphaModeCount > 0 {
		alphaModesTmp := unsafe.Slice((*CompositeAlphaMode)(caps.alphaModes), caps.alphaModeCount)
		ret.AlphaModes = make([]CompositeAlphaMode, caps.alphaModeCount)
		copy(ret.AlphaModes, alphaModesTmp)
	}

	return
}

func (p *Surface) Release() {
	C.wgpuSurfaceRelease(p.ref)
}

func (p *Surface) Present() {
	C.wgpuSurfacePresent(p.ref)
}
