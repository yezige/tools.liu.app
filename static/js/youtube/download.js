import { docCookies } from '../cookie.js'

const toggleFast = function(showfast) {
  const display = showfast ? 'flex' : 'none'
  const nodisplay = showfast ? 'none' : 'flex'
  document.querySelectorAll('.item.down-fast').forEach((item) => {
    item.style.display = display
  })
  document.querySelectorAll('.item.down-standard').forEach((item) => {
    item.style.display = nodisplay
  })
  setDownOpt()
}

const toggleAudio = function(showaudio) {
  const display = showaudio ? 'flex' : 'none'
  const nodisplay = showaudio ? 'none' : 'flex'
  document.querySelectorAll('.download_content').forEach((row) => {
    if (row.querySelector('[data-audio="1"]')) {
      row.style.display = display
    } else {
      row.style.display = nodisplay
    }
  })
  setDownOpt()
}
// 监听下载选项事件
const listenerDownOpt = function() {
  const to_fast = document.querySelector('#to-fast')
  const has_audio = document.querySelector('#has-audio')
  // 控制切换展示极速下载
  to_fast.addEventListener(
    'click',
    function() {
      toggleFast(document.querySelector('#to-fast').checked)
    },
    { passive: true }
  )

  // 控制切换展示有音频
  has_audio.addEventListener(
    'click',
    function() {
      toggleAudio(document.querySelector('#has-audio').checked)
    },
    { passive: true }
  )
}
listenerDownOpt()

// 设置下载选项设置
const setDownOpt = function() {
  const to_fast = document.querySelector('#to-fast')
  const has_audio = document.querySelector('#has-audio')

  docCookies.setItem('downOpt', JSON.stringify({
    to_fast: to_fast.checked,
    has_audio: has_audio.checked
  }))
}

// 企取得下载选项设置
const getDownOpt = function() {
  const downOpt = docCookies.getItem('downOpt') || '{}'
  return JSON.parse(downOpt)
}

// 展示下载选项设置
const showDownOpt = function() {
  const to_fast = document.querySelector('#to-fast')
  const has_audio = document.querySelector('#has-audio')
  const downOpt = getDownOpt()
  if (downOpt.to_fast) {
    to_fast.click()
  }
  if (downOpt.has_audio) {
    has_audio.click()
  }
}
showDownOpt()
