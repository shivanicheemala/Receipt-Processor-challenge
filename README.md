# Receipt Processor

## Overview

This project is a web service that processes receipts and calculates points based on a set of predefined rules. The API allows for submitting receipts and querying the points awarded to a specific receipt.

## Docker Setup

### Build Docker Image
Build the Docker image for the Receipt Processor.

docker build -t receipt-processor .

### Run Docker Container
Start the container on your local machine on port 8080.

docker run -p 8080:8080 receipt-processor

### try hitting the API as shown in example or as below

$headers = @{
    "Content-Type" = "application/json"
}
$body = @{
    "retailer" = "Target"
    "purchaseDate" = "2022-01-01"
    "purchaseTime" = "13:01"
    "items" = @(
        @{
            "shortDescription" = "Mountain Dew 12PK"
            "price" = "6.49"
        },
        @{
            "shortDescription" = "Emils Cheese Pizza"
            "price" = "12.25"
        },
        @{
            "shortDescription" = "Knorr Creamy Chicken"
            "price" = "1.26"
        },
        @{
            "shortDescription" = "Doritos Nacho Cheese"
            "price" = "3.35"
        },
        @{
            "shortDescription" = "Klarbrunn 12-PK 12 FL OZ"
            "price" = "12.00"
        }
    )
    "total" = "35.35"
}
$bodyJson = $body | ConvertTo-Json
$response = Invoke-WebRequest -Uri http://localhost:8080/receipts/process -Method Post -Headers $headers -Body $bodyJson
$response.Content

### Run the following command using receipt-id you get from above

curl http://localhost:8080/receipts/<receipt_id>/points





