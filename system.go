package ecs_cpp_style

import (
	"reflect"
)

type System interface {
	AddEntity(entity Entity)
	RemoveEntity(entity Entity)
	HasEntity(entity Entity) bool
}

type SystemManager struct {
	SystemSignatures map[string]*Signature
	Systems          map[string]System
}

func NewSystemManager() *SystemManager {

	return &SystemManager{
		SystemSignatures: make(map[string]*Signature),
		Systems:          make(map[string]System),
	}
}

func (sm *SystemManager) RegisterSystem(sys System) {
	t := reflect.TypeOf(sys).Elem()
	name := t.Name()
	sm.Systems[name] = sys
}

func (sm *SystemManager) SetSignature(sys System, signature *Signature) {
	t := reflect.TypeOf(sys).Elem()
	name := t.Name()
	sm.SystemSignatures[name] = signature
}

func (sm *SystemManager) EntitySignatureChanged(entity Entity, signature *Signature) {
	for sysType, system := range sm.Systems {

		sysSig := sm.SystemSignatures[sysType]

		// ok := signature.Contains(sysSig)
		ok := sysSig.Contains(signature)

		if ok {
			if !system.HasEntity(entity) {
				system.AddEntity(entity)
			}
		} else {
			system.RemoveEntity(entity)
		}
	}
}

func (sm *SystemManager) EntityDestroyed(entity Entity) {
	for _, sys := range sm.Systems {
		sys.RemoveEntity(entity)
	}
}
