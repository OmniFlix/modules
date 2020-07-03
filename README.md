For those of you that landed here for  the first time, please check out the [freeflix-media-hub](https://github.com/FreeFlixMedia/freeflix-media-hub) repository to setup the chain.

In this document, you can go through the structure of both the ```nfts``` and the ```xnfts``` modules within. Please don't mind the plural convention, the production version will have a singular naming convention.

This document in itself is extremely 'inspired' from prev. Code with Us resources, especially the [Scavenger Hunt](https://tutorials.cosmos.network/scavenge/tutorial/01-background.html) tutorial 

## [Modules Documenation](docs)
### Table of Contents
1. [NFT Module](docs/nfts)
   1. [Messages](./docs/nfts/Msgs.md)
   2. [Codec](docs/nfts/Codec.md)
   3. [Alias](docs/nfts/Alias.md)
   4. [Keeper](docs/nfts/Keeper.md)
   5. [Prefixes](docs/nfts/Prefixes.md)
   6. [Iterators](docs/nfts/Iterators.md)
   7. [Handler](docs/nfts/Handler.md)
   8. [Events](docs/nfts/Events.md)
   9. [Querier](docs/nfts/Querier.md)
   10. [Cli](docs/nfts/Cli.md)
   11. [Module](docs/nfts/Module.md)
2. [XNFT Module](docs/xnfts)
   1. [Messages](./docs/xnfts/Msgs.md)
   2. [Codec](docs/xnfts/Codec.md)
   3. [Alias](docs/xnfts/Alias.md)
   4. [Keeper](docs/xnfts/Keeper.md)
   5. [Packets](./docs/xnfts/Packets.md)
   6. [Prefixes](docs/xnfts/Prefixes.md)
   7. [Handler](docs/xnfts/Handler.md)
   8.  [Events](docs/xnfts/Events.md)
   9.  [Relay](docs/xnfts/Relay.md)
   10. [Cli](docs/xnfts/Cli.md)
   11. [Module](docs/xnfts/Module.md)

## App Creation
In order use the **NFT & XNFT** Modules, you need to have the App. you can follow [App Creation](docs/App_creation.md) documentation to create your own app with custom configuration.

## Modules intergartion

By [App Creation](#app-creation), you can able to create scaffold of basic app, now you need to [Integrate modules](docs/Integration_modules.md) into App. 