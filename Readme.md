```
cd cli  
go build  
$env:NODE_ID = "3000"
cli createwallet  
```

```
output : Your new address: [address]  
```

```
cli createblockchain -address [address]  
cli getbalance -address [address]  
cli createwallet  
```

```
output : Your new address: [receiverAddress]  
```

```
cli send -from [address] -to [receiverAddress] -amount 6 -mine   
cli getbalance -address [address]  
```