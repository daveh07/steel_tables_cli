// Package models defines data structures for steel section properties.
package models

import "encoding/json"

// SteelProperty defines the structure for a single steel section property.
type SteelProperty struct {
	Section  string      `json:"Section"`
	Grade    int         `json:"Grade"`
	Weight   float64     `json:"Weight"`
	D        float64     `json:"d"`
	Bf       float64     `json:"bf"`
	Tf       float64     `json:"tf"`
	Tw       float64     `json:"tw"`
	R1       interface{} `json:"r1"`
	D1       float64     `json:"d1"`
	Tw1      interface{} `json:"tw__1"`
	Tf1      interface{} `json:"tf__1"`
	Ag       float64     `json:"Ag"`
	Ix       float64     `json:"Ix"`
	Zx       float64     `json:"Zx"`
	Sx       float64     `json:"Sx"`
	Rx       float64     `json:"rx"`
	Iy       float64     `json:"Iy"`
	Zy       float64     `json:"Zy"`
	Sy       float64     `json:"Sy"`
	Ry       float64     `json:"ry"`
	J        float64     `json:"J"`
	Iw       interface{} `json:"Iw"`
	Flange   interface{} `json:"flange"`
	Web      interface{} `json:"web"`
	Kf       interface{} `json:"kf"`
	CNS      interface{} `json:"-"`
	Zex      float64     `json:"Zex"`
	CNS2     interface{} `json:"-"`
	Zey      float64     `json:"Zey"`
	TwoTf    interface{} `json:"2tf"`
	Zy5      float64     `json:"Zy5"`
	TanAlpha float64     `json:"Tan Alpha"`
	AlphaB   interface{} `json:"αb"`
	Fu       interface{} `json:"Fu"`
	R2       interface{} `json:"r2"`
	ZeyD     float64     `json:"ZeyD"`
	In       float64     `json:"In"`
	Ip       float64     `json:"Ip"`
	ZexC     float64     `json:"ZexC"`
	X5       interface{} `json:"x5"`
	Y5       float64     `json:"y5"`
	NL       float64     `json:"nL"`
	PB       float64     `json:"pB"`
	PT       interface{} `json:"pT"`
	Residual string      `json:"Residual"`
	Type     interface{} `json:"Type"`
}

// UnmarshalJSON handles JSON fields with commas in their names.
func (sp *SteelProperty) UnmarshalJSON(data []byte) error {
	type Alias SteelProperty
	aux := &struct{ *Alias }{Alias: (*Alias)(sp)}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	var rawMap map[string]interface{}
	if err := json.Unmarshal(data, &rawMap); err != nil {
		return err
	}

	if val, ok := rawMap["C,N,S"]; ok {
		sp.CNS = val
	}
	if val, ok := rawMap["C,N,S__1"]; ok {
		sp.CNS2 = val
	}

	return nil
}
