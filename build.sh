VERSION="0.1"

./node_modules/.bin/webpack
csplit index_template.html /INLINE_SCRIPT_HERE/

{
  cat xx00
  echo '<script>'
  cat ui.js
  cat dist/bundle.js
  echo '</script>'
  tail -n +2 xx01
} > out.html

rm xx00 xx01
name=memwallet_${VERSION}_SHA256_$(sha256sum out.html  | cut -d' ' -f1).html
mv out.html $name

sed "s/REDIR_PAGE/$name/" redir_template.html > index.html

