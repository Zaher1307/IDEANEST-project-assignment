# IDEANEST-project-assignment

## Design Choices

### Project Structure

- Decoupled packages for different layers of the application
  - `cmd` contains the main package for application interface (REST API)
  - `internal` contains private packages for core application logic (`business`) decoupled from the application interface and from the data access layer (`database`)
- This project structure allows greater flexibility when it comes to extending the application or adding another clients/databases because the `business` package operates separately from the database or its clients.

```
.
├── cmd
│   ├── main.go
│   ├── middlewares.go
│   └── handlers.go
├── internal
│   ├── auth
│   │   └── auth.go
│   ├── business
│   │   └── business.go
│   ├── types
│   │   └── types.go
│   └── database
│       └── database.go
├── docker-compose.yaml
├── Dockerfile
├── go.mod
├── go.sum
└── README.md
```

### Code Architecture

Code Architecture follows from the project structure, I started with the idea of high loose-coupling in mind. As mentioned before, this provides great flexibility in adding new clients (e.g. gRPC alongside REST) and different database solutions (e.g. postgres).<br><br>With that said, here is a brief description of functionality encompassed in each file:

- `cmd/main.go` : contains the main function of the application that starts a `gin` web server and attatches handlers to endpoints
- `cmd/middlewares.go` : contains the authentication middleware for protected endpoins and any future middlewares
- `cmd/handlers.go` : contains handlers code for interfacing with the REST client and refine the http request data to be passed to core application logic.
- `internal/auth/auth.go`: contains authentication logic that handles creating/revoking tokens.
- `internal/business/business.go`: contains core application logic, this is the true API of the application, which can be used by different clients.
- `internal/types/types.go`: contains types for the core functionality, these types are used all across the application code to keep consistency and to decouple how the application operates on data from how data is stored in whatever backing database, so that when trying to use different database, all application code won't need to change.
- `internal/database/database.go`: contains the data access layer for the application, its main job is to operate as an interface to the database and to smoothly handle the conversion between core application types and whatever format these types are actually stored in the database.

### Database Schema

Data entities are straightforward, the application consists of 2 entities:

- **_User_**:
  - Name (string)
  - Email (string)
  - Password (string)
  - Orgs (array [ ] )
- **_Organization_**:
  - Name (string)
  - Description (string)
  - OrgMembers (array [ ] )

```
├─ User
│    ├── Name
│    ├── Email
│    ├── Password
│    └── Orgs[]
│
└─ Organization
     ├── Name
     ├── Description
     └── OrgMembers[]
          ├── Name
          ├── Email
          └── AccessLevel
```

#### Notes on the schema:

- To avoid a full scan of the database when reading all organizations of a user, I added an index-like field in the User collection to cache IDs of their organizations.

## Running the application

```bash
$ git clone https://github.com/Zaher1307/IDEANEST-project-assignment.git
$ cd IDEANEST-project-assignment
$ docker compose up
```

### Notes

1. It was not clear that refresh token is to be revoked after a certain amount of time, also the fact that we have a dedicated endpoint to revoke the refresh token made it clear that we don’t want to revoke the refresh token automatically.

**Action**: I assumed that the endpoint “POST /refresh-token” only generates a new access token from a given refresh token.

---

2. It was not very clear what are the data fields in the application entities (e.g. User), for example what are the access levels and their privileges? Can any organization member add a new user to that organization?

**Action**: I assumed that there are 2 privilege levels (Admin, Member) for each organization, so that each organization has only one Admin, its original creator. And this means that only that Admin is privileged to successfully invoke the following endpoints

- `POST /organization/{organization_id}/invite`
- `PUT /organization/{organization_id}`
- `DELETE /organization/{organization_id}`

---

3. Given the previous 2 notes and decisions made upon them, It’s meaningless to have an endpoint to read everything in the system (read all organizations and their members) because it requires the user to be a member of all organizations.

**Action**: I assumed that the endpoint `GET /organization` only reads all organizations that the authorized user is a member of.
