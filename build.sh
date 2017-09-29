./node_modules/.bin/webpack
csplit index.html /INLINE_SCRIPT_HERE/

{
  cat xx00
  echo '<script>'
  cat ui.js
  cat dist/bundle.js
  echo '</script>'
  tail -n +2 xx01
} > out.html


