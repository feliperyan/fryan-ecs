package ecs_cpp_style

import "fmt"

type Entity int

type EntityManager struct {
	availableEntities Queue[Entity]
	signatures        []*Signature
	entityCount       int
	max_ents          int
}

func NewEntityManager(max_ents int) *EntityManager {

	em := &EntityManager{
		availableEntities: make(Queue[Entity], 0),
		signatures:        make([]*Signature, max_ents),
		entityCount:       0,
		max_ents:          max_ents,
	}

	for i := 0; i < max_ents; i++ {
		em.availableEntities.enqueue(Entity(i))
	}

	return em
}

func (em *EntityManager) CreateEntity() Entity {
	if em.entityCount > em.max_ents {
		panic("too many entities")
	}

	em.entityCount++

	ent, ok := em.availableEntities.dequeue()
	if !ok {
		fmt.Println("create entities failed to deque a new one")
	}
	return ent
}

func (em *EntityManager) DestroyEntity(entity Entity) {
	if int(entity) > em.max_ents {
		panic("out of range")
	}

	em.signatures[entity] = NewSignature() // TODO: should reset this to nil
	em.availableEntities.enqueue(entity)
	em.entityCount--
}

func (em *EntityManager) SetSignature(entity Entity, signature *Signature) {
	if int(entity) > em.max_ents {
		panic("out of range")
	}

	if em.signatures[entity] == nil {
		em.signatures[entity] = NewSignature()
	}

	em.signatures[entity] = signature
}

func (em *EntityManager) GetSignature(entity Entity) *Signature {
	if int(entity) > em.max_ents {
		panic("out of range")
	}

	if em.signatures[entity] == nil {
		em.signatures[entity] = NewSignature()
	}

	return em.signatures[entity]
}
