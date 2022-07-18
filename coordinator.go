package ecs_cpp_style

type Coordinator struct {
	em *EntityManager
	sm *SystemManager
	cm *ComponentManager
}

func NewCoordinator(ents int) *Coordinator {
	return &Coordinator{
		em: NewEntityManager(ents),
		sm: NewSystemManager(),
		cm: NewComponentManager(),
	}
}

func CreateNewEntity(coord *Coordinator) Entity {
	return coord.em.CreateEntity()
}

func EraseEntity(coord *Coordinator, entity Entity) {
	coord.em.DestroyEntity(entity)
	EntityDestroyed(coord.cm, entity)
	coord.sm.EntityDestroyed(entity)
}

func RegisterNewComponentType[T any](coord *Coordinator) {
	RegisterComponent[T](coord.cm)
}

func AddNewComponent[T any](coord *Coordinator, entity Entity, comp T) {
	AddComponent[T](coord.cm, entity, comp)

	sig := coord.em.GetSignature(entity) //
	sig.Set(int(GetComponentType[T](coord.cm)), true)
	coord.em.SetSignature(entity, sig)

	coord.sm.EntitySignatureChanged(entity, sig)
}

func RemoveExistingComponent[T any](coord *Coordinator, entity Entity) {
	RemoveComponent(coord.cm, entity)

	sig := coord.em.GetSignature(entity) //
	sig.Set(int(GetComponentType[T](coord.cm)), false)
	coord.em.SetSignature(entity, sig)

	coord.sm.EntitySignatureChanged(entity, sig)
}

func GetExistingComponent[T any](coord *Coordinator, entity Entity) *T {

	c, err := GetComponent[T](coord.cm, entity)
	if err != nil {
		return nil
	}

	return &c
}

func RegisterNewSystem(coord *Coordinator, system System) {
	coord.sm.RegisterSystem(system)
}
func SetSystemSignature(coord *Coordinator, system System, sig DumbSignature) {
	coord.sm.SetSignature(system, sig)
}
