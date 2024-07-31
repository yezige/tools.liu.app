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

/**
 * 下载文件带进度封装
 * @param {url: string} param0 
 * @param {*} callback 
 * @returns 
 */
const fetchFileWithProgress = async ({ url }, callback) => {
  callback = callback ? callback : function () {}
  if (url instanceof URL || typeof url === 'string') {
  } else {
    throw new Error('不支持的url类型')
  }

  // 定义
  let response = await fetch(url)
  const reader = response.body.getReader()
  const contentLength = +response.headers.get('Content-Length')
  if (!contentLength) {
    throw new Error('计算长度错误')
  }
  console.log('content-length', contentLength)

  // 读取数据，显示进度
  let receivedLength = 0 // 当前接收到了这么多字节
  let chunks = [] // 接收到的二进制块的数组（包括 body）
  try {
    while (true) {
      // 当最后一块下载完成时，done 值为 true
      // value 是块字节的 Uint8Array
      const { done, value } = await reader.read()

      if (done) {
        break
      }

      chunks.push(value)
      receivedLength += value.length

      callback({ total: contentLength, received: receivedLength })
    }
  } catch (e) {}

  // 将块连接到单个 Uint8Array
  let chunksAll = new Uint8Array(receivedLength)
  let position = 0
  for (let chunk of chunks) {
    chunksAll.set(chunk, position)
    position += chunk.length
  }

  return chunksAll
}

export { fetchHead, fetchFileWithProgress }
