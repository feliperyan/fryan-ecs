package ecs_cpp_style

import "fmt"

type Entity int

type EntityManager struct {
	availableEntities Queue[Entity]
	signatures        []DumbSignature
	entityCount       int
	max_ents          int
}

func NewEntityManager(max_ents int) *EntityManager {

	em := &EntityManager{
		availableEntities: make(Queue[Entity], 0),
		signatures:        make([]DumbSignature, max_ents),
		entityCount:       0,
		max_ents:          max_ents,
	}

	for i := 0; i < max_ents; i++ {
		em.availableEntities.enqueue(Entity(i))
		em.signatures[i] = make(DumbSignature, MAX_COMPONENTS)
	}

	return em
}

func (em *EntityManager) CreateEntity() Entity {
	if em.entityCount > em.max_ents {
		panic("too many entities")
	}

	em.entityCount++

	ent, ok := em.availableEntities.dequeue()
	fmt.Println(ok)
	return ent
}

func (em *EntityManager) DestroyEntity(entity Entity) {
	if int(entity) > em.max_ents {
		panic("out of range")
	}

	em.signatures[entity].Reset()
	em.availableEntities.enqueue(entity)
	em.entityCount--
}

func (em *EntityManager) SetSignature(entity Entity, signature DumbSignature) {
	if int(entity) > em.max_ents {
		panic("out of range")
	}

	em.signatures[entity] = signature
}

func (em *EntityManager) GetSignature(entity Entity) DumbSignature {
	if int(entity) > em.max_ents {
		panic("out of range")
	}

	return em.signatures[entity]
}
