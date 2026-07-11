![Go Tests](https://github.com/lelledev/upaygo/workflows/Go/badge.svg)
![CodeQL Advanced](https://github.com/lelledev/upaygo/workflows/CodeQL%20Advanced/badge.svg)

[![Run in Postman](https://run.pstmn.io/button.svg)](https://app.getpostman.com/run-collection/bf7491736c431cd822f6)

# uPay in Golang

Payment Gateway Microservice in Golang

## PSD2 SCA

**EU SCA law will be on duty after 14th September 2019**

### Updates
 
- 13/08/2019 - For **UK cards** the [SCA implementation deadline is March 2021](https://www.fca.org.uk/news/press-releases/fca-agrees-plan-phased-implementation-strong-customer-authentication)

## Feature

- SCA ready with [Stripe Payment Intents](https://stripe.com/docs/payments/payment-intents)
    - Off-session intents
    - Separation of auth and capture
        - New intent
        - Confirm intent
        - Capture/Delete intent
- No database infrastructure needed
- Stripe API keys configuration per currency

## Installation

```bash
cp config.json.dist config.json
vi config.json # Add your config values

# API doc
swag init

# Enable pre-commit hooks (Git 2.54+ config-based hooks; no copy into .git/hooks)
git config --local include.path ../.github/hooks/gitconfig
# Verify registration (list only; does not execute hooks):
git hook list --show-scope pre-commit
# Run the configured pre-commit commands to confirm they work:
git hook run pre-commit

go run main.go -config=config.json
```


## How to use

*Example of a checkout web page*

1) User to insert card information with Stripe Elements
2) Create an intent with JS SDK
3) Confirm the intent
4) Does the intent requires 3D Secure (intent status and next_action param)
    1) No, Point 5)    
    2) Yes, Stripe Elements will open the a 3D Secure popup
5) Do you after checkout domain logic
6) Any error during your checkout process?
    1) Yes, cancel the intent
    2) No, capture the intent

## Stripe client

All Stripe calls go through a per-currency client built by
`appconfig.GetStripeClientByCurrency` (the stripe-go `stripe.Client` API), so
the secret key is bound to the client instead of shared global state. The
deprecated global key (`stripe.Key`) and the legacy resource packages
(`stripe-go/v82/paymentintent`, `stripe-go/v82/customer`) must not be used:
the pre-commit hook and CI reject them.

## Tests

```bash
go test ./... -failfast -tags=unit
go test ./... -failfast -tags=stripe -config=ABS_PATH/config.json
```

### APIs

- Swagger */swagger/index.html*
