go test -buildmode=exe -run TestDeployChain
pause
go test -buildmode=exe -run TestDeployContractOnly
pause
go test -buildmode=exe -run TestDeployContractAndSpawn
