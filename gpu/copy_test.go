package gpu

import (
	"nimble-cube/core"
	"testing"
)

func TestCopy(t *testing.T) {
	LockCudaThread()

	cell := 1e-9
	mesh := core.NewMesh(2, 4, 8, cell, cell, cell)
	N := mesh.NCell()
	F := 100
	a := core.MakeChan("a", "", mesh)
	b := MakeChan("b", "", mesh)
	c := core.MakeChan("c", "", mesh)

	up := NewUploader(a.NewReader(), b)
	down := NewDownloader(b.NewReader(), c)

	go up.Run()
	go down.Run()

	go func() {
		for f := 0; f < F; f++ {
			list := a.WriteNext(N)
			for i := range list {
				list[i] = float32(i)
			}
			a.WriteDone()
		}
	}()

	C := c.NewReader()
	for f := 0; f < F; f++ {
		list := C.ReadNext(N)
		for i := range list {
			if list[i] != float32(i) {
				t.Error("expected:", float32(i), "got:", list[i])
			}
		}
		C.ReadDone()
	}
}
