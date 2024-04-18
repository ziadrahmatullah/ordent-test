# Ordent Test - Ecommerce
## Postman Documentation

For the postman documentation, you can clik this link : https://documenter.getpostman.com/view/31472691/2sA3Bn7CjZ

## System Design Documentation

For system design documentation, you can clik this link : 
https://whimsical.com/ordent-ecommerce-XQ6EDpZDMcuYnQVnod94vy@FNpptVQ1B1cQufatFQE8oqqUUA6k9

## How to run
1. Make sure that you have already installed docker
2. Copy the .env.example for adjust the application
    ```terminal
    cp .env.example .env

    ```
3. Make sure to change the DB_PORT at .env or you can just stop your postgresql system
    ```terminal
    sudo systemctl stop postgres

    ```
4. Run docker compose
    ```terminal
    docker compose up

    ```
5. Run the server
    ```terminal
    docker run ordent-test-rest

    ```
6. You can try every endpoint on the postman

