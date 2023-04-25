Assumption: if the relayer paused for a while, that pool state not synced in time.

```js

// Swap request during the relayer halt
pool1.swapX4Y(10)
pool1.swapX4Y(10)
pool1.swapX4Y(10)
pool1.swapX4Y(10)
pool1.swapX4Y(10)
pool1.swapX4Y(10)
pool1.swapX4Y(10)
pool1.swapX4Y(10)
pool1.swapX4Y(10)

pool2.swapY4X(10)
pool2.swapY4X(10)
pool2.swapY4X(10)
pool2.swapY4X(10)
pool2.swapY4X(10)
pool2.swapY4X(10)
pool2.swapY4X(10)

// sync the state after relayer restart
pool1.swapY4X(10)
pool1.swapY4X(10)
pool1.swapY4X(10)
pool1.swapY4X(10)
pool1.swapY4X(10)
pool1.swapY4X(10)
pool1.swapY4X(10)

pool2.swapX4Y(10)
pool2.swapX4Y(10)
pool2.swapX4Y(10)
pool2.swapX4Y(10)
pool2.swapX4Y(10)
pool2.swapX4Y(10)
pool2.swapX4Y(10)
pool2.swapX4Y(10)
pool2.swapX4Y(10)

// synced

// check if the out amounts from pool1 on A chain and pool2 on chain B 
console.log(pool1.swapX4Y(10), pool2.swapX4Y(10))
console.log(pool1.swapY4X(10), pool2.swapY4X(10))
```

```sh
0.005811342500464889 0.005811272070786799
16948.228154912824 16948.432618577033
```

As we can see, it has a little different. but not much.
