/**
 * HEAD 封装
 * @example await fetchHead({url})
 * @param { {url:string, opt:*} } param0
 */
const fetchHead = async ({ url, opt = {} }) => {
  const option = Object.assign(
    {
      method: 'HEAD',
      redirect: 'follow'
    },
    opt
  )
  console.log('fetchHead', url)
  const response = await fetch(url, option)
  console.log('fetchHead', response.status)
  response.headers.forEach((v, k) => {
    console.log('fetchHead', k, v)
  })
  return response
}

const fetchFileWithProgress = async (file, callback) => {
  callback = callback ? callback : function () {}
  if (file instanceof URL || typeof file === 'string') {
  } else {
    return new Uint8Array()
  }

  // 计算打下
  let head = await fetchHead({ url: downUrl })
  if (head.status != 200) {
    return new Uint8Array()
  }
  const contentLength = ~~head.headers.get('content-length')
  if (!contentLength) {
    return new Uint8Array()
  }
  console.log('content-length', contentLength)

  // 定义
  let response = await fetch(file)
  const reader = response.body.getReader()

  // 读取数据，显示进度
  let receivedLength = 0 // 当前接收到了这么多字节
  let chunks = [] // 接收到的二进制块的数组（包括 body）
  while (true) {
    // 当最后一块下载完成时，done 值为 true
    // value 是块字节的 Uint8Array
    const { done, value } = await reader.read()

    if (done) {
      break
    }

    chunks.push(value)
    receivedLength += value.length

    console.log(`Received ${receivedLength} of ${contentLength}`)
    callback({ total: contentLength, received: receivedLength })
  }

  // 将块连接到单个 Uint8Array
  let chunksAll = new Uint8Array(receivedLength) // (4.1)
  let position = 0
  for (let chunk of chunks) {
    chunksAll.set(chunk, position) // (4.2)
    position += chunk.length
  }

  return chunksAll
}

export { fetchHead, fetchFileWithProgress }
