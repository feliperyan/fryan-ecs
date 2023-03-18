## Entity Manager

1. Holds a queue of entity Ids. 
2. Entity creation consumes an Id
3. Destroying an entity returns its Id to the queue
4. Holds a map of `EntityId` to `Signature` to represent an Entity's archetype and therefore, allow 


## Component Manager

1. Allows "registration" of new `Components` types which creates a separate `ComponentArray` for each component type.
2. Component registration returns an Id for the component which is then used to create a `Signature` for an archetype for an `Entity` as per above.
3. Also allows adding `Components` to specific `Entities`.
4. Holds a map of `String` to `ComponentArray` where the key is a reflected type name for the Component type.

## System Manager

1. Holds a map of `System` to `Signature` which represents what minimum components an `Entity` archetype must have in order for a `System` to be relevant to that `Entity`.
2. Similar to `ComponentManager`, holds a map of `String` to `System` where the key is a reflected type name for the System type.
3. Decides which `Entities` map to each `Systems` based on `func EntitySignatureChanged` which must be called everytime an `Entity` changes its archetype.

## Coordinator

Makes sense of the above 3 abstractions and provides quality of life functions.

### Set up
1. Create a new coordinator.
2. Register all component types.
3. Create and register Systems.

### Running
1. Create entities.
2. Create necessary components for each entity and add them to the entity.
3. Loop over all systems.
4. For each system, get a pointer to each componentArray needed.
5. For each system loop over the entity ids and get the corresponding components from the componentArray pointers.
6. Apply system to the components.
