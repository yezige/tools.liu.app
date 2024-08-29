import { setSearch, startProgress } from '/static/js/util.js'

document.addEventListener('DOMContentLoaded', function () {
  setSearch()

  for (const row of document.querySelectorAll('.start_progress')) {
    row.addEventListener('click', function() {
      startProgress()
    })
  }
})

