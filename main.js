var EC = require("elliptic").ec;
var EdDSA = require("elliptic").eddsa;
var pbkdf2 = require("pbkdf2");
var scrypt = require("scrypt-js");
var bs58 = require("bs58");
var shajs = require("sha.js");
var Ripemd160 = require("ripemd160");
var keccak = require("keccak");
var BN = require("bn.js");
var argon2 = require("argon2-browser");
var Worker = require("./argon.worker.js");

var s256 = new EC("secp256k1");
var ed25519 = new EdDSA("ed25519");

function sha256(s) {
  return shajs("sha256")
    .update(s)
    .digest();
}

function ripemd160(s) {
  return new Ripemd160().update(s).digest();
}

function keccak256(s) {
  return keccak("keccak256")
    .update(s)
    .digest();
}
function keccak256(s) {
  return keccak("keccak256")
    .update(s)
    .digest();
}

function reduce32(s) {
  return new BN(s, "le")
    .mod(
      new BN(
        "1000000000000000000000000000000014def9dea2f79cd65812631a5cf5d3ed",
        16
      )
    )
    .toArrayLike(Buffer, "le", 32);
}

function b58checkencode(version, unbufferedKey) {
  var buffer = Buffer.from(unbufferedKey);
  buffer = Buffer.concat([Buffer.alloc(1, version), buffer]);
  var hash = sha256(sha256(buffer));
  buffer = Buffer.concat([buffer, hash.slice(0, 4)]);
  var encoded = bs58.encode(buffer);
  return encoded;
}

function keyToBitcoinish(key, version) {
  var pubkey = new Buffer(
    s256.keyFromPrivate(key).getPublic(false, "hex"),
    "hex"
  );
  return {
    private: b58checkencode(version + 0x80, key),
    public: b58checkencode(version, ripemd160(sha256(pubkey)))
  };
}

function keyToBitcoin(key) {
  return keyToBitcoinish(key, 0);
}

function keyToLitecoin(key) {
  return keyToBitcoinish(key, 48);
}

function keyToEthereum(key) {
  var pubkey = new Buffer(
    s256.keyFromPrivate(key).getPublic(false, "hex"),
    "hex"
  );
  return {
    private: key.toString("hex"),
    public:
      "0x" +
      keccak256(pubkey.slice(1))
        .slice(12)
        .toString("hex")
  };
}

function keyToMonero(seed) {
  var private_spend = reduce32(seed);
  var private_view = reduce32(keccak256(private_spend));

  // Hack
  var kp = ed25519.keyFromSecret();
  kp._privBytes = Array.from(private_spend);
  var public_spend = new Buffer(kp.pubBytes());
  var kp = ed25519.keyFromSecret();
  kp._privBytes = Array.from(private_view);
  var public_view = new Buffer(kp.pubBytes());

  var address_buf = Buffer.concat([
    Buffer.alloc(1, 0x12),
    public_spend,
    public_view
  ]);
  address_buf = Buffer.concat([
    address_buf,
    keccak256(address_buf).slice(0, 4)
  ]);
  var address = "";
  for (var i = 0; i < 8; i++) {
    address += bs58.encode(address_buf.slice(i * 8, i * 8 + 8));
  }
  address += bs58.encode(address_buf.slice(64, 69));

  return {
    private_spend: private_spend.toString("hex"),
    private_view: private_view.toString("hex"),
    public_spend: public_spend.toString("hex"),
    public_view: public_view.toString("hex"),
    public: address
  };
}
/*
argon2.hash({
    // required
    pass: 'password',
    salt: 'salt',
    // optional
    time: 1, // the number of iterations
    mem: 1024, // used memory, in KiB
    hashLen: 24, // desired hash length
    parallelism: 1, // desired parallelism (will be computed in parallel only for PNaCl)
    type: argon2.ArgonType.Argon2d, // or argon2.ArgonType.Argon2i
    distPath: '' // asm.js script location, without trailing slash
})
// result
.then(res => {
    console.log(res.hash) // hash as Uint8Array
    console.log(res.hashHex) // hash as hex-string
    console.log(res.encoded) // encoded hash, as required by argon2
})
// or error
.catch(err => {
    console.log(err.message) // error message as string, if available
    console.log(err.code) // numeric error code
})
*/
function allocateArray(strOrArr) {
  var arr =
    strOrArr instanceof Uint8Array || strOrArr instanceof Array
      ? strOrArr
      : intArrayFromString(strOrArr);
  return allocate(arr, "i8", ALLOC_NORMAL);
}
/*	salt		The salt to use, at least 8 characters
	-i		Use Argon2i (this is the default)
	-d		Use Argon2d instead of Argon2i
	-id		Use Argon2id instead of Argon2i
	-t N		Sets the number of iterations to N (default = 3)
	-m N		Sets the memory usage of 2^N KiB (default 12)
	-k N		Sets the memory usage of N KiB (default 4096)
	-p N		Sets parallelism to N threads (default 1)
	-l N		Sets hash output length to N bytes (default 32)
	-e		Output only encoded hash
	-r		Output only the raw bytes of the hash
	-v (10|13)	Argon2 version (defaults to the most recent version, currently 13)
	-h		Print ./argon2 usage

v*/

function now() {
  return performance ? performance.now() : Date.now();
}
function warpwallet(password, salt, power, hashSuffix, callback) {
  /*var arg = {};
    var t_cost = arg && arg.time || 10;
    var m_cost = arg && arg.mem || 4096;
    var parallelism = arg && arg.parallelism || 1;
    var pwd = allocateArray(arg && arg.pass || 'password');
    var pwdlen = arg && arg.pass ? arg.pass.length : 8;
    var salt = allocateArray(arg && arg.salt || 'somesalt');
    var saltlen = arg && arg.salt ? arg.salt.length : 8;
    var hash = allocate(new Array(arg && arg.hashLen || 32), 'i8', ALLOC_NORMAL);
    var hashlen = arg && arg.hashLen || 32;
    var encoded = allocate(new Array(512), 'i8', ALLOC_NORMAL);
    var encodedlen = 512;
    var argon2_type = arg && arg.type || 0;
    var version = 0x13;
*/
  var password_buffer = Buffer.from(password, "utf-8");
  var salt_buffer = Buffer.from(salt, "utf-8");
  var x1 = Buffer.alloc(1, hashSuffix);
  var x2 = Buffer.alloc(1, hashSuffix + 1);

  var workerData = {
    pass: password,
    salt: salt,
    hashSuffix: hashSuffix
  };
  const worker = new Worker();

  worker.postMessage(workerData);
  worker.onmessage = function(event) {
    console.log(event.data);

    if (event.data.hashSuffix != 4) {
      callback(event.data.status, event.data.bufferkey);
    } else {
      callback(event.data.status, event.data.strkey);
    }
  };
  //,,,
  console.log("no worker", argon2.ArgonType.Argon2i);
  callback(0.1, null);
  //cheating a bit here to show some progress and not make the user leave in dispair
  setTimeout(function() {
    callback(0.2, null);
  }, 5000);
  setTimeout(function() {
    callback(0.4, null);
  }, 15000);
  setTimeout(function() {
    callback(0.6, null);
  }, 30000);
  setTimeout(function() {
    callback(0.7, null);
  }, 35000);
  /*        argon2.hash({
            // required
            pass: password,
            salt: salt,
            // optional
            time: 262144, // the number of iterations
            mem: 8, // used memory, in KiB
            hashLen: 32, // desired hash length
            parallelism: 1, // desired parallelism (will be computed in parallel only for PNaCl)
            type: argon2.ArgonType.Argon2i, // or argon2.ArgonType.Argon2i
            distPath: '' // asm.js script location, without trailing slash
        })
        // result
            .then(res => {
                callback(0.3, null);
                console.log(res.hash) // hash as Uint8Array
                console.log(res.hashHex) // hash as hex-string
                console.log(res.encoded) // encoded hash, as required by argon2
                var key1 = res.hash;
                pbkdf2.pbkdf2(Buffer.concat([password_buffer, x2]), Buffer.concat([salt_buffer, x2]), Math.pow(2, 16), 32, 'sha256', function(err, key2) {
                    for (var i = 0; i < 32; i++) {
                        key2[i] = key2[i] ^ key1[i];
                    }
                    //console.log(key2.toString('hex'));
                    callback(1, key2);
                });
            })
        // or error
            .catch(err => {
                console.log(err.message) // error message as string, if available
                console.log(err.code) // numeric error code
            })*/
  //callback(progress, null);
}

var currencies = {
  bitcoin: {
    fn: keyToBitcoin,
    hashSuffix: 1
  },
  litecoin: {
    fn: keyToLitecoin,
    hashSuffix: 2
  },
  monero: {
    fn: keyToMonero,
    hashSuffix: 3
  },
  ethereum: {
    fn: keyToEthereum,
    hashSuffix: 4
  }
};

function generateWallet(passphrase, salt, currency, callback) {
  warpwallet(passphrase, salt, 18, currencies[currency].hashSuffix, function(
    progress,
    result
  ) {
    if (result) {
      var wallet = currencies[currency].fn(result);
      callback(1, wallet);
    } else {
      callback(progress, null);
    }
  });
}

module.exports = {
  generateWallet: generateWallet
};

//warpwallet('hello', 'a@b.c', 10);
//warpwallet('hello', 'a@b.c', 18);
