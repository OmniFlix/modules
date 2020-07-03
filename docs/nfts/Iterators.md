

### Iterators

Sometimes you will want to access a `TweetNFT` directly by their key. That's why we have the methods `MintTweetNFT` and `GetTweetNFTByID`. However, sometimes you will want to get every `TweetNFT` at once. To do this we use an Iterator called `KVStorePrefixIterator`. This utility comes from the `sdk` and iterates over a key store. If you provide a prefix, it will only iterate over the keys that contain that prefix. Since we have prefixes defined for our `TweetNFT`, we can use them here to only return our desired data types.

---
Now that you've seen the `Keeper` where every `TweetNFT` that is stored, we need to connect the messages to this storage. This process is called *handling* the messages and is done inside the `Handler`.