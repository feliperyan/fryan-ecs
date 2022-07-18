package ecs_cpp_style

import (
	"fmt"
	"reflect"
)

type System interface {
	AddEntity(entity Entity)
	RemoveEntity(entity Entity)
}

type SystemManager struct {
	SystemSignatures map[string]DumbSignature
	Systems          map[string]System
}

func NewSystemManager() *SystemManager {
	return &SystemManager{
		SystemSignatures: map[string]DumbSignature{},
		Systems:          map[string]System{},
	}
}

func (sm *SystemManager) RegisterSystem(sys System) {
	t := reflect.TypeOf(sys).Elem()
	name := t.Name()
	sm.Systems[name] = sys
}

func (sm *SystemManager) SetSignature(sys System, signature DumbSignature) {
	t := reflect.TypeOf(sys).Elem()
	name := t.Name()
	sm.SystemSignatures[name] = signature
}

func (sm *SystemManager) EntitySignatureChanged(entity Entity, signature DumbSignature) {
	for sysType, system := range sm.Systems {

		sysSig := sm.SystemSignatures[sysType]

		ok, err := sysSig.Contains(signature)
		if err != nil {
			panic(fmt.Sprintf("Issue with signatures: %v", err))
		}

		if ok {
			system.AddEntity(entity)
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
