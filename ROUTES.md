# Echo Boilerplate

A simple Boilerplate for [Echo](https://github.com/labstack/echo)

## Routes

### Web

- **[GET] `/health-check`**: Check server

  ```bash
  http GET localhost:3001/health-check
  ```

- **[GET] `/metrics`**: Prometheus metrics
  ```bash
  http GET localhost:3001/metrics
  ```

### API

- **[POST] `/api/v1/login`**: Authentication

  ```bash
  http POST localhost:3001/api/v1/login username=test@gmail.com password=00000000
  ```

  Response:

  ```json
  {
    "id": "2a40080f-6077-4273-9075-1c5503ac95eb",
    "username": "test@gmail.com",
    "lastname": "Test",
    "firstname": "Toto",
    "created_at": "2021-03-08T20:43:28.345Z",
    "updated_at": "2021-03-08T20:43:28.345Z",
    "token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJjcmVhdGVkQXQiOiIyMDIxLTAzLTA4VDIwOjQzOjI4LjM0NVoiLCJleHAiOjE2MTYxMDAyMTUsImZpcnN0bmFtZSI6IkZhYmllbiIsImlhdCI6MTYxNTIzNjIxNSwiaWQiOjEsImxhc3RuYW1lIjoiQmVsbGFuZ2VyIiwibmJmIjoxNjE1MjM2MjE1LCJ1c2VybmFtZSI6InZhbGVudGlsQGdtYWlsLmNvbSJ9.RL_1C2tYqqkXowEi8Np-y3IH1qQLl8UVdFNWswcBcIOYB6W4T-L_RAkZeVK04wtsY4Hih2JE1KPcYqXnxj2FWg",
    "expires_at": "2021-03-18T21:43:35.641Z"
  }
  ```

- **[POST] `/api/v1/users`**: User creation

  ```bash
  http POST localhost:3001/api/v1/users "Authorization: Bearer eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiSUQiLCJ1c2VybmFtZSI6IlVzZXJuYW1lIiwibGFzdG5hbWUiOiJGaXJzdG5hbWUiLCJmaXJzdG5hbWUiOiJMYXN0bmFtZSIsImF1ZCI6IkNsaWVudCIsImV4cCI6MTY0OTg1MTcxNSwiaWF0IjoxNjQ5ODUwMjc1LCJpc3MiOiJBUEkiLCJuYmYiOjE2NDk4NTAyNzUsInN1YiI6IkFQSSBhdXRoZW50aWNhdGlvbiJ9.whL0nJDvbFkQji0gxbjiXrXkpK_Fm9xnlu_LEeXmD-yj21g22pvjs1hE-2ISQVIgzkxvmWLbIWN-Ki_T3p-RnQ" lastname=Test firstname=Toto username=test@gmail.com password=00000000
  ```

  Response:

  ```json
  {
    "id": "cb13cc29-13bb-4b84-bf30-17da00ec7400",
    "username": "test@gmail.com",
    "lastname": "Test",
    "firstname": "Toto",
    "created_at": "2021-03-09T21:05:35.564747+01:00",
    "updated_at": "2021-03-09T21:05:35.564747+01:00"
  }
  ```
