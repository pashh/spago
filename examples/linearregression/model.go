// Copyright 2019 spaGO Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"io"
	"log"
	"saientist.dev/spago/pkg/mat"
	"saientist.dev/spago/pkg/ml/ag"
	"saientist.dev/spago/pkg/ml/nn"
)

type LinearRegression struct {
	W *nn.Param `type:"weights"`
}

func NewLinearRegression(in, out int) *LinearRegression {
	return &LinearRegression{
		W: nn.NewParam(mat.NewEmptyDense(out, in)),
	}
}

func (m *LinearRegression) ForEachParam(callback func(param *nn.Param)) {
	nn.ForEachParam(m, callback)
}

func (m *LinearRegression) Serialize(w io.Writer) (int, error) {
	return nn.Serialize(m, w)
}

func (m *LinearRegression) Deserialize(r io.Reader) (int, error) {
	return nn.Deserialize(m, r)
}

type Processor struct {
	opt   []interface{}
	model *LinearRegression
	g     *ag.Graph
	w     ag.Node
}

func (m *LinearRegression) NewProc(g *ag.Graph, opt ...interface{}) nn.Processor {
	p := &Processor{
		model: m,
		opt:   opt,
		g:     g,
		w:     g.NewWrap(m.W),
	}
	p.init(opt)
	return p
}

func (p *Processor) init(opt []interface{}) {
	if len(opt) > 0 {
		log.Fatal("linearregression: invalid init options")
	}
}

func (p *Processor) Model() nn.Model {
	return p.model
}

func (p *Processor) Graph() *ag.Graph {
	return p.g
}

func (p *Processor) RequiresFullSeq() bool {
	return false
}

func (p *Processor) Reset() {
	p.init(p.opt)
}

func (p *Processor) Forward(xs ...ag.Node) []ag.Node {
	ys := make([]ag.Node, len(xs))
	for i, x := range xs {
		ys[i] = nn.Linear(p.g, p.w, x)
	}
	return ys
}