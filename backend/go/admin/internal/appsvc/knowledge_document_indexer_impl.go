package appsvc

import (
	"admin/internal/ai/agent/knowledge_pipeline"
	"admin/internal/ai/indexer"
	"admin/internal/config"
	"admin/internal/fiberc/res"
	"admin/internal/services/milvusc"
	"admin/internal/services/objectstore"
	"admin/internal/services/orm/models"
	"admin/internal/services/orm/query"
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/compose"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
	"gorm.io/datatypes"
)

type knowledgeDocumentIndexer struct {
	q       *query.Query
	conf    *config.Config
	storage objectstore.Engine
	runner  KnowledgeIndexRunner
	cleaner VectorStoreCleaner
	now     func() time.Time
}

func NewKnowledgeDocumentIndexer(q *query.Query, conf *config.Config, storage objectstore.Engine) KnowledgeDocumentIndexer {
	if q == nil {
		q = query.Q
	}
	if storage == nil {
		storage = objectstore.Current()
	}
	if conf == nil {
		conf = &config.Config{}
	}
	return &knowledgeDocumentIndexer{
		q:       q,
		conf:    conf,
		storage: storage,
		runner:  &einoKnowledgeIndexRunner{},
		cleaner: &milvusVectorStoreCleaner{},
		now:     time.Now,
	}
}

func (k *knowledgeDocumentIndexer) ImportFile(ctx context.Context, input ImportKnowledgeFileInput) (*models.KnowledgeDocument, error) {
	collection, err := k.q.KnowledgeCollection.Where(k.q.KnowledgeCollection.ID.Eq(input.CollectionID)).First()
	if err != nil {
		return nil, res.FailMsg("知识库集合不存在")
	}
	if err = validateCollectionEmbedding(collection); err != nil {
		return nil, err
	}

	asset, err := k.q.FileAsset.Where(k.q.FileAsset.ID.Eq(input.FileAssetID)).First()
	if err != nil {
		return nil, res.FailMsg("文件不存在")
	}
	if asset.Status != models.FileAssetStatusActive {
		return nil, res.FailMsg("文件尚未完成上传")
	}
	if !isSupportedKnowledgeFile(asset.Extension, asset.ContentType) {
		return nil, res.FailMsg("暂不支持该文件类型")
	}

	body, err := k.storage.GetObject(ctx, asset.Bucket, asset.ObjectKey)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	data, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	title := strings.TrimSpace(input.Title)
	if title == "" {
		title = strings.TrimSuffix(asset.OriginalName, path.Ext(asset.OriginalName))
	}
	contentType := normalizeDocumentContentType(input.ContentType, asset.ContentType, asset.Extension)
	source := fmt.Sprintf("minio://%s/%s", asset.Bucket, asset.ObjectKey)
	metadata := datatypes.JSONMap{
		"fileAssetID": asset.ID,
		"objectKey":   asset.ObjectKey,
	}

	doc := &models.KnowledgeDocument{
		CollectionID:  input.CollectionID,
		DocumentID:    fmt.Sprintf("file-%d", asset.ID),
		Title:         title,
		Content:       string(data),
		ContentType:   contentType,
		Source:        source,
		VectorStatus:  "pending",
		IsEnabled:     models.KnowledgeDocument{}.IsEnabled,
		Remark:        models.KnowledgeDocument{}.Remark,
		Metadata:      metadata,
		IndexingError: "",
	}
	doc.IsEnabled.IsEnabled = true
	doc.Remark.Remark = input.Remark
	doc.OperatorID.CreatedBy.CreatedBy = input.OperatorID
	doc.OperatorID.UpdatedBy.UpdatedBy = input.OperatorID

	if err = k.q.KnowledgeDocument.Create(doc); err != nil {
		return nil, err
	}
	return k.q.KnowledgeDocument.Where(k.q.KnowledgeDocument.ID.Eq(doc.ID)).First()
}

func (k *knowledgeDocumentIndexer) IndexDocument(ctx context.Context, documentID uint64) error {
	doc, err := k.q.KnowledgeDocument.Where(k.q.KnowledgeDocument.ID.Eq(documentID)).First()
	if err != nil {
		return res.FailMsg("文档不存在")
	}
	collection, err := k.q.KnowledgeCollection.Where(k.q.KnowledgeCollection.ID.Eq(doc.CollectionID)).First()
	if err != nil {
		return res.FailMsg("知识库集合不存在")
	}
	return k.indexDocument(ctx, doc, collection, extensionForContentType(doc.ContentType))
}

func (k *knowledgeDocumentIndexer) DeleteDocumentVectors(ctx context.Context, documentID uint64) error {
	doc, err := k.q.KnowledgeDocument.Where(k.q.KnowledgeDocument.ID.Eq(documentID)).First()
	if err != nil {
		return res.FailMsg("文档不存在")
	}
	collection, err := k.q.KnowledgeCollection.Where(k.q.KnowledgeCollection.ID.Eq(doc.CollectionID)).First()
	if err != nil {
		return res.FailMsg("知识库集合不存在")
	}
	return k.cleaner.DeleteByDocument(ctx, collection.CollectionName, doc.ID, stableDocumentSource(doc))
}

func (k *knowledgeDocumentIndexer) indexDocument(ctx context.Context, doc *models.KnowledgeDocument, collection *models.KnowledgeCollection, extension string) error {
	if err := validateCollectionEmbedding(collection); err != nil {
		return k.markIndexFailed(ctx, doc.ID, err)
	}

	source := stableDocumentSource(doc)
	if doc.Source != source {
		if _, err := k.q.KnowledgeDocument.Where(k.q.KnowledgeDocument.ID.Eq(doc.ID)).UpdateSimple(
			k.q.KnowledgeDocument.Source.Value(source),
		); err != nil {
			return err
		}
	}

	filePath, err := writeTempKnowledgeFile(doc.Content, extension)
	if err != nil {
		return k.markIndexFailed(ctx, doc.ID, err)
	}
	defer os.Remove(filePath)

	if err = k.cleaner.DeleteByDocument(ctx, collection.CollectionName, doc.ID, source); err != nil {
		return err
	}

	metadata := buildKnowledgeIndexMetadata(doc, source)
	ids, err := k.runner.Index(ctx, k.conf, collection, filePath, metadata)
	if err != nil {
		return k.markIndexFailed(ctx, doc.ID, err)
	}

	_, err = k.q.KnowledgeDocument.Where(k.q.KnowledgeDocument.ID.Eq(doc.ID)).UpdateSimple(
		k.q.KnowledgeDocument.VectorStatus.Value("indexed"),
		k.q.KnowledgeDocument.VectorID.Value(strings.Join(ids, ",")),
		k.q.KnowledgeDocument.LastIndexedAt.Value(k.now().UnixMilli()),
		k.q.KnowledgeDocument.IndexingError.Value(""),
	)
	return err
}

func (k *knowledgeDocumentIndexer) markIndexFailed(ctx context.Context, documentID uint64, err error) error {
	_, updateErr := k.q.KnowledgeDocument.Where(k.q.KnowledgeDocument.ID.Eq(documentID)).UpdateSimple(
		k.q.KnowledgeDocument.VectorStatus.Value("failed"),
		k.q.KnowledgeDocument.VectorID.Value(""),
		k.q.KnowledgeDocument.LastIndexedAt.Value(int64(0)),
		k.q.KnowledgeDocument.IndexingError.Value(sanitizeIndexingError(err)),
	)
	if updateErr != nil {
		return updateErr
	}
	return err
}

func buildKnowledgeIndexMetadata(doc *models.KnowledgeDocument, source string) map[string]any {
	metadata := map[string]any{
		"_source":        source,
		"document_db_id": strconv.FormatUint(doc.ID, 10),
		"collection_id":  strconv.FormatUint(doc.CollectionID, 10),
		"document_id":    doc.DocumentID,
		"title":          doc.Title,
		"chunk_index":    strconv.Itoa(doc.ChunkIndex),
	}
	if doc.Metadata != nil {
		if fileAssetID, ok := doc.Metadata["fileAssetID"]; ok {
			metadata["file_asset_id"] = fmt.Sprintf("%v", fileAssetID)
		}
	}
	return metadata
}

func sanitizeIndexingError(err error) string {
	if err == nil {
		return ""
	}
	msg := strings.TrimSpace(err.Error())
	if len(msg) > 512 {
		return msg[:512]
	}
	return msg
}

func normalizeDocumentContentType(input string, assetContentType string, extension string) string {
	if trimmed := strings.TrimSpace(input); trimmed != "" {
		return trimmed
	}
	switch {
	case strings.Contains(assetContentType, "markdown"):
		return "markdown"
	case strings.Contains(assetContentType, "html"):
		return "html"
	default:
		switch strings.ToLower(extension) {
		case ".md", ".markdown":
			return "markdown"
		case ".html", ".htm":
			return "html"
		default:
			return "text"
		}
	}
}

func extensionForContentType(contentType string) string {
	switch strings.ToLower(strings.TrimSpace(contentType)) {
	case "markdown":
		return ".md"
	case "html":
		return ".html"
	default:
		return ".txt"
	}
}

func stableDocumentSource(doc *models.KnowledgeDocument) string {
	if strings.TrimSpace(doc.Source) != "" {
		return doc.Source
	}
	if strings.TrimSpace(doc.DocumentID) != "" {
		return "manual://" + doc.DocumentID
	}
	return fmt.Sprintf("manual://%d", doc.ID)
}

func validateCollectionEmbedding(collection *models.KnowledgeCollection) error {
	if collection == nil {
		return res.FailMsg("知识库集合不存在")
	}
	return nil
}

func isSupportedKnowledgeFile(extension string, contentType string) bool {
	switch strings.ToLower(strings.TrimSpace(extension)) {
	case ".md", ".markdown", ".txt", ".html", ".htm":
		return true
	}
	return strings.HasPrefix(strings.ToLower(strings.TrimSpace(contentType)), "text/")
}

func writeTempKnowledgeFile(content string, extension string) (string, error) {
	pattern := "knowledge-*"
	if extension != "" {
		pattern += extension
	}
	file, err := os.CreateTemp("", pattern)
	if err != nil {
		return "", err
	}
	defer file.Close()
	if _, err = file.WriteString(content); err != nil {
		return "", err
	}
	return file.Name(), nil
}

type einoKnowledgeIndexRunner struct{}

func (r *einoKnowledgeIndexRunner) Index(ctx context.Context, conf *config.Config, collection *models.KnowledgeCollection, sourcePath string, metadata map[string]any) ([]string, error) {
	confCopy := *conf
	confCopy.AI = conf.AI
	confCopy.AI.Embedding = conf.AI.Embedding
	confCopy.AI.Knowledge = conf.AI.Knowledge
	if collection != nil {
		confCopy.AI.Knowledge.Collection = collection.CollectionName
		confCopy.AI.Knowledge.IndexMetricType = collection.MetricType
		confCopy.AI.Knowledge.IndexType = collection.IndexType
	}
	pipeline, err := knowledge_pipeline.BuildKnowledgeIndexingWithOptions(ctx, &confCopy, &knowledge_pipeline.BuildOptions{
		Metadata: metadata,
	})
	if err != nil {
		return nil, err
	}
	return pipeline.Invoke(ctx, document.Source{URI: sourcePath}, compose.WithCallbacks())
}

type milvusVectorStoreCleaner struct{}

func (c *milvusVectorStoreCleaner) DeleteByDocument(ctx context.Context, collectionName string, documentDBID uint64, source string) error {
	cli := milvusc.Client
	if cli == nil || collectionName == "" {
		return nil
	}

	documentExpr := fmt.Sprintf(`metadata["document_db_id"] == "%d"`, documentDBID)
	filter := documentExpr
	if strings.TrimSpace(source) != "" {
		filter = fmt.Sprintf(`(%s) or (metadata["_source"] == %s)`, documentExpr, mustJSONString(source))
	}

	queryResult, err := cli.Query(ctx, milvusclient.NewQueryOption(collectionName).
		WithFilter(filter).
		WithOutputFields("id"))
	if err != nil {
		return err
	}
	if queryResult.Len() == 0 {
		return nil
	}

	idColumn := queryResult.GetColumn("id")
	idsToDelete := make([]string, 0, idColumn.Len())
	for i := 0; i < idColumn.Len(); i++ {
		id, getErr := idColumn.GetAsString(i)
		if getErr == nil && id != "" {
			idsToDelete = append(idsToDelete, id)
		}
	}
	if len(idsToDelete) == 0 {
		return nil
	}

	_, err = cli.Delete(ctx, milvusclient.NewDeleteOption(collectionName).
		WithStringIDs("id", idsToDelete))
	return err
}

func mustJSONString(input string) string {
	raw, _ := sonic.MarshalString(input)
	return raw
}

func (k *knowledgeDocumentIndexer) EnsureCollection(ctx context.Context, collection *models.KnowledgeCollection) error {
	confCopy := *k.conf
	confCopy.AI = k.conf.AI
	confCopy.AI.Embedding = k.conf.AI.Embedding
	confCopy.AI.Knowledge = k.conf.AI.Knowledge

	confCopy.AI.Knowledge.Collection = collection.CollectionName
	if collection.MetricType != "" {
		confCopy.AI.Knowledge.IndexMetricType = collection.MetricType
	}
	if collection.IndexType != "" {
		confCopy.AI.Knowledge.IndexType = collection.IndexType
	}

	_, err := indexer.New(ctx, &confCopy)
	return err
}

func (k *knowledgeDocumentIndexer) DropCollection(ctx context.Context, collectionName string) error {
	cli := milvusc.Client
	if cli == nil || collectionName == "" {
		return nil
	}
	return cli.DropCollection(ctx, milvusclient.NewDropCollectionOption(collectionName))
}
