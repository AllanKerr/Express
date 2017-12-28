# Teardown Command

The teardown command allows developers to delete previously deployed application containers.

## Command

The teardown command is made available to developers through the `services` executable.

```
services teardown <name>
```

## Parameters

1. ***Name.*** The only parameter required for the teardown command is `name`. The name must match the name passed to the `deploy` command. The list of all deployed application containers can be found using the `list` command.

## Design

[The design of the teardown command can be found here.](./teardown-command-design.md)
