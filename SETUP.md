# Setup steps

## 0. On each node, deploy the latest version of criu and cuda-checkpoint binary
### 0.1 Setup NVIDIA driver 
```bash
sudo ubuntu-driver install 570-server
```

### 0.2 Setup CRIU
```bash
# Add ppa
sudo add-apt-repository ppa:checkpoint-restore/criu
sudo apt update -y
sudo apt install criu -y
```
- To setup custom build for gpu support, grab it from the release of this repo:
  - https://github.com/pedestrianlove/criu
  - and install it with `sudo dpkg -i criu_<version>.deb`.

## 1. Install microk8s on 2 nodes and join them together
```bash
sudo snap install microk8s --classic
```

### 1.1 Setup Strimzi(kafka) on the cluster
```bash
microk8s kubectl create namespace kafka
microk8s kubectl create -f "https://strimzi.io/install/latest?namespace=kafka" -n kafka
microk8s kubectl apply -f yaml/kafka/kafka-ephemeral-single.yaml -n kafka
```

### 1.2 Add your user to the microk8s group
```bash
sudo usermod -aG microk8s $USER
```

### 1.3 Enable addons in microk8s
```bash
microk8s enable cert-manager gpu metrics-server
```

## 2. Deploy the Checkpoint Operator
```bash
# Generate manifests
make manifest
# Install CRDs
make install
# Deploy the operator (built using github action in this repo)
make deploy IMG=ghcr.io/pedestrianlove/checkpoint-operator:latest
```

## 3. Deploy your own workload in the cluster 
- Take PyTorch training as an example
```bash
microk8s kubectl apply -f config/samples/train_pytorch.yaml
```

## 4. Deploy the migration CRD spec (TBD)
- Currently not working at all.
```bash
microk8s kubectl apply -f config/samples/migrate_pytorch.yaml
```
