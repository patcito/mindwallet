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

onmessage = function(e) {
  console.log(e.data);
  var data = e.data;
  var password = data.pass;
  var salt = data.salt;
  var hashSuffix = data.hashSuffix;

  var password_buffer = Buffer.from(password, "utf-8");
  var salt_buffer = Buffer.from(salt, "utf-8");
  var x1 = Buffer.alloc(1, hashSuffix);
  var x2 = Buffer.alloc(1, hashSuffix + 1);

  console.log("yes worker", argon2.ArgonType.Argon2i);

  argon2
    .hash({
      // required
      pass: password,
      salt: salt,
      // optional
      time: 262144, // the number of iterations
      mem: 8, // used memory, in KiB
      hashLen: 32, // desired hash length
      parallelism: 1, // desired parallelism (will be computed in parallel only for PNaCl)
      type: argon2.ArgonType.Argon2i, // or argon2.ArgonType.Argon2i
      distPath: "./" // asm.js script location, without trailing slash
    })
    .then(res => {
      postMessage({ status: 0.9, data: null });
      console.log(res.hash); // hash as Uint8Array
      console.log(res.hashHex); // hash as hex-string
      console.log(res.encoded); // encoded hash, as required by argon2
      var key1 = res.hash;
      pbkdf2.pbkdf2(
        Buffer.concat([password_buffer, x2]),
        Buffer.concat([salt_buffer, x2]),
        Math.pow(2, 16),
        32,
        "sha256",
        function(err, key2) {
          for (var i = 0; i < 32; i++) {
            key2[i] = key2[i] ^ key1[i];
          }
          console.log("key2", key2, "key2 tostring hex", key2.toString("hex"));
          var objData = {
            status: 1,
            hashSuffix: hashSuffix,
            strkey: key2.toString("hex"),
            bufferkey: key2
          };
          //postMessage({status: 1, data: key2.toString('hex'), hashSuffix: hashSuffix}, key2)
          console.log("type of ", typeof objData.bufferkey);
          postMessage(objData);
        }
      );
    })
    // or error
    .catch(err => {
      console.log(err.message); // error message as string, if available
      console.log(err.code); // numeric error code
    });
};
/*
            */
