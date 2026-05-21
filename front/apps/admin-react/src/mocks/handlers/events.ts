import { http, HttpResponse } from 'msw'
import { url } from '.'

export const eventsHandlers = [
  http.get(url('/api/events'), async () => {
    const encoder = new TextEncoder()
    let count = 0

    const stream = new ReadableStream({
      start(controller) {
        const interval = setInterval(() => {
          count++
          const data = JSON.stringify({ count })
          controller.enqueue(encoder.encode(`event: count\ndata: ${data}\n\n`))

          if (count >= 100) {
            clearInterval(interval)
            controller.close()
          }
        }, 30000)
      },
    })

    return new HttpResponse(stream, {
      headers: {
        'Content-Type': 'text/event-stream',
        'Cache-Control': 'no-cache',
        'Connection': 'keep-alive',
      },
    })
  }),
]
