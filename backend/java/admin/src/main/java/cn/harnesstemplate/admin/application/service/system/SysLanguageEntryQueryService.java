package cn.harnesstemplate.admin.application.service.system;

import cn.harnesstemplate.admin.domain.entity.SysLanguageEntry;
import cn.harnesstemplate.admin.infrastructure.persistence.EasyQuerySupport;
import cn.harnesstemplate.admin.infrastructure.persistence.EntityProxies.SysLanguageEntryProxy;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.util.List;

@Service
public class SysLanguageEntryQueryService {

    private final EasyQuerySupport easyQuery;

    public SysLanguageEntryQueryService(EasyQuerySupport easyQuery) { this.easyQuery = easyQuery; }

    public List<SysLanguageEntry> findByTypeId(Long typeId) {
        return easyQuery.queryable(SysLanguageEntryProxy.createTable())
                .where(entry -> {
                    entry.sysLanguageTypeId().eq(typeId);
                    entry.deletedAt().eq(0L);
                })
                .orderBy(entry -> {
                    entry.sortOrder().asc();
                    entry.id().asc();
                })
                .toList();
    }

    public SysLanguageEntry findById(Long id) {
        return easyQuery.queryable(SysLanguageEntryProxy.createTable())
                .where(entry -> {
                    entry.id().eq(id);
                    entry.deletedAt().eq(0L);
                })
                .firstOrNull();
    }

    public SysLanguageEntry save(SysLanguageEntry e) {
        return easyQuery.insert(e);
    }

    public void update(SysLanguageEntry e) {
        easyQuery.expressionUpdatable(SysLanguageEntryProxy.createTable())
                .setColumns(entry -> {
                    entry.entryCode().set(e.getEntryCode());
                    entry.entryValue().set(e.getEntryValue());
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
        easyQuery.expressionUpdatable(SysLanguageEntryProxy.createTable())
                .setColumns(entry -> entry.deletedAt().set(System.currentTimeMillis()))
                .where(entry -> {
                    entry.id().eq(id);
                    entry.deletedAt().eq(0L);
                })
                .executeRows();
    }
}
