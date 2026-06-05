package cn.harnesstemplate.admin.application.service.knowledge;

import cn.harnesstemplate.admin.domain.entity.KnowledgeCollection;
import cn.harnesstemplate.admin.infrastructure.persistence.EasyQuerySupport;
import cn.harnesstemplate.admin.infrastructure.persistence.EntityProxies.KnowledgeCollectionProxy;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.util.List;

@Service
public class KnowledgeCollectionQueryService {

    private final EasyQuerySupport easyQuery;

    public KnowledgeCollectionQueryService(EasyQuerySupport easyQuery) { this.easyQuery = easyQuery; }

    public List<KnowledgeCollection> listAll() {
        return easyQuery.queryable(KnowledgeCollectionProxy.createTable())
                .where(collection -> collection.deletedAt().eq(0L))
                .orderBy(collection -> collection.id().asc())
                .toList();
    }

    public KnowledgeCollection findById(Long id) {
        return easyQuery.queryable(KnowledgeCollectionProxy.createTable())
                .where(collection -> {
                    collection.id().eq(id);
                    collection.deletedAt().eq(0L);
                })
                .firstOrNull();
    }

    public KnowledgeCollection save(KnowledgeCollection c) {
        return easyQuery.insert(c);
    }

    public void update(KnowledgeCollection c) {
        easyQuery.expressionUpdatable(KnowledgeCollectionProxy.createTable())
                .setColumns(collection -> {
                    collection.collectionName().set(c.getCollectionName());
                    collection.displayName().set(c.getDisplayName());
                    collection.metricType().set(c.getMetricType());
                    collection.indexType().set(c.getIndexType());
                    collection.isEnabled().set(c.getIsEnabled());
                    collection.remark().set(c.getRemark());
                    collection.updatedAt().set(LocalDateTime.now());
                })
                .where(collection -> {
                    collection.id().eq(c.getId());
                    collection.deletedAt().eq(0L);
                })
                .executeRows();
    }

    public void delete(Long id) {
        easyQuery.expressionUpdatable(KnowledgeCollectionProxy.createTable())
                .setColumns(collection -> collection.deletedAt().set(System.currentTimeMillis()))
                .where(collection -> {
                    collection.id().eq(id);
                    collection.deletedAt().eq(0L);
                })
                .executeRows();
    }
}
