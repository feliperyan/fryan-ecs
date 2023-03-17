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
	addComponent(coord.cm, entity, comp)

	sig := coord.em.GetSignature(entity)

	sig.Set(int(getComponentType[T](coord.cm)))

	coord.em.SetSignature(entity, sig)

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
