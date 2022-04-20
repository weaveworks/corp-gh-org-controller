load('ext://restart_process', 'docker_build_with_restart')

DIRNAME = os.path.basename(os. getcwd())
IMG='controller:latest'
CONTROLLERGEN='crd:trivialVersions=true rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases;'

def yaml():
    return local('cd config/manager; kustomize edit set image controller=' + IMG + '; cd ../..; kustomize build config/default')

def manifests():
    return 'controller-gen ' + CONTROLLERGEN

def generate():
    return 'controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./...";'

def vetfmt():
    return 'go vet ./...; go fmt ./...'

def binary():
    return 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -o bin/manager main.go'

deps = ['controllers', 'main.go', 'api']

local(manifests() + generate())

local_resource('crd', manifests() + 'kustomize build config/crd | kubectl apply -f -', deps=["api"])

k8s_yaml(yaml())

local_resource('watch', generate() + binary(), deps=deps, ignore=['*/*/zz_generated.deepcopy.go'])

local_resource('sample', 'kubectl apply -f ./config/samples', deps=["./config/samples"], resource_deps=[DIRNAME + "-controller-manager"])

docker_build_with_restart(IMG, '.',
 dockerfile='tilt.docker',
 entrypoint='/manager',
 only=['./bin/manager'],
 live_update=[
       sync('./bin/manager', '/manager'),
   ]
)
