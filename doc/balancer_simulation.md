# Balancer Protocol Analysis Report

## Purpose

The purpose of this report is to evaluate the viability of withdrawing all assets from a Balancer pool with additional incentives to facilitate the pool's liquidity. This is crucial for a liquidity provider looking to fully divest from a pool. Our findings reveal deviations from expectations based on assumptions outlined in Balancer's white paper.

## Weighted Pool Simulation

### Scenario 1

Let's create a scenario where a pool is comprised of two tokens: X and Y. The initial quantities of these tokens are `X = 1000, Y = 1000`. Thus, the Initial Liquidity Provider (LP) token supply is `2000 (1000X + 1000Y)`.

Suppose another user wants to deposit the same amount of tokens into this pool (`Xd = 1000, Yd = 1000`).

According to Balancer's multi-deposit formula (`Dk = ((Ps - Pi)/Ps - 1) * Bk`), we need to determine the quantity of LP tokens to mint for this deposit. The variables in the formula are defined as follows:

- Dk: The deposit amount of token k
- Bk: The amount of token k in the pool
- Pi: Newly issued LP token amount
- Ps: Current total supply of LP tokens

From this formula, we derive the quantity of newly issued LP tokens as `Pi = Ps * Dk/Bk`.

Thus, the minted LP tokens are:

- For X tokens: `Pix = 2000 * 1000/1000 = 2000`
- For Y tokens: `Piy = 2000 * 1000/1000 = 2000`

As a result, the total supply of LP tokens is `Ps = 2000 + 2000 + 2000 = 6000` with `Bx = 2000, By = 2000`.

Let's now assume that these tokens are withdrawn again.

From Balancer's withdrawal formula (`Ak = (1 −(Ps-Pr)/Ps) * Bk`), we can deduce the amount withdrawn as `Ak = Bk*Pr/Ps`.

The withdrawn quantities are:

- For X tokens: `Ax = 2000 * 2000/6000 = 666.67`
- For Y tokens: `Ay = 2000 * 2000/6000 = 666.67`

The resulting quantities suggest a withdrawal of `666.67` for each token type, instead of the deposited `1000`. This represents a loss of `33.34%` for each token.

In this scenario, the loss appears to decrease significantly if the deposit amount is substantially smaller than the initial LP. However, additional analysis is required for a comprehensive understanding of the Balancer protocol's implications on liquidity providers' investments.

### Scenario 2

In this step, we simulate the decrease in loss according to the deposit amount. We'll use the same pool, but this time we'll only deposit a small amount, specifically `5%` of initial liquidity.

- For X tokens: `Pix = 2000 * 50/1000 = 100`
- For Y tokens: `Piy = 2000 * 50/1000 = 100`

As a result, the total supply of LP tokens is `Ps = 2000 + 100 + 100 = 2200` with `Bx = 2050, By = 2050`.

Let's now assume that these tokens are withdrawn again.

From Balancer's withdrawal formula (`Ak = (1 −(Ps-Pr)/Ps) * Bk`), we can deduce the amount withdrawn as `Ak = Bk*Pr/Ps`.

The withdrawn quantities are:

- For X tokens: `Ax = 1050 * 100/2200 = 47.72`
- For Y tokens: `Ay = 1050 * 100/2200 = 47.72`

The resulting quantities suggest a withdrawal of `47.72` for each token type, instead of the deposited `50`. This represents a loss of `4.56%` for each token. The lower deposit in this scenario has led to a significantly reduced loss compared to Scenario 1, confirming our initial hypothesis that losses decrease with smaller deposit amounts relative to the initial LP.

### Scenario 3

In this case, we will work with a pool that is not balanced in its initial token distribution: `X = 2000000, Y = 1000`. This makes our initial LP token supply `2000000 (2000000X + 1000Y)`.

If a user wants to deposit `20%` of the initial token amount into this pool, the deposit amounts are:

- For X tokens: `Xd = 2000000 * 0.20 = 400000`
- For Y tokens: `Yd = 1000 * 0.20 = 200`

Applying Balancer's multi-deposit formula, we calculate the quantity of LP tokens to mint for this deposit.

- For X tokens: `Pix = 2000000 * 400000/2000000 = 400000`
- For Y tokens: `Piy = 2000000 * 200/1000 = 400000`

The total supply of LP tokens is `Ps = 2000000 + 400000 + 400000 = 2800000`, with `Bx = 2400000, By = 1200`.

Assuming these tokens are withdrawn:

- For X tokens: `Ax = 2400000 * 400000/2800000 = 342857.14`
- For Y tokens: `Ay = 1200 * 400000/2800000 = 171.43`

We have withdrawal amounts of `342857.14` for X tokens and `171.43` for Y tokens, compared to the deposited `400000` and `200` respectively. This represents a loss of `14.29%` for X tokens and `14.29%` for Y tokens, showing that the loss percentage is the same for both token types, despite the significant difference in their quantities.

This scenario further demonstrates that losses decrease with smaller deposit amounts relative to the initial LP, even in cases where the token distribution in the pool is heavily skewed towards one token.

### Scenario 4 (Revised Again)

Consider an uneven initial distribution for a pool (`X = 2000000, Y = 1000`), leading to an initial LP token supply of `2001000 (2000000X + 1000Y)`.

Now, suppose a second user deposits `10%` of the original assets into this pool:

- For X tokens: `Xd = 2000000 * 0.10 = 200000`
- For Y tokens: `Yd = 1000 * 0.10 = 100`

Applying Balancer's multi-deposit formula, the newly minted LP tokens are:

- For X tokens: `Pix = 2001000 * 200000/2000000 = 200100`
- For Y tokens: `Piy = 2001000 * 100/1000 = 200100`

The total supply of LP tokens now is `Ps = 2001000 + 200100 + 200100 = 2401200`, with `Bx = 2200000, By = 1100`.

Now, if the original depositor attempts to withdraw all assets using all the initial and subsequent LP tokens (`2000000 + 200100` for X and `1000 + 200100` for Y), the quantities withdrawn are:

- For X tokens: `Ax = 2200000 * (2000000 + 200100)/2401200 = 2015750`
- For Y tokens: `Ay = 1100 * (1000 + 200100)/2401200 = 92`

The remaining quantities in the pool will be:

- For X tokens: `Bx_remaining = Bx - Ax = 2200000 - 2015750 = 184250`
- For Y tokens: `By_remaining = By - Ay = 1100 - 92 = 1008`

So, if the original depositor uses all the LP tokens in existence for withdrawal, all the assets will be withdrawn, leaving the pool empty.

# Results

- The loss amount in the Balancer protocol does not depend on the ratio of the assets in the pool, but instead depends on the ratio of the current asset amount and the new deposit amount.
- This loss will consistently occur when a user makes a withdrawal, thereby preventing potential attackers from destabilizing the pool through rapid deposit/withdrawal operations.
- It seems plausible that this loss can be compensated through other incentives or swap fees. Given the loss is a constant value, full compensation could be achievable if a user maintains their deposit in the pool over a long-term period.
- Therefore, only long-term deposits can potentially generate profits from providing liquidity in the pool.
- This inherent property of the Balancer protocol also provides resistance against malicious pool attackers.
- Additionally, it is noteworthy that even when utilizing all LP tokens, it's impossible to withdraw the original funds entirely from the pool. This indicates that a withdrawal action cannot collapse the pool, which enhances the protocol's robustness.
