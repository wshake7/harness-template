import { createFileRoute } from '@tanstack/react-router'
import { Splitter } from 'antd'

export const Route = createFileRoute('/_app/system/language')({
  staleTime: 1000 * 60 * 2,
  component: LanguageManagement,
})

function LanguageManagement() {
  const [selectedType, setSelectedType] = useState<LanguageType>()
  const [refreshKey] = useState(0)

  const handleSelectType = useCallback((record: LanguageType) => {
    setSelectedType(record)
  }, [])

  const handleDeleteSelectedType = useCallback(() => {
    setSelectedType(undefined)
  }, [])

  const handleUpdateSelectedType = useCallback((record: LanguageType) => {
    setSelectedType(record)
  }, [])

  const handleClearSelectedType = useCallback(() => {
    setSelectedType(undefined)
  }, [])

  return (
    <Splitter>
      <Splitter.Panel defaultSize="50%" min="25%" max="60%">
        <LanguageTypePanel
          selectedType={selectedType}
          onSelectType={handleSelectType}
          onDeleteSelectedType={handleDeleteSelectedType}
          onUpdateSelectedType={handleUpdateSelectedType}
        />
      </Splitter.Panel>
      <Splitter.Panel>
        <LanguageEntryPanel
          selectedType={selectedType}
          refreshKey={refreshKey}
          onClearType={handleClearSelectedType}
        />
      </Splitter.Panel>
    </Splitter>
  )
}
