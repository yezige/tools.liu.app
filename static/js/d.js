import { setFancybox, setToggle, loadComments, delEmptyThtd, setMenu, setLazyload, setVideoMove, setProgress } from '/static/js/util.js'

document.addEventListener('DOMContentLoaded', function () {
  console.log('DOMContentLoaded')
  // lozad.js 需要手动开启，lazysizes.js 不需要
  setLazyload()
  setProgress()
  setFancybox()
})

setMenu()
setVideoMove()
delEmptyThtd()
setToggle()
