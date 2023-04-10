package ecs_cpp_style

import (
	"fmt"
	"reflect"
)

type ComponentType uint8

type IComponentArray interface {
	EntityDestroyed(entity Entity)
}

type ComponentArray[T any] struct {
	Components    []T
	EntityToIndex map[Entity]int
	IndexToEntity map[int]Entity
	size          int
}

func NewComponentArray[T any](maxEntities int) *ComponentArray[T] {
	return &ComponentArray[T]{
		Components:    make([]T, 0, maxEntities),
		EntityToIndex: map[Entity]int{},
		IndexToEntity: map[int]Entity{},
		size:          0,
	}
}

func (ca *ComponentArray[T]) Insert(entity Entity, component T) {

	if ca.size < len(ca.Components) {
		ca.Components[ca.size] = component
		ca.EntityToIndex[entity] = ca.size
		ca.IndexToEntity[ca.size] = entity

	} else {
		ca.Components = append(ca.Components, component)
		ca.IndexToEntity[ca.size] = entity
		ca.EntityToIndex[entity] = ca.size
	}

	ca.size++
}

func (ca *ComponentArray[T]) remove(entity Entity) {

	indexOfRemovedEntity := ca.EntityToIndex[entity]

	// Replace remove entity with last entity in the slice
	ca.Components[indexOfRemovedEntity] = ca.Components[ca.size-1]

	// Update maps

	entityOfLastElement := ca.IndexToEntity[ca.size-1]
	ca.EntityToIndex[entityOfLastElement] = indexOfRemovedEntity
	ca.IndexToEntity[indexOfRemovedEntity] = entityOfLastElement

	delete(ca.EntityToIndex, entity)
	delete(ca.IndexToEntity, ca.size-1)

	ca.size--
}

func (ca *ComponentArray[T]) Get(entity Entity) (*T, bool) {
	i, ok := ca.EntityToIndex[entity]
	if ok {
		return &ca.Components[i], true
	}
	var a T
	return &a, false

}

func (ca *ComponentArray[T]) EntityDestroyed(entity Entity) {
	_, ok := ca.EntityToIndex[entity]
	if ok {
		ca.remove(entity)
	}
}

// -----

type ComponentManager struct {
	componentTypes    map[string]ComponentType
	componentArrays   map[string]IComponentArray
	nextComponentType ComponentType
}

func newComponentManager() *ComponentManager {
	return &ComponentManager{
		componentTypes:    make(map[string]ComponentType),
		componentArrays:   make(map[string]IComponentArray),
		nextComponentType: 1,
	}
}

func registerComponent[T any](cm *ComponentManager) int {
	t := reflect.TypeOf((*T)(nil)).Elem()
	name := t.Name()

	cm.componentTypes[name] = cm.nextComponentType

	compArray := &ComponentArray[T]{
		Components:    make([]T, 0),
		EntityToIndex: make(map[Entity]int),
		IndexToEntity: make(map[int]Entity),
		size:          0,
	}
	cm.componentArrays[name] = compArray

	cm.nextComponentType++

	return int(cm.nextComponentType) - 1
}

func getComponentType[T any](cm *ComponentManager) ComponentType {
	t := reflect.TypeOf((*T)(nil)).Elem()
	name := t.Name()

	return cm.componentTypes[name]
}

func addComponent[T any](cm *ComponentManager, entity Entity, component T) {
	t := reflect.TypeOf((*T)(nil)).Elem()
	name := t.Name()

	array, ok := cm.componentArrays[name]
	if !ok {
		panic("addComponent can't add to a Component which has not been registered")
	}
	array.(*ComponentArray[T]).Insert(entity, component)
}

func removeComponent[T any](cm *ComponentManager, entity Entity) {
	t := reflect.TypeOf((*T)(nil)).Elem()
	name := t.Name()

	array := cm.componentArrays[name]
	array.(*ComponentArray[T]).remove(entity)
}

func getComponent[T any](cm *ComponentManager, entity Entity) (*T, error) {
	t := reflect.TypeOf((*T)(nil)).Elem()
	name := t.Name()

	array := cm.componentArrays[name]
	comp, ok := array.(*ComponentArray[T]).Get(entity)
	if ok {
		return comp, nil
	}

	var a T
	return &a, fmt.Errorf("could not find Component of type %v for Entity %v", name, entity)
}

func entityDestroyed(cm *ComponentManager, entity Entity) {
	for _, compArray := range cm.componentArrays {
		compArray.EntityDestroyed(entity)
	}
}
