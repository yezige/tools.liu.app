const timeout = (ms) => {
  return new Promise((resolve, reject) => {
    setTimeout(() => {
      resolve()
    }, ms)
  })
}
const setFancybox = function () {
  if (!CONFIG.page.site.fancybox) {
    return false
  }
  document.querySelectorAll('.content :not(a) > img, .content :not(.nofancybox) > img, .content > img').forEach((element) => {
    var $image = $(element)
    var imageLink = $image.attr('data-src') || $image.attr('src')
    var $imageWrapLink = $image.wrap(`<a class="fancybox fancybox.image" href="${imageLink}" itemscope itemtype="http://schema.org/ImageObject" itemprop="url"></a>`).parent('a')
    if ($image.is('.post-gallery img')) {
      $imageWrapLink.attr('data-fancybox', 'gallery').attr('rel', 'gallery')
    } else if ($image.is('.group-picture img')) {
      $imageWrapLink.attr('data-fancybox', 'group').attr('rel', 'group')
    } else {
      $imageWrapLink.attr('data-fancybox', 'default').attr('rel', 'default')
    }

    var imageTitle = $image.attr('title') || $image.attr('alt')
    if (imageTitle) {
      $imageWrapLink.append(`<p class="image-caption">${imageTitle}</p>`)
      // Make sure img title tag will show correctly in fancybox
      $imageWrapLink.attr('title', imageTitle).attr('data-caption', imageTitle)
    }

    if ($image.hasClass('inline')) {
      $imageWrapLink.addClass('inline')
    }
    if ($image.hasClass('center')) {
      $imageWrapLink.addClass('center')
    }
    if ($image.hasClass('right')) {
      $imageWrapLink.addClass('right')
    }
    if ($image.attr('width')) {
      $imageWrapLink.attr('width', $image.attr('width'))
    }
  })

  $.fancybox.defaults.hash = false
  $('.fancybox').fancybox({
    loop: true,
    helpers: {
      overlay: {
        locked: false
      }
    }
  })
}

const setToggle = function () {
  const menuToggle = document.querySelector('.toggle')
  const menu = document.querySelector('.menu')
  const header = document.querySelector('.showbody header')

  menuToggle.addEventListener('click', (e) => {
    menuToggle.classList.toggle('active')
    menu.classList.toggle('active')
    header.classList.toggle('active')
  })
  header.addEventListener('click', (e) => {
    e.stopPropagation()
  })
  document.addEventListener('click', () => {
    menuToggle.classList.remove('active')
    menu.classList.remove('active')
    header.classList.remove('active')
  })
}

const loadComments = function (element, callback) {
  if (!CONFIG.page.site.disqus.lazyload || !element) {
    callback()
    return
  }
  let intersectionObserver = new IntersectionObserver((entries, observer) => {
    let entry = entries[0]
    if (entry.isIntersecting) {
      callback()
      observer.disconnect()
    }
  })
  intersectionObserver.observe(element)
  return intersectionObserver
}

const delEmptyThtd = function () {
  document.querySelectorAll('thead th').forEach((element) => {
    if (element.innerHTML == '') {
      element.remove()
    }
  })
}

const setMenu = function () {
  document.querySelectorAll('.menu a').forEach((element) => {
    element.style.display = 'inline-block'
    if (element.getAttribute('href') == location.pathname) {
      element.classList.add('active')
    }
  })
}

const searchHandler = function () {
  const inputvalue = document.getElementById('search').value
  if (!inputvalue) return false

  location.href = '/youtube/search?q=' + inputvalue
}

// 搜索事件
const setSearch = function () {
  document.getElementById('search').addEventListener('keydown', (event) => {
    if (event.key === 'Enter') {
      searchHandler()
    }
  })
  document.querySelector('.search_submit').addEventListener('click', searchHandler)
}

// 显示加载中
const showLoading = function (els) {
  const loading_html = `
    <svg class="loading" width="120" height="120" viewBox="0 0 120 120" fill="none" xmlns="http://www.w3.org/2000/svg">
      <circle class="path" cx="60.5" cy="59.5" r="32.5" stroke="url(#paint0_linear_311_2)" stroke-width="4" stroke-linejoin="round" />
      <defs>
        <linearGradient id="paint0_linear_311_2" x1="28" y1="27" x2="93" y2="92" gradientUnits="userSpaceOnUse">
          <stop offset="0.283303" stop-color="#FFB74D" />
          <stop offset="0.482322" stop-color="#E57373" />
          <stop offset="0.753156" stop-color="#42A5F5" />
          <stop offset="0.890625" stop-color="#BDBDBD" />
        </linearGradient>
      </defs>
    </svg>`
  for (const el of els) {
    el.classList.add('loading_box')
    let loading_dom = el.querySelector('.loading')
    if (!loading_dom) {
      loading_dom = new DOMParser().parseFromString(loading_html, 'text/html').querySelector('.loading')
      el.append(loading_dom)
    }
    loading_dom.style.display = 'block'
  }
}

// 隐藏加载中
const hideLoading = function (els) {
  if (!els || !els.length) {
    return
  }
  for (const el of els) {
    const loading_dom = el.closest('.loading_box').querySelector('.loading')
    if (loading_dom) {
      loading_dom.remove()
    }
  }
}

// lozad.js 需要手动开启，lazysizes.js 不需要
const setLazyload = async function () {
  lozad(document.querySelectorAll('.lazyload'), {
    threshold: 0.1,
    enableAutoReload: true,
    load: async function (el) {
      if (el.getAttribute('data-src')) {
        el.src = el.getAttribute('data-src')
      }
      el.onload = function () {
        el.classList.remove('loading_box')
      }
    },
    loaded: function (el) {}
  }).observe()
}

const moveLeft = function (event) {
  const videoBox = event.currentTarget.closest('.popular_videos_box')
  const videos = videoBox.querySelector('.videos')
  if (!videoBox) {
    return false
  }

  let scrollLeft = videos.scrollLeft - videoBox.offsetWidth
  videos.scrollTo({
    left: scrollLeft,
    behavior: 'smooth'
  })

  if (scrollLeft <= 0) {
    videoBox.querySelector('.move_left_btn').classList.remove('active')
  }
}
const moveRight = function (event) {
  const videoBox = event.currentTarget.closest('.popular_videos_box')
  const videos = videoBox.querySelector('.videos')
  if (!videoBox) {
    return false
  }

  let scrollLeft = videos.scrollLeft + videoBox.offsetWidth
  videos.scrollTo({
    left: scrollLeft,
    behavior: 'smooth'
  })

  if (scrollLeft > 0) {
    videoBox.querySelector('.move_left_btn').classList.add('active')
  }
}
const setVideoMove = function () {
  document.querySelectorAll('.move_left_btn').forEach((element) => {
    element.addEventListener('click', moveLeft)
  })
  document.querySelectorAll('.move_right_btn').forEach((element) => {
    element.classList.add('active')
    element.addEventListener('click', moveRight)
  })
}
const setVideoLoading = function () {
  for (const row of document.querySelectorAll('.video_link')) {
    row.addEventListener('click', function () {
      showLoading([row])
    })
  }
}
// 顶部加载进度
const startProgress = function () {
  if (document.body.clientWidth < 1024) {
    return false
  }
  NProgress.start()
}
const doneProgress = function () {
  NProgress.done()
}
const setProgress = async function () {
  if (document.body.clientWidth < 1024) {
    return false
  }
  NProgress.configure({ showSpinner: false })
  NProgress.start()
  window.addEventListener('load', function () {
    NProgress.done()
  })
  await timeout(2000)
  NProgress.done()
}

const showMask = function (data) {
  let conf = {
    title: '提示',
    content: '',
    ok: '确认',
    cancle: '关闭',
    style: '',
    callback_init: function () {},
    callback_ok: function () {},
    callback_cancle: function () {}
  }
  conf = Object.assign(conf, data)

  const html = `<div class="mask_box" style="display: none;" id="mask_box">
  <div class="mask"></div>
  <div class="mask_panel" style="${conf.style}">
    <div class="mask_body">
      <div class="mask_title">${conf.title}</div>
      <div class="mask_content">${conf.content}</div>
    </div>
    <div class="mask_footer">
      <div class="mask_btn">${conf.ok}</div>
      <div class="mask_btn_cancle">${conf.cancle}</div>
    </div>
  </div>
</div>`
  if (!document.getElementById('mask_box')) {
    const mask = document.createElement('div')
    mask.innerHTML = html
    document.querySelector('body').append(mask)
    // mask 事件
    mask.querySelector('.mask_box .mask_btn').addEventListener('click', (e) => {
      if (conf.callback_ok()) {
        hideMask(e.target)
      }
    })
    mask.querySelector('.mask_box .mask_btn_cancle').addEventListener('click', (e) => {
      hideMask(e.target, conf.callback_cancle)
    })
    // 初始化回调
    conf.callback_init(mask)
  }
  const el = document.getElementById('mask_box')
  el.style.display = 'block'
}

const hideMask = function (e, callback) {
  callback = callback || function () {}
  let el = e ? e.closest('.mask_box') : document.querySelector('.mask_box')
  el.style.display = 'none'
  callback()
}

const setEles = function (dom, query, callback) {
  dom = dom || document
  callback = callback || function () {}
  dom.querySelectorAll(query).forEach((res) => {
    callback(res)
  })
}

const getParams = function (url) {
  return new URLSearchParams(url || window.location.search)
}

const setInputBox = function (mask) {
  // 设置假input事件
  setEles(mask, '.input_box ._input_tel', (ele) => {
    // 输入文字时展示到方块处
    ele.addEventListener('input', (e) => {
      const block = mask.querySelector('.input_box .input_block.active')
      ele.focus()

      // 输入数字操作
      if (['1', '2', '3', '4', '5', '6', '7', '8', '9', '0'].indexOf(e.data) === -1) {
        return false
      }
      setEles(mask, '.input_box .input_block.active', (ele) => {
        ele.innerHTML = e.data
      })
      const next = block.nextElementSibling

      if (!next || !next.classList.contains('input_block')) {
        return false
      }
      setEles(mask, '.input_box .input_block', (ele) => {
        ele.classList.remove('active')
      })
      next.classList.add('active')
    })

    // 删除操作
    ele.addEventListener('keyup', (e) => {
      const block = mask.querySelector('.input_box .input_block.active')
      ele.focus()

      // 删除操作
      if (e.key == 'Backspace') {
        // 清空方块内容
        setEles(mask, '.input_box .input_block.active', (ele) => {
          ele.innerHTML = ''
        })
        const next = block.previousElementSibling
        if (!next || !next.classList.contains('input_block')) {
          return false
        }
        setEles(mask, '.input_box .input_block', (ele) => {
          ele.classList.remove('active')
        })
        next.classList.add('active')
      }
    })
  })
  // 方块点击事件
  setEles(mask, '.input_box .input_block', (ele) => {
    ele.addEventListener('click', (e) => {
      const curr_block = e.currentTarget
      const _input = mask.querySelector('.input_box ._input_tel')
      // 闪动光标
      setEles(mask, '.input_box .input_block', (ele) => {
        ele.classList.remove('active')
      })
      curr_block.classList.add('active')
      // 唤起输入法
      _input.focus()
    })
  })
}

const ajax = (options) => {
  const option = {
    method: 'POST',
    data: {}
  }

  const xhr = new XMLHttpRequest()
  // 取得配置项
  let opt = Object.assign(option, options)
  let res = { success: false, msg: '请求异常' }
  return new Promise((resolve, reject) => {
    if (!opt.url) {
      res.msg = 'url 不能为空'
      resolve(res)
    }
    xhr.onreadystatechange = function () {
      if (xhr.readyState == XMLHttpRequest.DONE) {
        if (xhr.status !== 200) {
          resolve(res)
        }
        let data = JSON.parse(xhr.responseText)
        if (typeof data.success != 'undefined') {
          resolve(data)
        } else {
          resolve(res)
        }
      }
    }
    xhr.open(opt.method, opt.url)
    xhr.send(getRequestBody(opt))
  })
}
// 拼接请求参数
const getRequestBody = (opt) => {
  let databody = ''
  if (opt.method.toUpperCase() == 'POST') {
    if (opt.json) {
      xhr.setRequestHeader('Content-Type', 'application/json')
      databody = JSON.stringify(opt.data)
      console.log(databody)
      return databody
    }
    let databody_arr = []
    xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded')
    for (const [k, v] of Object.entries(opt.data)) {
      if (typeof v == 'object') {
        for (let [vk, vv] of Object.entries(v)) {
          databody_arr.push(`${k}[${vk}]=${vv}`)
        }
      } else {
        databody_arr.push(`${k}=${v}`)
      }
    }
    databody = databody_arr.join('&')
  }
  return databody
}

export {
  setFancybox,
  setToggle,
  loadComments,
  delEmptyThtd,
  setMenu,
  setSearch,
  searchHandler,
  setLazyload,
  setVideoMove,
  showLoading,
  hideLoading,
  setVideoLoading,
  timeout,
  setProgress,
  startProgress,
  doneProgress,
  showMask,
  setEles,
  setInputBox,
  getParams,
  ajax
}
