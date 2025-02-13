# Puzzles

## About
Puzzles is a distributed, simple e-commerce backend written in Golang designed to demonstrate knowledge of distributed service architecture. Strictly following OpenAPI 3.0 specifications, it's a good demonstration of my skills in a 5 week crunch period.

## Getting Started
### Prerequisites
- Docker
- Docker Compose

## Running the Application
To run the application in a containerized environment, follow these steps:
1. Clone the repository:
```bash
git clone https://github.com/robertjshirts/puzzles.git
```
2. Navigate to the project directory
```bash
cd Puzzles
```
3. Start the services using Docker Compose
```bash
docker compose up
```
Note: If `docker compose up` doesn't work, try it with a dash: `docker-compose up`.

4. Wait for the services to build and start. This may take a few moments.

### Service Dashboards
Once the services are up and running, you access the following dashboards:
- Consul (Service Registry): http://localhost:8500
- Traefik (API Gateway): http://localhost:8080
- RabbitMQ: http://localhost:15672 (username: guest, password: guest)

## Documentation
To interact with the containers, send requests to localhost:80, and refer to the following doucmentation in the repo:
- Catalog Service: ./service/catalog/catalog-api.yaml
- Basket Service: ./service/basket/basket-api.yaml
- Order Service: ./service/order/order-api.yaml
  
