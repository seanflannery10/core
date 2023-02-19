# Load ENV Settings
load('ext://dotenv', 'dotenv')
dotenv()

POSTGRES_USER = os.getenv('POSTGRES_USER')
POSTGRES_PASSWORD = os.getenv('POSTGRES_PASSWORD')
POSTGRES_DB = os.getenv('POSTGRES_DB')

# Tests
load('ext://tests/golang', 'test_go')
test_go('test-core-cmd', './cmd/...', './cmd')
test_go('test-core-internal', './internal/...', './internal')

# Build App
local_resource(
  'core-compile',
  'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -gcflags="all=-N -l" -o ./bin/core ./cmd/core',
   deps=['./cmd/core/', './internal/'],
)

# Run App
dockerfile='''
FROM alpine
COPY /bin/ /
'''

load('ext://restart_process', 'docker_build_with_restart')
docker_build_with_restart(
  'core-image',
  '.',
  entrypoint='/core',
  dockerfile_contents=dockerfile,
  only=['./bin/'],
  live_update=[sync('./bin/', '/')],
)

core = '''
apiVersion: v1
kind: ConfigMap
metadata:
  name: core
  labels:
    app: core
data:
  DB_DSN: 'postgres://{USER}:{PASS}@postgres:5432/{DB}?sslmode=disable'
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: core
  labels:
    app: core
spec:
  selector:
    matchLabels:
      app: core
  template:
    metadata:
      labels:
        app: core
    spec:
      containers:
        - name: core
          image: core-image
          envFrom:
            - configMapRef:
                name: core
          ports:
            - containerPort: 4000
'''.format(USER=POSTGRES_USER, PASS=POSTGRES_PASSWORD, DB=POSTGRES_DB)

k8s_yaml(blob(core))
k8s_resource('core', port_forwards='4000', resource_deps=['postgres', 'core-compile'])

# Run App Migrations
migrations_dockerfile='''
FROM amacneil/dbmate
COPY /sql/migrations/ /db/migrations/
'''

docker_build(
  'core-migrations-image',
  '.',
  dockerfile_contents=migrations_dockerfile,
  only=['./sql/migrations/'],
)

core_migrations = '''
apiVersion: v1
kind: ConfigMap
metadata:
  name: core-migrations
  labels:
    app: core-migrations
data:
  DATABASE_URL: 'postgres://{USER}:{PASS}@postgres:5432/{DB}?sslmode=disable'
---
apiVersion: batch/v1
kind: Job
metadata:
  name: core-migrations
  labels:
    app: core-migrations
spec:
  template:
    spec:
      restartPolicy: Never
      containers:
      - name: core-migrations
        image: core-migrations-image
        command: ["/bin/sh", "-c", 'dbmate down; dbmate up']
        envFrom:
          - configMapRef:
              name: core-migrations
'''.format(USER=POSTGRES_USER, PASS=POSTGRES_PASSWORD, DB=POSTGRES_DB)

k8s_yaml(blob(core_migrations))
k8s_resource('core-migrations', resource_deps=['postgres'])

# Run Postgres
postgres = '''
apiVersion: v1
kind: ConfigMap
metadata:
  name: postgres
  labels:
    app: postgres
data:
  POSTGRES_USER: {USER}
  POSTGRES_PASSWORD: {PASS}
  POSTGRES_DB: {DB}
  PGUSER: {USER}
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: postgres
  labels:
    app: postgres
spec:
  selector:
    matchLabels:
      app: postgres
  template:
    metadata:
      labels:
        app: postgres
    spec:
      containers:
        - name: postgres
          image: postgres:15
          args:
            - postgres
            - -c
            - log_statement=all
          envFrom:
            - configMapRef:
                name: postgres
          ports:
            - containerPort: 5432
          startupProbe:
            exec:
              command:
                - /bin/sh
                - -c
                - exec pg_isready -h localhost
            periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: postgres
  labels:
    app: postgres
spec:
  ports:
    - port: 5432
      protocol: TCP
  selector:
    app: postgres
'''.format(USER=POSTGRES_USER, PASS=POSTGRES_PASSWORD, DB=POSTGRES_DB)

k8s_yaml(blob(postgres))
k8s_resource('postgres', port_forwards='5432')
