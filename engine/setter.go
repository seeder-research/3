package engine

import (
	"code.google.com/p/mx3/cuda"
	"code.google.com/p/mx3/data"
)

// quantity that is not stored, but can output to (set) a buffer
type setterQuant struct {
	autosave
	setFn func(dst *data.Slice, good bool) // calculates quantity and stores in dst
}

// constructor
func setter(nComp int, m *data.Mesh, name, unit string, setFunc func(dst *data.Slice, good bool)) setterQuant {
	return setterQuant{newAutosave(nComp, name, unit, m), setFunc}
}

// calculate quantity, save in dst, notify output may be needed
func (b *setterQuant) set(dst *data.Slice, goodstep bool) {
	b.setFn(dst, goodstep)
	if goodstep && b.needSave() {
		goSaveCopy(b.autoFname(), dst, Time)
		b.saved()
	}
}

// get the quantity, recycle will be true (q needs to be recycled)
func (b *setterQuant) GetGPU() (q *data.Slice, recycle bool) {
	buffer := cuda.GetBuffer(b.nComp, b.mesh)
	b.set(buffer, false)
	return buffer, true // must recycle
}

// get the quantity, recycle will be true, q will be on GPU
func (b *setterQuant) Get() (*data.Slice, bool) {
	return b.GetGPU()
}

func (p *setterQuant) Save() {
	save(p)
}

func (p *setterQuant) SaveAs(fname string) {
	saveAs(p, fname)
}
