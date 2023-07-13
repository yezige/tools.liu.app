import { setSearch, startProgress } from '../util.js'

setSearch()

for (const row of document.querySelectorAll('.start_progress')) {
  row.addEventListener('click', function() {
    startProgress()
  })
}
