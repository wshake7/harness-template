package cn.harnesstemplate.admin.application.service.system;

import cn.harnesstemplate.admin.domain.entity.SysLanguageType;
import cn.harnesstemplate.admin.infrastructure.persistence.EasyQuerySupport;
import cn.harnesstemplate.admin.infrastructure.persistence.EntityProxies.SysLanguageTypeProxy;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.util.List;

@Service
public class SysLanguageTypeQueryService {

    private final EasyQuerySupport easyQuery;

    public SysLanguageTypeQueryService(EasyQuerySupport easyQuery) { this.easyQuery = easyQuery; }

    public List<SysLanguageType> listAll() {
        return easyQuery.queryable(SysLanguageTypeProxy.createTable())
                .where(type -> type.deletedAt().eq(0L))
                .orderBy(type -> {
                    type.sortOrder().asc();
                    type.id().asc();
                })
                .toList();
    }

    public SysLanguageType findById(Long id) {
        return easyQuery.queryable(SysLanguageTypeProxy.createTable())
                .where(type -> {
                    type.id().eq(id);
                    type.deletedAt().eq(0L);
                })
                .firstOrNull();
    }

    public SysLanguageType save(SysLanguageType t) {
        return easyQuery.insert(t);
    }

    public void update(SysLanguageType t) {
        easyQuery.expressionUpdatable(SysLanguageTypeProxy.createTable())
                .setColumns(type -> {
                    type.typeCode().set(t.getTypeCode());
                    type.typeName().set(t.getTypeName());
                    type.isDefault().set(t.getIsDefault());
                    type.isEnabled().set(t.getIsEnabled());
                    type.sortOrder().set(t.getSortOrder());
                    type.updatedAt().set(LocalDateTime.now());
                })
                .where(type -> {
                    type.id().eq(t.getId());
                    type.deletedAt().eq(0L);
                })
                .executeRows();
    }

    public void delete(Long id) {
        easyQuery.expressionUpdatable(SysLanguageTypeProxy.createTable())
                .setColumns(type -> type.deletedAt().set(System.currentTimeMillis()))
                .where(type -> {
                    type.id().eq(id);
                    type.deletedAt().eq(0L);
                })
                .executeRows();
    }
}
