//go:build linux && !android && !wayland

package wgpuext_glfw // import "github.com/birdbrain-holdings/go-webgpu/wgpuext/glfw"

import (
	"unsafe"

	"github.com/birdbrain-holdings/go-webgpu/wgpu"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func GetSurfaceDescriptor(w *glfw.Window) *wgpu.SurfaceDescriptor {
	return &wgpu.SurfaceDescriptor{
		XlibWindow: &wgpu.SurfaceDescriptorFromXlibWindow{
			Display: unsafe.Pointer(glfw.GetX11Display()),
			Window:  uint32(w.GetX11Window()),
		},
	}
}
