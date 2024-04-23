/**
 * 多线程下载
 * 参考原文
 * https://gist.github.com/semlinker/837211c039e6311e1e7629e5ee5f0a42
 */

import { isError, isSuccess } from './util.js'
/**
 * HEAD 封装
 * @example await fetchHead({url})
 * @param { {url:string, opt:*} } param0
 */
const fetchHead = async ({ url, opt = {} }) => {
  const option = Object.assign(
    {
      method: 'HEAD'
    },
    opt
  )
  console.log('fetchHead', url)
  const response = await fetch(url, option)
  response.headers.forEach((v, k) => {
    console.log('fetchHead', k, v)
  })
  return response
}

/**
 * Download 封装
 * @example await fetchDownload({url})
 * @param { {url: string, name: string, chunkSize: int, poolLimit: int} } param0 {链接，文件名唯一标识，块大小，线程池大小}
 * @returns { {success: boolean} }
 */
const fetchDownload = async ({ url, name, chunkSize, poolLimit = 1 }) => {
  try {
    // 计算长度
    const head = await fetchHead({ url })
    if (head.success) {
    }
    const contentLength = ~~head.headers.get('content-length')
    if (!contentLength) {
      throw new Error('文件长度为空')
    }
    console.log('length', contentLength)

    // 拆分的块数
    const chunks = typeof chunkSize === 'number' ? Math.ceil(contentLength / chunkSize) : 1
    const results = await asyncPool(poolLimit, [...new Array(chunks).keys()], async (i) => {
      let start = i * chunkSize
      let end = i + 1 == chunks ? contentLength - 1 : (i + 1) * chunkSize - 1
      try {
        return getPartContent(url, start, end, i)
      } catch (err) {
        throw new Error(err)
      }
    })

    // 组装
    return isSuccess('成功', assemble(results))
  } catch (err) {
    console.error(err)
    return isError(err)
  }
}

/**
 * 多线程执行
 * @param {int} poolLimit 线程池限制
 * @param {array} array 块个数
 * @param {function} iteratorFn 回调方法
 * @returns
 */
const asyncPool = async (poolLimit, array, iteratorFn) => {
  const ret = [] // 存储所有的异步任务
  const executing = [] // 存储正在执行的异步任务
  for (const item of array) {
    // 调用iteratorFn函数创建异步任务
    const p = Promise.resolve().then(() => iteratorFn(item, array))
    ret.push(p) // 保存新的异步任务

    // 当poolLimit值小于或等于总任务个数时，进行并发控制
    if (poolLimit <= array.length) {
      // 当任务完成后，从正在执行的任务数组中移除已完成的任务
      const e = p.then(() => executing.splice(executing.indexOf(e), 1))
      executing.push(e) // 保存正在执行的异步任务
      if (executing.length >= poolLimit) {
        await Promise.race(executing) // 等待较快的任务执行完成
      }
    }
  }
  return Promise.all(ret)
}

/**
 * 获取分段
 * @param {string} url 链接
 * @param {int} start 开始位置
 * @param {int} end 结束位置
 * @param {int} i 序号
 * @returns
 */
const getPartContent = (url, start, end, i) => {
  return new Promise((resolve, reject) => {
    try {
      console.log(`download-start-${i + 1}`)
      fetch(url, {
        headers: {
          range: `bytes=${start}-${end}`
        }
      }).then(async (res) => {
        if (res.status != 206) {
          throw new Error(`取得块 ${i + 1} 异常`)
        }
        console.log(`download-end-${i + 1}`)
        resolve(new Uint8Array(await res.blob()))
      })
    } catch (err) {
      reject(new Error(err))
    }
  })
}

/**
 * 组装
 * @param {Uint8Array[]} arrays
 * @returns {Blob}
 */
const assemble = (arrays) => {
  if (!arrays.length) return null
  let totalLength = arrays.reduce((acc, value) => acc + value.length, 0)
  let result = new Uint8Array(totalLength)
  let length = 0
  for (let array of arrays) {
    result.set(array, length)
    length += array.length
  }
  return new Blob([result])
}

const getURL = (url) => {
  return new URL(url)
}

export { fetchDownload }
