function setKeyResult(name, show, label, value){
  i('cnt-' + name).style.display = show ? 'block' : 'none';
  i('label-' + name).textContent = show ? label : '';
  i('result-' + name).value = show ? value : '';
}

function titleCase(string) {
  return string.charAt(0).toUpperCase() + string.slice(1);
}

function generate() {
  var passphrase = document.getElementById('passphrase').value;
  var salt = document.getElementById('salt').value;
  var currency = document.querySelector('[name=coin]:checked').value;
  var progress_blk = i('progress');
  var progress_bar = i('progress-bar');

  for(var el of document.querySelectorAll('.form input')) el.disabled = true;
  i('btn').style.display = 'none';

  mw.generateWallet(passphrase, salt, currency, function(progress, result) {
    progress_blk.style.display = !result ? 'block' : 'none';
    progress_bar.style.width = progress*100 + '%';
    if(result){
      for(var el of document.querySelectorAll('.form input')) el.disabled = false;
      i('wallet-title').style.display = 'block';
      i('wallet-title').textContent = 'Your ' + titleCase(currency) + ' wallet',

      setKeyResult('public', true, 'Public ' + titleCase(currency) + ' address', result.public);
      setKeyResult('private', currency != 'monero', 'Private key', result.private);

      setKeyResult('private-spend', currency == 'monero', 'Private Spend key', result.private_spend);
      setKeyResult('public-spend',  currency == 'monero', 'Public Spend key', result.public_spend);
      setKeyResult('public-view',   currency == 'monero', 'Public View key', result.public_view);
      setKeyResult('private-view',  currency == 'monero', 'Private View key', result.private_view);
    }
  });
}

function i(id){
  return document.getElementById(id);
}

function updateButton(){
  var passphrase = i('passphrase').value;
  var salt = i('salt').value;
  var button = i('btn');
  resetState();
  if(passphrase.length == 0 || salt.length <=7) {
    button.disabled = true;
    button.textContent = 'Generate';
  } else if (passphrase.length < 12) {
    button.disabled = false;
    button.textContent = 'Consider a larger passpharase';
  } else {
    button.disabled = false;
    button.textContent = 'Generate';
  }
}

function resetState(){
  if(passphrase.length == 0) {
    button.disabled = true;
    button.textContent = 'Generate';
  }   setKeyResult('public',        false);
  setKeyResult('private',       false);
  setKeyResult('private-spend', false);
  setKeyResult('public-spend',  false);
  setKeyResult('public-view',   false);
  setKeyResult('private-view',  false);

  i('btn').style.display = 'block';
  i('wallet-title').style.display = 'none';
}

function toggle(id) {
  var el = i(id);
  el.style.display = el.style.display == 'block' ? 'none' : 'block';
}

function updateOnline() {
  document.getElementById('online-alert').style.display = navigator.onLine ? 'block' : 'none';
}

function init() {
  window.addEventListener('online', updateOnline);
  window.addEventListener('offline', updateOnline);
  updateOnline();

  if(location.hash.length > 3) {
    var check = document.querySelector('input[name="coin"][value="' + location.hash.replace(/\W/g, '') + '"]');
    if(check) check.checked = true;
  }
}
init();
