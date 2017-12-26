# Services Design and Architecture

## Overview

## Goals

## Architecture

This project consists of two main components: an authorization and identity service and a gateway to expose deployed applications.

*TODO: Add architecture diagram*

### [Gateway](./gateway/overview.md)

The Services gateway consists of two components: the gateway and the gateway controller. The gateway is exposed to the public and responsible for loading balancing, routing and protection. The gateway controller is responsible for automatically configuring the gateway and its associated services when applications are deployed, updated or teared down.

### [Authorization and Identity Service]()
