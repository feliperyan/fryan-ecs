## Entity Manager

1. Holds a queue of entity Ids. 
2. Entity creation consumes an Id
3. Destroying an entity returns its Id to the queue
4. Each entity has a signature to represent its combination of Components.


## Component Manager

1. 