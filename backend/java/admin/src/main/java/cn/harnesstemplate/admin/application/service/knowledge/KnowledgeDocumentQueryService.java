package cn.harnesstemplate.admin.application.service.knowledge;

import cn.harnesstemplate.admin.domain.entity.KnowledgeDocument;
import cn.harnesstemplate.admin.infrastructure.persistence.EasyQuerySupport;
import cn.harnesstemplate.admin.infrastructure.persistence.EntityProxies.KnowledgeDocumentProxy;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.util.List;

@Service
public class KnowledgeDocumentQueryService {

    private final EasyQuerySupport easyQuery;

    public KnowledgeDocumentQueryService(EasyQuerySupport easyQuery) { this.easyQuery = easyQuery; }

    public List<KnowledgeDocument> findByCollectionId(Long collectionId) {
        return easyQuery.queryable(KnowledgeDocumentProxy.createTable())
                .where(document -> {
                    document.collectionId().eq(collectionId);
                    document.deletedAt().eq(0L);
                })
                .orderBy(document -> document.id().asc())
                .toList();
    }

    public KnowledgeDocument findById(Long id) {
        return easyQuery.queryable(KnowledgeDocumentProxy.createTable())
                .where(document -> {
                    document.id().eq(id);
                    document.deletedAt().eq(0L);
                })
                .firstOrNull();
    }

    public KnowledgeDocument save(KnowledgeDocument d) {
        return easyQuery.insert(d);
    }

    public void update(KnowledgeDocument d) {
        easyQuery.expressionUpdatable(KnowledgeDocumentProxy.createTable())
                .setColumns(document -> {
                    document.title().set(d.getTitle());
                    document.content().set(d.getContent());
                    document.contentType().set(d.getContentType());
                    document.vectorStatus().set(d.getVectorStatus());
                    document.vectorId().set(d.getVectorId());
                    document.metadata().set(d.getMetadata());
                    document.indexingError().set(d.getIndexingError());
                    document.lastIndexedAt().set(d.getLastIndexedAt());
                    document.isEnabled().set(d.getIsEnabled());
                    document.remark().set(d.getRemark());
                    document.updatedAt().set(LocalDateTime.now());
                })
                .where(document -> {
                    document.id().eq(d.getId());
                    document.deletedAt().eq(0L);
                })
                .executeRows();
    }

    public void updateIndexingStatus(Long id, String vectorStatus, String vectorId, String error) {
        easyQuery.expressionUpdatable(KnowledgeDocumentProxy.createTable())
                .setColumns(document -> {
                    document.vectorStatus().set(vectorStatus);
                    document.vectorId().set(vectorId);
                    document.indexingError().set(error);
                    document.lastIndexedAt().set(System.currentTimeMillis());
                    document.updatedAt().set(LocalDateTime.now());
                })
                .where(document -> document.id().eq(id))
                .executeRows();
    }

    public void delete(Long id) {
        easyQuery.expressionUpdatable(KnowledgeDocumentProxy.createTable())
                .setColumns(document -> document.deletedAt().set(System.currentTimeMillis()))
                .where(document -> {
                    document.id().eq(id);
                    document.deletedAt().eq(0L);
                })
                .executeRows();
    }
}
