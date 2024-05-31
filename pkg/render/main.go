package render

import "github.com/williamflynt/topolith/pkg/world"

type RenderedReturnType string

type OnRenderFunction = func(world.World)
type UnhookFunction = func()

type Renderer interface {
	Render(w world.World) ([]byte, RenderedReturnType, error)
	OnRender(f OnRenderFunction) UnhookFunction
}

// TODO: Implement a renderer.
