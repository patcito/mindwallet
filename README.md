## MindWallet

MindWallet is a deterministic cryptocurrency address generator heavily based on [MemWallet](https://github.com/dvdbng/memwallet) but using argon2 scrypt instead as hashing function,
it's like [WarpWallet](https://keybase.io/warp/), but it works for Ethereum, Litecoin, Monero and Bitcoin. You never need to save or store your private key anywhere. Just pick a really good password - many random words, for example - and never use it for anything else.

Given the same Passphrase and Salt, MindWallet will always generate the same address and private key, so you only need to remember your password to access your funds.

For more information on why this is safer than a regular brainwallet, see [WarpWallet](https://keybase.io/warp/)'s help, MindWallet is a re-implementation of WarpWallet, but it works for other currencies thanks to MemWallet who makes up for most of the code. WarpWallet and MemWallet use the same algorithm, so WarpWallet and MemWallet will generate the same Bitcoin address for a given Passphrase and salt.

Another difference between MindWallet and MemWallet is that MindWallet makes use of a web worker to make it faster, this forces you to make use of a web server.

This repo contains an implementation of MindWallet in JavaScript and Go.

### Instructions to run the html:

$ python -m SimpleHTTPServer

Then open your browser and go to http://localhost:8000

### Instructions to run the go executable:

First [install golang](https://golang.org/dl/).

Second install go dep:

On Mac:

$ brew install dep

Anywhere else:

$ go get -u github.com/golang/dep/cmd/dep

Then run:

$ dep ensure

Now build the binary:

$ go build mindwallet.go

Run it like that:

$ ./mindwallet bitcoin my_password my_salt

You can add more memory cost at the end by adding one the following options: strong, very_strong or ridiculously_strong. These can take minutes
to execute but as this is something you won't be re-uinsg a lot, it may not be a big deal.

Have fun!
