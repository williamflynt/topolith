# Ideation

This project is designed to help me do something I've needed to do many times - map out complex systems.
Usually without existing diagrams, and from verbal interviews with various team members.

## Current Approach

I prefer to use C4 diagrams for more polished presentation materials.
The focus on structure, levels, and verb-style relations help quickly map systems in a way that's comprehensible.
I'd typically define a `.puml` file for each level of C4 diagram I wanted, and end up rendering and saving the images.

The problem is that those are hard to update quickly as I learn, and they aren't interactive.
They don't easily provide the ability to dive into individual systems/components, either.

## `topolith` Data Structure

At the core of this `topolith` is a data structure than independently defines systems and their components.
There isn't a concept of depth, but rather of bounds.
Each component (`Item`) stands on its own in the `World`, and a separate `Tree` helps us bound components into systems.
We create the "lines" between components as separate entities - `Rel`.

These choices help us separate the data from the surrounding behaviors.
Functions operate on top of the data structures, and are defined separately.

### Nouns

I'm thinking about three main nouns:

* `World` - holder of everything that describes the diagram.
* `Item` - any system or component or database or person.
* `Rel` - relationship between `Item` of any level.

## Behaviors in `topolith`

In general, we're striving for something inspired by functional programming purity.
Side effects will be executed by hooking into specific points in the code.
Persistence, notifications, rendering, and other side effects hook into the code and are defined independently.
We can think of this like a plugin system in some ways.

## Undo and Time Travel

One interesting piece of the brainstorm idea here is "undo-ability".
Within the `World`, we will track the `Command` stream that created it.
Each `Command` has a `Command.Dual()` that produces the "undo" `Command` as output.
If we treat the `World.History` as an immutable log of commands, we get high confidence undo/redo capability.s

As a result of that choice we get time lapse capability almost for free, along with audit stream, etc...

## User Interfaces

Tracking commands also makes creating interfaces to the ultimate application much simpler.
We can convert user inputs from our CLI, web UI, or email server into `Command` objects.
The rest just happens - we execute the command and return the resulting `World` or `Item` or `Rel`.

If we can serialize commands as text, we can also store the `World` as a simple text file - the series of commands to rebuild it from scratch.

### Rendering Relationships

The most important outcome of `topolith` is a rendering of the `World` that people can understand and explore.

Given that we'll structure all `Item` in the `World` into a `World.Tree`, we need the concept of "rolling up" our `Rel`.
For example, if `SystemA.Client` calls `SystemB.Server`, a fully collapsed `World` would show `SystemA -calls-> SystemB`.
Expanding `SystemB` would update the rendering to `SystemA -calls-> SystemB.Server`.

### CLI

This is a brainstorm of the way I'd like to structure a CLI around our concepts.
I find that a CLI is a fast way to input and update information, so I want it to be a first class consideration.
A well-explored CLI also clarifies some of the concepts around an application!

For these commands, imagine we've already typed `topolith` and entered something like a REPL.
Anything after `;;;` is a comment - the remainder of the line isn't part of a command.

**NOTE:** The command history isn't represented for `World` in the output below.
This was a brainstorming choice to save horizontal space in this specific document.
In actual implementation, the `World` includes `History`.

```
world myWorld  ;;; Create a new World named myWorld.
> World{myWorld}
item database  ;;; Create a new Item with ID `database`.
> Item{database}
item database  ;;; Retrieve the Item with ID `database`.
> Item{database}
item? server  ;;; Retrieve the Item with ID `server` or nothing.
>
world  ;;; Retrieve the World.
> World{myWorld {database Item{database}}}
item server
> Item{server}
rel server database Uses  ;;; Create a Rel from `server` to `database` with verb "Uses".
> Rel{server database Uses}
rel? database server ;;; Retrieve the Rel from `database` to `server` or nothing.
> 
rel server database Writes  ;;; Update the verb to "Writes".
> Rel{server database Writes}
world
> World{myWorld {database Item{database} server Item{server}} {server {Rel{server database Writes}}}}
item server --delete  ;;; Remove the `Item`.
>
world
> World{myWorld {database Item{database}}}
item client && rel client database "Reads once"  ;;; Support for spaces in verbs using double or single quotes.
> Rel{client database "Reads once"}
item server && item client in server  ;;; Create a `server` and make the `client` a component of it.
> Tree{database server{client}}
in? database server  ;;; Check if `database` is part of the `server`.
> false
item notReal in client  ;;; Operations with more than one item do not create automatically.
> Error{"Item{notReal} not found in World{myWorld}", 412}
item! monitor in client  ;;; But we can specify creation with a bang.
> Tree{database server{client{monitor}}
in? monitor server  ;;; The `in?` query traverses the Tree.
> true
rel monitor database Pings
> Rel{monitor database Pings}
rel? server database  ;;; Retrieve a list of all `Rel` between `Item`.
> {Rel{client database "Reads once"} Rel{monitor database Pings}}
rel? server database --strict  ;;; Unless we specify strict mode.
>
world write myWorld.txt  ;;; Save our World and get the filename back.
> "/home/myUser/topolith/myWorld.txt"
world load myWorld.txt  ;;; Load the world.
> World{myWorld {database Item{database} server Item{server} client Item{client} monitor Item{monitor}} {monitor {Rel{monitor database Pings}} client {Rel{client database "Reads once"}}} Tree{database server{client{monitor}}}
world --pretty  ;;; Pretty print the World.
> World{
>   myWorld
> 
>   {
>       database Item{database}
>       server Item{server}
>       client Item{client}
>       monitor Item{monitor}
>   }
>
>   {
>       monitor {
>           Rel{monitor database Pings}
>       }
>       client {
>           Rel{client database "Reads once"}
>       }
>   }
>
>   Tree{
>       database
>       server{
>           client{
>               monitor
>           }
>       }
>   }
```

## Code

I want to code this up in Go for the core application itself, and the CLI.
It would be great if I could define some generation around my data structures to get the application.
Then I could use the core data structures in a generated grammar definition for the CLI and generate the CLI.
