package ECSCppStyle

import "reflect"

type Coordinator struct {
	em *EntityManager
	sm *SystemManager
	cm *ComponentManager
}

func GetTotalEntities(coord *Coordinator) int {
	return coord.em.entityCount
}

func NewCoordinator(ents int) *Coordinator {
	return &Coordinator{
		em: NewEntityManager(ents),
		sm: NewSystemManager(),
		cm: newComponentManager(),
	}
}

func CreateNewEntity(coord *Coordinator) Entity {
	return coord.em.CreateEntity()
}

func EraseEntity(coord *Coordinator, entity Entity) {
	coord.em.DestroyEntity(entity)
	entityDestroyed(coord.cm, entity)
	coord.sm.EntityDestroyed(entity)
}

func RegisterNewComponentType[T any](coord *Coordinator) int {
	return registerComponent[T](coord.cm)
}

func AddNewComponent[T any](coord *Coordinator, entity Entity, comp T) {

	// Adds the Component to the right ComponentArray for the Entity
	addComponent(coord.cm, entity, comp)

	// Updates the Entity archetype
	sig := coord.em.GetSignature(entity)
	if sig == nil {
		sig = NewSignature()
	}
	sig.Set(int(getComponentType[T](coord.cm)))
	coord.em.SetSignature(entity, sig)

	// Informs all Systems that an entity has changed its archetype. Ie: more systems may now apply to this Entity
	coord.sm.EntitySignatureChanged(entity, sig)
}

func RemoveExistingComponent[T any](coord *Coordinator, entity Entity) {
	removeComponent[T](coord.cm, entity)

	sig := coord.em.GetSignature(entity) //
	sig.Unset(int(getComponentType[T](coord.cm)))
	coord.em.SetSignature(entity, sig)

	coord.sm.EntitySignatureChanged(entity, sig)
}

func GetExistingComponent[T any](coord *Coordinator, entity Entity) *T {

	c, err := getComponent[T](coord.cm, entity)
	if err != nil {
		return nil
	}

	return c
}

func RegisterNewSystem(coord *Coordinator, system System) {
	coord.sm.RegisterSystem(system)
}

func SetSystemSignature(coord *Coordinator, system System, sig *Signature) {
	coord.sm.SetSignature(system, sig)
}

func GetComponentType[T any](coord *Coordinator) ComponentType {
	return getComponentType[T](coord.cm)
}

func GetComponentArrayForComponentType[T any](coord *Coordinator) *ComponentArray[T] {
	t := reflect.TypeOf((*T)(nil)).Elem()
	name := t.Name()

	list, ok := coord.cm.componentArrays[name]
	if !ok {
		panic("Component Manager has no array for that component")
	}

	return list.(*ComponentArray[T])
}

func GetEntitySignature(coord *Coordinator, ent Entity) *Signature {
	return coord.em.GetSignature(ent)
}
