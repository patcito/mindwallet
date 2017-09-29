var EC = require('elliptic').ec;
var EdDSA = require('elliptic').eddsa;
var pbkdf2 = require('pbkdf2');
var scrypt = require('scrypt-js');
var bs58 = require('bs58');
var shajs = require('sha.js')
var Ripemd160 = require('ripemd160');
var keccak = require('keccak');
var BN = require('bn.js');

var s256 = new EC('secp256k1');
var ed25519 = new EdDSA('ed25519');

function sha256(s) {
  return shajs('sha256').update(s).digest();
}

function ripemd160(s) {
  return (new Ripemd160()).update(s).digest();
}

function keccak256(s) {
  return keccak('keccak256').update(s).digest();
}
function keccak256(s) {
  return keccak('keccak256').update(s).digest();
}

function reduce32(s) {
  return (new BN(s, 'le').mod(new BN('1000000000000000000000000000000014def9dea2f79cd65812631a5cf5d3ed', 16))).toArrayLike(Buffer, 'le', 32);
}

function b58checkencode(version, buffer) {
  buffer = Buffer.concat([Buffer.alloc(1, version), buffer])
  var hash = sha256(sha256(buffer));
  buffer = Buffer.concat([buffer, hash.slice(0, 4)]);
  var encoded = bs58.encode(buffer);
  return encoded;
}

function keyToBitcoinish(key, version) {
  var pubkey = new Buffer(s256.keyFromPrivate(key).getPublic(false, 'hex'), 'hex');
  return {
    private: b58checkencode(version + 0x80, key),
    public: b58checkencode(version, ripemd160(sha256(pubkey))),
  };
}

function keyToBitcoin(key) {
  return keyToBitcoinish(key, 0);
}

function keyToLitecoin(key) {
  return keyToBitcoinish(key, 48);
}

function keyToEthereum(key) {
  var pubkey = new Buffer(s256.keyFromPrivate(key).getPublic(false, 'hex'), 'hex');
  return {
    private: key.toString('hex'),
    public: '0x' + keccak256(pubkey.slice(1)).slice(12).toString('hex')
  };
}

function keyToMonero(seed) {
  var private_spend = reduce32(seed);
  var private_view = reduce32(keccak256(private_spend));

  // Hack
  var kp = ed25519.keyFromSecret()
  kp._privBytes = Array.from(private_spend);
  var public_spend = new Buffer(kp.pubBytes());
  var kp = ed25519.keyFromSecret()
  kp._privBytes = Array.from(private_view);
  var public_view = new Buffer(kp.pubBytes());


  var address_buf = Buffer.concat([Buffer.alloc(1, 0x12), public_spend, public_view])
  address_buf = Buffer.concat([address_buf, keccak256(address_buf).slice(0,4)]);
  var address = ''
  for (var i = 0; i < 8; i++) {
    address += bs58.encode(address_buf.slice(i*8, i*8+8));
  }
  address += bs58.encode(address_buf.slice(64, 69));

  return {
    private_spend: private_spend.toString('hex'),
    private_view: private_view.toString('hex'),
    public_spend: public_spend.toString('hex'),
    public_view: public_view.toString('hex'),
    public: address
  }
}



function warpwallet(password, salt, power, hashSuffix, callback) {
  var password_buffer = Buffer.from(password, 'utf-8');
  var salt_buffer = Buffer.from(salt, 'utf-8');
  var x1 = Buffer.alloc(1, hashSuffix);
  var x2 = Buffer.alloc(1, hashSuffix + 1);

  scrypt(Buffer.concat([password_buffer, x1]), Buffer.concat([salt_buffer, x1]), Math.pow(2, power), 8, 1, 32, function(error, progress, key1) {
    if(key1) {
      pbkdf2.pbkdf2(Buffer.concat([password_buffer, x2]), Buffer.concat([salt_buffer, x2]), Math.pow(2, 16), 32, 'sha256', function(err, key2) {
        for (var i = 0; i < 32; i++) {
          key2[i] = key2[i] ^ key1[i];
        }
        //console.log(key2.toString('hex'));
        callback(1, key2);
      });
    }
    callback(progress, null);
  });
}

var currencies = {
  bitcoin: {
    fn: keyToBitcoin,
    hashSuffix: 1,
  },
  litecoin: {
    fn: keyToLitecoin,
    hashSuffix: 2,
  },
  monero: {
    fn: keyToMonero,
    hashSuffix: 3,
  },
  ethereum: {
    fn: keyToEthereum,
    hashSuffix: 4
  }
}

function generateWallet(passphrase, salt, currency, callback) {
  warpwallet(passphrase, salt, 18, currencies[currency].hashSuffix, function(progress, result) {
    if(result) {
      var wallet = currencies[currency].fn(result);
      callback(1, wallet)
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
