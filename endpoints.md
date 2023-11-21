Request an ID
POST /v1/purchase
Requests an UUID.

Body:
N/A

Response:
POST /v1/purchase
(HTTP OK)
{
    "id": "fba670cc-ffb5-431f-9bc9-639cf48d5224"
}

====

Record a Purchase
PUT /v1/purchase/:id
Repeated calls to this endpoint will be not honored.

Body:
{
    "description": "abcde",
    "value": 1.23,
    "date": "2023-11-21"
}

Response:
PUT /v1/purchase/83174a7b-3b26-4c78-9469-57cab608ad24
(HTTP NoContent)

====

Retrieve Purchase
GET /v1/purchase/:id

Response:
GET /v1/purchase/83174a7b-3b26-4c78-9469-57cab608ad24
(HTTP OK):
{
    "date": "2023-11-21",
    "description": "abcde",
    "id": "83174a7b-3b26-4c78-9469-57cab608ad24",
    "value": 1.23
}


====

Retrieve Converted Purchase
GET /v1/purchase/:id?currency=:currency
:currency must be as specified in https://fiscaldata.treasury.gov/datasets/treasury-reporting-rates-exchange/treasury-reporting-rates-of-exchange, field Country-Currency.

Response:
GET /v1/purchase/83174a7b-3b26-4c78-9469-57cab608ad24?currency=Canada-Dollar
(HTTP OK):
{
    "converted_value:": 1.63,
    "currency": "Dollar",
    "date": "2023-11-21",
    "description": "abcde",
    "id": "83174a7b-3b26-4c78-9469-57cab608ad24",
    "rate": 1.326,
    "rate_date": "2023-06-30",
    "value": 1.23
}
