const setFancybox = function () {
  if (!CONFIG.site.fancybox) {
    return false
  }
  document.querySelectorAll('.content :not(a) > img, .content > img').forEach((element) => {
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

const loadComments = function (element, callback) {
  if (!CONFIG.site.disqus.lazyload || !element) {
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
    var $th = $(element)
    if ($th.html() == '') {
      $th.remove()
    }
  })
}

const setMenu = function () {
  $('.menu a').show().each(function () {
    if ($(this).attr('href') == location.pathname) {
      $(this).addClass('active')
    }
  })
}

const searchHandler = function () {
  const inputvalue = document.getElementById('search').value
  if (!inputvalue) return false

  location.href = '/youtube/search?q=' + inputvalue
}

const setSearch = function () {
  document.getElementById('search').addEventListener('keydown', (event) => {
    if (event.key === 'Enter') {
      searchHandler()
    }
  })
  document.querySelector('.search_submit').addEventListener('click', searchHandler)
}

// lozad.js 需要手动开启，lazysizes.js 不需要
const setLazyload = function () {
  lozad(document.querySelectorAll('img.lazyload')).observe()
}

const moveLeft = function (event) {
  const videoBox = event.currentTarget.closest('.popular_videos_box')
  const videos = videoBox.querySelector('.videos')
  if (!videoBox){
    return false
  }

  videos.scrollLeft -= videoBox.offsetWidth

  if (videos.scrollLeft <= 0) {
    videoBox.querySelector('.move_left_btn').classList.remove('active')
  }
}
const moveRight = function (event) {
  const videoBox = event.currentTarget.closest('.popular_videos_box')
  const videos = videoBox.querySelector('.videos')
  if (!videoBox){
    return false
  }

  videos.scrollLeft += videoBox.offsetWidth

  if (videos.scrollLeft > 0) {
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
export { setFancybox, setToggle, loadComments, delEmptyThtd, setMenu, setSearch, searchHandler, setLazyload, setVideoMove }
