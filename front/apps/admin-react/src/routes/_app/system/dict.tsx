import { createFileRoute } from '@tanstack/react-router'
import { Splitter } from 'antd'

export const Route = createFileRoute('/_app/system/dict')({
  staleTime: 1000 * 60 * 2,
  component: DictManagement,
})

function DictManagement() {
  const [selectedType, setSelectedType] = useState<DictType>()
  const [selectedEntryIds, setSelectedEntryIds] = useState<number[]>([])
  const [refreshKey, setRefreshKey] = useState(0)

  const handleSelectType = useCallback((record: DictType) => {
    setSelectedType(record)
    setSelectedEntryIds([])
  }, [])

  const handleDeleteSelectedType = useCallback(() => {
    setSelectedType(undefined)
    setSelectedEntryIds([])
  }, [])

  const handleUpdateSelectedType = useCallback((record: DictType) => {
    setSelectedType(record)
  }, [])

  const handleBatchCopyEntries = useCallback(async (entryIds: number[], targetTypeId: number) => {
    try {
      await DictApi.entryBatchCopy({ entryIds, targetTypeId })
      gMessage.success(`成功复制 ${entryIds.length} 项`)
      setSelectedEntryIds([])
      setRefreshKey(k => k + 1)
    }
    catch {
      gMessage.error('复制失败')
    }
  }, [])

  return (
    <Splitter>
      <Splitter.Panel defaultSize="40%" min="25%" max="75%">
        <DictTypePanel
          selectedType={selectedType}
          onSelectType={handleSelectType}
          onDeleteSelectedType={handleDeleteSelectedType}
          onUpdateSelectedType={handleUpdateSelectedType}
          onBatchCopyEntries={handleBatchCopyEntries}
        />
      </Splitter.Panel>
      <Splitter.Panel>
        <DictEntryPanel
          selectedType={selectedType}
          selectedEntryIds={selectedEntryIds}
          onSelectionChange={setSelectedEntryIds}
          onClearType={handleDeleteSelectedType}
          refreshKey={refreshKey}
        />
      </Splitter.Panel>
    </Splitter>
  )
}
