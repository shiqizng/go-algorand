local algotest = require("algotest")
acct1=algotest.makeAccount()
print("acct1: ", acct1)
txn, contractAddr = algotest.createAppFromConfig(acct1,"contract1")
print("contract addr: ",contractAddr)
-- print(algotest.createAppFromConfig("AAAAA","contract2"))
print("starting private network")
algotest.startPrivateNetwork()
appID=txn:submit("533e921d0920c6f82d2874daefa70f76") -- walletID.
print("created app id:", appID)
acct1="O6P5JCFLRQVQUZAQ3QGNGOGMICETX2CH7LJ2DWZKZOF5YE3YF7JQWJQPXQ"
expected={balance=999000,apps={{id=26,extrapages=1}},totalapps=1}
algotest.assertAccountStates(acct1,expected)
