# Nintendo Badge Arcade secure

The NEX secure server for Nintendo Badge Arcade.

## Requirements  

This server requires protocols and methods that *currently* aren't upstream on the Pretendo Network's `nex-protocols-go` repository. The `go.mod` points to [my fork](https://github.com/DaniElectra/nex-protocols-go) with those methods.

## Current status

Currently, the secure server is **incomplete** and **untested**. For example: the ability to create a new account is missing, and there might be issues related to the `Shop` protocol.  

## Credits

The codebase of the server is based on Pretendo Network's [Super Mario Maker secure](https://github.com/PretendoNetwork/super-mario-maker-secure) server.

