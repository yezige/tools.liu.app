import { setFancybox, setToggle, loadComments, delEmptyThtd, setMenu, setLazyload, setVideoMove} from './util.js'

$(function() {
  delEmptyThtd()
  setFancybox()
  setToggle()
  setMenu()
  // lozad.js 需要手动开启，lazysizes.js 不需要
  setLazyload()
  setVideoMove()
})
window.loadComments = loadComments
