# Balancer Protocol Analysis Report

## Purpose

The purpose of this report is to evaluate the viability of withdrawing all assets from a Balancer pool with additional incentives to facilitate the pool's liquidity. This is crucial for a liquidity provider looking to fully divest from a pool. Our findings reveal deviations from expectations based on assumptions outlined in Balancer's white paper.

## Weighted Pool Simulation

### Scenario 1

Let's create a scenario where a pool is comprised of two tokens: X and Y. The initial quantities of these tokens are `X = 1000, Y = 1000, Wx = 0.5, Wy = 0.5`. Thus, the Initial Liquidity Provider (LP) token supply is `2000 (1000X + 1000Y)`.

Suppose another user wants to deposit the same amount of tokens into this pool (`Xd = 1000, Yd = 1000`).

According to Balancer's multi-deposit formula (`Dk = ((Ps - Pi)/Ps - 1)/Wk * Bk`), we need to determine the quantity of LP tokens to mint for this deposit. The variables in the formula are defined as follows:

- Dk: The deposit amount of token k
- Bk: The amount of token k in the pool
- Pi: Newly issued LP token amount
- Ps: Current total supply of LP tokens

From this formula, we derive the quantity of newly issued LP tokens as `Pi = Ps*Wk * Dk/Bk`.

Thus, the minted LP tokens are:

- For X tokens: `Pix = 2000 * 0.5 * 1000/1000 = 1000`
- For Y tokens: `Piy = 2000 * 0.5* 1000/1000 = 1000`

As a result, the total supply of LP tokens is `Ps = 2000 + 1000 + 1000 = 4000` with `Bx = 2000, By = 2000`.

Let's now assume that these tokens are withdrawn again.

From Balancer's withdrawal formula (`Ak = (1 −(Ps-Pr)/Ps)/Wk * Bk`), we can deduce the amount withdrawn as `Ak = Bk*Pr/Ps`.

The withdrawn quantities are:

- For X tokens: `Ax = 2000 * 1000/4000/0.5 = 1000`
- For Y tokens: `Ay = 2000 * 1000/4000/0.5 = 1000`

### Scenario 2

In this step, we simulate the decrease in loss according to the deposit amount. We'll use the same pool, but this time we'll only deposit a small amount, specifically `5%` of initial liquidity.

- For X tokens: `Pix = 2000 *0.5* 50/1000 = 50`
- For Y tokens: `Piy = 2000 *0.5* 50/1000 = 50`

As a result, the total supply of LP tokens is `Ps = 2000 + 50 + 50 = 2100` with `Bx = 1050, By = 1050`.

Let's now assume that these tokens are withdrawn again.

From Balancer's withdrawal formula (`Ak = (1 −(Ps-Pr)/Ps)/Wk * Bk`), we can deduce the amount withdrawn as `Ak = Bk*Pr/Ps/Wk`.

The withdrawn quantities are:

- For X tokens: `Ax = 1050 * 50/2100/0.5 = 50`
- For Y tokens: `Ay = 1050 * 50/2100/0.5 = 50`

### Scenario 3

In this case, we will work with a pool that is not balanced in its initial token distribution: `X = 2000000, Y = 1000, Wx = 0.2, Wy = 0.8`. This makes our initial LP token supply `2001000 (2000000X + 1000Y)`.

If a user wants to deposit `20%` of the initial token amount into this pool, the deposit amounts are:

- For X tokens: `Xd = 2000000 * 0.20 = 400000`
- For Y tokens: `Yd = 1000 * 0.20 = 200`

  Applying Balancer's multi-deposit formula, we calculate the quantity of LP tokens to mint for this deposit.
  `Pi = Ps*Wk * Dk/Bk`

- For X tokens: `Pix = 2001000 *0.2*400000/2000000 = 80040`
- For Y tokens: `Piy = 2001000 *0.8*200/1000 = 320160`

The total supply of LP tokens is `Ps = 2001000 + 80040 + 320160  = 2401200`, with `Bx = 2400000, By = 1200`.

Assuming these tokens are withdrawn:
`Ak = Bk*Pr/Ps/Wk`.

- For X tokens: `Ax = 2400000 * 80040/2401200/0.2 = 400000`
- For Y tokens: `Ay = 1200 * 320160/2401200/0.8 = 200`

### Scenario 4 (Revised Again)

Consider an uneven initial distribution for a pool (`X = 2000000, Y = 1000, Wx = 0.2, Wy=0.8`), leading to an initial LP token supply of `400800(2000000*0.2 + 1000*0.8)`.

Now, suppose a second user deposits `10%` of the original assets into this pool:

- For X tokens: `Xd = 2000000 * 0.10 = 200000`
- For Y tokens: `Yd = 1000 * 0.10 = 100`

Applying Balancer's multi-deposit formula, the newly minted LP tokens are:

- For X tokens: `Pix = 400800 *0.2*200000/2000000 = 8016`
- For Y tokens: `Piy = 400800 *0.8*100/1000 = 32064`

The total supply of LP tokens now is `Ps = 400800 + 8016 + 32064 = 440880`, with `Bx = 2200000, By = 1100`.

Now, if the original depositor attempts to withdraw all assets using all the initial and subsequent LP tokens (`2000000*0.2 + 40020` for X and `1000*0.8 + 160080` for Y), the quantities withdrawn are:

- For X tokens: `Ax = 2200000 * (2000000*0.2 + 40020)/2201100/0.2 = 2199000`
- For Y tokens: `Ay = 1100 * (1000*0.8 + 160080)/2201100/0.8 = 100`

The remaining quantities in the pool will be:

- For X tokens: `Bx_remaining = Bx - Ax = 2200000 - 2199000 = 1000`
- For Y tokens: `By_remaining = By - Ay = 1100 - 100 = 1000`



