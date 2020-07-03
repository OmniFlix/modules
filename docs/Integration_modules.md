## App Integration

If you use our [Scaffold]() to create App, you are able to find TODOs in `app.go`. Follow those TODOs to config the Modules in `app.go` file.

### app.go
- #### Importing Modules
```go=
    // TODO : Import of nft & xnft
    "github.com/FreeFlixMedia/modules/nfts"
    "github.com/FreeFlixMedia/modules/xnfts"
```

- #### Adding Modules in NewBasicMager
```go=
    // TODO: Add nft & xnft module(s) AppModuleBasic
    nfts.AppModuleBasic{},
    xnfts.AppModuleBasic{},
```

- #### Adding Keeper and ScopedKeeper 
```go=
    // TODO: Add nft & xnft Keeper
    nftKeeper nfts.Keeper
    xnftKeeper xnfts.Keeper
    
     // TODO: Add scoped xnft Keeper
    scopedXNFTKeeper capability.ScopedKeeper
```

- #### Adding StoreKeys
```go=
    // TODO: Add the keys that module requires
    nfts.StoreKey, xnfts.StoreKey
```

- #### Adding ScopedXNFTKeeper
```go=
    // TODO: Add scopedXNFTKeeper
    scopedXNFTKeeper := app.capabilityKeeper.ScopeToModule(xnfts.ModuleName)
```
- #### Adding Module Keeper
```go=
	// TODO: initialize nft & xnft Keepers
    app.nftKeeper = nfts.NewKeeper(app.cdc, keys[nfts.StoreKey])
    app.xnftKeeper = xnfts.NewKeeper(app.cdc, keys[xnfts.StoreKey], app.nftKeeper, app.bankKeeper,app.ibcKeeper.ChannelKeeper, &app.ibcKeeper.PortKeeper, scopedXNFTKeeper)
    xnftModule := xnfts.NewAppModule(app.xnftKeeper)
```

- #### Adding Router
```go=
    // TODO: Add xnft Route
    ibcRouter.AddRoute(xnfts.ModuleName, xnftModule)
```
- #### Intialiting modules
```go=
    // TODO: Add nft & xnft module(s)
    nfts.NewAppModule(app.nftKeeper),
    xnftModule,
```

- #### Init gensis modules
```go=
    // TODO: Init  nft & xnft module(s)
    nfts.StoreKey,
    xnfts.StoreKey,
```

- #### Adding Scoped Keeper
```go=
    //TODO: Add ScopedXNFTKeeper
    app.scopedXNFTKeeper= scopedXNFTKeeper
```