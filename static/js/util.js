const timeout = (ms) => {
  return new Promise((resolve, reject) => {
    setTimeout(() => {
      resolve()
    }, ms)
  })
}
const setFancybox = function() {
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

const setToggle = function() {
  const menuToggle = document.querySelector('.toggle')
  const menu = document.querySelector('.menu')
  const header = document.querySelector('.showcase header')

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

const loadComments = function(element, callback) {
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

const delEmptyThtd = function() {
  document.querySelectorAll('thead th').forEach((element) => {
    var $th = $(element)
    if ($th.html() == '') {
      $th.remove()
    }
  })
}

const setMenu = function() {
  $('.menu a')
    .show()
    .each(function() {
      if ($(this).attr('href') == location.pathname) {
        $(this).addClass('active')
      }
    })
}

const searchHandler = function() {
  const inputvalue = document.getElementById('search').value
  if (!inputvalue) return false

  location.href = '/youtube/search?q=' + inputvalue
}

// 搜索事件
const setSearch = function() {
  document.getElementById('search').addEventListener('keydown', (event) => {
    if (event.key === 'Enter') {
      searchHandler()
    }
  })
  document.querySelector('.search_submit').addEventListener('click', searchHandler)
}

// 显示加载中
const showLoading = function(els) {
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
const hideLoading = function(els) {
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
const setLazyload = async function() {
  await timeout(200)
  lozad(document.querySelectorAll('.lazyload'), {
    loaded: function(el) {
      if (!el.closest('.loading_box')) {
        return
      }
      el.closest('.loading_box')
        .querySelector('.loading')
        .remove()
    }
  }).observe()
}

const moveLeft = function(event) {
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
const moveRight = function(event) {
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
const setVideoMove = function() {
  document.querySelectorAll('.move_left_btn').forEach((element) => {
    element.addEventListener('click', moveLeft)
  })
  document.querySelectorAll('.move_right_btn').forEach((element) => {
    element.classList.add('active')
    element.addEventListener('click', moveRight)
  })
}
const setVideoLoading = function() {
  for (const row of document.querySelectorAll('.video_link')) {
    row.addEventListener('click', function() {
      showLoading([row])
    })
  }
}
// 顶部加载进度
const startProgress = function() {
  if (document.body.clientWidth < 1024) {
    return false
  }
  NProgress.start()
}
const doneProgress = function() {
  NProgress.done()
}
const setProgress = async function() {
  if (document.body.clientWidth < 1024) {
    return false
  }
  NProgress.configure({ showSpinner: false })
  NProgress.start()
  window.addEventListener('load', function() {
    NProgress.done()
  })
  await timeout(2000)
  NProgress.done()
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
  doneProgress
}
