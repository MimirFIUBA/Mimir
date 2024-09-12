# MIMIR - Triggers
This is a trigger framework design.

## Overview
The main idea of the framework is to provide a flexible way to react to events and act accordingly.

## Trigger Observer
The trigger observer works as an interface for all triggers. The trigger provides a method `Update(event Event)` that should be called when we want the trigger to react.

### Implementations
#### Trigger
The first implementation of the trigger observer. The trigger contains a condition and actions. Once the update event is called it will check that the condition is satisfied and will execute all actions.

#### Time Trigger
This trigger is mainly meant to be used as a timeout. It also has a condition and actions, but works slightly different. For it to start running it provides a method `Start()`. Once this method is executed a timer will start ticking and will execute all the actions in intervals of time equal to the `Duration` attribute of the Trigger.
If the Trigger is updated by the `Update` method it will evaluate the condition and if it's true it will restart the ticker.

#### FrequencyTrigger


## Condition
Condition is the interface for the different evaluations of the triggers.

### Implementations
#### Boolean Conditions
- **Compare condition:** compares a value with a reference based on an opreator `>, >=, <, <=, ==, !=`.
- **AndCondition:** applies the operator `and` to all of its contained conditions.
- **OrCondition:** applies the operator `or` to all of its contained conditions.
- **TrueCondition:** it will allways evaluate to true.

#### Receive value condition
It will evaluate to true once it receives a value. It is mainly used for timeout triggers.

#### CustomCondition
It provides the possibility to define a custom function to evaluate.

#### BetweenCondition
Evaluates if the value is between a max and min.

#### DeltaCondition
Evaluates if the current value has a difference of delta from the previous one.

#### AverageCondition
Under development. It should evaluate an average from a set of values.



##### TODO:
- [ ] COMPARE CONDITION that compares a value from an event with another value of another event and not a reference.
- [ ] Check that receive value is restarted once it is evaluated.
- [ ] False condition? (no le encontre uso por ahora)
- [ ] Average condition

## Action
Action provides an interface for users to define it's own actions. An action should implement the method `Execute()`.


