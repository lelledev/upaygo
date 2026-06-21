![Go Tests](https://github.com/lelledev/upaygo/workflows/Go/badge.svg)
![Code scanning - action](https://github.com/lelledev/upaygo/workflows/Code%20scanning%20-%20action/badge.svg)

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

# If you want to contribute
cp .github/hooks/pre-commit .git/hooks/pre-commit
# Open and change absolute config path

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

## Tests

```bash
go test ./... -failfast -tags=unit
go test ./... -failfast -tags=stripe -config=ABS_PATH/config.json
```

### APIs

- Swagger */swagger/index.html*
