# Nintendo Badge Arcade secure

The NEX secure server for Nintendo Badge Arcade.

## Current status

Currently, the secure server is able to create a new account, but trying to access the catchers after the first time makes the console *crash*. Other game functionality hasn't been tested yet.  

Also, the server currently doesn't work with Pretendo Network, as there are some internal game checks that report `IncompatibleBOSSData` error (at least when checking the internal game dialogue names).  

## Credits

The codebase of the server is based on Pretendo Network's [Super Mario Maker secure](https://github.com/PretendoNetwork/super-mario-maker-secure) server.

