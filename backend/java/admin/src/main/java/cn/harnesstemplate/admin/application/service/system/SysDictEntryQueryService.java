package cn.harnesstemplate.admin.application.service.system;

import cn.harnesstemplate.admin.domain.entity.SysDictEntry;
import cn.harnesstemplate.admin.infrastructure.persistence.EasyQuerySupport;
import cn.harnesstemplate.admin.infrastructure.persistence.EntityProxies.SysDictEntryProxy;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.util.Collections;
import java.util.List;

@Service
public class SysDictEntryQueryService {

    private final EasyQuerySupport easyQuery;

    public SysDictEntryQueryService(EasyQuerySupport easyQuery) { this.easyQuery = easyQuery; }

    public List<SysDictEntry> findByTypeId(Long typeId) {
        return easyQuery.queryable(SysDictEntryProxy.createTable())
                .where(entry -> {
                    entry.sysDictTypeId().eq(typeId);
                    entry.deletedAt().eq(0L);
                })
                .orderBy(entry -> {
                    entry.sortOrder().asc();
                    entry.id().asc();
                })
                .toList();
    }

    public List<SysDictEntry> findByTypeIds(List<Long> typeIds) {
        if (typeIds == null || typeIds.isEmpty()) return Collections.emptyList();
        return easyQuery.queryable(SysDictEntryProxy.createTable())
                .where(entry -> {
                    entry.sysDictTypeId().in(typeIds);
                    entry.isEnabled().eq(true);
                    entry.deletedAt().eq(0L);
                })
                .orderBy(entry -> {
                    entry.sortOrder().asc();
                    entry.id().asc();
                })
                .toList();
    }

    public SysDictEntry findById(Long id) {
        return easyQuery.queryable(SysDictEntryProxy.createTable())
                .where(entry -> {
                    entry.id().eq(id);
                    entry.deletedAt().eq(0L);
                })
                .firstOrNull();
    }

    public SysDictEntry save(SysDictEntry e) {
        return easyQuery.insert(e);
    }

    public void update(SysDictEntry e) {
        easyQuery.expressionUpdatable(SysDictEntryProxy.createTable())
                .setColumns(entry -> {
                    entry.entryLabel().set(e.getEntryLabel());
                    entry.entryValue().set(e.getEntryValue());
                    entry.labelComponent().set(e.getLabelComponent());
                    entry.languageCode().set(e.getLanguageCode());
                    entry.isEnabled().set(e.getIsEnabled());
                    entry.sortOrder().set(e.getSortOrder());
                    entry.remark().set(e.getRemark());
                    entry.updatedAt().set(LocalDateTime.now());
                })
                .where(entry -> {
                    entry.id().eq(e.getId());
                    entry.deletedAt().eq(0L);
                })
                .executeRows();
    }

    public void delete(Long id) {
        easyQuery.expressionUpdatable(SysDictEntryProxy.createTable())
                .setColumns(entry -> entry.deletedAt().set(System.currentTimeMillis()))
                .where(entry -> {
                    entry.id().eq(id);
                    entry.deletedAt().eq(0L);
                })
                .executeRows();
    }
}
