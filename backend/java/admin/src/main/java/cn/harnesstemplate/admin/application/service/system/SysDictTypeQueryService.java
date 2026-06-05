package cn.harnesstemplate.admin.application.service.system;

import cn.harnesstemplate.admin.domain.entity.SysDictType;
import cn.harnesstemplate.admin.infrastructure.persistence.EasyQuerySupport;
import cn.harnesstemplate.admin.infrastructure.persistence.EntityProxies.SysDictTypeProxy;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.util.Collections;
import java.util.List;

@Service
public class SysDictTypeQueryService {

    private final EasyQuerySupport easyQuery;

    public SysDictTypeQueryService(EasyQuerySupport easyQuery) { this.easyQuery = easyQuery; }

    public List<SysDictType> listAll() {
        return easyQuery.queryable(SysDictTypeProxy.createTable())
                .where(type -> type.deletedAt().eq(0L))
                .orderBy(type -> {
                    type.sortOrder().asc();
                    type.id().asc();
                })
                .toList();
    }

    public SysDictType findById(Long id) {
        return easyQuery.queryable(SysDictTypeProxy.createTable())
                .where(type -> {
                    type.id().eq(id);
                    type.deletedAt().eq(0L);
                })
                .firstOrNull();
    }

    public SysDictType save(SysDictType t) {
        return easyQuery.insert(t);
    }

    public void update(SysDictType t) {
        easyQuery.expressionUpdatable(SysDictTypeProxy.createTable())
                .setColumns(type -> {
                    type.typeCode().set(t.getTypeCode());
                    type.typeName().set(t.getTypeName());
                    type.isEnabled().set(t.getIsEnabled());
                    type.sortOrder().set(t.getSortOrder());
                    type.remark().set(t.getRemark());
                    type.updatedAt().set(LocalDateTime.now());
                })
                .where(type -> {
                    type.id().eq(t.getId());
                    type.deletedAt().eq(0L);
                })
                .executeRows();
    }

    public List<SysDictType> findByCodes(List<String> codes) {
        if (codes == null || codes.isEmpty()) return Collections.emptyList();
        return easyQuery.queryable(SysDictTypeProxy.createTable())
                .where(type -> {
                    type.typeCode().in(codes);
                    type.isEnabled().eq(true);
                    type.deletedAt().eq(0L);
                })
                .toList();
    }

    public void delete(Long id) {
        easyQuery.expressionUpdatable(SysDictTypeProxy.createTable())
                .setColumns(type -> type.deletedAt().set(System.currentTimeMillis()))
                .where(type -> {
                    type.id().eq(id);
                    type.deletedAt().eq(0L);
                })
                .executeRows();
    }
}
