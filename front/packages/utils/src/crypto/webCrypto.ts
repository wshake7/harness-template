export async function generateAesKey() {
  const key = await window.crypto.subtle.generateKey(
    { name: 'AES-GCM', length: 256 },
    true,
    ['encrypt', 'decrypt'],
  )
  const keyRaw = await window.crypto.subtle.exportKey('raw', key)
  return {
    key,
    keyBase64: arrayBufferToBase64(keyRaw),
  }
}

export async function aesEncrypt(
  key: CryptoKey,
  aad: string,
  data?: any,
) {
  const encoder = new TextEncoder()
  const iv = window.crypto.getRandomValues(new Uint8Array(12))

  const plainBytes = data
    ? encoder.encode(JSON.stringify(data))
    : new Uint8Array()

  const encryptedBuffer = await crypto.subtle.encrypt(
    {
      name: 'AES-GCM',
      iv,
      additionalData: encoder.encode(aad),
      tagLength: 128,
    },
    key,
    plainBytes,
  )

  const encrypted = new Uint8Array(encryptedBuffer)
  const tagLength = 16

  if (encrypted.length < tagLength) {
    throw new Error('encrypted data too short')
  }

  const ciphertext = encrypted.slice(0, encrypted.length - tagLength)
  const tag = encrypted.slice(encrypted.length - tagLength)
  const tagIv = new Uint8Array(tag.length + iv.length)
  tagIv.set(tag, 0)
  tagIv.set(iv, tag.length)

  const combined = new Uint8Array(encrypted.length + iv.length)
  combined.set(encrypted, 0)
  combined.set(iv, encrypted.length)

  const toBase64 = (buf: Uint8Array) => btoa(String.fromCharCode(...buf))

  return {
    CiphertextRaw: ciphertext,
    Ciphertext: toBase64(ciphertext),
    TagIvRaw: tagIv,
    TagIv: toBase64(tagIv),
    CombinedRaw: combined,
    Combined: toBase64(combined),
  }
}

export async function aesDecrypt(
  combinedBase64: string,
  key: CryptoKey,
  aad: string,
) {
  const encoder = new TextEncoder()
  const data = new Uint8Array(base64ToArrayBuffer(combinedBase64))
  const ivLength = 12
  const tagLength = 16

  if (data.length < ivLength + tagLength) {
    throw new Error('invalid data length')
  }

  const iv = data.slice(data.length - ivLength)
  const sealed = data.slice(0, data.length - ivLength)

  const decrypted = await crypto.subtle.decrypt(
    {
      name: 'AES-GCM',
      iv,
      additionalData: encoder.encode(aad),
      tagLength: 128,
    },
    key,
    sealed,
  )

  return new TextDecoder().decode(decrypted)
}

export async function aesDecryptCiphertextAndTag(
  ciphertextBase64: string,
  tagIvBase64: string,
  key: CryptoKey,
  aad: string,
) {
  const encoder = new TextEncoder()
  const ciphertext = new Uint8Array(base64ToArrayBuffer(ciphertextBase64))
  const tagIv = new Uint8Array(base64ToArrayBuffer(tagIvBase64))
  const tagLength = 16
  const ivLength = 12

  if (tagIv.length !== tagLength + ivLength) {
    throw new Error('tagIv length invalid')
  }

  const tag = tagIv.slice(0, tagLength)
  const iv = tagIv.slice(tagLength)
  const sealed = new Uint8Array(ciphertext.length + tagLength)
  sealed.set(ciphertext, 0)
  sealed.set(tag, ciphertext.length)

  const decrypted = await crypto.subtle.decrypt(
    {
      name: 'AES-GCM',
      iv,
      additionalData: encoder.encode(aad),
      tagLength: 128,
    },
    key,
    sealed,
  )

  return new TextDecoder().decode(decrypted)
}

export async function rsaEncrypt(data: string, key: CryptoKey) {
  const encoder = new TextEncoder()
  const encrypted = await window.crypto.subtle.encrypt(
    { name: 'RSA-OAEP' },
    key,
    encoder.encode(data),
  )
  return arrayBufferToBase64(encrypted)
}

export function base64ToArrayBuffer(base64: string) {
  const binaryString = window.atob(base64)
  const bytes = new Uint8Array(binaryString.length)
  for (let i = 0; i < binaryString.length; i++) {
    bytes[i] = binaryString.charCodeAt(i)
  }
  return bytes.buffer
}

export function arrayBufferToBase64(buffer: ArrayBuffer | Uint8Array) {
  const bytes = new Uint8Array(buffer)
  let binary = ''
  for (let i = 0; i < bytes.byteLength; i++) {
    binary += String.fromCharCode(bytes[i])
  }
  return window.btoa(binary)
}

export function uriSort(
  obj: Record<string, any>,
  filterFn?: (key: string) => boolean,
): string {
  const fn = filterFn ?? (() => true)

  const keys = Object.keys(obj)
    .filter((key) => {
      const value = obj[key]
      return fn(key) && value !== '' && value !== undefined && value !== null
    })
    .sort()

  return keys
    .map(key => `${key}=${String(obj[key])}`)
    .join('&')
}
