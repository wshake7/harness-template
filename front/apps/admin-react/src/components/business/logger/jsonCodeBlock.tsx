import { useEffect, useMemo, useState } from 'react'

interface JsonHighlighter {
  codeToHtml: (code: string, options: { lang: string, theme: string }) => string
}

interface HighlightedHtml {
  source: string
  html: string
}

let jsonHighlighterPromise: Promise<JsonHighlighter> | undefined

function formatJsonContent(value?: string) {
  const content = value?.trim()
  if (!content) {
    return ''
  }

  try {
    return JSON.stringify(JSON.parse(content), null, 2)
  }
  catch {
    return content
  }
}

function getJsonHighlighter() {
  jsonHighlighterPromise ??= Promise
    .all([
      import('shiki/core'),
      import('shiki/engine/javascript'),
      import('shiki/langs/json.mjs'),
      import('shiki/themes'),
    ])
    .then(async ([{ createHighlighterCore }, { createJavaScriptRegexEngine }, json, themes]) => {
      const githubLight = await themes.bundledThemes['github-light']()
      const highlighter = await createHighlighterCore({
        themes: [githubLight.default],
        langs: [json.default],
        engine: createJavaScriptRegexEngine({ forgiving: true }),
      })
      return {
        codeToHtml: (code: string, options: { lang: string, theme: string }) => highlighter.codeToHtml(code, options),
      }
    })

  return jsonHighlighterPromise
}

function PlainCodeBlock({ children }: { children: string }) {
  return (
    <pre
      style={{
        margin: 0,
        padding: 12,
        whiteSpace: 'pre-wrap',
        overflowWrap: 'anywhere',
        fontSize: 12,
        lineHeight: 1.6,
      }}
    >
      {children}
    </pre>
  )
}

export function JsonCodeBlock({ value }: { value?: string }) {
  const formattedValue = useMemo(() => formatJsonContent(value), [value])
  const [highlighted, setHighlighted] = useState<HighlightedHtml | null>(null)

  useEffect(() => {
    if (!formattedValue) {
      return
    }

    let disposed = false
    getJsonHighlighter()
      .then(highlighter => highlighter.codeToHtml(formattedValue, {
        lang: 'json',
        theme: 'github-light',
      }))
      .then((html) => {
        if (!disposed) {
          setHighlighted({ source: formattedValue, html })
        }
      })
      .catch(() => {
        if (!disposed) {
          setHighlighted(null)
        }
      })

    return () => {
      disposed = true
    }
  }, [formattedValue])

  if (!formattedValue) {
    return '-'
  }

  const highlightedHtml = highlighted?.source === formattedValue ? highlighted.html : ''

  return (
    <div
      style={{
        maxHeight: 240,
        maxWidth: '100%',
        overflow: 'auto',
        border: '1px solid var(--ant-color-border-secondary)',
        borderRadius: 6,
        background: 'var(--ant-color-bg-layout)',
      }}
    >
      <style>
        {`
          .api-log-json-code pre {
            margin: 0 !important;
            padding: 12px !important;
            white-space: pre-wrap !important;
            overflow-wrap: anywhere !important;
            background: transparent !important;
            font-size: 12px !important;
            line-height: 1.6 !important;
          }
          .api-log-json-code code {
            white-space: pre-wrap !important;
          }
        `}
      </style>
      {highlightedHtml
        ? (
            <div
              className="api-log-json-code"
              // Shiki owns this HTML. Callers only pass raw text, never trusted or prebuilt HTML.
              // eslint-disable-next-line react-dom/no-dangerously-set-innerhtml
              dangerouslySetInnerHTML={{ __html: highlightedHtml }}
            />
          )
        : <PlainCodeBlock>{formattedValue}</PlainCodeBlock>}
    </div>
  )
}
