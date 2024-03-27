import { docCookies } from '../cookie.js'
import { showMask, setInputBox, setEles, getParams, ajax } from '../util.js'
import { FFmpeg } from '/static/node_modules/@ffmpeg/ffmpeg/dist/esm/index.js'
import { fetchFile, toBlobURL, downloadWithProgress } from '/static/node_modules/@ffmpeg/util/dist/esm/index.js'

const TXCLOUD_HOST = 'https://cloud1-5giq10fn52e7fc0e-1307628865.ap-shanghai.app.tcloudbase.com/cloud-function/httpFunction'

const toggleFast = function (showfast) {
  const fast_display = showfast ? 'flex' : 'none'
  const standard_display = showfast ? 'none' : 'flex'
  document.querySelectorAll('.item.down-fast').forEach((item) => {
    item.style.display = fast_display
    showfast ? item.classList.add('active') : item.classList.remove('active')
  })
  document.querySelectorAll('.item.down-standard').forEach((item) => {
    item.style.display = standard_display
    showfast ? item.classList.remove('active') : item.classList.add('active')
  })
  setDownOpt()
}

const toggleAudio = function (showaudio) {
  const audio_display = showaudio ? 'flex' : 'flex'
  const other_display = showaudio ? 'none' : 'flex'
  document.querySelectorAll('.download_content').forEach((row) => {
    if (row.querySelector('[data-audio="1"]')) {
      row.style.display = audio_display
    } else {
      row.style.display = other_display
    }
  })
  setDownOpt()
}
// 监听下载选项事件
const listenerDownOpt = function () {
  const to_fast = document.querySelector('#to-fast')
  const has_audio = document.querySelector('#has-audio')
  // 控制切换展示极速下载
  to_fast.addEventListener(
    'click',
    function () {
      toggleFast(document.querySelector('#to-fast').checked)
    },
    { passive: true }
  )

  // 控制切换展示有音频
  has_audio.addEventListener(
    'click',
    function () {
      toggleAudio(document.querySelector('#has-audio').checked)
    },
    { passive: true }
  )
}
listenerDownOpt()

// 设置下载选项设置
const setDownOpt = function () {
  const to_fast = document.querySelector('#to-fast')
  const has_audio = document.querySelector('#has-audio')

  docCookies.setItem(
    'downOpt',
    JSON.stringify({
      to_fast: to_fast.checked,
      has_audio: has_audio.checked
    })
  )
}

// 取得下载选项设置
const getDownOpt = function () {
  const downOpt = docCookies.getItem('downOpt') || '{}'
  return JSON.parse(downOpt)
}

// 展示下载选项设置
const showDownOpt = function () {
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

const load = async () => {
  const baseURL = 'https://unpkg.com/@ffmpeg/core@0.12.6/dist/esm'
  const ffmpeg = new FFmpeg()
  ffmpeg.on('log', ({ message }) => {
    console.log(message)
  })
  // toBlobURL is used to bypass CORS issue, urls with the same
  // domain can be used directly.
  await ffmpeg.load({
    coreURL: await toBlobURL(`${baseURL}/ffmpeg-core.js`, 'text/javascript'),
    wasmURL: await toBlobURL(`${baseURL}/ffmpeg-core.wasm`, 'application/wasm')
  })
  return ffmpeg
}

const ffmpegToDownload = async (data) => {
  const error_ele = document.querySelector('.vip_down_link_box .error')
  const progress_ele = document.querySelector('.vip_down_link_box .progress')
  error_ele.innerHTML = ''
  progress_ele.innerHTML = ''
  try {
    const ffmpeg = await load()
    ffmpeg.on('progress', ({ progress, time }) => {
      progress_ele.innerHTML = `${progress * 100} % (transcoded time: ${time / 1000000} s)`
    })
    await ffmpeg.writeFile(
      data.video_name,
      await downloadWithProgress(data.video_url, ({ total, received }) => {
        if (!total) {
          return false
        }
        progress_ele.innerHTML = `${(received / total).toFixed(2) * 100} %`
      })
    )
    await ffmpeg.writeFile(
      data.audio_name,
      await downloadWithProgress(data.audio_url, ({ total, received }) => {
        if (!total) {
          return false
        }
        progress_ele.innerHTML = `${(received / total).toFixed(2) * 100} %`
      })
    )
    await ffmpeg.exec(['-i', data.video_name, '-i', data.audio_name, '-f mp4', data.output_name])
    const ffdata = await ffmpeg.readFile(data.output_name)
    const downEle = document.getElementById('vip_down_link')
    downEle.src = URL.createObjectURL(new Blob([ffdata.buffer], { type: 'video/mp4' }))
    // 设置为加载完成，隐藏loading，显示下载链接
    setVipDownloadLoad('loaded')
  } catch (error) {
    error_ele.innerHTML = error
    // 设置为加载错误，隐藏loading、下载链接，显示错误
    setVipDownloadLoad('error')
  }
}

// VIP下载提交事件
const doVipDownload = async () => {
  if (!(await checkVipCode())) {
    return false
  }
  // 设置为加载中，隐藏数字输入框显示loading
  setVipDownloadLoad('loading')
  const id = getParams().get('id')
  // 按钮对象
  const vip = document.querySelector('.item.down-vip.active')
  if (!vip) {
    return false
  }
  // 视频对象
  const video = vip.closest('.download_content')
  const quality = video.getAttribute('data-quality')
  const itag = video.getAttribute('data-itag')
  const ext = video.getAttribute('data-ext')
  const video_name = `${id}_${quality}_${itag}`
  const video_name_ext = `${video_name}.${ext}`
  let video_url
  if (video.querySelector('.item.down.active')) {
    video_url = video.querySelector('.item.down.active a').getAttribute('href')
  } else {
    video_url = video.querySelector('.item.down.down-standard a').getAttribute('href')
  }
  // 音频对象
  let audio
  for (const row of document.querySelectorAll('.download_content[data-audio="1"]')) {
    audio = row
  }
  const audio_ext = audio.getAttribute('data-ext')
  const audio_name_ext = `${video_name}_audio.${audio_ext}`
  // 音频链接
  let audio_url
  if (audio.querySelector('.item.down.active')) {
    audio_url = audio.querySelector('.item.down.active a').getAttribute('href')
  } else {
    audio_url = audio.querySelector('.item.down.down-standard a').getAttribute('href')
  }
  const data = {
    id,
    video_url,
    audio_url,
    video_name: video_name_ext,
    audio_name: audio_name_ext,
    output_name: `${video_name}.mp4`
  }
  console.log('down-info', data)
  await ffmpegToDownload(data)
}

// 校验码是否正确
const checkVipCode = async () => {
  const error_ele = document.querySelector('.vip_code_box .error')
  error_ele.innerHTML = ''
  let pass_ele = document.querySelector('.input_box').innerText.replaceAll('\n', '').replaceAll(' ', '')
  if (!pass_ele) {
    error_ele.innerHTML = '密码不能为空'
    return false
  }
  const res = await ajax({
    url: `${TXCLOUD_HOST}?functionName=orderCouponUsed&couponCode=${pass_ele}`,
    method: 'GET'
  })
  if (!res.success) {
    error_ele.innerHTML = res.msg
    return false
  }
  error_ele.innerHTML = ''
  return true
}

// 设置VIP下载所需事件
const setVipDownload = async () => {
  document.querySelectorAll('.item.down-vip').forEach((item) => {
    item.addEventListener('click', function (e) {
      setEles(document, '.item.down-vip', (ele) => {
        ele.classList.remove('active')
      })
      e.currentTarget.classList.add('active')
      setVipDownloadLoad('init')
      showMask({
        title: '下载音视频合成版',
        content: `<div class="vip_download">
  <span class="statement">通过支付 ¥1.99，获得音频+视频合成后的结果，券码8小时内有效不限次数，退款无忧（使用后仍支持退款）</span>
  <div class="vip_pic_box">
    <a href="https://aii.ren/u/zcpQwubV" target="_blank"><span class="vip_pic_img"></span></a>
  </div>
  <div class="vip_code_box active">
    <div class="input_box">
      <!-- 假的，只用来唤起输入法输入，监听事件等 -->
      <div class="_input">
        <input class="_input_tel" type="tel">
      </div>
      <span class="input_block" id="input_block_0"></span>
      <span class="input_block" id="input_block_1"></span>
      <span class="input_block" id="input_block_2"></span>
      <span class="input_block" id="input_block_3"></span>
      <span class="input_block" id="input_block_4"></span>
      <span class="input_block" id="input_block_5"></span>
    </div>
    <div class="error"></div>
  </div>
  <div class="vip_down_link_box loading">
    <div class="loading_box"></div>
    <a id="vip_down_link" href="javascript:void(0);">点击下载</a>
    <div class="progress"></div>
    <div class="error"></div>
  </div>
</div>`,
        ok: '开始',
        cancle: '关闭',
        style: 'width: 260px; height: 440px;',
        callback_init: function (mask) {
          setInputBox(mask)
          setVipDownloadLoad('init')
        },
        callback_ok: function () {
          const down_link_box = document.querySelector('.vip_down_link_box')
          if (down_link_box.classList.contains('loading')) {
            return false
          }
          doVipDownload()
          return false
        },
        callback_cancle: function () {}
      })
    })
  })
}
setVipDownload()

const setVipDownloadLoad = (status) => {
  const code_box = document.querySelector('.vip_code_box')
  const down_link_box = document.querySelector('.vip_down_link_box')
  if (!down_link_box) {
    return false
  }
  if (status == 'init') {
    code_box.classList.add('active')
    down_link_box.classList.remove('active')
  } else {
    code_box.classList.remove('active')
    down_link_box.classList.add('active')
  }

  if (status == 'loading') {
    down_link_box.classList.add('loading')
  } else {
    down_link_box.classList.remove('loading')
  }

  if (status == 'loaded') {
    down_link_box.classList.add('loaded')
  } else {
    down_link_box.classList.remove('loaded')
  }

  if (status == 'error') {
    down_link_box.classList.add('error')
  } else {
    down_link_box.classList.remove('error')
  }
}
