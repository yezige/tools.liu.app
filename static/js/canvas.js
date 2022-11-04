// From https://codepen.io/Mamboleoo/pen/obWGYr

const dark = window.matchMedia('(prefers-color-scheme: dark)').matches,
  dpr = 1, // window.devicePixelRatio,
  canvas = document.createElement('canvas'),
  ctx = canvas.getContext('2d'),
  canvas_box = document.querySelector('header .logo'),
  canvas_xy = canvas.getBoundingClientRect(),
  font_size = 1.5,
  background = dark ? '#121212' : '#f5f5f5', // 背景色
  colors = [
    '#d81b60', // 红
    '#f4511e', // 橙
    '#1976d2', // 蓝
    '#ffb300', // 黄
    '#6d4c41', // 咖
  ], // 颜色列表
  text = CONFIG.page.site.logo || 'Icons',
  width = 120, // canvas 宽度
  height = 43, // canvas 高度
  mouse_diameter = Math.round(height / 3), // 鼠标直径大小
  speed = 100 // 移动速度

let particles = [],
  amount = 0,
  mouse = { x: 0, y: 0 },
  radius = dpr * 2,
  max_radius = dpr * 2.5

var ww = (canvas.width = width)
var wh = (canvas.height = height)

function Particle(x, y) {
  this.x = Math.random() * ww
  this.y = Math.random() * wh
  this.dest = {
    x: x,
    y: y,
  }

  // 从radius至max_radius之间随机一个球大小
  const r = parseFloat((max_radius * Math.random()).toFixed(1))
  this.r = r < radius ? radius: r

  this.vx = (Math.random() - 0.5) * 20
  this.vy = (Math.random() - 0.5) * 20
  this.accX = 0
  this.accY = 0
  this.friction = Math.random() * 0.05 + 0.94

  this.color = colors[Math.floor(Math.random() * 6)]
}
// 400 为移动速度
Particle.prototype.render = function() {
  this.accX = (this.dest.x - this.x) / speed
  this.accY = (this.dest.y - this.y) / speed
  this.vx += this.accX
  this.vy += this.accY
  this.vx *= this.friction
  this.vy *= this.friction

  this.x += this.vx
  this.y += this.vy

  ctx.fillStyle = this.color
  ctx.beginPath()
  ctx.arc(this.x, this.y, this.r, Math.PI * 2, false)
  ctx.fill()

  var a = this.x - mouse.x
  var b = this.y - mouse.y

  // a*a + b*b = c*c
  var distance = Math.sqrt(a * a + b * b)
  if (distance < radius * mouse_diameter) {
    this.accX = (this.x - mouse.x) / 50
    this.accY = (this.y - mouse.y) / 50
    this.vx += this.accX
    this.vy += this.accY
  }
}
function initScene() {
  if ($ && $(canvas_box).is(':hidden')) {
    return false
  }
  mouse = { x: -999, y: -999 }
  canvas_box.innerHTML = '' // 可以不清空
  canvas_box.appendChild(canvas)
  canvas.style.letterSpacing = '3px'
  ctx.scale(dpr, dpr)

  ww = canvas.width = width
  wh = canvas.height = height

  ctx.clearRect(0, 0, canvas.width, canvas.height)

  ctx.font = `${font_size}em Roboto`
  // ctx.textAlign = "center";
  ctx.fillText(text, canvas.width / 2 - ctx.measureText(text).width / 2, canvas.height / 2 + font_size * 16 / 2 - 3)

  var data = ctx.getImageData(0, 0, ww, wh).data
  ctx.clearRect(0, 0, canvas.width, canvas.height)
  ctx.globalCompositeOperation = 'screen'

  particles = []
  // 定义密集程度Math.round(90/50) = 2，过大的分辨率建议结果=10，目前写死三像素的步进，当宽过大时Math.round(ww/50)可能会=4，导致球数量不够
  for (var i = 0; i < ww; i += 3) {
    for (var j = 0; j < wh; j += 3) {
      if (data[(i + j * ww) * 4 + 3] > 50) {
        particles.push(new Particle(i, j))
      }
    }
  }
  amount = particles.length
}
function onMouseMove(e) {
  let x = e.offsetX // clientX是相对于窗口的位置 screenX是相对于屏幕的位置
  let y = e.offsetY
  // 防止停在边缘
  if (x >= width - 15 || x <= 15) {
    x = -999
  }
  if (y >= height - 15 || y <= 15) {
    y = -999
  }
  mouse.x = x
  mouse.y = y
}

function onTouchMove(e) {
  if (e.touches.length > 0) {
    mouse.x = e.touches[0].offsetX
    mouse.y = e.touches[0].offsetY
  }
}

function onTouchEnd(e) {
  mouse.x = -9999
  mouse.y = -9999
}

function onMouseClick() {
  radius++
  if (radius === 5) {
    radius = 0
  }
}

function render(a) {
  requestAnimationFrame(render)
  ctx.clearRect(0, 0, canvas.width, canvas.height)
  for (var i = 0; i < amount; i++) {
    particles[i].render()
  }
}

// window.addEventListener('resize', initScene)
canvas.addEventListener('mousemove', onMouseMove, {passive: true})
canvas.addEventListener('touchmove', onTouchMove, {passive: true})
canvas.addEventListener('click', onMouseClick, {passive: true})
canvas.addEventListener('touchend', onTouchEnd, {passive: true})

if (false && FontFace) {
  // 字体加载完后再渲染
  const canvasFont = new FontFace('Inconsolata', 'url(./static/font/Inconsolata.woff2)')
  canvasFont
    .load()
    .then((font) => {
      document.fonts.add(font)
    })
    .then(() => {
      initScene()
    })
} else {
  initScene()
}
requestAnimationFrame(render)
