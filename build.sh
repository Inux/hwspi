export GOOS=linux
export GOARCH=arm
export GOARM=6

go build ...hwspi
sudo -E go install ...hwspi
