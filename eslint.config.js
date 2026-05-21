import antfu from '@antfu/eslint-config'
import { tanstackConfig } from '@tanstack/eslint-config'

export default antfu({
  ...tanstackConfig,
  typescript: true,
  formatters: true,
  react: true,
  isInEditor: false,
  lessOpinionated: true,
  rules: {
    'pnpm/yaml-no-unused-catalog-item': 'off',
    'yaml/quotes': 'off',
    'no-console': 'off',
    'eslint-comments/no-unlimited-disable': 'off',
    'react-refresh/only-export-components': 'off',
    'unused-imports/no-unused-vars': 'off',
    'style/max-statements-per-line': 'off',
    'e18e/prefer-static-regex': 'off',
  },
  ignores: [
    '.github',
    '.vitepress/dist',
    '.vitepress/cache',
    'node_modules',
    'public',
    '**/*.d.ts',
    '.eslint-config-inspector',
    '.vscode',
    '.husky',
    '**/*.md',
  ],
})
