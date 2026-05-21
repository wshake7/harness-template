import type { MenuDataItem, ProSettings } from '@ant-design/pro-components'
import type { FormListFieldData } from 'antd'
import type { ComponentType, CSSProperties } from 'react'
import type { ResourceMenuNode } from '~/api/business/sysResourceMenu'

import type { ChangePwdFormValues } from '~/components/business/account/changePwdModal'
import { CloseCircleOutlined, CloseOutlined, ExportOutlined, LockOutlined, LogoutOutlined, ReloadOutlined } from '@ant-design/icons'

import {
  PageContainer,
  ProCard,
  ProConfigProvider,
  ProLayout,
  SettingDrawer,
} from '@ant-design/pro-components'

import {
  createFileRoute,
  Outlet,
  useNavigate,
  useRouter,
  useRouterState,
} from '@tanstack/react-router'

import { ConfigProvider, Dropdown } from 'antd'

import useApp from 'antd/es/app/useApp'

import { memo, startTransition, useCallback, useEffect, useMemo, useReducer, useRef, useState } from 'react'
import { ResourceMenuApi } from '~/api/business/sysResourceMenu'
import ChangePwdModal from '~/components/business/account/changePwdModal'
import { TAB_REFRESH_INTERVAL } from '~/config/tabs'
import { useResourceMenuStore } from '~/stores/resourceMenu'
import { renderAntIcon } from '~/utils/antIcons'

export const Route = createFileRoute('/_app')({
  component: AppLayout,
})

interface CachedTabPane {
  lastHiddenAt?: number
  version: number
}

interface CachedTabPaneNavigateAction {
  type: 'NAVIGATE'
  payload: {
    previousPath: string
    pathname: string
    currentMenuTab: { key: string } | null
  }
}

interface CachedTabPaneRemoveAction {
  type: 'REMOVE'
  payload: {
    key: string
  }
}

interface CachedTabPaneRemoveAllAction {
  type: 'REMOVE_ALL'
}

interface CachedTabPaneRefreshAction {
  type: 'REFRESH'
  payload: {
    key: string
  }
}

type CachedTabPaneAction
  = | CachedTabPaneNavigateAction
    | CachedTabPaneRemoveAction
    | CachedTabPaneRemoveAllAction
    | CachedTabPaneRefreshAction

function sortDynamicMenuNode(a: ResourceMenuNode, b: ResourceMenuNode) {
  const orderDiff = (a.sortOrder ?? 0) - (b.sortOrder ?? 0)
  return orderDiff || a.id - b.id
}

function toDynamicMenuItems(nodes: ResourceMenuNode[], router: ReturnType<typeof useRouter>): MenuDataItem[] {
  const routesByPath = router.routesByPath as unknown as Record<string, unknown>

  return nodes
    .toSorted(sortDynamicMenuNode)
    .filter(node => !node.hidden)
    .map((node) => {
      const children = toDynamicMenuItems(node.children ?? [], router)
      const isUrl = node.isUrl || /^https?:\/\//.test(node.path)
      const routeExists = isUrl || Boolean(routesByPath[node.path])
      if (!routeExists && children.length === 0) {
        return null
      }

      return {
        path: node.path,
        name: node.name,
        icon: renderAntIcon(node.icon),
        isUrl,
        menuType: node.menuType,
        routes: children.length > 0 ? children : undefined,
      } as MenuDataItem
    })
    .filter(Boolean) as MenuDataItem[]
}

function flattenMenuItems(items: MenuDataItem[]) {
  const map = new Map<string, MenuDataItem>()
  const walk = (nodes: MenuDataItem[]) => {
    for (const node of nodes) {
      if (node.path) {
        map.set(node.path, node)
      }
      const routes = (node as any).routes as MenuDataItem[] | undefined
      if (routes?.length) {
        walk(routes)
      }
    }
  }
  walk(items)
  return map
}

function cachedTabPaneReducer(
  state: Record<string, CachedTabPane>,
  action: CachedTabPaneAction,
): Record<string, CachedTabPane> {
  if (action.type === 'NAVIGATE') {
    const { previousPath, pathname, currentMenuTab } = action.payload
    const now = Date.now()
    const next = { ...state }

    if (previousPath !== pathname && next[previousPath]) {
      next[previousPath] = {
        ...next[previousPath],
        lastHiddenAt: now,
      }
    }

    if (currentMenuTab) {
      const currentPane = next[currentMenuTab.key]
      if (!currentPane) {
        next[currentMenuTab.key] = { version: 0 }
      }
      else if (currentPane.lastHiddenAt && now - currentPane.lastHiddenAt > TAB_REFRESH_INTERVAL) {
        next[currentMenuTab.key] = {
          version: currentPane.version + 1,
        }
      }
      else {
        next[currentMenuTab.key] = {
          ...currentPane,
          lastHiddenAt: undefined,
        }
      }
    }

    return next
  }

  if (action.type === 'REMOVE' && action.payload.key) {
    const next = { ...state }
    delete next[action.payload.key]
    return next
  }

  if (action.type === 'REMOVE_ALL') {
    return {}
  }

  if (action.type === 'REFRESH' && action.payload.key) {
    const pane = state[action.payload.key] ?? { version: 0 }
    return {
      ...state,
      [action.payload.key]: {
        version: pane.version + 1,
      },
    }
  }

  return state
}

const CachedTabPaneContent = memo(({ path }: { path: string }) => {
  const router = useRouter()
  const routesByPath = router.routesByPath as unknown as Record<string, { options: { component?: ComponentType } }>
  const route = routesByPath[path]
  const Component = route?.options.component as ComponentType | undefined

  if (!Component) {
    return <Outlet />
  }

  return <Component />
})

function getCachedTabPaneStyle(active: boolean): CSSProperties {
  return {
    inset: active ? undefined : 0,
    opacity: active ? 1 : 0,
    pointerEvents: active ? 'auto' : 'none',
    position: active ? 'relative' : 'absolute',
    transform: active ? 'translate3d(0, 0, 0)' : 'translate3d(0, 8px, 0)',
    transition: 'opacity 180ms ease, transform 220ms cubic-bezier(0.22, 1, 0.36, 1)',
    width: '100%',
    willChange: 'opacity, transform',
    zIndex: active ? 1 : 0,
  }
}

function AppLayout() {
  const [config, setConfig] = useState<Partial<ProSettings>>({
    fixSiderbar: true,
    layout: 'mix',
    fixedHeader: true,
    menu: {
      defaultOpenAll: false,
    },
  })

  const app = useApp()
  const navigate = useNavigate()
  const router = useRouter()
  const accountToken = useAccountStore(state => state.token)

  const pathname = useRouterState({
    select: s => s.location.pathname,
  })
  const [interactivePathname, setInteractivePathname] = useReducer((_state: string, nextPathname: string) => nextPathname, pathname)

  const items = useMenuTabsStore(state => state.items)
  const add = useMenuTabsStore(state => state.add)
  const remove = useMenuTabsStore(state => state.remove)
  const removeAll = useMenuTabsStore(state => state.removeAll)
  const dynamicMenuTree = useResourceMenuStore(state => state.dynamicMenuTree)
  const setDynamicMenuTree = useResourceMenuStore(state => state.setDynamicMenuTree)
  const menuItems = useMemo(
    () => dynamicMenuTree.length > 0 ? toDynamicMenuItems(dynamicMenuTree, router) : [],
    [dynamicMenuTree, router],
  )
  const menuItemMap = useMemo(() => flattenMenuItems(menuItems), [menuItems])
  const previousPathRef = useRef(pathname)
  const [cachedTabPanes, dispatch] = useReducer(cachedTabPaneReducer, {})

  useEffect(() => {
    setInteractivePathname(pathname)
  }, [pathname])

  useEffect(() => {
    if (!accountToken) {
      setDynamicMenuTree([])
      return
    }
    let disposed = false
    ResourceMenuApi.menuTree()
      .send()
      .then((res) => {
        if (disposed) {
          return
        }
        const nodes = res.data ?? []
        setDynamicMenuTree(nodes)
      })
      .catch(() => {})
    return () => {
      disposed = true
    }
  }, [accountToken, setDynamicMenuTree])

  const currentMenuTab = useMemo(() => {
    const dynamicMenu = menuItemMap.get(pathname) as (MenuDataItem & { menuType?: ResourceMenuNode['menuType'] }) | undefined

    if (!dynamicMenu || dynamicMenu.menuType === 'CATALOG') { return null }

    return {
      key: pathname,
      label: dynamicMenu.name,
      icon: dynamicMenu.icon,
    }
  }, [menuItemMap, pathname])

  useEffect(() => {
    if (!currentMenuTab) { return }
    if (!items.some(item => item.key === currentMenuTab.key)) {
      add(currentMenuTab)
    }
  }, [currentMenuTab, items, add])

  useEffect(() => {
    dispatch({
      type: 'NAVIGATE',
      payload: {
        previousPath: previousPathRef.current,
        pathname,
        currentMenuTab,
      },
    })

    previousPathRef.current = pathname
  }, [currentMenuTab, pathname])

  const renderedTabs = useMemo(() => {
    const tabMap = new Map(items.map(item => [item.key, item]))
    if (currentMenuTab) {
      tabMap.set(currentMenuTab.key, currentMenuTab)
    }
    return [...tabMap.values()]
  }, [currentMenuTab, items])

  const closeTab = useCallback((key: string) => {
    dispatch({
      type: 'REMOVE',
      payload: { key },
    })
    remove(key)
  }, [remove])

  const closeAllTabs = useCallback(() => {
    setInteractivePathname('/')
    dispatch({ type: 'REMOVE_ALL' })
    removeAll()
    startTransition(() => {
      navigate({ to: '/' })
    })
  }, [navigate, removeAll])

  const refreshTab = useCallback((key: string) => {
    dispatch({
      type: 'REFRESH',
      payload: { key },
    })
  }, [])

  const openTabInNewWindow = useCallback((key: string) => {
    const url = new URL(key, window.location.origin)
    window.open(url.toString(), '_blank', 'noopener,noreferrer')
  }, [])

  const tabList = useMemo(() => {
    return renderedTabs.map(item => ({
      ...item,
      label: (
        <Dropdown
          trigger={['contextMenu']}
          menu={{
            items: [
              {
                key: 'refresh',
                icon: <ReloadOutlined />,
                label: '刷新',
              },
              {
                key: 'close',
                icon: <CloseOutlined />,
                label: '关闭',
              },
              {
                key: 'closeAll',
                icon: <CloseCircleOutlined />,
                label: '全部关闭',
              },
              {
                key: 'openNewWindow',
                icon: <ExportOutlined />,
                label: '新窗口打开',
              },
            ],
            onClick: ({ key, domEvent }) => {
              domEvent.stopPropagation()
              if (key === 'openNewWindow') {
                openTabInNewWindow(item.key)
              }
              else if (key === 'refresh') {
                refreshTab(item.key)
              }
              else if (key === 'close') {
                closeTab(item.key)
              }
              else if (key === 'closeAll') {
                closeAllTabs()
              }
            },
          }}
        >
          <span onContextMenu={event => event.stopPropagation()}>{item.label}</span>
        </Dropdown>
      ),
    }))
  }, [closeAllTabs, closeTab, openTabInNewWindow, refreshTab, renderedTabs])

  /** ---------------- handlers ---------------- */

  function onMenuClick(item: MenuDataItem & { isUrl: boolean, onClick: () => void }) {
    if (!item.path) {
      return
    }
    if (item.isUrl || /^https?:\/\//.test(item.path)) {
      window.open(item.path, '_blank', 'noopener,noreferrer')
      return
    }
    if (!(router.routesByPath as unknown as Record<string, unknown>)[item.path]) {
      return
    }
    setInteractivePathname(item.path)
    startTransition(() => {
      navigate({ to: item.path })
    })
  }

  function onTabClick(key: string) {
    setInteractivePathname(key)
    startTransition(() => {
      navigate({ to: key })
    })
  }

  async function submitChangePwd(values?: ChangePwdFormValues, error?: FormListFieldData) {
    if (error || !values) {
      app.message.error('\u63D0\u4EA4\u5931\u8D25,\u8BF7\u68C0\u67E5\u8F93\u5165')
      return false
    }

    try {
      await AccountApi.changePwd({
        oldPwd: values.oldPwd,
        newPwd: values.newPwd,
      })
    }
    catch {
      return false
    }
    app.message.success('娣囶喗鏁肩€靛棛鐖滈幋鎰')
    return true
  }

  if (typeof document === 'undefined') {
    return <div />
  }

  return (
    <div
      id="test-pro-layout"
      style={{
        height: '100vh',
        overflow: 'auto',
      }}
    >
      <ProConfigProvider hashed={false}>
        <ConfigProvider
          getTargetContainer={() => document.getElementById('test-pro-layout') || document.body}
        >
          <ProLayout
            route={{ routes: menuItems }}
            location={{ pathname: interactivePathname }}
            token={{
              header: {
                // colorBgMenuItemSelected: 'rgba(0,0,0,0.04)',
              },
            }}
            menu={{
              collapsedShowGroupTitle: true,
            }}
            title="Wshake"
            menuItemRender={(item, dom) => (
              <div
                className="w-full"
                onClick={() => {
                  onMenuClick(item)
                }}
              >
                {dom}
              </div>
            )}
            avatarProps={{
              src: 'https://gw.alipayobjects.com/zos/antfincdn/efFD%24IOql2/weixintupian_20170331104822.jpg',
              size: 'small',
              title: '\u4E03\u5996\u5996',
              render: (_, dom) => (
                <Dropdown
                  menu={{
                    items: [
                      {
                        key: 'changePwd',
                        label: (
                          <ChangePwdModal username="\u4E03\u5996\u5996" onSubmit={submitChangePwd}>
                            <div>
                              <LockOutlined className="ant-dropdown-menu-item-icon" />
                              <span>{'\u4FEE\u6539\u5BC6\u7801'}</span>
                            </div>
                          </ChangePwdModal>
                        ),
                      },
                      {
                        key: 'logout',
                        icon: <LogoutOutlined />,
                        label: '\u9000\u51FA\u767B\u5F55',
                        onClick: async () => {
                          AccountApi.logout()
                        },
                      },
                    ],
                  }}
                >
                  {dom}
                </Dropdown>
              ),
            }}
            {...config}
          >
            <PageContainer
              fixedHeader
              title={false}
              tabList={tabList}
              tabProps={{
                activeKey: interactivePathname,
                hideAdd: true,
                onChange: onTabClick,
                tabBarStyle: {
                  margin: 0,
                },
                onEdit(e, action) {
                  if (action === 'remove') {
                    if (typeof e === 'string') {
                      closeTab(e)
                    }
                  }
                },
                type: 'editable-card',
              }}
            >
              <ProCard style={{ minHeight: 1000, overflow: 'hidden' }}>
                <div style={{ position: 'relative' }}>
                  {renderedTabs.length > 0
                    ? renderedTabs.map((item) => {
                        const pane = cachedTabPanes[item.key] ?? { version: 0 }
                        const active = item.key === interactivePathname

                        return (
                          <div
                            aria-hidden={!active}
                            key={`${item.key}:${pane.version}`}
                            style={getCachedTabPaneStyle(active)}
                          >
                            <CachedTabPaneContent path={item.key} />
                          </div>
                        )
                      })
                    : <Outlet />}
                </div>
              </ProCard>
            </PageContainer>

            <SettingDrawer
              pathname={pathname}
              enableDarkTheme
              settings={config}
              onSettingChange={setConfig}
            />
          </ProLayout>
        </ConfigProvider>
      </ProConfigProvider>
    </div>
  )
}
