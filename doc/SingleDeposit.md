Compare single side and two side deposit:

```ts

const initialX = 2_000_000 // USDT
const initialY = 1000 // ETH
const fee = 0.003
const pool1 = new WeightedPool(initialX, initialY, fee)
const pool2 = new WeightedPool(initialX, initialY, fee)
const pool3 = new WeightedPool(initialX, initialY, fee)
const pool4 = new WeightedPool(initialX, initialY, fee)

const depositX = initialX * 30 / 100
const depositY = initialY * 30 / 100

console.log("deposit: ", depositX, depositY)

pool1.deposit(depositX, depositY)
console.log('Market price after add:', pool1.x, pool1.y, pool1.marketPrice())
pool2.deposit(depositX * 2, 0)
console.log('Market price after add:', pool2.x, pool2.y, pool2.marketPrice())

const a1 = pool1.swapX4Y(depositX), a2 = pool2.swapX4Y(depositX)
const b1 = pool1.swapY4X(depositY), b2 = pool2.swapY4X(depositY)

console.log(a1, a2, a1-a2,  `${((a1/a2 - 1)*100).toFixed(2)}%`)
console.log(b1, b2, b1-b2, `${((b1/b2 - 1)*100).toFixed(2)}%`)

pool3.deposit(depositX, depositY)
console.log('Market price after add:', pool3.x, pool3.y, pool3.marketPrice())
pool4.deposit(0, depositY * 2)
console.log('Market price after add:', pool4.x, pool4.y, pool4.marketPrice())

const c1 = pool3.swapX4Y(depositX), c2 = pool3.swapX4Y(depositX)
const d1 = pool4.swapY4X(depositY), d2 = pool4.swapY4X(depositY)

console.log(c1, c2, c1-c2,  `${((c1/c2 - 1)*100).toFixed(2)}%`)
console.log(d1, d2, d1-d2, `${((d1/d2 - 1)*100).toFixed(2)}%`)
```

```sh
% node singleDeposit.js
deposit:  600000 300
Market price after add: 2600000 1300 2000
Market price after add: 3200000 1000 3200
243.1555249828029 157.49565583697546 85.65986914582743 54.39%
705472.5599939243 995127.2748816465 -289654.71488772216 -29.11%
Market price after add: 2600000 1300 2000
Market price after add: 2000000 1600 1250
243.1555249828029 166.5273324610914 76.6281925217115 46.02%
314991.3116739507 229272.17663466535 85719.13503928535 37.39%
```

conclusion:
when the amount of single deposit increase, the difference increases as well. we need set a upper limit. say single deposit amount / pools.amount <= 1%


