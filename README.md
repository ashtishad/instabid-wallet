## instabid-wallet

[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![LinkedIn][linkedin-shield]][linkedin-url]

A digital wallet with fast money transfer capabilities and a real-time commodity auction platform. Made with GoLang,
microservices, domain-driven design, and hexagonal architecture.

### Built With

![go][go]
![postgres][postgres]
![aws][aws]
![docker][docker]
![redis][redis]
![apache-kafka][apache-kafka]
![github-actions][github-actions]

<!-- GETTING STARTED -->

### Getting Started

###### Clone using ssh protocol `git@github.com:ashtishad/instabid-wallet.git`

#### Environment-variables

Change environment variables in Makefile, if empty then default values listed here will be loaded, check
app_helpers.go -> SanityCheck()

- API_HOST      `[IP Address of the machine]` : `127.0.0.1`
- USER_API_PORT `[Port of the user api]` : `8000`
- DB_USER       `[Database username]` : `postgres`
- DB_PASSWD     `[Database password]`: `postgres`
- DB_ADDR       `[IP address of the database]` : `127.0.0.1`
- DB_PORT       `[Port of the database]` : `5432`
- DB_NAME       `[Name of the database]` : `instabid`
- GIN_MODE      `[Name of the gin mode]` : `debug`

#### Postgres-Database-Setup

* Run docker compose: Bring the container up with `docker compose up`. Configurations are in `compose.yaml` file.
* (optional) Remove databases and volumes:
  ``` 
  docker compose down
  docker volume rm instabid-wallet_data
  ```

#### Run-the-application

* Run the application with `make run` command from project root. or, if you want to run it from IDE, please set
  environment variables by executing commands mentioned in Makefile on your terminal.

<p align="right"><a href="#instabid-wallet">↑ Top</a></p>

<!-- Project Structure -->

### Project Structure

```
├── user-api                 <-- user-api microservice.
├── .github/workflows        <-- Github CI workflows(Build, Test, Lint).
├── config                   <-- Database initialization script with docker compose.
├── db/migrations            <-- Postgres DB migrations scripts for golang-migrate.
├── lib                      <-- Common setup, configs used across all services.
├── compose.yaml             <-- Docker services setup(databases)
├── golangci.yml             <-- Config for golangci-lint. 
├── Makefile                 <-- Builds the whole app with exporting environment variables.
├── main.go                  <-- Start all server concurrently, init logger, init db, env port check, graceful shutdown.
├── readme.md                <-- Readme for the whole app.

```

<p align="right"><a href="#instabid-wallet">↑ Top</a></p>

<!-- Data Flow (Hexagonal architecture) -->

### Data Flow (Hexagonal architecture)

    Incoming : Client --(JSON)-> REST Handlers --(DTO)-> Service --(Domain Object)-> RepositoryDB

    Outgoing : RepositoryDB --(Domain Object)-> Service --(DTO)-> REST Handlers --(JSON)-> Client

<p align="right"><a href="#instabid-wallet">↑ Top</a></p>

### Routes Planned

#### Auth-API(:8001)

* POST /login: Authenticate and log a user in and generate JWT access token.
* GET /verify: Verify the user's authentication token.
* POST /logout: Log out the currently authenticated user.
* POST /reset-password: Reset the password for a given user.

#### User-API(:8000)

* GET /users: Retrieve a paginated list of users with their profiles, with optional filtering.
* POST /users/: Register a new user.
* GET /users/:user_id: Fetch details for a specific user by ID.
* POST /users/:user_id: Create profile details for a specific user by ID.
* PUT /users/:user_id: Update the details of a specific user by ID.
* DELETE /users/:user_id: Delete a specific user by ID.
* GET /users/:user_id/profile: Fetch the profile details of a specific user by ID.
* PUT /users/:user_id/profile: Update the profile details of a specific user by ID.

#### Wallet-API(:8002)

* POST /users/:user_id/wallet: Create a new wallet for a specific user by ID.
* GET /users/:user_id/wallet: Fetch the wallet details of a specific user by ID.
* POST /users/:user_id/wallet/transfer: Transfer money from the wallet of a specific user.
* GET /users/:user_id/wallet/transactions: Retrieve the transaction history for a specific user's wallet.

#### Auction-API(:8003)

* POST /auctions: Create a new auction.
* GET /auctions: Retrieve a list of all auctions.
* GET /auctions/:auction_id: Fetch details of a specific auction by ID.
* POST /auctions/:auction_id/bid: Place a bid on a specific auction by ID.
* GET /auctions/:auction_id/bids: Retrieve a list of bids for a specific auction by ID.

#### Admin-API(:8004)

* PUT /users/:user_id/role: Assign or modify the role of a specific user.
* GET /users/inactive: Retrieve a list of all inactive users.
* PUT /users/:user_id/activate: Activate a specific user by ID.
* PUT /users/:user_id/deactivate: Deactivate a specific user by ID.

#### MISC

* GET /health (Health check endpoint for monitoring and maintenance.)

<p align="right"><a href="#instabid-wallet">↑ Top</a></p>

<!-- CONTACT -->

### Contact

Ashef Tishad - [@ashef](https://www.linkedin.com/in/ashef/)

Project Link: [https://github.com/ashtishad/instabid-wallet](https://github.com/ashtishad/instabid-wallet)

<p align="right"><a href="#instabid-wallet">↑ Top</a></p>

<!-- Credits -->

### Credits

Readme template [Readme Template](https://github.com/othneildrew/Best-README-Template)

Badges and Icons [Shields.io](https://shields.io/)

<p align="right"><a href="#instabid-wallet">↑ Top</a></p>


<!-- MARKDOWN LINKS & IMAGES -->
<!-- Github -->

[forks-shield]: https://img.shields.io/github/forks/ashtishad/instabid-wallet?logo=github&style=for-the-badge

[forks-url]: https://github.com/ashtishad/instabid-wallet/network/members

[stars-shield]: https://img.shields.io/github/stars/ashtishad/instabid-wallet?logo=github&style=for-the-badge

[stars-url]: https://github.com/ashtishad/instabid-wallet/stargazers

<!-- Social -->

[linkedin-shield]: https://img.shields.io/badge/-LinkedIn-black.svg?style=for-the-badge&logo=linkedin&colorB=555

[linkedin-url]: https://www.linkedin.com/in/ashef/

<!-- Language -->

[go]: https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white

<!-- Database -->

[postgres]: https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white

[elastic-search]: https://img.shields.io/badge/Elastic_Search-005571?style=for-the-badge&logo=elasticsearch&logoColor=white

[redis]: https://img.shields.io/badge/redis-%23DD0031.svg?&style=for-the-badge&logo=redis&logoColor=white

<!-- Cloud -->

[docker]: https://img.shields.io/badge/Docker-2CA5E0?style=for-the-badge&logo=docker&logoColor=white

[aws]: https://img.shields.io/badge/Amazon_AWS-FF9900?style=for-the-badge&logo=amazonaws&logoColor=white

[github-actions]: https://img.shields.io/badge/GitHub_Actions-2088FF?style=for-the-badge&logo=github-actions&logoColor=white

<!-- Libraries -->

[apache-kafka]: https://img.shields.io/badge/Apache_Kafka-231F20?style=for-the-badge&logo=apache-kafka&logoColor=white

[jwt]: https://img.shields.io/badge/JWT-000000?style=for-the-badge&logo=JSON%20web%20tokens&logoColor=white

[swagger]: https://img.shields.io/badge/Swagger-85EA2D?style=for-the-badge&logo=Swagger&logoColor=white

<!-- Blogs -->

[medium]: https://img.shields.io/badge/Medium-12100E?style=for-the-badge&logo=medium&logoColor=white

[sponsor]: https://img.shields.io/badge/sponsor-30363D?style=for-the-badge&logo=GitHub-Sponsors&logoColor=#white

